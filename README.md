# URL Uptime Monitor

A Go backend service that pings registered URLs on a schedule and notifies users via HMAC-signed webhooks when a URL transitions between up and down.

## Concept

One sentence: **It pings URLs you register and alerts you via webhook when they go up or down.**

Think Pingdom / UptimeRobot, stripped down. The domain is intentionally boring — the code and architecture are the point, not the problem.

## Technologies

| Tech | Role |
|------|------|
| **Go** | Language |
| **PostgreSQL** | Monitor records + check history |
| **Cloud Run** | Hosts the HTTP API and worker endpoints |
| **Cloud Build** | CI/CD — GitHub push triggers container build + deploy |
| **Cloud Scheduler** | Fires every 1 minute → hits the dispatch endpoint |
| **Cloud Tasks** | Dispatcher enqueues one check task per monitor; workers run in parallel |
| **Webhooks + HMAC** | On status transition, POST to the user's webhook with HMAC signature |
| **Bearer API keys** (SHA-256 hashed) | Auth for the management API |
| **Structured logging** (`slog`) | Every check emits a log line Cloud Run indexes |

## Build in 3 Layers

Each layer is independently shippable. Finish and commit each before moving on — avoid the "half-finished monorepo" trap.

### Layer 1 — Base CRUD service

- `POST /monitors` — register a URL
- `GET /monitors` — list registered monitors
- `GET /monitors/:id` — detail view
- Bearer API-key auth (SHA-256 hashed keys in DB, middleware)
- PostgreSQL with migrations (goose or golang-migrate)
- Dockerfile
- Cloud Run deploy
- Cloud Build trigger on push to main

**Outcome:** a working (but dumb) CRUD service deployed on Cloud Run. Stop and commit.

### Layer 2 — Scheduling core

- `POST /dispatch` — internal endpoint; scans monitors due for a check, enqueues one Cloud Task per monitor
- `POST /check` — Cloud Tasks worker endpoint, OIDC-authed; performs the HTTP check, writes result row, updates current status
- Cloud Scheduler triggers `/dispatch` every 1 minute

**Outcome:** the service actually does its job. Stop and commit.

### Layer 3 — Notifications

- `PATCH /monitors/:id/webhook` — register a callback URL + shared secret
- On status transition (up → down or down → up), enqueue a Cloud Task to `/notify`
- `/notify` POSTs to the user's webhook with an `X-Signature: sha256=...` header
- Cloud Tasks handles retry semantics automatically on 5xx responses

### Optional Layer 4 — Polish

- **GCS**: store response body snapshots for failed checks; return signed URLs for inspection
- Minimal HTML `GET /monitors/:id` page showing the last 24h of check history

## What This Project Is Not

- Not a task/todo API (redundant with iRolls work)
- Not a URL shortener or pastebin (too generic, no async/scheduled work to showcase)
- Not a chat/forum (scope trap)

## Resume Bullet (Goal State)

Once Layer 3 ships, the resume entry should read:

```
URL Uptime Monitor | Source Code        Go | PostgreSQL | Cloud Run | Cloud Tasks | Cloud Scheduler
• Built a distributed health-check service in Go on Cloud Run, scheduling
  per-URL checks via Cloud Scheduler + Cloud Tasks to isolate work units
  and enable parallel retries.
• Implemented HMAC-signed webhook notifications on status transitions,
  delivered asynchronously through Cloud Tasks with automatic retry on failure.
```

## Alternatives (Same Architecture, Different Domain)

The architecture transfers if the domain stops being motivating:

- **RSS feed aggregator** — Scheduler fetches feeds on interval instead of pinging URLs; everything else identical.
- **Webhook relay** — Drop Cloud Scheduler entirely; service receives webhooks, verifies HMAC, forwards to registered destinations with retry via Cloud Tasks.

Pick whichever domain keeps momentum.
