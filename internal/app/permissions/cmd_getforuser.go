package permissions

import "context"

type ForUser struct {
	Role              Role
	TenantPermissions TenantPermissions
	UserPermissions   UserPermissions
	Resources         Resources
}

func (s *Service) GetForUser(ctx context.Context, tenantID, userID string, resources []string) (ForUser, error) {
	return ForUser{}, nil
}
