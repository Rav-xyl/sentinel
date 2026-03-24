package sidecars

import (
	"fmt"
	"log"
)

// DBType represents the type of database sidecar
type DBType string

const (
	Postgres DBType = "postgres"
	Redis    DBType = "redis"
)

// Provision requests a local container or process to spin up a database for a project
func Provision(projectID string, dbType DBType) (string, error) {
	// In a full implementation, this would use Docker SDK or systemd to spin up the DB
	log.Printf("[Sidecar] Provisioning %s database for project %s...\n", dbType, projectID)
	
	// Mock connection string return
	connStr := fmt.Sprintf("%s://user:pass@localhost:5432/%s_db", dbType, projectID)
	return connStr, nil
}
