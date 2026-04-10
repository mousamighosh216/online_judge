package main

import (
	"log"
	"net/http"

	"oj/internal/api/routes"
)

func main() {
	router := routes.SetupRoutes()

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
