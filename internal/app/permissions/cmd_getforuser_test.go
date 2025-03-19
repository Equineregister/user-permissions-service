//go:build test
// +build test

package permissions_test

import (
	"context"
	"testing"

	"github.com/Equineregister/user-permissions-service/pkg/config"
)

func TestGetForUser_LambdaGetUserPermissions(t *testing.T) {
	ctx := context.Background()

	t.Run("No resources in request", func(t *testing.T) {
		svc, _ := NewTestEnv(ctx, t, config.LambdaGetUserPermissions)

		_, err := svc.GetForUser(ctx, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

	})
}
