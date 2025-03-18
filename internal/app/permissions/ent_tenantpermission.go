package permissions

import "fmt"

type TenantPermissions []TenantPermission

type TenantPermission Permission

func (tp TenantPermission) Format(f string) (string, error) {
	switch f {
	case FormatPermission:
		return fmt.Sprintf(f, tp.Resource, tp.Action), nil
	default:
		return "", fmt.Errorf("invalid format: %s", f)
	}
}
