package main

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

// runMigrations applies all pending goose migrations against the given DB.
// Migration SQL files are embedded in the binary — no filesystem dependency.
func runMigrations(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "db/migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
