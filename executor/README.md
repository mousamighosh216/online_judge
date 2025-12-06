# Docker-based Executor Template

This minimal executor uses Docker images as per-language sandboxes.
It supports C, C++, and Python out of the box.

## Requirements

- Docker installed and running (Docker Desktop on Windows)
- Go 1.18+ for running the worker (or build the binary with `go build`)

## How to run tests

1. From `executor/` run:
   ```bash
   cd executor
   go run cmd/executor-worker/main.go c /path/to/code.c 2
