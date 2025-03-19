package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func rollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		slog.Error("Failed to rollback transaction", "error", err.Error())
	}
}
