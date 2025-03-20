package main

import (
	"context"
	"os"

	"github.com/Equineregister/user-permissions-service/internal/adapters/secondary/postgres"
	"github.com/Equineregister/user-permissions-service/internal/app/permissions"
	"github.com/Equineregister/user-permissions-service/internal/pkg/application"
	"github.com/Equineregister/user-permissions-service/internal/pkg/contextkey"
	"github.com/Equineregister/user-permissions-service/pkg/config"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/exp/slog"
)

// Request represents the input event structure
type Request struct {
	TenantID  string   `json:"tenantId"`
	UserID    string   `json:"userId"`
	Resources []string `json:"resources"`
}

// Response represents the output structure
type Response struct {
	TenantID          string     `json:"tenantId"`
	UserID            string     `json:"userId"`
	Roles             []string   `json:"roles"`
	TenantPermissions []string   `json:"tenantPermissions"`
	UserPermissions   []string   `json:"userPermissions"`
	UserResources     []Resource `json:"userResources"`
}

type Resource struct {
	ResourceID   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
}

type handler struct {
	repo    permissions.Reader
	service *permissions.Service
}

// Handler is the Lambda function handler
func (h *handler) handle(ctx context.Context, request Request) (Response, error) {
	slog.Debug("received request", "tenantId", request.TenantID, "userId", request.UserID)

	ctx = context.WithValue(ctx, contextkey.CtxKeyTenantID, request.TenantID)
	ctx = context.WithValue(ctx, contextkey.CtxKeyUserID, request.UserID)

	forUser, err := h.service.GetForUser(ctx, request.Resources)
	if err != nil {
		slog.Error("error getting permissions for user", "error", err.Error())
		return Response{}, err
	}

	return response(request.TenantID, request.UserID, forUser), nil
}

func response(tenantID, userID string, forUser *permissions.ForUser) Response {
	resp := Response{
		TenantID: tenantID,
		UserID:   userID,
	}
	if forUser == nil {
		resp.Roles = []string{}
		resp.TenantPermissions = []string{}
		resp.UserPermissions = []string{}
		resp.UserResources = []Resource{}
		return resp
	}

	resp.UserResources = make([]Resource, len(forUser.Resources))
	for i, r := range forUser.Resources {
		resp.UserResources[i] = Resource{
			ResourceID:   r.ID,
			ResourceType: r.Type,
		}
	}

	resp.Roles = forUser.Roles.StringSlice()
	resp.TenantPermissions = forUser.TenantPermissions.StringSlice()
	resp.UserPermissions = forUser.UserPermissions.StringSlice()
	return resp
}

func main() {
	application.InitLogger()

	ctx := context.Background()

	cfg := config.New()
	if err := cfg.Load(ctx, config.LambdaGetUserPermissions); err != nil {
		slog.Error("Error loading config", "error", err.Error())
		os.Exit(1)
	}

	repo := postgres.NewPermissionsRepo(cfg)
	service := permissions.NewService(repo)

	h := &handler{
		repo:    repo,
		service: service,
	}

	lambda.Start(h.handle)
}
