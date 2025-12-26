# Go URL Shortener 🧩

A lightweight URL shortener written in Go to explore distributed systems fundamentals: routing, concurrency, observability, resilience, and deployment.

👉 I document what I learn from each pull request in [**LEARNINGS.md**](./LEARNINGS.md).

---

## At a glance

- REST API built with `chi`, request timeouts, and graceful shutdown.
- SQLite persistence with embedded Goose migrations; data stored in `data/app.db`.
- In-memory LRU cache for hot redirect lookups.
- Idempotent create flow when clients reuse `X-Request-ID`.
- Observability: Zap JSON logs, Prometheus metrics at `/metrics`, and a local Grafana stack.
- Testing: table-driven unit tests with `sqlmock` plus integration tests via `httptest`.

## 🚀 Quick Start

Run the server:
```bash
go run ./cmd/server
```

Override the port:
```bash
PORT=9090 go run ./cmd/server
```

## API

Endpoints:
- `POST /v1/shorten` create a short URL
- `GET /v1/r/{code}` redirect to original URL
- `GET /health` and `GET /readyz` for liveness/readiness
- `GET /metrics` Prometheus scrape

Example:
```bash
curl -X POST http://localhost:8080/v1/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

Response:
```json
{"short":"https://short.example/abc123"}
```

Full OpenAPI spec: `api/openapi.yaml`

## Observability (local)

Start the monitoring stack (Prometheus + Grafana):
```bash
docker compose up
```

- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

## Tests

Run unit tests:
```bash
go test ./...
```

Run integration (black-box) tests:
```bash
go test -tags=integration ./...
```
