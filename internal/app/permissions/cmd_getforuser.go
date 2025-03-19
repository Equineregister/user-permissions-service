package permissions

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type ForUser struct {
	Roles             Roles
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

	chUserPermissions := make(chan UserPermissions, 1)
	eg.Go(func() error {
		up, err := s.repo.GetUserPermissions(ctxEg, resources)
		if err != nil {
			return err
		}
		chUserPermissions <- up
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

	close(chTenantPermissions)
	close(chUserPermissions)
	close(chUserExtraPermissions)
	close(chUserRevokedPermissions)
	close(chRoles)
	close(chResources)

	// Combine the user permissions with the extra permissions and remove the revoked permissions
	// This is done in the service layer to keep the repository layer simple.
	up := <-chUserPermissions
	extra := <-chUserExtraPermissions
	revoked := <-chUserRevokedPermissions

	// Add extra permissions
	up = append(up, extra...)

	// Remove revoked permissions. Do this after adding the extra permissions to ensure that the revoked permissions are removed.

	for _, r := range revoked {
		for i, p := range up {
			if p.ID == r.ID {
				up = append(up[:i], up[i+1:]...)
				break
			}
		}
	}

	return &ForUser{
		Resources:         <-chResources,
		Roles:             <-chRoles,
		TenantPermissions: <-chTenantPermissions,
		UserPermissions:   up,
	}, nil
}
