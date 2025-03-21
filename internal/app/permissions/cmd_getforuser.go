package permissions

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type ForUser struct {
	Roles              Roles
	RevokedPermissions UserRevokedPermissions
	ExtraPermissions   UserExtraPermissions
	Resources          Resources
	RoleMap            TenantRoleMap
}

func (s *Service) GetForUser(ctx context.Context, resources []string) (*ForUser, error) {

	eg, ctxEg := errgroup.WithContext(ctx)

	chTenantRoleMap := make(chan TenantRoleMap, 1)
	eg.Go(func() error {
		trm, err := s.repo.GetTenantRoleMap(ctxEg, resources)
		if err != nil {
			return err
		}
		chTenantRoleMap <- trm
		return nil
	})

	chUserExtraPermissions := make(chan UserExtraPermissions, 1)
	chUserRevokedPermissions := make(chan UserRevokedPermissions, 1)
	eg.Go(func() error {
		ex, rv, err := s.repo.GetUserPermissionsExtraAndRevoked(ctxEg, resources)
		if err != nil {
			return err
		}
		chUserExtraPermissions <- ex
		chUserRevokedPermissions <- rv
		return nil
	})

	chRoles := make(chan Roles, 1)
	eg.Go(func() error {
		roles, err := s.repo.GetUserRoles(ctxEg)
		if err != nil {
			return err
		}
		chRoles <- roles
		return nil
	})

	chResources := make(chan Resources, 1)
	eg.Go(func() error {
		resources, err := s.repo.GetUserResources(ctxEg, resources)
		if err != nil {
			return err
		}
		chResources <- resources
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("get for user: %w", err)
	}

	close(chTenantRoleMap)
	close(chUserExtraPermissions)
	close(chUserRevokedPermissions)
	close(chRoles)
	close(chResources)

	return &ForUser{
		ExtraPermissions:   <-chUserExtraPermissions,
		RevokedPermissions: <-chUserRevokedPermissions,
		Resources:          <-chResources,
		Roles:              <-chRoles,
		RoleMap:            <-chTenantRoleMap,
	}, nil
}
