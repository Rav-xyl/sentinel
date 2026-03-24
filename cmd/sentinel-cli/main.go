package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const apiURL = "http://localhost:8081/api/routes"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		listRoutes()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Sentinel CLI 🛡️

Usage:
  sentinel-cli <command>

Commands:
  list    List all active routing rules from the Sentinel core
  help    Show this help message`)
}

func listRoutes() {
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("Failed to connect to Sentinel API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API returned error status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var routes map[string]string
	if err := json.Unmarshal(body, &routes); err != nil {
		log.Fatalf("Failed to parse routes: %v", err)
	}

	fmt.Println("🌐 Active Sentinel Routes:")
	fmt.Println("-------------------------")
	if len(routes) == 0 {
		fmt.Println("No routes currently configured.")
		return
	}

	for host, target := range routes {
		fmt.Printf(" %s  -->  %s\n", host, target)
	}
}
