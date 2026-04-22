package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	dbmigrations "erp-job/db/migrations"
)

const createSchemaMigrationsQuery = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version VARCHAR(255) NOT NULL PRIMARY KEY,
	applied_at DATETIME(6) NOT NULL
);
`

func ApplyMigrations(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, createSchemaMigrationsQuery); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	applied, err := loadAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	names, err := dbmigrations.UpFiles()
	if err != nil {
		return fmt.Errorf("list migration files: %w", err)
	}
	sort.Strings(names)

	migrationFS := dbmigrations.FS()
	for _, name := range names {
		if applied[name] {
			continue
		}

		script, err := fs.ReadFile(migrationFS, name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if err := applyMigration(ctx, db, name, string(script)); err != nil {
			return err
		}
	}

	return nil
}

func loadAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("load applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return applied, nil
}

func applyMigration(ctx context.Context, db *sql.DB, version, script string) error {
	if _, err := db.ExecContext(ctx, strings.TrimSpace(script)); err != nil {
		return fmt.Errorf("apply migration %s: %w", version, err)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`, version, time.Now().UTC()); err != nil {
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	return nil
}
