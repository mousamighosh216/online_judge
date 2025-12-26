package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"judge_project/api/internal/config"
	"judge_project/api/internal/db"
	httphandler "judge_project/api/internal/http"
	"judge_project/api/internal/queue"
	"judge_project/api/internal/submissions"
)

func main() {
	// 1️⃣ Initialize Postgres
	ctx := context.Background()
	cfg := config.Load()

	pg, err := db.NewPostgres(ctx, "postgres://ojuser:tempPASS@postgres:5432/oj?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	// 2️⃣ Initialize RedisQueue
	redisQueue := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPass)

	// 3️⃣ Retry loop to wait for Postgres
	for i := 0; i < 10; i++ {
		if err := pg.Pool.Ping(context.Background()); err == nil {
			break
		}
		log.Println("Waiting for Postgres...")
		time.Sleep(time.Second)
	}

	// 4️⃣ Retry loop to wait for Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 10; i++ {
		if err := redisQueue.Client.Ping(ctx).Err(); err == nil {
			break
		}
		log.Println("Waiting for Redis...")
		time.Sleep(time.Second)
	}

	// 5️⃣ Initialize SubmissionService
	subSvc := submissions.NewService(pg, redisQueue)

	// 6️⃣ Initialize HTTP handler
	h := httphandler.New(subSvc)

	// 7️⃣ Start HTTP server
	log.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
