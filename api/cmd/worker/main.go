package main

import (
	"context"
	"log"
	"time"

	"judge_project/api/internal/config"
	"judge_project/api/internal/db"
	"judge_project/api/internal/executor"
	"judge_project/api/internal/queue"
	"judge_project/api/internal/submissions"
)

func main() {
	ctx := context.Background()

	// Load env config
	cfg := config.Load()

	// Postgres
	pg, err := db.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect postgres:", err)
	}
	defer pg.Pool.Close()

	// Redis
	rq := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPass)

	// Business logic
	subSvc := &submissions.Service{
		DB:    pg,
		Queue: rq,
	}

	// Docker executor
	exec := executor.New(cfg.WorkDir)

	log.Println("Worker started. Waiting for submissions...")

	for {
		// 1️⃣ Get submission ID from Redis
		subID, err := rq.DequeueSubmission(ctx)
		if err != nil {
			log.Println("queue error:", err)
			time.Sleep(time.Second)
			continue
		}

		log.Println("Processing submission:", subID)

		// 2️⃣ Load submission
		sub, err := subSvc.GetSubmissionByID(ctx, subID)
		if err != nil {
			log.Println("failed to load submission:", err)
			continue
		}

		// ⚠️ TEMP values (we’ll fetch real testcases next)
		result, err := exec.RunSubmission(ctx, executor.Submission{
			ID:              sub.ID,
			Language:        "python", // TEMP — will map from language_id
			SourceCode:      sub.SourceCode,
			InputData:       "1 2\n",
			ExpectedOutput:  "3\n",
			TimeLimitMillis: 2000,
		})
		if err != nil {
			log.Println("execution error:", err)
			continue
		}

		log.Println("Execution result:", result.Status)

		// 3️⃣ Update DB
		err = subSvc.UpdateSubmissionStatus(
			ctx,
			sub.ID,
			result.Status,
			&result.TimeMs,
			nil,
		)
		if err != nil {
			log.Println("failed to update submission:", err)
		}
	}
}
