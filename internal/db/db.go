package db

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func OpenAndMigrate(dbPath string) (*sql.DB, error) {
	// connect to SQLite
	dsn := sqliteDSN(dbPath)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// set goose to use embedded migrations
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}

	// run migrations
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return sqlDB, nil
}

func sqliteDSN(dbPath string) string {
	if dbPath == ":memory:" {
		return "file::memory:?cache=shared&_pragma=foreign_keys(1)"
	}

	sep := "?"
	if strings.Contains(dbPath, "?") {
		sep = "&"
	}
	if strings.HasPrefix(dbPath, "file:") {
		return dbPath + sep + "_pragma=foreign_keys(1)"
	}
	return "file:" + dbPath + sep + "_pragma=foreign_keys(1)"
}
