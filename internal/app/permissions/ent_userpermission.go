package permissions

type UserPermissions []UserPermission

type UserPermission Permission

func (up UserPermission) Format(f string) (string, error) {
	return Permission(up).Format(f)
}
