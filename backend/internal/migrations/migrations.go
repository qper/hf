package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type Result struct {
	Applied  int
	Version  uint
	NoChange bool
}

func RunWithDSN(dsn string) (Result, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return Result{}, fmt.Errorf("open database: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	return RunWithDB(db, dsn)
}

func RunWithDB(db *sql.DB, dsn string) (Result, error) {
	if db == nil {
		return Result{}, errors.New("database handle is nil")
	}
	if err := db.Ping(); err != nil {
		return Result{}, fmt.Errorf("connect to database: %w", err)
	}

	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return Result{}, fmt.Errorf("create migration source: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return Result{}, fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return Result{}, fmt.Errorf("create migration instance: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	prevVersion, _, prevErr := m.Version()
	if prevErr != nil && !errors.Is(prevErr, migrate.ErrNilVersion) {
		return Result{}, fmt.Errorf("get current migration version: %w", prevErr)
	}

	err = m.Up()
	switch {
	case err == nil:
		version, _, versionErr := m.Version()
		if versionErr != nil {
			return Result{}, fmt.Errorf("get migration version after apply: %w", versionErr)
		}
		applied := 0
		if prevErr != nil || prevVersion == 0 {
			applied = int(version)
		} else {
			applied = int(version - prevVersion)
		}
		return Result{Applied: applied, Version: version}, nil
	case errors.Is(err, migrate.ErrNoChange):
		version, _, versionErr := m.Version()
		if versionErr != nil {
			return Result{}, fmt.Errorf("get migration version after no change: %w", versionErr)
		}
		return Result{Version: version, NoChange: true}, nil
	default:
		return Result{}, fmt.Errorf("run migrations: %w", err)
	}
}
