package permissions

type Permission struct {
	Name string
	ID   string
}

func (p Permission) String() string {
	return p.Name
}
