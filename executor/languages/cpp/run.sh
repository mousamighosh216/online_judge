#!/bin/sh
# run.sh - executed inside box
ulimit -t 2
ulimit -v 65536
#run the compiled binary
./main