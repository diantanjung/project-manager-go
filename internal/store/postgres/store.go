package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func Open(ctx context.Context, databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("opening database pool: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return db, nil
}
