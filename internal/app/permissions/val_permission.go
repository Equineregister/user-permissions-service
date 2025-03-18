package permissions

import "fmt"

type Permission struct {
	Resource string
	Action   string
}

const (
	FormatPermission = `%s:%s`
)

func (p Permission) Format(f string) (string, error) {
	switch f {
	case FormatPermission:
		return fmt.Sprintf(f, p.Resource, p.Action), nil
	default:
		return "", fmt.Errorf("invalid format: %s", f)
	}
}
