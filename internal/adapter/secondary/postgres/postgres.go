package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Equineregister/user-permissions-service/internal/pkg/contextkey"
	"github.com/jackc/pgx/v5"
)

func whereSQL(where []string) string {
	if len(where) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(where, " AND ")
}

func rollback(ctx context.Context, tx pgx.Tx) error {
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func extractUserID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(contextkey.CtxKeyUserID).(string)
	if ok {
		return tenantID, true
	}

	return "", false
}
