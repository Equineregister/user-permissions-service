package main

import (
	"encoding/json"
	"testing"

	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/pkg/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// wantJSONLogged is a flag to enable JSON logging of the expected output, useful for OPA testing.
	wantJSONLogged = false
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
				TenantID:           "test_tenant",
				UserID:             "cba1470a-58b6-444f-a763-31b309f087e2",
				Roles:              []string{},
				ExtraPermissions:   []string{},
				RevokedPermissions: []string{},
				UserResources:      []Resource{},
				RoleGraph:          rego.RoleGraph{},
			},
		},
		{
			name: "Empty forUser",
			args: args{
				tenantID: "test_tenant",
				userID:   "4817f881-0081-4a96-a8c1-7da5b743c2ec",
				forUser: &permissions.ForUser{
					Roles:              nil,
					ExtraPermissions:   nil,
					RevokedPermissions: nil,
					Resources:          nil,
					RoleMap:            nil,
				},
			},
			want: Response{
				TenantID:           "test_tenant",
				UserID:             "4817f881-0081-4a96-a8c1-7da5b743c2ec",
				Roles:              []string{},
				ExtraPermissions:   []string{},
				RevokedPermissions: []string{},
				UserResources:      []Resource{},
				RoleGraph:          rego.RoleGraph{},
			},
		},
		{
			name: "Populated forUser",
			args: args{
				tenantID: "test_tenant",
				userID:   "2cdabaf2-24fb-4c90-961f-b92f129f895e",
				forUser: &permissions.ForUser{
					Roles: permissions.Roles{
						{Name: "customer service", ID: "eb1386c5-6a18-43e3-9176-b7ffa927ecc2"},
						{Name: "discount decider", ID: "b5622eba-1c4c-42de-803d-778261c61b79"},
					},
					ExtraPermissions: permissions.UserExtraPermissions{
						{Name: "customers:update", ID: "2c1c7083-e97e-456c-8e33-db4b1ae4ef49"},
					},
					RevokedPermissions: permissions.UserRevokedPermissions{
						{Name: "discounts:delete", ID: "06597dcc-26dd-4f48-99be-95ccce65bca9"},
					},
					Resources: permissions.Resources{
						{ID: "90a12308-003c-4b90-957e-59ad1f3e5b7a", Type: "customers"},
						{ID: "6b7ef64b-8f4f-47e2-9cc6-ebeb0075904b", Type: "discounts"},
					},
					RoleMap: permissions.TenantRoleMap{
						{Name: "customer manager", ID: "e4d9424c-8249-4888-9812-35f5aefb02b2"}: {
							Permissions: permissions.TenantPermissions{
								{Name: "customers:delete", ID: "af02ce53-955d-4c5c-89db-9a9e0c9f11fa"},
								{Name: "customers:update", ID: "2c1c7083-e97e-456c-8e33-db4b1ae4ef49"},
							},
							Inherits: []permissions.Role{
								{Name: "customer service", ID: "eb1386c5-6a18-43e3-9176-b7ffa927ecc2"},
							},
						},
						{Name: "customer service", ID: "eb1386c5-6a18-43e3-9176-b7ffa927ecc2"}: {
							Permissions: permissions.TenantPermissions{
								{Name: "customers:create", ID: "9bb579ce-e5a0-4917-a8c8-f9916e5ee9aa"},
								{Name: "customers:read", ID: "1349adb7-052b-4879-b843-29621d96966f"},
							},
							Inherits: []permissions.Role{
								{Name: "discount decider", ID: "b5622eba-1c4c-42de-803d-778261c61b79"},
							},
						},
						{Name: "discount decider", ID: "b5622eba-1c4c-42de-803d-778261c61b79"}: {
							Permissions: permissions.TenantPermissions{
								{Name: "discounts:create", ID: "2a16bbbb-40cb-4bfb-8e99-fa521ab74508"},
								{Name: "discounts:delete", ID: "06597dcc-26dd-4f48-99be-95ccce65bca9"},
								{Name: "discounts:read", ID: "49539ec8-f619-47b4-8101-299ed590e5fb"},
							},
							Inherits: []permissions.Role{},
						},
					},
				},
			},
			want: Response{
				TenantID:           "test_tenant",
				UserID:             "2cdabaf2-24fb-4c90-961f-b92f129f895e",
				Roles:              []string{"customer service", "discount decider"},
				ExtraPermissions:   []string{"customers:update"},
				RevokedPermissions: []string{"discounts:delete"},
				UserResources: []Resource{
					{ResourceID: "90a12308-003c-4b90-957e-59ad1f3e5b7a", ResourceType: "customers"},
					{ResourceID: "6b7ef64b-8f4f-47e2-9cc6-ebeb0075904b", ResourceType: "discounts"},
				},
				RoleGraph: rego.RoleGraph{
					"customer manager": rego.Role{Permissions: []string{"customers:delete", "customers:update"},
						Inherits: []string{"customer service"},
					},
					"customer service": rego.Role{Permissions: []string{"customers:create", "customers:read"},
						Inherits: []string{"discount decider"},
					},
					"discount decider": rego.Role{
						Permissions: []string{"discounts:create", "discounts:delete", "discounts:read"},
						Inherits:    []string{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response(tt.args.tenantID, tt.args.userID, tt.args.forUser)
			assert.Equal(t, tt.want, got)

			if wantJSONLogged {
				resp, err := json.Marshal(got)
				require.NoError(t, err)
				t.Log(string(resp))
			}

		})
	}
}
