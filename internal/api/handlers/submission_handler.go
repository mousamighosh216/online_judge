package handlers

import (
	"encoding/json"
	"net/http"

	"oj/internal/api/services"
)

func CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req services.SubmissionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	id, err := services.CreateSubmission(req)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func GetSubmission(w http.ResponseWriter, r *http.Request) {
	// TODO: parse ID from URL
}