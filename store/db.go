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

// GetAllHosts retrieves a list of all registered hostnames in the database
func (db *DB) GetAllHosts() ([]string, error) {
	rows, err := db.Query("SELECT host FROM routes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []string
	for rows.Next() {
		var host string
		if err := rows.Scan(&host); err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return hosts, nil
}

// GetAllRoutes retrieves all registered routes from the database
func (db *DB) GetAllRoutes() (map[string]string, error) {
	rows, err := db.Query("SELECT host, target FROM routes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := make(map[string]string)
	for rows.Next() {
		var host, target string
		if err := rows.Scan(&host, &target); err != nil {
			return nil, err
		}
		routes[host] = target
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return routes, nil
}

// SetRoute adds or updates a routing rule
func (db *DB) SetRoute(host, target string) error {
	_, err := db.Exec("INSERT INTO routes (host, target) VALUES (?, ?) ON CONFLICT(host) DO UPDATE SET target=excluded.target", host, target)
	return err
}
