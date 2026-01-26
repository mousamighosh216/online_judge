package executor

import (
	"context"
	"testing"
)

func TestPythonExecution(t *testing.T) {
	exec := New("/tmp")

	res, err := exec.RunSubmission(context.Background(), Submission{
		Language:   "python",
		SourceCode: "print(1+2)",
		InputData:  "",
	})

	if err != nil {
		t.Fatal(err)
	}

	if res.Output != "3\n" {
		t.Fatalf("expected 3, got %q", res.Output)
	}
}
