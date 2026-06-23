# Uptime Monitor

A cloud-native uptime monitoring service built on GCP as a learning project — the domain is intentionally simple so the focus stays on infrastructure, deployment, and architecture.

## What It Does

Registers URLs, schedules periodic health checks via Cloud Scheduler and Cloud Tasks, and records results. Each registered site gets a Cloud Task enqueued per polling cycle; the worker hits the URL and records status and latency.

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/sites` | List registered sites |
| `POST` | `/sites/add` | Register a site |
| `POST` | `/sites/remove` | Remove a site |
| `POST` | `/sites/poll/enqueue` | Enqueue a poll task per site (called by Cloud Scheduler) |
| `POST` | `/sites/poll/worker` | Worker endpoint called by Cloud Tasks per site |

## Stack

| Layer | Technology |
|-------|-----------|
| Language | Go |
| Router | chi |
| Database | PostgreSQL (Cloud SQL) |
| DB Access | pgx + sqlc (type-safe generated queries) |
| Migrations | goose |
| Container | Docker (multi-stage, hardened, non-root) |
| Hosting | Cloud Run |
| Queue | Cloud Tasks |
| Scheduler | Cloud Scheduler |
| IaC | Terraform |
| CI/CD | Cloud Build (GitOps — push to main triggers build + deploy + migrate) |
| Secrets | Secret Manager |

## Architecture

```
Cloud Scheduler (1 min)
        │
        ▼
POST /sites/poll/enqueue
        │  (one task per registered site)
        ▼
Cloud Tasks Queue
        │
        ▼
POST /sites/poll/worker
        │
        ▼
HTTP GET → target URL → record result
```

## Infrastructure

All GCP resources are managed by Terraform:

- **Cloud Run** — hosts the Go service
- **Cloud SQL (Postgres 18)** — stores sites and check results; uses IAM database authentication (no passwords for the app)
- **Cloud Tasks** — queue with rate limits, exponential backoff, and automatic retries
- **Cloud Scheduler** — triggers the enqueue endpoint on a 1-minute cron
- **Secret Manager** — holds migration credentials and Terraform state config
- **Artifact Registry** — stores Docker images
- **IAM** — least-privilege service accounts for runtime and CI/CD

## CI/CD Pipeline

Every push to `main` runs:

1. Build Docker image
2. Push to Artifact Registry
3. `terraform apply` — provisions or updates all infrastructure
4. Run `goose up` migrations via Cloud SQL Proxy

Terraform state is stored remotely in GCS.
