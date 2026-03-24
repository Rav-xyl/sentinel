package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Job represents a single deployment task
type Job struct {
	RepoURL    string
	Branch     string
	ProjectID  string
	BuildCmd   string
	StartCmd   string
	Env        []string
	Workspace  string
}

// NewJob creates a new deployment job
func NewJob(repoURL, branch, projectID string) *Job {
	return &Job{
		RepoURL:   repoURL,
		Branch:    branch,
		ProjectID: projectID,
		Workspace: filepath.Join("deployments", projectID),
	}
}

// Execute handles the full Clone -> Build lifecycle
func (j *Job) Execute() error {
	// 1. Prepare Workspace
	if err := os.MkdirAll(j.Workspace, 0755); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// 2. Clone Repository
	log.Printf("[Builder] Cloning %s (%s) into %s...\n", j.RepoURL, j.Branch, j.Workspace)
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "--branch", j.Branch, j.RepoURL, ".")
	cloneCmd.Dir = j.Workspace
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("clone failed: %s %w", string(output), err)
	}

	// 3. Build (if command provided)
	if j.BuildCmd != "" {
		log.Printf("[Builder] Running build: %s\n", j.BuildCmd)
		buildCmd := exec.Command("sh", "-c", j.BuildCmd)
		if os.PathSeparator == '\\' {
			buildCmd = exec.Command("cmd", "/C", j.BuildCmd)
		}
		buildCmd.Dir = j.Workspace
		buildCmd.Env = append(os.Environ(), j.Env...)
		if output, err := buildCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build failed: %s %w", string(output), err)
		}
	}

	return nil
}

// Start spawns the application process
func (j *Job) Start() (*exec.Cmd, error) {
	if j.StartCmd == "" {
		return nil, fmt.Errorf("no start command provided")
	}

	log.Printf("[Builder] Starting application: %s\n", j.StartCmd)
	cmd := exec.Command("sh", "-c", j.StartCmd)
	if os.PathSeparator == '\\' {
		cmd = exec.Command("cmd", "/C", j.StartCmd)
	}
	cmd.Dir = j.Workspace
	cmd.Env = append(os.Environ(), j.Env...)
	
	// Day 12: Connect stdout/stderr for log streaming
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	return cmd, nil
}
