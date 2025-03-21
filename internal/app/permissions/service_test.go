//go:build test
// +build test

package permissions_test

import (
	"context"
	"os"
	"testing"

	"github.com/Equineregister/user-permissions-service/internal/adapters/secondary/postgres"
	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/internal/pkg/application"
	"github.com/Equineregister/user-permissions-service/pkg/testdatabase"
)

const (
	TestTenantID = "test_tenant"
)

func NewTestEnv(ctx context.Context, t *testing.T) (*permissions.Service, permissions.ReaderWriter) {
	t.Helper()

	os.Setenv("LOG_LEVEL", "debug")
	application.InitLogger()

	db, err := testdatabase.NewTestDatabase(ctx, postgres.Migrations, postgres.Tenants)
	if err != nil {
		t.Fatalf("failed to create new test database: %s", err.Error())
	}
	tp := postgres.NewTenantPoolFromSuppliedPool(ctx, TestTenantID, db.DB)

	repo := postgres.NewPermissionsRepoWithTenantPool(tp)
	service := permissions.NewService(repo)

	return service, repo
}
