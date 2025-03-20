package permissions

import "context"

type Reader interface {
	GetTenantPermissions(ctx context.Context, resources []string) (TenantPermissions, error)
	GetUserPermissions(ctx context.Context, resources []string) (UserPermissions, error)
	GetUserPermissionsExtraAndRevoked(ctx context.Context, resources []string) (UserExtraPermissions, UserRevokedPermissions, error)
	GetUserResources(ctx context.Context, resources []string) (Resources, error)
	GetUserRoles(ctx context.Context) (Roles, error)
	GetTenantRoles(ctx context.Context) (Roles, error)
	GetTenantRoleMap(ctx context.Context, resources []string) (TenantRoleMap, error)
}

type Writer interface {
	// Nothing yet.
}

type ReaderWriter interface {
	Reader
	Writer
}
