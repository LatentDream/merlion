package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"merlion/internal/store/local/operations"
)

const ENV_DB_PATH_OVERWRITTEN = "MERLION_DB_PATH"

// GetDBPath determines the appropriate path for the SQLite database file.
// It prioritizes the MERLION_DB_PATH environment variable if set.
// Otherwise, it defaults to a path within the user's home directory (`~/.merlion/notes.db`).
func GetDBPath() (string, error) {
	if dbPath := os.Getenv(ENV_DB_PATH_OVERWRITTEN); dbPath != "" {
		log.Infof("Using database path from MERLION_DB_PATH: %s", dbPath)
		// Ensure directory for custom path also exists
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0750); err != nil { // 0750 for rwx for user, rx for group
			return "", fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
		}
		return dbPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	defaultPath := filepath.Join(homeDir, ".merlion", "notes.db")
	log.Infof("Using default database path: %s", defaultPath)

	dbDir := filepath.Dir(defaultPath)
	if err := os.MkdirAll(dbDir, 0750); err != nil { // 0750 for rwx for user, rx for group
		return "", fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
	}
	return defaultPath, nil
}

// InitDB initializes the SQLite database connection.
// It determines the database path, opens the connection, pings it to ensure liveness,
// and applies any pending database migrations.
func InitDB() (*sql.DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine database path: %w", err)
	}

	var db *sql.DB
	db, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database at %s: %w", dbPath, err)
	}

	// Ping the database to ensure the connection is live.
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Infof("Successfully connected to database: %s", dbPath)

	if err := ApplyMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Info("Database migrations applied successfully.")
	return db, nil
}
