package db

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open opens a GORM database connection using DATABASE_URL from the environment.
func Open() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return db, nil
}

// OpenWithSchema opens a GORM connection scoped to a Postgres schema.
// stdlib.OpenDB creates a fresh driver instance per call - no global registry,
// fully isolated. Safe for concurrent tests.
func OpenWithSchema(schema string) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	// RuntimeParams are sent at Postgres connection startup before any query.
	connConfig.RuntimeParams["search_path"] = schema + ", public"

	// stdlib.OpenDB uses a fresh driver instance - not the global RegisterConnConfig pool.
	sqlDB := stdlib.OpenDB(*connConfig)

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open database with schema: %w", err)
	}

	return db, nil
}
