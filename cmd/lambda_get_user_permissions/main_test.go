package main

import (
	"testing"

	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/stretchr/testify/assert"
)

func Test_response(t *testing.T) {
	type args struct {
		tenantID string
		userID   string
		forUser  *permissions.ForUser
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		{
			name: "Nil forUser",
			args: args{
				tenantID: "test_tenant",
				userID:   "cba1470a-58b6-444f-a763-31b309f087e2",
				forUser:  nil,
			},
			want: Response{
				TenantID:          "test_tenant",
				UserID:            "cba1470a-58b6-444f-a763-31b309f087e2",
				Roles:             []string{},
				TenantPermissions: []string{},
				UserPermissions:   []string{},
				UserResources:     []Resource{},
			},
		},
		{
			name: "Empty forUser",
			args: args{
				tenantID: "test_tenant",
				userID:   "4817f881-0081-4a96-a8c1-7da5b743c2ec",
				forUser: &permissions.ForUser{
					Roles:             nil,
					TenantPermissions: nil,
					UserPermissions:   nil,
					Resources:         nil,
				},
			},
			want: Response{
				TenantID:          "test_tenant",
				UserID:            "4817f881-0081-4a96-a8c1-7da5b743c2ec",
				Roles:             []string{},
				TenantPermissions: []string{},
				UserPermissions:   []string{},
				UserResources:     []Resource{},
			},
		},
		{
			name: "Populated forUser",
			args: args{
				tenantID: "test_tenant",
				userID:   "2cdabaf2-24fb-4c90-961f-b92f129f895e",
				forUser: &permissions.ForUser{
					Roles: permissions.Roles{
						permissions.Role{Name: "admin", ID: "550e8400-e29b-41d4-a716-446655440000"},
						permissions.Role{Name: "sales manager", ID: "f47ac10b-58cc-4372-a567-0e02b2c3d479"},
					},
					TenantPermissions: permissions.TenantPermissions{
						permissions.TenantPermission{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
						permissions.TenantPermission{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
						permissions.TenantPermission{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
						permissions.TenantPermission{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
						permissions.TenantPermission{Name: "products:delete", ID: "acecdadf-f527-45bf-8123-353b7ee8dc6a"},
						permissions.TenantPermission{Name: "products:disable", ID: "cf7dc325-6bc9-44f5-aafb-fcdc694b111d"},
						permissions.TenantPermission{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
						permissions.TenantPermission{Name: "products:update", ID: "e12d692b-3a96-43aa-a966-dd3add99d312"},
					},
					UserPermissions: permissions.UserPermissions{
						permissions.UserPermission{Name: "invoices:create", ID: "2f9606d8-4bff-46e7-bd8f-ae9e476d3995"},
						permissions.UserPermission{Name: "invoices:delete", ID: "41c21275-b7d5-4031-b551-b5e293b85319"},
						permissions.UserPermission{Name: "invoices:read", ID: "8f20eca6-9859-4532-babb-65a528e1611e"},
						permissions.UserPermission{Name: "products:create", ID: "df6ae9bc-e957-41c1-a683-3773667c7628"},
						permissions.UserPermission{Name: "products:delete", ID: "acecdadf-f527-45bf-8123-353b7ee8dc6a"},
						permissions.UserPermission{Name: "products:read", ID: "62752f21-fbe2-4301-a72d-7dc8963e08e2"},
					},
					Resources: permissions.Resources{
						permissions.Resource{ID: "6b63b489-61cb-4087-8636-f10716bd724e", Type: "invoices"},
						permissions.Resource{ID: "568104df-6ff3-40be-b660-91e3160aa7e6", Type: "invoices"},
					},
				},
			},
			want: Response{
				TenantID:          "test_tenant",
				UserID:            "2cdabaf2-24fb-4c90-961f-b92f129f895e",
				Roles:             []string{"admin", "sales manager"},
				TenantPermissions: []string{"invoices:create", "invoices:delete", "invoices:read", "products:create", "products:delete", "products:disable", "products:read", "products:update"},
				UserPermissions:   []string{"invoices:create", "invoices:delete", "invoices:read", "products:create", "products:delete", "products:read"},
				UserResources: []Resource{
					{ResourceID: "6b63b489-61cb-4087-8636-f10716bd724e", ResourceType: "invoices"},
					{ResourceID: "568104df-6ff3-40be-b660-91e3160aa7e6", ResourceType: "invoices"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response(tt.args.tenantID, tt.args.userID, tt.args.forUser)
			assert.Equal(t, tt.want, got)
		})
	}
}
