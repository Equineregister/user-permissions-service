package postgres

import (
	"context"
	"fmt"

	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/internal/config"
	"github.com/Equineregister/user-permissions-service/internal/pkg/contextkey"
	"github.com/jackc/pgx/v5"
)

type PermissionsRepo struct {
	tenantPool *TenantPool
}

// NewPermissionsRepoWithTenantPool creates a new PermissionsRepo from the supplied TenantPool.
func NewPermissionsRepoWithTenantPool(tenantPool *TenantPool) *PermissionsRepo {
	return &PermissionsRepo{
		tenantPool: tenantPool,
	}
}

// NewPermissionsRepo creates a new PermissionsRepo from the supplied config.
func NewPermissionsRepo(cfg *config.Config) *PermissionsRepo {
	return &PermissionsRepo{
		tenantPool: NewTenantPool(cfg),
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
			p.permission_name ILIKE ANY (@permission_names::text[])
		ORDER BY
        	p.permission_name ASC
		`,
		pgx.NamedArgs{
			"permission_names": permissionNamesForResources(resources),
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
	userID, found := contextkey.UserID(ctx)
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

	// We must get the User's roles first, then get the permissions from those roles.
	roles, err := pr.getUserRoles(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
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

func (pr *PermissionsRepo) GetUserPermissionsExtraAndRevoked(ctx context.Context, resources []string) (permissions.UserExtraPermissions, permissions.UserRevokedPermissions, error) {
	userID, found := contextkey.UserID(ctx)
	if !found {
		return nil, nil, fmt.Errorf("user ID not found in context")
	}

	pool, err := pr.tenantPool.GetTenantConnection(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get tenant connection: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin: %w", err)
	}
	defer rollback(ctx, tx)

	// Get extra and revoked permissions
	extra, revoked, err := pr.getUserPermissionsExtraAndRevoked(ctx, tx, userID, resources)
	if err != nil {
		return nil, nil, fmt.Errorf("get user permissions extra and revoked: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit: %w", err)
	}

	return extra, revoked, nil
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
			p.permission_name ILIKE ANY (@permission_names::text[])
		ORDER BY
        	p.permission_name ASC
		`, pgx.NamedArgs{
		"role_ids":         userRoles.GetIDs(),
		"permission_names": permissionNamesForResources(resources),
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

func (pr *PermissionsRepo) getUserPermissionsExtraAndRevoked(ctx context.Context, tx pgx.Tx, userID string, resources []string) (permissions.UserExtraPermissions, permissions.UserRevokedPermissions, error) {

	rows, err := tx.Query(ctx, `
			SELECT 
				up.permission_id, 
				p.permission_name, 
				up.permission_type
			FROM 
				user_permissions up
			JOIN 
				permissions p ON up.permission_id = p.permission_id
			WHERE 
				up.user_id = @user_id
				AND
				p.permission_name ILIKE ANY (@permission_names::text[])
			ORDER BY
        		p.permission_name ASC
		`, pgx.NamedArgs{
		"user_id":          userID,
		"permission_names": permissionNamesForResources(resources),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("query user_permissions: %w", err)
	}
	defer rows.Close()

	extra := make(permissions.UserExtraPermissions, 0)
	revoked := make(permissions.UserRevokedPermissions, 0)
	for rows.Next() {

		var up permissions.UserPermission
		var permissionType string
		if err := rows.Scan(&up.ID, &up.Name, &permissionType); err != nil {
			return nil, nil, fmt.Errorf("scan user_permissions: %w", err)
		}

		switch {
		case permissionType == "extra":
			extra = append(extra, up)
		case permissionType == "revoked":
			revoked = append(revoked, up)
		default:
			return nil, nil, fmt.Errorf("unknown permission type: %s", permissionType)
		}
	}

	if rows.Err() != nil {
		return nil, nil, fmt.Errorf("rows user_permissions: %w", rows.Err())
	}

	return extra, revoked, nil
}

func (pr *PermissionsRepo) GetUserResources(ctx context.Context, resources []string) (permissions.Resources, error) {
	userID, found := contextkey.UserID(ctx)
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

	res, err := pr.getUserResources(ctx, tx, userID, resources)
	if err != nil {
		return nil, fmt.Errorf("get user resources: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return res, nil
}

func (pr *PermissionsRepo) getUserResources(ctx context.Context, tx pgx.Tx, userID string, resources []string) (permissions.Resources, error) {
	if len(resources) == 0 {
		resources = []string{"%"} // Match everything
	}

	rows, err := tx.Query(ctx, `
		SELECT 
			ur.resource_id, ur.resource_type_id, rt.resource_type_name
		FROM 
			user_resources ur
		JOIN 
			resource_types rt ON ur.resource_type_id = rt.resource_type_id
		WHERE 
			ur.user_id = @user_id
			AND
			rt.resource_type_name ILIKE ANY(@resource_types::text[])
		ORDER BY
        	rt.resource_type_name ASC, ur.resource_id ASC
		`, pgx.NamedArgs{
		"user_id":        userID,
		"resource_types": resources,
	})
	if err != nil {
		return nil, fmt.Errorf("query user_resources: %w", err)
	}
	defer rows.Close()

	var userResources permissions.Resources
	for rows.Next() {
		var ur permissions.Resource
		var resourceTypeID int
		if err := rows.Scan(&ur.ID, &resourceTypeID, &ur.Type); err != nil {
			return nil, fmt.Errorf("scan user_resources: %w", err)
		}
		userResources = append(userResources, ur)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows user_resources: %w", rows.Err())
	}

	return userResources, nil
}

func (pr *PermissionsRepo) GetUserRoles(ctx context.Context) (permissions.Roles, error) {
	userID, found := contextkey.UserID(ctx)
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
	directRoles, err := pr.getUserDirectRoles(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get direct roles: %w", err)
	}
	if len(directRoles) == 0 {
		return nil, nil
	}

	var allChildRoles permissions.Roles
	currentRoles := directRoles
	for {
		childRoles, err := pr.getChildRoles(ctx, tx, currentRoles.GetIDs())
		if err != nil {
			return nil, fmt.Errorf("get child roles: %w", err)
		}
		if len(childRoles) == 0 {
			break
		}
		allChildRoles = append(allChildRoles, childRoles...)
		currentRoles = childRoles
	}
	// Maintain the order of the roles, don't sort them.
	return append(directRoles, allChildRoles...), nil
}

func (pr *PermissionsRepo) getUserDirectRoles(ctx context.Context, tx pgx.Tx, userID string) (permissions.Roles, error) {
	rows, err := tx.Query(ctx, `
		SELECT 
			ur.role_id, r.role_name
		FROM 
			user_roles ur
		JOIN 
			roles r ON ur.role_id = r.role_id
		WHERE 
			ur.user_id = @user_id
		ORDER BY
			r.role_name ASC
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
		if err := rows.Scan(&ur.ID, &ur.Name); err != nil {
			return nil, fmt.Errorf("scan user_roles: %w", err)
		}
		userRoles = append(userRoles, ur)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows user_roles: %w", rows.Err())
	}

	return userRoles, nil
}

func (pr *PermissionsRepo) getChildRoles(ctx context.Context, tx pgx.Tx, assignedRoles []string) (permissions.Roles, error) {

	rows, err := tx.Query(ctx, `
		SELECT 
			rh.child_role_id, r.role_name
		FROM 
			role_hierarchy rh
		JOIN 
			roles r ON rh.child_role_id = r.role_id
		WHERE 
			rh.parent_role_id = ANY(@assigned_roles)
		ORDER BY
			r.role_name ASC
		`, pgx.NamedArgs{
		"assigned_roles": assignedRoles,
	})
	if err != nil {
		return nil, fmt.Errorf("query role_hierarchy: %w", err)
	}
	defer rows.Close()

	var childRoles permissions.Roles
	for rows.Next() {
		var ur permissions.Role
		if err := rows.Scan(&ur.ID, &ur.Name); err != nil {
			return nil, fmt.Errorf("scan user_roles: %w", err)
		}
		childRoles = append(childRoles, ur)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows user_roles: %w", rows.Err())
	}

	return childRoles, nil
}

func (pr *PermissionsRepo) GetTenantRoles(ctx context.Context) (permissions.Roles, error) {
	pool, err := pr.tenantPool.GetTenantConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tenant connection: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer rollback(ctx, tx)

	roles, err := pr.getTenantRoles(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("get tenant roles: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return roles, nil
}

func (pr *PermissionsRepo) getTenantRoles(ctx context.Context, tx pgx.Tx) (permissions.Roles, error) {

	allRoles, err := pr.getAllRoles(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("get direct roles: %w", err)
	}
	if len(allRoles) == 0 {
		return nil, nil
	}

	return allRoles, nil
}

func (pr *PermissionsRepo) getAllRoles(ctx context.Context, tx pgx.Tx) (permissions.Roles, error) {
	rows, err := tx.Query(ctx, `
		SELECT 
			r.role_id, r.role_name
		FROM 
			roles r
		ORDER BY
			r.role_name ASC
		`)
	if err != nil {
		return nil, fmt.Errorf("query roles: %w", err)
	}
	defer rows.Close()

	var roles permissions.Roles
	for rows.Next() {
		var r permissions.Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, fmt.Errorf("scan roles: %w", err)
		}
		roles = append(roles, r)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows roles: %w", rows.Err())
	}

	return roles, nil
}

func (pr *PermissionsRepo) GetTenantRoleMap(ctx context.Context, resources []string) (permissions.TenantRoleMap, error) {
	pool, err := pr.tenantPool.GetTenantConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tenant connection: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer rollback(ctx, tx)

	roles, err := pr.getTenantRoles(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("get tenant role map: %w", err)
	}

	rolemap := permissions.TenantRoleMap{}

	for _, role := range roles {
		rolePermissions, err := pr.getRoleTenantPermissions(ctx, tx, role.ID, resources)
		if err != nil {
			return nil, fmt.Errorf("get role permissions: %w", err)
		}

		childRoles, err := pr.getChildRoles(ctx, tx, []string{role.ID})
		if err != nil {
			return nil, fmt.Errorf("get role permissions: %w", err)
		}

		rolemap[role] = permissions.TenantMappedRole{
			Permissions: permissions.TenantPermissions(rolePermissions),
			Inherits:    childRoles,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return rolemap, nil
}

func (pr *PermissionsRepo) getRoleTenantPermissions(ctx context.Context, tx pgx.Tx, roleID string, resources []string) (permissions.TenantPermissions, error) {

	rows, err := tx.Query(ctx, `
		SELECT 
			rp.permission_id, p.permission_name
		FROM 
			role_permissions rp
		JOIN 
			permissions p ON rp.permission_id = p.permission_id
		WHERE 
			rp.role_id = @role_id
			AND
			p.permission_name ILIKE ANY (@permission_names::text[])
		ORDER BY
			p.permission_name ASC
		`, pgx.NamedArgs{
		"role_id":          roleID,
		"permission_names": permissionNamesForResources(resources),
	})
	if err != nil {
		return nil, fmt.Errorf("query role_permissions: %w", err)
	}
	defer rows.Close()

	var tenantPermissions permissions.TenantPermissions
	for rows.Next() {
		var tp permissions.TenantPermission
		if err := rows.Scan(&tp.ID, &tp.Name); err != nil {
			return nil, fmt.Errorf("scan role_permissions: %w", err)
		}
		tenantPermissions = append(tenantPermissions, tp)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows role_permissions: %w", rows.Err())
	}

	return tenantPermissions, nil
}

func permissionNamesForResources(resources []string) []string {
	if len(resources) == 0 {
		return []string{"%"} // Match everything
	}
	permissionNames := make([]string, len(resources))
	for i, resource := range resources {
		permissionNames[i] = resource + ":%"
	}
	return permissionNames
}
