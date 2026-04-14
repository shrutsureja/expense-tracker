package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"expense-tracker/internal/database/migrations"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	// Ensure schema_migrations table exists for tracking
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	migrationFiles := []struct {
		version  int
		filename string
	}{
		{1, "001_initial_schema.up.sql"},
	}

	for _, mf := range migrationFiles {
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", mf.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking migration %d: %w", mf.version, err)
		}
		if exists > 0 {
			continue
		}

		sqlBytes, err := migrations.FS.ReadFile(mf.filename)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", mf.filename, err)
		}

		log.Printf("Applying migration %d: %s", mf.version, mf.filename)
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("executing migration %s: %w", mf.filename, err)
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", mf.version); err != nil {
			return fmt.Errorf("recording migration %d: %w", mf.version, err)
		}
	}

	return nil
}
