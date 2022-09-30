// Package dbschema contains the database schema, migrations and seeding data.
package dbschema

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
	"github.com/phbpx/gobeers/internal/sys/database"
)

//go:embed sql
var migrations embed.FS

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	// Load the migrations from the embedded filesystem.
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("invalid source instance: %w", err)
	}

	// Create the database driver for the migrations.
	target, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("invalid target postgres instance, %w", err)
	}

	// Create the migration instance.
	m, err := migrate.NewWithInstance("httpfs", source, "postgres", target)
	if err != nil {
		return err
	}

	// Run the migrations.
	if err := m.Up(); err != nil {
		// If the error is not a "no change" error, return it.
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
