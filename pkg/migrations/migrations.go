package migrations

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createMigrationsTable = `
	CREATE TABLE IF NOT EXISTS _migrations (
    file_name text NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
`
const selectMigrations = `
	SELECT file_name FROM _migrations
	ORDER BY file_name DESC 
	LIMIT 1;
`
const insertMigration = `
	INSERT INTO _migrations (file_name) VALUES ($1);
`

const (
	migrationsDir      = "migrations"
	tenantsDir         = "migrations/tenants/test"
	testTenantFilename = "0001-foundation.sql"
)

func Migrate(ctx context.Context, db *pgxpool.Pool, fsys fs.FS) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("rollback tx", "error", err.Error())
		}
	}()
	_, err = tx.Exec(ctx, createMigrationsTable)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	rows, err := tx.Query(ctx, selectMigrations)
	if err != nil {
		return fmt.Errorf("select migrations: %w", err)
	}
	defer rows.Close()

	var appliedMigrations []string
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		appliedMigrations = append(appliedMigrations, file)
	}
	if rows.Err() != nil {
		return fmt.Errorf("iterate over rows: %w", rows.Err())
	}

	latestMigration := ""
	switch len(appliedMigrations) {
	case 0:
		break
	case 1:
		latestMigration = appliedMigrations[0]
	default:
		return fmt.Errorf("unexpected number of rows: %d", len(appliedMigrations))
	}

	files, err := fs.ReadDir(fsys, migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir (%s): %w", migrationsDir, err)
	}

	files = filterFiles(files)
	for _, file := range files {
		// skip to the one after the latest
		if file.Name() <= latestMigration {
			continue
		}

		bytes, err := fs.ReadFile(fsys, fmt.Sprintf("%s/%s", migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		query := string(bytes)
		if len(strings.TrimSpace(query)) == 0 {
			return fmt.Errorf("empty query in file: %s", file.Name())
		}

		_, err = tx.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("exec query for file: %s caused: %w", file.Name(), err)
		}

		_, err = tx.Exec(ctx, insertMigration, file.Name())
		if err != nil {
			return fmt.Errorf("insert migration: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func filterFiles(files []fs.DirEntry) []fs.DirEntry {
	var filteredFiles []fs.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		filteredFiles = append(filteredFiles, file)
	}
	return filteredFiles
}

func LoadTestTenantData(ctx context.Context, db *pgxpool.Pool, fsys fs.FS) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("rollback tx", "error", err.Error())
		}
	}()

	bytes, err := fs.ReadFile(fsys, fmt.Sprintf("%s/%s", tenantsDir, testTenantFilename))
	if err != nil {
		return fmt.Errorf("read tenants file (%s): %w", testTenantFilename, err)
	}
	query := string(bytes)
	if len(strings.TrimSpace(query)) == 0 {
		return fmt.Errorf("empty query in file: %s", testTenantFilename)
	}

	_, err = tx.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("exec query for file: %s caused: %w", testTenantFilename, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
