package services

// "oj/internal/api/repository"

type SubmissionRequest struct {
	SourceCode string `json:"source_code"`
	Language   string `json:"language"`
}

func CreateSubmission(req SubmissionRequest) (string, error) {
	return "hey", nil
}
