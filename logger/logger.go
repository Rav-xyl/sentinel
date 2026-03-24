package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Stream represents a log stream for a specific project
type Stream struct {
	mu      sync.Mutex
	ProjectID string
	Writer  io.Writer
}

// NewStream initializes a log stream that writes to both stdout and a file
func NewStream(projectID string) (*Stream, error) {
	logPath := fmt.Sprintf("deployments/%s/build.log", projectID)
	
	// Ensure directory exists
	if err := os.MkdirAll(fmt.Sprintf("deployments/%s", projectID), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// MultiWriter allows logging to terminal and file simultaneously
	writer := io.MultiWriter(os.Stdout, file)

	return &Stream{
		ProjectID: projectID,
		Writer:    writer,
	}, nil
}

// Log writes a formatted message to the stream
func (s *Stream) Log(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Fprintf(s.Writer, "[%s] %s\n", s.ProjectID, message)
}
