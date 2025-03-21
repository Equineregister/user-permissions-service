package permissions

type UserExtraPermissions UserPermissions

func (ueps UserExtraPermissions) StringSlice() []string {
	return UserPermissions(ueps).StringSlice()
}

type UserRevokedPermissions UserPermissions

func (urps UserRevokedPermissions) StringSlice() []string {
	return UserPermissions(urps).StringSlice()
}

type UserPermissions []UserPermission

func (ups UserPermissions) StringSlice() []string {
	names := make([]string, len(ups))
	for i, up := range ups {
		names[i] = up.String()
	}
	return names
}

type UserPermission Permission

func (up UserPermission) String() string {
	return Permission(up).String()
}
