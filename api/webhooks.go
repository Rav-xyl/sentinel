package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Rav-xyl/sentinel/builder"
	"github.com/Rav-xyl/sentinel/logger"
	"github.com/Rav-xyl/sentinel/vault"
)

// handleGitHubWebhook processes incoming push events from GitHub
func (s *Server) handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "failed to read body"}`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !verifySignature(secret, signature, body) {
			log.Println("[Webhook] Invalid GitHub signature detected")
			http.Error(w, `{"error": "invalid signature"}`, http.StatusUnauthorized)
			return
		}
	}

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

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "accepted", "message": "deployment queued"}`))

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	branch = strings.TrimPrefix(branch, "refs/tags/")

	// Asynchronously process the deployment
	go s.processDeployment(payload.Repository.Name, payload.Repository.CloneURL, branch)
}

func (s *Server) processDeployment(projectName, repoURL, branch string) {
	log.Printf("[CI/CD] Starting deployment pipeline for %s...\n", projectName)

	// 1. Setup Logger
	stream, err := logger.NewStream(projectName)
	if err != nil {
		log.Printf("[CI/CD Error] Failed to create logger for %s: %v\n", projectName, err)
		return
	}
	stream.Log("Deployment Pipeline Initiated")

	// 2. Load Secrets
	envConfig, err := vault.LoadEnv(projectName)
	if err != nil {
		stream.Log(fmt.Sprintf("Warning: Failed to load secrets: %v", err))
	} else {
		stream.Log("Secrets loaded from Vault")
	}

	// 3. Create Builder Job
	// In a real system, BuildCmd and StartCmd would be read from a sentinel.yaml in the repo.
	// For this sprint, we assume a standard Node.js app pattern.
	job := builder.NewJob(repoURL, branch, projectName)
	job.BuildCmd = "npm install && npm run build" // Default assumption
	job.StartCmd = "npm start"
	if envConfig != nil {
		job.Env = envConfig.Vars
	}

	// 4. Execute Build
	stream.Log("Cloning and Building...")
	if err := job.Execute(); err != nil {
		stream.Log(fmt.Sprintf("Build Failed: %v", err))
		return
	}
	stream.Log("Build Successful!")

	// 5. Start Application (Find an open port in production, hardcoding for sprint demo)
	targetPort := "3000" // Example
	job.StartCmd = fmt.Sprintf("PORT=%s npm start", targetPort)
	
	cmd, err := job.Start()
	if err != nil {
		stream.Log(fmt.Sprintf("Failed to start application: %v", err))
		return
	}
	
	// Keep track of the cmd process to kill it on redeploy (omitted for brevity, requires a ProcessManager map)
	_ = cmd

	stream.Log(fmt.Sprintf("Application started on port %s", targetPort))

	// 6. Update Routing Table dynamically
	host := fmt.Sprintf("%s.example.com", projectName) // Default domain pattern
	target := fmt.Sprintf("http://127.0.0.1:%s", targetPort)
	
	if err := s.db.SetRoute(host, target); err != nil {
		stream.Log(fmt.Sprintf("Failed to update router: %v", err))
		return
	}

	stream.Log(fmt.Sprintf("Routing updated: %s -> %s", host, target))
	stream.Log("DEPLOYMENT COMPLETE 🚀")
}

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
