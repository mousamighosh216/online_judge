package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ExecResult struct {
	Status   string
	Stdout   string
	Stderr   string
	ExitCode int
}

// Run is the main entry point for running code
func Run(lang, codePath string, timeLimitSec int) (*ExecResult, error) {
	// Create a temporary working directory
	tmpDir, err := os.MkdirTemp("", "exec_job_*")
	if err != nil {
		return nil, err
	}
	// copy code file into tmpDir
	ext := filepath.Ext(codePath)
	dest := filepath.Join(tmpDir, "code"+ext)
	if err := copyFile(codePath, dest); err != nil {
		return nil, err
	}

	// Choose language-specific runner
	switch lang {
	case "c":
		return runC(tmpDir, timeLimitSec)
	case "cpp":
		return runCPP(tmpDir, timeLimitSec)
	case "python":
		return runPython(tmpDir, timeLimitSec)
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// -------------------- Language Runners --------------------

func runC(workdir string, timeLimitSec int) (*ExecResult, error) {
	codeFile := filepath.Join(workdir, "code.c")
	binary := filepath.Join(workdir, "a.out")

	// Compile
	compileCmd := exec.Command("gcc", "-O2", "-std=c17", "-o", binary, codeFile)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr
	if err := compileCmd.Run(); err != nil {
		return &ExecResult{
			Status:   "CE",
			Stdout:   "",
			Stderr:   compileStderr.String(),
			ExitCode: 1,
		}, nil
	}

	// Run with timeout
	return runWithTimeout(binary, timeLimitSec)
}

func runCPP(workdir string, timeLimitSec int) (*ExecResult, error) {
	codeFile := filepath.Join(workdir, "code.cpp")
	binary := filepath.Join(workdir, "a.out")

	compileCmd := exec.Command("g++", "-O2", "-std=c++17", "-o", binary, codeFile)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr
	if err := compileCmd.Run(); err != nil {
		return &ExecResult{
			Status:   "CE",
			Stdout:   "",
			Stderr:   compileStderr.String(),
			ExitCode: 1,
		}, nil
	}

	return runWithTimeout(binary, timeLimitSec)
}

func runPython(workdir string, timeLimitSec int) (*ExecResult, error) {
	codeFile := filepath.Join(workdir, "code.py")
	return runWithTimeout("python", timeLimitSec, codeFile)
}

// -------------------- Helper --------------------

func runWithTimeout(cmdName string, timeoutSec int, args ...string) (*ExecResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdName, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	res := &ExecResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if ctx.Err() == context.DeadlineExceeded {
		res.Status = "TO"
		res.ExitCode = 124
		return res, fmt.Errorf("time limit exceeded")
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			res.ExitCode = exitErr.ExitCode()
			res.Status = "RTE"
			return res, nil
		}
		return nil, err
	}

	res.Status = "OK"
	res.ExitCode = 0
	return res, nil
}
