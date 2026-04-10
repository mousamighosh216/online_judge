// NewSubmission(), Validate(), SetVerdict(v Verdict)

package models

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Status string
type Verdict string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"

	VerdictAC  Verdict = "AC"
	VerdictWA  Verdict = "WA"
	VerdictTLE Verdict = "TLE"
	VerdictMLE Verdict = "MLE"
	VerdictRE  Verdict = "RE"
	VerdictCE  Verdict = "CE"
)

// SupportedLanguages is the set of languages the system can execute.
// Populated at startup from config/languages.json.
var SupportedLanguages = map[string]bool{
	"cpp":    true,
	"python": true,
	"java":   true,
}

// MaxSourceBytes is the hard cap on submitted source code size.
const MaxSourceBytes = 64 * 1024 // 64 KB

type TestResult struct {
	TestCaseIndex int     `json:"index"`
	Verdict       Verdict `json:"verdict"`
	TimeMs        int64   `json:"time_ms"`
	MemKb         int64   `json:"mem_kb"`
	Stdout        string  `json:"stdout,omitempty"`
	Stderr        string  `json:"stderr,omitempty"`
}

type Submission struct {
	mu sync.Mutex // guards Status and FinalVerdict

	ID           string       `json:"id"`
	ProblemID    string       `json:"problem_id"`
	UserID       string       `json:"user_id"`
	Language     string       `json:"language"`
	SourceCode   string       `json:"source_code"`
	Status       Status       `json:"status"`
	FinalVerdict Verdict      `json:"final_verdict,omitempty"`
	Results      []TestResult `json:"results,omitempty"`
	MaxTimeMs    int64        `json:"max_time_ms,omitempty"`
	MaxMemKb     int64        `json:"max_mem_kb,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// NewSubmission constructs a Submission with a fresh UUID, default
// status "queued", and both timestamps set to now (UTC).
// The caller must still set ProblemID, UserID, Language, SourceCode
// before persisting.
func NewSubmission(problemID, userID, language, sourceCode string) *Submission {
	now := time.Now().UTC()
	return &Submission{
		ID:         uuid.New().String(),
		ProblemID:  problemID,
		UserID:     userID,
		Language:   language,
		SourceCode: sourceCode,
		Status:     StatusQueued,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Validate checks the submission fields before the record is persisted
// or enqueued. It does NOT hit the database — problem existence is
// checked by the service layer via the repository.
//
// Returns a joined error listing every violation, or nil if valid.
func (s *Submission) Validate(problemExists bool) error {
	var errs []error

	if s.ProblemID == "" {
		errs = append(errs, errors.New("problem_id is required"))
	}
	if !problemExists {
		errs = append(errs, errors.New("problem does not exist"))
	}

	if s.UserID == "" {
		errs = append(errs, errors.New("user_id is required"))
	}

	if !SupportedLanguages[s.Language] {
		errs = append(errs, errors.New("unsupported language: "+s.Language))
	}

	if len(s.SourceCode) == 0 {
		errs = append(errs, errors.New("source_code is empty"))
	}
	if len(s.SourceCode) > MaxSourceBytes {
		errs = append(errs, errors.New("source_code exceeds 64 KB limit"))
	}

	return errors.Join(errs...)
}

// SetVerdict atomically transitions the submission to StatusDone and
// records the final verdict. It is safe to call from the worker
// goroutine while the HTTP handler may be reading Status concurrently.
func (s *Submission) SetVerdict(v Verdict) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = StatusDone
	s.FinalVerdict = v
	s.UpdatedAt = time.Now().UTC()
}
