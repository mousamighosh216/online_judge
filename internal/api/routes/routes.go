package routes

import (
	"net/http"

	"oj/internal/api/handlers"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/submissions", handlers.CreateSubmission)
	mux.HandleFunc("/submissions/", handlers.GetSubmission)

	return mux
}