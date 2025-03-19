package permissions

type TenantPermissions []TenantPermission

func (tps TenantPermissions) StringSlice() []string {
	var names []string
	for _, tp := range tps {
		names = append(names, tp.String())
	}
	return names
}

type TenantPermission Permission

func (tp TenantPermission) String() string {
	return Permission(tp).String()
}
