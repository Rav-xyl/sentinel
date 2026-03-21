package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// handleGitHubWebhook processes incoming push events from GitHub
func (s *Server) handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "failed to read body"}`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Verify the signature if a secret is configured
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !verifySignature(secret, signature, body) {
			log.Println("[Webhook] Invalid GitHub signature detected")
			http.Error(w, `{"error": "invalid signature"}`, http.StatusUnauthorized)
			return
		}
	}

	// Parse the event payload
	var payload struct {
		Ref        string `json:"ref"`
		Repository struct {
			Name     string `json:"name"`
			CloneURL string `json:"clone_url"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, `{"error": "invalid payload"}`, http.StatusBadRequest)
		return
	}

	// Acknowledge receipt immediately before starting the build
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "accepted", "message": "deployment queued"}`))

	// TODO (Day 10): Trigger the actual clone and build process asynchronously
	log.Printf("[Webhook] Received push event for %s on branch %s. Queuing deployment...\n", payload.Repository.Name, payload.Ref)
}

// verifySignature checks the HMAC hex digest against the provided GitHub signature
func verifySignature(secret, headerSignature string, body []byte) bool {
	const signaturePrefix = "sha256="
	if !strings.HasPrefix(headerSignature, signaturePrefix) {
		return false
	}

	signature := headerSignature[len(signaturePrefix):]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
