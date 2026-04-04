# 🚀 Online Judge System (Go-based)

A scalable, modular **Online Judge (OJ)** system inspired by modern execution engines like Judge0, but designed with **true multi-test execution, better abstraction, and production-ready architecture**.

---

# 🧠 Overview

This project is a **distributed system** that safely executes untrusted user code and evaluates it against multiple test cases.

### Core Flow

```
Client
 → API (Go)
   → Database (Postgres)
   → Queue (Redis)
     → Scheduler
       → Worker
         → Executor Pipeline
           → Docker Sandbox
           → Compile + Run
           → Judge
   ← Result stored
 ← Client fetches result
```

---

# 🏗️ Directory Structure

```
online-judge/
│
├── cmd/
│   ├── api/                 # API entrypoint
│   ├── worker/              # Worker entrypoint
│   ├── scheduler/           # Scheduler entrypoint
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── repository/
│   │   └── routes/
│   │
│   ├── scheduler/
│   │   ├── queue/
│   │   └── scheduler.go
│   │
│   ├── worker/
│   │   ├── consumer/
│   │   └── worker.go
│   │
│   ├── executor/
│   │   ├── pipeline.go
│   │   ├── workspace.go
│   │   └── result.go
│   │
│   ├── sandbox/
│   │   ├── docker_runner.go
│   │   └── sandbox.go
│   │
│   ├── languages/
│   │   ├── base.go
│   │   ├── cpp.go
│   │   ├── python.go
│   │   └── java.go
│   │
│   ├── judge/
│   │   ├── comparator.go
│   │   ├── verdict.go
│   │   └── testcase_runner.go
│   │
│   ├── models/
│   │   ├── submission.go
│   │   ├── problem.go
│   │   └── testcase.go
│   │
│   ├── infra/
│   │   ├── db/
│   │   ├── redis/
│   │   └── logger/
│
├── config/
│   ├── languages.json
│   ├── limits.json
│
├── docker/
│   ├── Dockerfile.executor
│   └── docker-compose.yml
│
└── README.md
```

---

# ⚙️ Component Breakdown (Functions + Responsibilities)

---

## 🔹 `cmd/api/main.go`

### Responsibilities:

* Boot API server
* Initialize DB + Redis
* Register routes

---

## 🔹 `internal/api/handlers/submission_handler.go`

### Functions:

```
CreateSubmission()
GetSubmission()
ListSubmissions()
```

### Responsibility:

* Handle HTTP requests
* Validate input
* Call service layer

---

## 🔹 `internal/api/services/submission_service.go`

### Functions:

```
CreateSubmissionRecord()
EnqueueSubmission()
ValidateSubmission()
```

### Responsibility:

* Business logic
* Push job to queue

---

## 🔹 `internal/api/repository/submission_repo.go`

### Functions:

```
InsertSubmission()
UpdateSubmissionStatus()
GetSubmissionByID()
```

### Responsibility:

* Database operations

---

## 🔹 `internal/scheduler/scheduler.go`

### Functions:

```
DispatchJobs()
ApplyRateLimiting()
PrioritizeSubmissions()
```

### Responsibility:

* Control execution order
* Prevent overload

---

## 🔹 `internal/worker/consumer/consumer.go`

### Functions:

```
ConsumeQueue()
RetryFailedJobs()
AcknowledgeJob()
```

### Responsibility:

* Pull jobs from Redis
* Pass to worker

---

## 🔹 `internal/worker/worker.go`

### Functions:

```
ProcessSubmission()
```

### Responsibility:

* Orchestrate execution lifecycle

---

## 🔥 `internal/executor/pipeline.go` (CORE)

### Functions:

```
ExecuteSubmission()
Compile()
RunTestcases()
AggregateResults()
Cleanup()
```

### Responsibility:

* Core execution pipeline

---

## 🔹 `internal/executor/workspace.go`

### Functions:

```
CreateWorkspace()
WriteSourceCode()
WriteInputFile()
CleanupWorkspace()
```

### Responsibility:

* Manage temporary files

---

## 🔹 `internal/sandbox/docker_runner.go`

### Functions:

```
RunContainer()
BuildDockerCommand()
ApplyLimits()
```

### Responsibility:

* Execute code inside Docker
* Enforce:

  * memory limits
  * CPU limits
  * process limits

---

## 🔹 `internal/languages/base.go`

### Functions:

```
CompileCommand()
RunCommand()
FileName()
```

---

## 🔹 `internal/languages/cpp.go`

### Functions:

```
CompileCommand() → g++
RunCommand() → ./binary
```

---

## 🔹 `internal/languages/python.go`

### Functions:

```
RunCommand() → python3 script.py
```

---

## 🔹 `internal/judge/testcase_runner.go`

### Functions:

```
RunAllTestcases()
RunSingleTestcase()
```

### Responsibility:

* Execute per test case
* Early exit on failure

---

## 🔹 `internal/judge/comparator.go`

### Functions:

```
CompareExact()
CompareTrimmed()
CompareFloating()
```

---

## 🔹 `internal/judge/verdict.go`

### Functions:

```
GetVerdict()
MapExitCode()
```

---

## 🔹 `internal/models/submission.go`

### Structure:

```
Submission {
  ID
  SourceCode
  Language
  Status
  Results[]
}
```

---

## 🔹 `internal/models/testcase.go`

```
TestCase {
  Input
  ExpectedOutput
  TimeLimit
  MemoryLimit
}
```

---

# 🐳 Docker Execution (Sandbox)

Each submission runs inside a container:

```
docker run --rm \
  --memory=128m \
  --cpus=0.5 \
  --pids-limit=64 \
  --network=none \
  executor-image
```

---

# ⚖️ Judge0 vs This System

| Feature           | Judge0                | This System            |
| ----------------- | --------------------- | ---------------------- |
| Execution Model   | 1 test per submission | Multi-test batch       |
| Compilation       | per test ❌            | once per submission ✅  |
| Scheduler         | FIFO only ❌           | priority-based ✅       |
| Language handling | JSON config           | Go abstraction         |
| Pipeline          | shell scripts         | structured Go pipeline |
| Sandbox           | isolate               | Docker (extendable)    |
| Performance       | moderate              | optimized              |
| Flexibility       | high                  | high + structured      |

---

# 🧠 Key Differences

### Judge0:

* Simple execution API
* Stateless per submission
* Easier to integrate

### This System:

* Full Online Judge
* Handles:

  * multiple test cases
  * scheduling
  * fairness
  * extensibility

---

# 🚀 Tech Stack

* **Language:** Go (Golang)
* **Queue:** Redis
* **Database:** PostgreSQL
* **Sandbox:** Docker (Linux cgroups)
* **Orchestration:** Docker Compose

---

# ⚠️ Known Bottlenecks

* High concurrency → CPU contention
* Disk I/O (workspace creation)
* Docker startup latency
* Queue backlog under heavy load

---

# 🎯 Future Improvements

* Container pooling
* Distributed workers
* Custom checkers
* Subtask scoring
* Kubernetes deployment

---

# 💣 Final Note

This project is not just “run code”.

It is:

> **A distributed, adversarial-safe execution system with deterministic evaluation**

---

# 📌 Getting Started

```
docker-compose up --build
```

---

