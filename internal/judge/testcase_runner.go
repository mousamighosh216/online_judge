package judge

import (
	"fmt"

	"github.com/mousamighosh216/oj/internal/executor"
	"github.com/mousamighosh216/oj/internal/models"
)

// TestCaseMap keys each test case by a human-readable identifier
// e.g. "tc_1", "tc_2" — preserving insertion order via a slice of keys.
type TestCaseMap struct {
	Keys  []string
	Cases map[string]models.TestCase
}

// NewTestCaseMap builds an ordered map from a slice of TestCase.
// Keys are generated as "tc_1", "tc_2", ... matching the problem's
// test case index so failures are traceable back to the original set.
func NewTestCaseMap(testcases []models.TestCase) *TestCaseMap {
	m := &TestCaseMap{
		Keys:  make([]string, 0, len(testcases)),
		Cases: make(map[string]models.TestCase, len(testcases)),
	}
	for i, tc := range testcases {
		key := fmt.Sprintf("tc_%d", i+1)
		m.Keys = append(m.Keys, key)
		m.Cases[key] = tc
	}
	return m
}

// RunResult holds the outcome of a single test case execution.
type RunResult struct {
	Key    string // which test case failed, e.g. "tc_3"
	Result models.TestResult
	Passed bool
}

// RunAllTestcases iterates the TestCaseMap in insertion order.
// On the first non-AC verdict it breaks the loop and returns:
//   - results collected so far (including the failing one)
//   - the key of the failing test case
//   - a non-nil error describing the failure
//
// If all test cases pass, failedKey is empty and err is nil.
func RunAllTestcases(
	tcMap *TestCaseMap,
	ws *executor.Workspace,
	lang models.Language,
) (results []RunResult, failedKey string, err error) {

	results = make([]RunResult, 0, len(tcMap.Keys))

	for _, key := range tcMap.Keys {
		tc := tcMap.Cases[key]

		result, runErr := RunSingleTestcase(key, tc, ws, lang)
		results = append(results, result)

		if runErr != nil {
			// hard execution error (sandbox crash, OOM, etc.)
			return results, key, fmt.Errorf("test case %s: execution error: %w", key, runErr)
		}

		if !result.Passed {
			// wrong answer, TLE, MLE, RE — verdict already set on result
			return results, key, fmt.Errorf(
				"test case %s: got verdict %s (expected AC)",
				key, result.Result.Verdict,
			)
		}
	}

	return results, "", nil
}

// RunSingleTestcase runs one test case through the sandbox, compares
// its output, and returns a RunResult with the verdict populated.
// A non-nil error means the sandbox itself failed (not a wrong answer).
func RunSingleTestcase(
	key string,
	tc models.TestCase,
	ws *executor.Workspace,
	lang models.Language,
) (RunResult, error) {

	// write input file for this specific test case
	if err := ws.WriteInputFile(tc, key); err != nil {
		return RunResult{Key: key}, fmt.Errorf("writing input for %s: %w", key, err)
	}

	// run inside sandbox
	sandboxOut, err := ws.Sandbox.RunContainer(
		lang.RunCommand(ws),
		ws,
		tc.Limits,
	)
	if err != nil {
		// sandbox-level failure: OOM kill, seccomp violation, etc.
		return RunResult{Key: key}, fmt.Errorf("sandbox error on %s: %w", key, err)
	}

	// compare output
	passed := CompareTrimmed(sandboxOut.Stdout, tc.ExpectedOutput)

	// derive verdict from comparison + resource usage + exit code
	verdict := GetVerdict(GetVerdictInput{
		Matched:     passed,
		ExitCode:    sandboxOut.ExitCode,
		WallTimeMs:  sandboxOut.WallTimeMs,
		MemKb:       sandboxOut.MemKb,
		TimeLimitMs: tc.Limits.TimeLimitMs,
		MemLimitKb:  tc.Limits.MemLimitKb,
	})

	result := RunResult{
		Key:    key,
		Passed: verdict == models.VerdictAC,
		Result: models.TestResult{
			TestCaseIndex: tc.Index,
			Verdict:       verdict,
			TimeMs:        sandboxOut.WallTimeMs,
			MemKb:         sandboxOut.MemKb,
			Stdout:        sandboxOut.Stdout,
			Stderr:        sandboxOut.Stderr,
		},
	}

	return result, nil
}
