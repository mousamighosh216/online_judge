#!/bin/bash
# compile.sh - Executes C compilation inside the isolate sandbox.

# 1. Configuration variables (Passed from Go worker)
# These would typically be passed as arguments, but we'll hardcode them for demonstration.
# ARG1: Source file path, ARG2: Output executable name
SOURCE_FILE="${1:-main.c}"
OUTPUT_FILE="${2:-main}"

# 2. Isolate Initialization and Compilation
# The actual compilation command must be wrapped by 'isolate --run'
# We use /usr/bin/gcc and /usr/local/bin/isolate (absolute paths)
# We set time/mem limits for compilation itself (e.g., 5 seconds, 100MB)

/usr/local/bin/isolate --init --box-id 1 &> /dev/null

/usr/local/bin/isolate --run --box-id 1 \
    --time 5 --mem 102400 \
    --meta /tmp/compile_meta.txt \
    --dir=/app=/app:rw \
    -- /usr/bin/gcc "$SOURCE_FILE" -O2 -static -s -o "$OUTPUT_FILE" 2> /app/compile.stderr

# 3. Check Result Status
STATUS=$(grep status /tmp/compile_meta.txt | cut -d: -f2 | tr -d '[:space:]')

if [ "$STATUS" = "OK" ]; then
    echo "Compilation successful. Status: OK"
    # Exit with code 0 for the Go worker to proceed to run phase
    EXIT_CODE=0
else
    echo "Compilation failed. Status: $STATUS"
    # Exit with code 1 for the Go worker to report Compilation Error (CE)
    EXIT_CODE=1
fi

# 4. Cleanup and Exit
/usr/local/bin/isolate --cleanup --box-id 1 &> /dev/null
exit $EXIT_CODE