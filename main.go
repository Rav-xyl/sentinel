package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rav-xyl/sentinel/proxy"
	"github.com/Rav-xyl/sentinel/store"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	log.Println("🛡️ Sentinel Starting...")

	// Initialize the SQLite database
	db, err := store.InitDB("sentinel.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// Initialize the reverse proxy router with the database
	router := proxy.NewRouter(db)

	// Setup Let's Encrypt / ACME with dynamic host policy
	certManager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			// Fetch all allowed hosts from the database
			hosts, err := db.GetAllHosts()
			if err != nil {
				return fmt.Errorf("failed to retrieve hosts: %v", err)
			}
			
			// Check if the requested host is in our database
			for _, h := range hosts {
				if h == host {
					return nil // Host is allowed
				}
			}
			return fmt.Errorf("acme/autocert: host not configured in sentinel: %s", host)
		},
		Cache: autocert.DirCache("certs"),
	}

	// Start HTTP to HTTPS redirect server with strict timeouts
	go func() {
		log.Println("Listening for HTTP (Port 80) for redirects and ACME challenges...")
		httpServer := &http.Server{
			Addr:         ":80",
			Handler:      certManager.HTTPHandler(nil),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		if err := httpServer.ListenAndServe(); err != nil {
			log.Printf("HTTP Redirect Server Error: %v", err)
		}
	}()

	// Start the main HTTPS server with aggressive timeouts
	port := ":443"
	log.Printf("Listening for Secure HTTPS traffic on %s\n", port)

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Sentinel Secure Server Crashed: %v", err)
	}
}
