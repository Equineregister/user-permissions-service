//go:build test
// +build test

package permissions_test

import (
	"context"
	"testing"

	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/internal/pkg/contextkey"
	"github.com/Equineregister/user-permissions-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

const (
	userAdmin        = "032fb302-4aee-4a68-b426-0c6faf12081e"
	userSalesManager = "652f4d18-dd3d-40c0-874e-cbe3566abccf"
	userSalesPerson  = "2133479c-35a8-4a49-a682-2952d4772ecc"
)

var expectedTenantPermissionsNoResources = permissions.TenantPermissions{
	{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
	{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
	{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
	{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
	{Name: "products:delete", ID: "acecdadf-f527-45bf-8123-353b7ee8dc6a"},
	{Name: "products:disable", ID: "cf7dc325-6bc9-44f5-aafb-fcdc694b111d"},
	{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
	{Name: "products:update", ID: "e12d692b-3a96-43aa-a966-dd3add99d312"},
}

var testExpectations_NoResourcesInRequest = map[string]permissions.ForUser{
	userAdmin: {
		TenantPermissions: expectedTenantPermissionsNoResources,
		Roles: permissions.Roles{
			{Name: "admin", ID: "550e8400-e29b-41d4-a716-446655440000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
			{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
			{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
			{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
			{Name: "products:delete", ID: "acecdadf-f527-45bf-8123-353b7ee8dc6a"},
			{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
		},
	},
	userSalesManager: {
		TenantPermissions: expectedTenantPermissionsNoResources,
		Roles: permissions.Roles{
			{Name: "sales manager", ID: "f47ac10b-58cc-4372-a567-0e02b2c3d479"},
			{Name: "sales person", ID: "123e4567-e89b-12d3-a456-426614174000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
			{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
			{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
			{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
			{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
			{Name: "products:update", ID: "e12d692b-3a96-43aa-a966-dd3add99d312"},
		},
		Resources: permissions.Resources{
			{Type: "products", ID: "75248bd5-73a2-4507-9ab3-5418abd33a3c"},
		},
	},
	userSalesPerson: {
		TenantPermissions: expectedTenantPermissionsNoResources,
		Roles: permissions.Roles{
			permissions.Role{Name: "sales person", ID: "123e4567-e89b-12d3-a456-426614174000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
		},
		Resources: permissions.Resources{
			{Type: "invoices", ID: "6b63b489-61cb-4087-8636-f10716bd724e"},
			{Type: "invoices", ID: "568104df-6ff3-40be-b660-91e3160aa7e6"},
		},
	},
}

var expectedTenantPermissionsProductsResource = permissions.TenantPermissions{
	{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
	{Name: "products:delete", ID: "acecdadf-f527-45bf-8123-353b7ee8dc6a"},
	{Name: "products:disable", ID: "cf7dc325-6bc9-44f5-aafb-fcdc694b111d"},
	{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
	{Name: "products:update", ID: "e12d692b-3a96-43aa-a966-dd3add99d312"},
}

var testExpectations_WithProductsResourceInRequest = map[string]permissions.ForUser{
	userAdmin: {
		TenantPermissions: expectedTenantPermissionsProductsResource,
		Roles: permissions.Roles{
			{Name: "admin", ID: "550e8400-e29b-41d4-a716-446655440000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
			{Name: "products:delete", ID: "acecdadf-f527-45bf-8123-353b7ee8dc6a"},
			{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
		},
		Resources: nil,
	},
	userSalesManager: {
		TenantPermissions: expectedTenantPermissionsProductsResource,
		Roles: permissions.Roles{
			{Name: "sales manager", ID: "f47ac10b-58cc-4372-a567-0e02b2c3d479"},
			{Name: "sales person", ID: "123e4567-e89b-12d3-a456-426614174000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
			{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
			{Name: "products:update", ID: "e12d692b-3a96-43aa-a966-dd3add99d312"},
		},
		Resources: permissions.Resources{
			{Type: "products", ID: "75248bd5-73a2-4507-9ab3-5418abd33a3c"},
		},
	},
	userSalesPerson: {
		TenantPermissions: expectedTenantPermissionsProductsResource,
		Roles: permissions.Roles{
			{Name: "sales person", ID: "123e4567-e89b-12d3-a456-426614174000"},
		},
		UserPermissions: nil,
		Resources:       nil,
	},
}

var expectedTenantPermissionsInvoicesResource = permissions.TenantPermissions{
	{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
	{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
	{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
}

var testExpectations_WithInvoicesResourceInRequest = map[string]permissions.ForUser{
	userAdmin: {
		TenantPermissions: expectedTenantPermissionsInvoicesResource,
		Roles: permissions.Roles{
			{Name: "admin", ID: "550e8400-e29b-41d4-a716-446655440000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
			{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
			{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
		},
		Resources: nil,
	},
	userSalesManager: {
		TenantPermissions: expectedTenantPermissionsInvoicesResource,
		Roles: permissions.Roles{
			{Name: "sales manager", ID: "f47ac10b-58cc-4372-a567-0e02b2c3d479"},
			{Name: "sales person", ID: "123e4567-e89b-12d3-a456-426614174000"},
		},
		UserPermissions: permissions.UserPermissions{
			permissions.UserPermission{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
			permissions.UserPermission{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
			permissions.UserPermission{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
		},
		Resources: nil,
	},
	userSalesPerson: {
		TenantPermissions: expectedTenantPermissionsInvoicesResource,
		Roles: permissions.Roles{
			{Name: "sales person", ID: "123e4567-e89b-12d3-a456-426614174000"},
		},
		UserPermissions: permissions.UserPermissions{
			{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
		},
		Resources: permissions.Resources{
			{Type: "invoices", ID: "6b63b489-61cb-4087-8636-f10716bd724e"},
			{Type: "invoices", ID: "568104df-6ff3-40be-b660-91e3160aa7e6"},
		},
	},
}

func TestGetForUser_LambdaGetUserPermissions(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.CtxKeyTenantID, TestTenantID)

	svc, _ := NewTestEnv(ctx, t, config.LambdaGetUserPermissions)

	t.Run("No resources in request", func(t *testing.T) {

		for userID, expectations := range testExpectations_NoResourcesInRequest {

			ctx = context.WithValue(ctx, contextkey.CtxKeyUserID, userID)

			fu, err := svc.GetForUser(ctx, nil)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			assert.EqualValues(t, expectations.Roles, fu.Roles)
			assert.EqualValues(t, expectations.TenantPermissions, fu.TenantPermissions)
			assert.EqualValues(t, expectations.UserPermissions, fu.UserPermissions, "user ID: %s", userID)
			assert.EqualValues(t, expectations.Resources, fu.Resources, "user ID: %s", userID)
		}
	})

	t.Run("With Products resource in request", func(t *testing.T) {

		for userID, expectations := range testExpectations_WithProductsResourceInRequest {

			ctx = context.WithValue(ctx, contextkey.CtxKeyUserID, userID)

			fu, err := svc.GetForUser(ctx, []string{"products"})
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			assert.EqualValues(t, expectations.Roles, fu.Roles)
			assert.EqualValues(t, expectations.TenantPermissions, fu.TenantPermissions)
			assert.EqualValues(t, expectations.UserPermissions, fu.UserPermissions, "user ID: %s", userID)
			assert.EqualValues(t, expectations.Resources, fu.Resources, "user ID: %s", userID)
		}
	})
	t.Run("With Invoices resource in request", func(t *testing.T) {

		for userID, expectations := range testExpectations_WithInvoicesResourceInRequest {

			ctx = context.WithValue(ctx, contextkey.CtxKeyUserID, userID)

			fu, err := svc.GetForUser(ctx, []string{"INVOICES"})
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			assert.EqualValues(t, expectations.Roles, fu.Roles)
			assert.EqualValues(t, expectations.TenantPermissions, fu.TenantPermissions)
			assert.EqualValues(t, expectations.UserPermissions, fu.UserPermissions, "user ID: %s", userID)
			assert.EqualValues(t, expectations.Resources, fu.Resources, "user ID: %s", userID)
		}
	})
}
