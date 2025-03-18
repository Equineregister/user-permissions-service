package permissions

type Roles []Role

func (r Roles) GetIDs() []string {
	var ids []string
	for _, role := range r {
		ids = append(ids, role.ID)
	}
	return ids
}

type Role struct {
	ID       string
	Name     string
	CacheKey string
}

func (r Role) String() string {
	return string(r.Name)
}
