package permissions

type TenantPermissions []TenantPermission

func (tps TenantPermissions) StringSlice() []string {
	names := make([]string, len(tps))
	for i, tp := range tps {
		names[i] = tp.String()
	}
	return names
}

type TenantPermission Permission

func (tp TenantPermission) String() string {
	return Permission(tp).String()
}
