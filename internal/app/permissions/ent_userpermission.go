package permissions

import "fmt"

type UserPermissions []UserPermission

type UserPermission Permission

func (up UserPermission) Format(f string) (string, error) {
	switch f {
	case FormatPermission:
		return fmt.Sprintf(f, up.Resource, up.Action), nil
	default:
		return "", fmt.Errorf("invalid format: %s", f)
	}
}
