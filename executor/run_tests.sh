#!/bin/bash
echo "=== Running Executor Tests ==="

# Array of test cases: language, source file, timeout
tests=(
  "c tests/test_c.c 2"
  "cpp tests/test_cpp.cpp 2"
  "python tests/test_py.py 2"
)

for t in "${tests[@]}"; do
    IFS=' ' read -r lang file timeout <<< "$t"
    echo -n "Test: $lang Program ... "
    
    # Run the executor and capture JSON output
    result=$(go run ./cmd/executor-worker/main.go $lang $file $timeout 2>/dev/null)
    
    # Extract exit code and stdout from JSON
    exit_code=$(echo "$result" | grep -oP '(?<="exit_code": )\d+')
    stdout=$(echo "$result" | grep -oP '(?<="stdout": ")[^"]*')
    
    # Print PASS/FAIL
    if [ "$exit_code" -eq 0 ]; then
        echo "PASS"
        echo "Output: $stdout"
    else
        echo "FAIL"
        echo "Result: $result"
    fi
done
