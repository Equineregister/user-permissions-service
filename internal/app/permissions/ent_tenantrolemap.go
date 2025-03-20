package permissions

type TenantMappedRole struct {
	Permissions TenantPermissions
	Inherits    Roles
}

type TenantRoleMap map[Role]TenantMappedRole
