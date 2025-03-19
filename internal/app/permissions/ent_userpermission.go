package permissions

type UserExtraPermissions UserPermissions
type UserRevokedPermissions UserPermissions

type UserPermissions []UserPermission

func (ups UserPermissions) StringSlice() []string {
	var names []string
	for _, up := range ups {
		names = append(names, up.String())
	}
	return names
}

type UserPermission Permission

func (up UserPermission) String() string {
	return Permission(up).String()
}
