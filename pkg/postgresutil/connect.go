package postgresutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildDSN(host, user, password, dbname string, port int) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, password, dbname, port)
}

func Connect(ctx context.Context, port int, dbHost, dbUser, dbPass, dbName string) (*pgxpool.Pool, string, error) {

	dbAddr := BuildDSN(dbHost, dbUser, dbPass, dbName, port)
	conn, err := pgxpool.New(ctx, dbAddr)
	if err != nil {
		return nil, "", fmt.Errorf("open: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, "", fmt.Errorf("ping context: %w", err)
	}

	return conn, dbAddr, nil
}
