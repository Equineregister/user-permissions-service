package postgres

import (
	"context"
	"fmt"

	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/pkg/config"
	"github.com/jackc/pgx/v5"
)

type PermissionsRepo struct {
	tenantPool *TenantPool
}

func NewPermissionsRepo(cfg *config.Config) *PermissionsRepo {
	tenantPool := NewTenantPool(cfg)

	return &PermissionsRepo{
		tenantPool: tenantPool,
	}
}

func (pr *PermissionsRepo) GetTenantPermissions(ctx context.Context, resources []string) (permissions.TenantPermissions, error) {
	pool, err := pr.tenantPool.GetTenantConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tenant connection: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer rollback(ctx, tx)

	tp, err := pr.getTenantPermissions(ctx, tx, resources)
	if err != nil {
		return nil, fmt.Errorf("get tenant permissions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return tp, nil
}

func (pr *PermissionsRepo) getTenantPermissions(ctx context.Context, tx pgx.Tx, resources []string) (permissions.TenantPermissions, error) {

	rows, err := tx.Query(ctx, `
		SELECT
			p.permission_id, p.permission_name
		FROM 
			tenant_permissions tp
			JOIN permissions p ON tp.permission_id = p.permission_id
		WHERE 
			p.permission_name ILIKE ANY (ARRAY[@permission_names])
		ORDER BY
        	p.permission_name ASC
		`,
		pgx.NamedArgs{
			"permission_names": getPermissionNamesFromResources(resources),
		})
	if err != nil {
		return nil, fmt.Errorf("query tenant_permissions: %w", err)
	}
	defer rows.Close()

	tps := make(permissions.TenantPermissions, 0)
	for rows.Next() {
		var tp permissions.TenantPermission
		if err := rows.Scan(&tp.ID, &tp.Name); err != nil {
			return nil, fmt.Errorf("scan tenant_permissions: %w", err)
		}
		tps = append(tps, tp)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows tenant_permissions: %w", rows.Err())
	}
	return tps, nil
}

func (pr *PermissionsRepo) GetUserPermissions(ctx context.Context, resources []string) (permissions.UserPermissions, error) {
	userID, found := extractUserID(ctx)
	if !found {
		return nil, fmt.Errorf("user ID not found in context")
	}

	pool, err := pr.tenantPool.GetTenantConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tenant connection: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer rollback(ctx, tx)

	roles, err := pr.getUserRoles(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	if len(roles) == 0 {
		return nil, nil
	}

	up, err := pr.getUserPermissionsFromRoles(ctx, tx, roles, resources)
	if err != nil {
		return nil, fmt.Errorf("get user permissions from roles: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return up, nil
}

func (pr *PermissionsRepo) getUserPermissionsFromRoles(ctx context.Context, tx pgx.Tx, userRoles permissions.Roles, resources []string) (permissions.UserPermissions, error) {

	if len(userRoles) == 0 {
		return nil, nil
	}

	rows, err := tx.Query(ctx, `
		SELECT 
            rp.role_id, 
            p.permission_id, 
            p.permission_name
        FROM 
            role_permissions rp
        JOIN 
            permissions p ON rp.permission_id = p.permission_id
        WHERE 
			role_id = ANY(@role_ids)
			AND
			p.permission_name ILIKE ANY (ARRAY[@permission_names])
		`, pgx.NamedArgs{
		"role_ids":         userRoles.GetIDs(),
		"permission_names": getPermissionNamesFromResources(resources),
	})
	if err != nil {
		return nil, fmt.Errorf("query role_permissions: %w", err)
	}
	defer rows.Close()

	var rolePermissions permissions.UserPermissions
	for rows.Next() {
		var up permissions.UserPermission
		var rid string // Role ID - unused.
		if err := rows.Scan(&rid, &up.ID, &up.Name); err != nil {
			return nil, fmt.Errorf("scan role_permissions: %w", err)
		}
		rolePermissions = append(rolePermissions, up)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows role_permissions: %w", rows.Err())
	}

	return rolePermissions, nil

}

func (pr *PermissionsRepo) GetUserResources(ctx context.Context) (permissions.Resources, error) {
	return nil, nil
}

func (pr *PermissionsRepo) GetUserRoles(ctx context.Context) (permissions.Roles, error) {
	userID, found := extractUserID(ctx)
	if !found {
		return nil, fmt.Errorf("user ID not found in context")
	}

	pool, err := pr.tenantPool.GetTenantConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tenant connection: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer rollback(ctx, tx)

	roles, err := pr.getUserRoles(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return roles, nil
}

func (pr *PermissionsRepo) getUserRoles(ctx context.Context, tx pgx.Tx, userID string) (permissions.Roles, error) {
	rows, err := tx.Query(ctx, `
		SELECT 
			ur.role_id, r.role_name, ur.cache_key
		FROM 
			user_roles ur
		JOIN 
			roles r ON ur.role_id = r.role_id
		WHERE 
			ur.user_id = @user_id
		`, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("query user_roles: %w", err)
	}
	defer rows.Close()

	var userRoles permissions.Roles
	for rows.Next() {
		var ur permissions.Role
		if err := rows.Scan(&ur.ID, &ur.Name, &ur.CacheKey); err != nil {
			return nil, fmt.Errorf("scan user_roles: %w", err)
		}
		userRoles = append(userRoles, ur)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows user_roles: %w", rows.Err())
	}

	return userRoles, nil
}

func getPermissionNamesFromResources(resources []string) []string {
	if len(resources) == 0 {
		return []string{"%"} // Match everything
	}
	permissionNames := make([]string, len(resources))
	for i, resource := range resources {
		permissionNames[i] = resource + ":%"
	}
	return permissionNames
}
