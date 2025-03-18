package main

import (
	"context"
	"os"

	"github.com/Equineregister/user-permissions-service/internal/pkg/application"
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
	TenantPermissions []string   `json:"tenantPermissions"`
	UserPermissions   []string   `json:"userPermissions"`
	UserResources     []Resource `json:"userResources"`
}

type Resource struct {
	ResourceID   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
}

type handler struct {
}

// Handler is the Lambda function handler
func (h *handler) handle(ctx context.Context, request Request) (Response, error) {
	slog.Debug("received request", "tenantId", request.TenantID, "userId", request.UserID)

	return Response{
		TenantID: request.TenantID,
		UserID:   request.UserID,
		// TODO
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

	h := &handler{}

	lambda.Start(h.handle)
}
