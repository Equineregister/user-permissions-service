//go:build test
// +build test

package permissions_test

import (
	"context"
	"testing"

	"github.com/Equineregister/user-permissions-service/internal/pkg/contextkey"
	"github.com/Equineregister/user-permissions-service/pkg/config"
)

const (
	salesManager = "652f4d18-dd3d-40c0-874e-cbe3566abccf"
)

func TestGetForUser_LambdaGetUserPermissions(t *testing.T) {
	ctx := context.Background()

	ctx = context.WithValue(ctx, contextkey.CtxKeyTenantID, TestTenantID)

	t.Run("No resources in request", func(t *testing.T) {
		svc, _ := NewTestEnv(ctx, t, config.LambdaGetUserPermissions)

		ctx = context.WithValue(ctx, contextkey.CtxKeyUserID, salesManager)

		_, err := svc.GetForUser(ctx, nil)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

	})
}
