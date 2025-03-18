package postgres

import (
	"context"

	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/pkg/config"
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

func (pr *PermissionsRepo) GetTenantPermissionsForUser(ctx context.Context, tenantID, userID string, resources []string) (permissions.TenantPermissions, error) {
	return nil, nil
}

func (pr *PermissionsRepo) GetUserPermissions(ctx context.Context, userID string, resources []string) (permissions.UserPermissions, error) {
	return nil, nil
}

func (pr *PermissionsRepo) GetUserResources(ctx context.Context, userID string) (permissions.Resources, error) {
	return nil, nil
}

func (pr *PermissionsRepo) GetUserRole(ctx context.Context, userID string) (*permissions.Role, error) {
	return nil, nil
}
