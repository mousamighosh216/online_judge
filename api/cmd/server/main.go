package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"judge_project/api/internal/config"
	"judge_project/api/internal/db"
	httpHandlers "judge_project/api/internal/http"
	"judge_project/api/internal/queue"
	"judge_project/api/internal/submissions"
)

func main() {
	ctx := context.Background()

	// 1️⃣ Load environment config
	cfg := config.Load()

	// 2️⃣ Connect to Postgres
	pg, err := db.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect postgres:", err)
	}
	defer pg.Pool.Close()

	// 3️⃣ Connect to Redis
	rq := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPass)

	// 4️⃣ Initialize business service
	subSvc := &submissions.Service{
		DB:    pg,
		Queue: rq,
	}

	// 5️⃣ Build HTTP handler (router)
	handler := httpHandlers.New(subSvc)

	// 6️⃣ HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("API server running on http://localhost:8080")

	// 7️⃣ Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error:", err)
	}
}
