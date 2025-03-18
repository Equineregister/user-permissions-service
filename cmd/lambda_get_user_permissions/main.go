package main

import (
	"context"
	"os"

	"github.com/Equineregister/user-permissions-service/internal/adapter/secondary/postgres"
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
	Role              string     `json:"role"`
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

	tenantPerms := make([]string, len(forUser.TenantPermissions))
	for i, p := range forUser.TenantPermissions {
		tenantPerms[i] = p.Name
	}

	userPerms := make([]string, len(forUser.UserPermissions))
	for i, p := range forUser.UserPermissions {
		userPerms[i] = p.Name
	}

	resources := make([]Resource, len(forUser.Resources))
	for i, r := range forUser.Resources {
		resources[i] = Resource{
			ResourceID:   r.ID,
			ResourceType: r.Type,
		}
	}

	return Response{
		TenantID:          request.TenantID,
		UserID:            request.UserID,
		Role:              forUser.Role.String(),
		TenantPermissions: tenantPerms,
		UserPermissions:   userPerms,
		UserResources:     resources,
	}, nil
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
