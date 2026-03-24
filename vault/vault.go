package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config represents a collection of environment variables
type Config struct {
	Vars []string
}

// LoadEnv reads a .env file from the secrets directory and parses it into a slice of "KEY=VALUE"
func LoadEnv(projectID string) (*Config, error) {
	// Construct path to the protected secrets directory
	path := fmt.Sprintf("secrets/%s.env", projectID)
	
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Vars: []string{}}, nil // No env file is not an error
		}
		return nil, fmt.Errorf("failed to open env file: %w", err)
	}
	defer file.Close()

	var vars []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Validate KEY=VALUE format
		if !strings.Contains(line, "=") {
			continue
		}
		
		vars = append(vars, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading env file: %w", err)
	}

	return &Config{Vars: vars}, nil
}

// GetMap returns the variables as a map for easy manipulation
func (c *Config) GetMap() map[string]string {
	m := make(map[string]string)
	for _, v := range c.Vars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
