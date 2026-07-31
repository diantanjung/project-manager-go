package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"

	"project-manager-go/internal/config"
)

const defaultMigrationsSource = "file://migrations"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: go run ./cmd/db <ping|migrate-up|migrate-down|migrate-version|migrate-steps>")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	switch os.Args[1] {
	case "ping":
		return ping(cfg.DatabaseURL)
	case "migrate-up":
		return migrateUp(cfg.DatabaseURL)
	case "migrate-down":
		return migrateSteps(cfg.DatabaseURL, -1)
	case "migrate-version":
		return migrateVersion(cfg.DatabaseURL)
	case "migrate-steps":
		if len(os.Args) != 3 {
			return errors.New("usage: go run ./cmd/db migrate-steps <n>")
		}
		steps, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return fmt.Errorf("parsing migration steps: %w", err)
		}
		return migrateSteps(cfg.DatabaseURL, steps)
	default:
		return fmt.Errorf("unknown db command %q", os.Args[1])
	}
}

func ping(databaseURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("connecting database: %w", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute)

	var databaseName string
	if err := db.GetContext(ctx, &databaseName, "SELECT current_database()"); err != nil {
		return fmt.Errorf("querying database name: %w", err)
	}
	fmt.Printf("database ping ok: %s\n", databaseName)
	return nil
}

func migrateUp(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migrations up: %w", err)
	}
	fmt.Println("migrations up ok")
	return nil
}

func migrateSteps(databaseURL string, steps int) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migration steps: %w", err)
	}
	fmt.Printf("migration steps ok: %d\n", steps)
	return nil
}

func migrateVersion(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Println("migration version: none")
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading migration version: %w", err)
	}
	fmt.Printf("migration version: %d dirty=%t\n", version, dirty)
	return nil
}

func newMigrator(databaseURL string) (*migrate.Migrate, error) {
	sourceURL := os.Getenv("MIGRATIONS_SOURCE")
	if sourceURL == "" {
		sourceURL = defaultMigrationsSource
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("creating migrator: %w", err)
	}
	return m, nil
}

func closeMigrator(m *migrate.Migrate) {
	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		log.Printf("closing migration source: %v", sourceErr)
	}
	if databaseErr != nil {
		log.Printf("closing migration database: %v", databaseErr)
	}
}
