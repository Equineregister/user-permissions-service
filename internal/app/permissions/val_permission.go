package permissions

type Permission struct {
	Resource string
	Action   string
}

const (
	FormatPermission = `%s:%s`
)
