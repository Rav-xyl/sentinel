package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB encapsulates the SQLite database connection
type DB struct {
	*sql.DB
}

// InitDB initializes the SQLite database and creates necessary tables
func InitDB(filepath string) (*DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("[Store] SQLite database initialized at", filepath)
	return &DB{db}, nil
}

// createTables ensures the routes table exists
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS routes (
		host TEXT PRIMARY KEY,
		target TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}

// GetTarget retrieves the target URL for a given hostname
func (db *DB) GetTarget(host string) (string, error) {
	var target string
	err := db.QueryRow("SELECT target FROM routes WHERE host = ?", host).Scan(&target)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Not found, not an error
		}
		return "", err
	}
	return target, nil
}

// SetRoute adds or updates a routing rule
func (db *DB) SetRoute(host, target string) error {
	_, err := db.Exec("INSERT INTO routes (host, target) VALUES (?, ?) ON CONFLICT(host) DO UPDATE SET target=excluded.target", host, target)
	return err
}
