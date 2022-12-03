// Package database provides support for access the database.
package database

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib" // Calls init function.
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
)

// Config is the required properties to use the database.
type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

// String returns a string representation of the database configuration.
func (cfg *Config) String() string {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	return u.String()
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*bun.DB, error) {
	pgxConfig, err := pgx.ParseConfig(cfg.String())
	if err != nil {
		return nil, err
	}

	pgxConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	sqldb := stdlib.OpenDB(*pgxConfig)

	// Set connection pool settings.
	// https://bun.uptrace.dev/guide/running-bun-in-production.html
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(cfg.Name)))
	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *bun.DB) error {

	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// IsIntegrityViolation chec if the error code is one of the following:
//
//	"23000", "23001", "23502", "23503", "23505", "23514", "23P01"
//
// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
func IsIntegrityViolation(err error) bool {
	if pgErr, ok := err.(pgdriver.Error); ok && pgErr.IntegrityViolation() {
		return true
	}
	return false
}

// IsNoRowError checks if the error is caused by no row found in the database.
func IsNoRowError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
