package permissions

type Roles []Role

func (r Roles) GetIDs() []string {
	var ids []string
	for _, role := range r {
		ids = append(ids, role.ID)
	}
	return ids
}

func (r Roles) StringSlice() []string {
	names := make([]string, len(r))
	for i, role := range r {
		names[i] = role.Name
	}
	return names
}

type Role struct {
	ID   string
	Name string
}

func (r Role) String() string {
	return string(r.Name)
}
