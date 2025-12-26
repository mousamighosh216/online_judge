package submissions

import (
	"context"
	"errors"
	"time"

	"judge_project/api/internal/db"
	"judge_project/api/internal/queue"
)

//
// --------------------
// Models
// --------------------
//

// Service is the BUSINESS LOGIC layer for submissions
type Service struct {
	DB    *db.Postgres
	Queue *queue.RedisQueue
}

// Input for creating a submission
type CreateSubmissionInput struct {
	UserID     *int64
	ProblemID  int64
	LanguageID int16
	SourceCode string
}

// Submission entity (used by worker)
type Submission struct {
	ID         int64
	UserID     *int64
	ProblemID  int64
	LanguageID int16
	SourceCode string
	Status     string
}

//
// --------------------
// Business Logic
// --------------------
//

// CreateSubmission
// 1. Validate business rules
// 2. Insert into Postgres
// 3. Push submission ID to Redis queue
func (s *Service) CreateSubmission(
	ctx context.Context,
	in CreateSubmissionInput,
) (int64, error) {

	// ------------- RULE 1: Source code must not be empty
	if len(in.SourceCode) == 0 {
		return 0, errors.New("source code cannot be empty")
	}

	// ------------- RULE 2: Problem must exist
	var problemExists bool
	err := s.DB.Pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM problems WHERE id = $1)`,
		in.ProblemID,
	).Scan(&problemExists)
	if err != nil {
		return 0, err
	}
	if !problemExists {
		return 0, errors.New("problem does not exist")
	}

	// ------------- RULE 3: Language must exist & be active
	var languageActive bool
	err = s.DB.Pool.QueryRow(ctx,
		`SELECT is_active FROM languages WHERE id = $1`,
		in.LanguageID,
	).Scan(&languageActive)
	if err != nil {
		return 0, errors.New("language does not exist")
	}
	if !languageActive {
		return 0, errors.New("language is disabled")
	}

	// ------------- RULE 4: Insert submission
	var submissionID int64
	err = s.DB.Pool.QueryRow(ctx,
		`INSERT INTO submissions
		 (user_id, problem_id, language_id, source_code, status, created_at)
		 VALUES ($1, $2, $3, $4, 'queued', $5)
		 RETURNING id`,
		in.UserID,
		in.ProblemID,
		in.LanguageID,
		in.SourceCode,
		time.Now(),
	).Scan(&submissionID)
	if err != nil {
		return 0, err
	}

	// ------------- RULE 5: Push to Redis queue
	err = s.Queue.EnqueueSubmission(ctx, submissionID)
	if err != nil {
		return 0, err
	}

	return submissionID, nil
}

//
// --------------------
// Worker helpers
// --------------------
//

// GetSubmissionByID is used by worker
func (s *Service) GetSubmissionByID(
	ctx context.Context,
	id int64,
) (*Submission, error) {

	row := s.DB.Pool.QueryRow(ctx,
		`SELECT id, user_id, problem_id, language_id, source_code, status
		 FROM submissions
		 WHERE id = $1`,
		id,
	)

	var sub Submission
	err := row.Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ProblemID,
		&sub.LanguageID,
		&sub.SourceCode,
		&sub.Status,
	)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

// UpdateSubmissionStatus is used by worker after execution
func (s *Service) UpdateSubmissionStatus(
	ctx context.Context,
	id int64,
	status string,
	timeMs *int,
	memoryKb *int,
) error {

	_, err := s.DB.Pool.Exec(ctx,
		`UPDATE submissions
		 SET status = $1,
		     time_ms = $2,
		     memory_kb = $3,
		     finished_at = $4
		 WHERE id = $5`,
		status,
		timeMs,
		memoryKb,
		time.Now(),
		id,
	)

	return err
}

func NewService(db *db.Postgres, q *queue.RedisQueue) *Service {
	return &Service{
		DB:    db,
		Queue: q,
	}
}
