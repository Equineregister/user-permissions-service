package postgres

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Equineregister/user-permissions-service/internal/pkg/contextkey"
	"github.com/Equineregister/user-permissions-service/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantPool struct {
	cfg        *config.Config
	mu         sync.RWMutex
	tenantPool map[string]*pgxpool.Pool

	host string
	port int
}

func (tp *TenantPool) GetTenantConn(ctx context.Context) (*pgxpool.Pool, error) {
	tenantID, ok := tp.extractTenantID(ctx)
	if !ok {
		return nil, fmt.Errorf("tenantID not found in context: %v", tenantID)
	}

	tp.mu.RLock()
	tenant, ok := tp.tenantPool[tenantID]
	if ok {
		return tenant, nil
	}
	tp.mu.RUnlock()

	awsConfig := tp.cfg.GetAWSConfig()

	dsn := fmt.Sprintf("user=%s host=%s port=%d dbname=%s", "root", tp.host, tp.port, tenantID)
	pool, err := NewWithIAM(ctx, &awsConfig, dsn)
	if err != nil {
		return nil, fmt.Errorf("new with iam: %w", err)
	}

	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.tenantPool[tenantID] = pool

	return tenant, nil
}

func (tp *TenantPool) extractTenantID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(contextkey.CtxKeyTenantID).(string)
	if ok {
		return tenantID, true
	}

	return "", false
}

func NewTenantPool(cfg *config.Config) *TenantPool {
	tenantPool := make(map[string]*pgxpool.Pool)

	return &TenantPool{
		tenantPool: tenantPool,
	}
}

type RDSToken struct {
	Expires *time.Time
	Token   string
}

func NewWithIAM(ctx context.Context, awsConfig *aws.Config, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	var rdsToken RDSToken

	// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-connecting.html
	config.ConnConfig.RuntimeParams["application_name"] = "rds-pgx"

	config.BeforeConnect = func(ctx context.Context, pgxConfig *pgx.ConnConfig) error {
		if rdsToken.Expires == nil || rdsToken.Expires.Before(time.Now()) {
			token, expires, err := getRDSToken(
				ctx,
				awsConfig,
				awsConfig.Region,
				pgxConfig,
			)
			if err != nil {
				return fmt.Errorf("get rds token: %w", err)
			}

			rdsToken.Token = token
			rdsToken.Expires = expires
		}

		pgxConfig.Password = rdsToken.Token

		return nil
	}

	return pgxpool.NewWithConfig(ctx, config)
}

func getRDSToken(ctx context.Context, awsConfig *aws.Config, region string, pgxConfig *pgx.ConnConfig) (string, *time.Time, error) {
	token, err := auth.BuildAuthToken(
		ctx,
		fmt.Sprintf("%s:%d", pgxConfig.Host, pgxConfig.Port),
		region,
		pgxConfig.User,
		awsConfig.Credentials,
	)
	if err != nil {
		return "", nil, fmt.Errorf("build auth token: %w", err)
	}

	u, err := url.Parse(token)
	if err != nil {
		return "", nil, fmt.Errorf("parse token: %w", err)
	}

	issued, err := time.Parse(time.RFC3339, u.Query().Get("X-Amz-Date"))
	if err != nil {
		return "", nil, fmt.Errorf("parse time: %w", err)
	}

	ttl, err := time.ParseDuration(u.Query().Get("X-Amz-Expires") + "s")
	if err != nil {
		return "", nil, fmt.Errorf("parse duration: %w", err)
	}

	expires := issued.Add(ttl)

	return token, &expires, nil
}
