package permissions

import "context"

type Reader interface {
	GetTenantPermissions(ctx context.Context, resources []string) (TenantPermissions, error)
	GetUserPermissions(ctx context.Context, resources []string) (UserPermissions, error)
	GetUserResources(ctx context.Context) (Resources, error)
	GetUserRoles(ctx context.Context) (Roles, error)
}

type Writer interface{}

type ReaderWriter interface {
	Reader
	Writer
}
