package permissions

import "context"

type Reader interface {
	GetTenantPermissionsForUser(ctx context.Context, tenantID, userID string, resources []string) (TenantPermissions, error)
	GetUserPermissions(ctx context.Context, userID string, resources []string) (UserPermissions, error)
	GetUserResources(ctx context.Context, userID string) (Resources, error)
	GetUserRole(ctx context.Context, userID string) (*Role, error)
}

type Writer interface{}

type ReaderWriter interface {
	Reader
	Writer
}
