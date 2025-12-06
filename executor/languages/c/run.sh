#!/bin/bash
# run.sh: Executed inside the sandbox to run the compiled binary.
# NOTE: The actual time/memory limits should be enforced by the Go worker calling isolate!

# Arguments expected from the Go worker:
# $1 = Input file (e.g., /box/input.txt)
# $2 = Output file (e.g., /box/output.txt)

INPUT_FILE="${1:-/dev/null}"   # Use /dev/null if input is not provided
OUTPUT_FILE="${2:-/dev/null}" # Use /dev/null if output is not specified

# The ulimit commands are redundant if isolate is used for resource control, 
# but are included here for defense-in-depth.
ulimit -t 2 
ulimit -v 65536

# Run the compiled binary, redirecting standard input and standard output.
# The exit status of the binary will be captured by the outer 'isolate' process.
./main < "$INPUT_FILE" > "$OUTPUT_FILE"

# The script exits with the exit code of the compiled program (./main).