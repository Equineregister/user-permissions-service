package rego

import (
	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
)

type Role struct {
	Permissions []string `json:"permissions"`
	Inherits    []string `json:"inherits"`
}

type RoleGraph map[string]Role

// NewRoleGraph creates a new OPA ReGo RoleGraph from a TenantRoleMap
func NewRoleGraph(tenantRoleMap permissions.TenantRoleMap) RoleGraph {
	rg := make(RoleGraph)

	for role, mappedRole := range tenantRoleMap {
		rg[role.Name] = Role{
			Permissions: mappedRole.Permissions.StringSlice(),
			Inherits:    mappedRole.Inherits.StringSlice(),
		}
	}

	return rg
}
