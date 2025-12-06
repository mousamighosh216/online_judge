package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"executor/sandbox" // adjust import path if needed
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: executor <lang> <code_path> <time_limit_sec>")
		os.Exit(1)
	}

	lang := os.Args[1]
	codePath := os.Args[2]
	timeLimitSec, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Invalid time limit: %v\n", err)
		os.Exit(1)
	}

	// Run the code
	result, err := sandbox.Run(lang, codePath, timeLimitSec)
	if err != nil {
		// You can choose to log the error but still return JSON
		fmt.Printf("Error: %v\n", err)
	}

	// Wrap with optional wall time (not accurate, just for structure)
	output := map[string]interface{}{
		"status":    result.Status,
		"stdout":    result.Stdout,
		"stderr":    result.Stderr,
		"exit_code": result.ExitCode,
	}

	jsonData, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonData))
}
