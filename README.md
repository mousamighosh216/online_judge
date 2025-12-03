# online_judge

## 📂 Project Structure Details

<details>
<summary>Click to view the full directory tree</summary>

```text
oj/
├─ .github/
│  └─ workflows/
├─ infra/
│  ├─ k8s/
│  └─ docker/
├─ docker-compose.yml
├─ docker-compose.yml
├─ Dockerfile               # for the main API service
├─ docker-entrypoint.sh
├─ LICENSE
├─ README.md
├─ .env.example
├─ build/                   # CI build scripts and artifacts
│  └─ ci-scripts/
├─ api/                     # HTTP API service (Go)
│  ├─ cmd/
│  │  └─ judge0-api/        # main package entry
│  ├─ internal/
│  │  ├─ server/            # HTTP handlers, middleware
│  │  ├─ submissions/       # submission model + DB interactions
│  │  ├─ languages/         # language metadata (limits, compile/run commands)
│  │  ├─ workers/           # queueing, worker registry clients
│  │  └─ auth/              # API auth, rate limiting
│  ├─ pkg/                  # reusable packages (if needed)
│  ├─ configs/              # config structs, env parsing
│  └─ go.mod
├─ executor/                # code that runs user submissions (sandboxed)
│  ├─ cmd/                  # entry for executor worker
│  ├─ languages/            # language-specific runner wrappers
│  │  ├─ cpp/
│  │  ├─ runtime-python/
│  │  └─ java/
│  ├─ sandbox/              # sandbox driver (nsjail/firecracker/containers)
│  ├─ tests/                # harness for local executor tests
│  └─ go.mod
├─ sandbox-images/          # Dockerfiles or OCI images used for execution
│  ├─ base/                 # base images (with compilers, runtimes)
│  └─ slim/                 # minimal images for faster cold starts
├─ worker/                  # worker orchestration (queue consumers)
│  ├─ cmd/
│  ├─ handlers/             # how to run, collect results, store logs
│  └─ go.mod
├─ web/                     # optional front-end (React/Vite) or admin UI
│  ├─ public/
│  └─ src/
├─ db/                      # migrations, schema, seeds
│  ├─ migrations/
│  └─ schema.sql
├─ scripts/                 # helper scripts (setup, benchmark, admin)
│  ├─ setup_local.sh
│  └─ create_db_user.sh
├─ docs/                    # docs, API spec, architecture diagrams
│  ├─ architecture.md
│  └─ api.md
├─ tests/                   # integration / E2E tests (calls API executor)
│  └─ e2e/
├─ tools/                   # local dev tools (formatters, linters)
├─ logs/                    # example logs / rotation config (gitignored)
└─ .gitignore