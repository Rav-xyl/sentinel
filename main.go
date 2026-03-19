package main

import (
	"crypto/tls"
	"log"
	"net/http"

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

	// Hardcode a test route into the database for now
	err = db.SetRoute("localhost", "http://127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to add test route to DB: %v", err)
	}

	// Initialize the reverse proxy router with the database
	router := proxy.NewRouter(db)

	// Setup Let's Encrypt / ACME
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.com", "api.example.com"), // Will be dynamic later
		Cache:      autocert.DirCache("certs"),
	}

	// Start HTTP to HTTPS redirect server
	go func() {
		log.Println("Listening for HTTP (Port 80) for redirects and ACME challenges...")
		if err := http.ListenAndServe(":80", certManager.HTTPHandler(nil)); err != nil {
			log.Printf("HTTP Redirect Server Error: %v", err)
		}
	}()

	// Start the main HTTPS server
	port := ":443"
	log.Printf("Listening for Secure HTTPS traffic on %s\n", port)

	server := &http.Server{
		Addr:    port,
		Handler: router,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Sentinel Secure Server Crashed: %v", err)
	}
}
