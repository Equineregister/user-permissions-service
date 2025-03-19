package permissions

import "context"

type Reader interface {
	GetTenantPermissions(ctx context.Context, resources []string) (TenantPermissions, error)
	GetUserPermissions(ctx context.Context, resources []string) (UserPermissions, error)
	GetUserPermissionsExtraAndRevoked(ctx context.Context, resources []string) (UserExtraPermissions, UserRevokedPermissions, error)
	GetUserResources(ctx context.Context, resources []string) (Resources, error)
	GetUserRoles(ctx context.Context) (Roles, error)
}

type Writer interface{}

type ReaderWriter interface {
	Reader
	Writer
}
