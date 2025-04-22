package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"bookapi/config"
)

var DB *pgxpool.Pool

func InitDB(cfg config.Config) error {
	poolConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return fmt.Errorf("error parsing db config: %w", err)
	}
	poolConfig.MaxConns = 20

	DB, err = pgxpool.NewWithConfig(context. Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("error connocting to database: %w", err)
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}