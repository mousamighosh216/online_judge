package executor

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type Submission struct {
	ID              int64
	Language        string
	SourceCode      string
	InputData       string
	ExpectedOutput  string
	TimeLimitMillis int
}

type Result struct {
	Status string
	TimeMs int
	Output string
}

type Executor struct {
	WorkDir string
}

func New(workDir string) *Executor {
	return &Executor{WorkDir: workDir}
}

func (e *Executor) RunSubmission(ctx context.Context, sub Submission) (*Result, error) {
	start := time.Now()

	// DEV MODE (Windows / Mac)
	if runtime.GOOS != "linux" {
		return &Result{
			Status: "accepted",
			TimeMs: 1,
			Output: "dev-mode-output",
		}, nil
	}

	// Linux path (real execution later)
	tmpDir := filepath.Join(e.WorkDir, "sub-"+time.Now().Format("150405"))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := filepath.Join(tmpDir, "main.py")
	if err := os.WriteFile(sourceFile, []byte(sub.SourceCode), 0644); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "python3", sourceFile)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return &Result{
			Status: "runtime_error",
			TimeMs: elapsed(start),
			Output: string(out),
		}, nil
	}

	return &Result{
		Status: "accepted",
		TimeMs: elapsed(start),
		Output: string(out),
	}, nil
}

func elapsed(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
