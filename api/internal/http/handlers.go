package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"judge_project/api/internal/db"
	"judge_project/api/internal/submissions"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	SubmissionService *submissions.Service
	DB                *db.Postgres // or *sql.DB
	Redis             *redis.Client
}

// New creates the HTTP handler (router + handlers)
func New(subSvc *submissions.Service) http.Handler {
	h := &Handler{
		SubmissionService: subSvc,
	}

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/ready", h.ready)
	mux.HandleFunc("/submissions", h.createSubmission)
	mux.HandleFunc("/submissions/", h.getSubmission)

	return mux
}

// --------------------
// Handlers
// --------------------

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := map[string]string{
		"postgres": "ok",
		"redis":    "ok",
		"status":   "ready",
	}

	// 1️⃣ Check Postgres
	if err := h.SubmissionService.DB.Pool.Ping(ctx); err != nil {
		status["postgres"] = "down"
		status["status"] = "not_ready"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		// 2️⃣ Check Redis
		if err := h.SubmissionService.Queue.Client.Ping(ctx).Err(); err != nil {
			status["redis"] = "down"
			status["status"] = "not_ready"
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// POST /submissions
func (h *Handler) createSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProblemID  int64  `json:"problem_id"`
		LanguageID int16  `json:"language_id"`
		SourceCode string `json:"source_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	id, err := h.SubmissionService.CreateSubmission(
		r.Context(),
		submissions.CreateSubmissionInput{
			ProblemID:  req.ProblemID,
			LanguageID: req.LanguageID,
			SourceCode: req.SourceCode,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"submission_id": id,
		"status":        "queued",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// GET /submissions/{id}
func (h *Handler) getSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/submissions/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid submission id", http.StatusBadRequest)
		return
	}

	sub, err := h.SubmissionService.GetSubmissionByID(r.Context(), id)
	if err != nil {
		http.Error(w, "submission not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}
