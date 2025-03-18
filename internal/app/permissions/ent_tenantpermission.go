package permissions

type TenantPermissions []TenantPermission

type TenantPermission Permission

func (tp TenantPermission) Format(f string) (string, error) {
	return Permission(tp).Format(f)
}
