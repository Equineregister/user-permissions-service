package permissions

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type ForUser struct {
	Role              Role
	TenantPermissions TenantPermissions
	UserPermissions   UserPermissions
	Resources         Resources
}

func (s *Service) GetForUser(ctx context.Context, resources []string) (*ForUser, error) {

	eg, ctxEg := errgroup.WithContext(ctx)

	chTenantPermissions := make(chan TenantPermissions, 1)
	eg.Go(func() error {
		tp, err := s.repo.GetTenantPermissions(ctxEg, resources)
		if err != nil {
			return err
		}
		chTenantPermissions <- tp
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("get for user: %w", err)
	}
	close(chTenantPermissions)

	return &ForUser{
		TenantPermissions: <-chTenantPermissions,
	}, nil
}
