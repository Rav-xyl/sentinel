package api

import (
	"encoding/json"
	"net/http"

	"github.com/Rav-xyl/sentinel/store"
)

// Server holds the dependencies for the API
type Server struct {
	db *store.DB
}

// NewServer creates a new API server
func NewServer(db *store.DB) *Server {
	return &Server{db: db}
}

// RegisterRoutes registers all API endpoints on a given ServeMux
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/routes", s.handleRoutes)
}

func (s *Server) handleRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		routes, err := s.db.GetAllRoutes()
		if err != nil {
			http.Error(w, `{"error": "failed to fetch routes"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(routes)

	case http.MethodPost:
		var req struct {
			Host   string `json:"host"`
			Target string `json:"target"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid json payload"}`, http.StatusBadRequest)
			return
		}
		
		if req.Host == "" || req.Target == "" {
			http.Error(w, `{"error": "host and target are required"}`, http.StatusBadRequest)
			return
		}

		if err := s.db.SetRoute(req.Host, req.Target); err != nil {
			http.Error(w, `{"error": "failed to save route"}`, http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "success"}`))

	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
