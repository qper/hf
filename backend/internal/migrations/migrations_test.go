package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestRunMigrationsIdempotent(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	result, err := RunWithDB(db, dsn)
	if err != nil {
		t.Fatalf("first run migrations: %v", err)
	}
	if result.Version == 0 {
		t.Fatalf("expected migration version > 0, got %d", result.Version)
	}

	result, err = RunWithDB(db, dsn)
	if err != nil {
		t.Fatalf("second run migrations: %v", err)
	}
	if !result.NoChange {
		t.Fatalf("expected no change on second run, got %+v", result)
	}
}

func TestRunMigrationsFailsFastOnUnavailableDB(t *testing.T) {
	_, err := RunWithDSN("postgres://127.0.0.1:1/test?sslmode=disable")
	if err == nil {
		t.Fatal("expected error for unavailable database")
	}
	if got := err.Error(); got == "" || fmt.Sprintf("%v", err) == "" {
		t.Fatal("expected descriptive error message")
	}
}
