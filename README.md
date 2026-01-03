# Go URL Shortener

A high performance URL shortener written in Go to explore distributed systems fundamentals: observability, scalability, deployment, and testing.

I document what I learn from each pull request in [**LEARNINGS.md**](./LEARNINGS.md).

---

## At a glance

- REST API built with `chi`, request timeouts, and graceful shutdown.
- SQLite persistence with embedded Goose migrations; data stored in `data/app.db`.
- In-memory LRU cache for hot redirect lookups.
- Idempotent create flow when clients reuse `X-Request-ID`.
- Observability: Zap JSON logs, Prometheus metrics at `/metrics`, and a local Grafana stack.
- Testing: table-driven unit tests with `sqlmock` plus integration tests via `httptest`.

## Quick Start

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

[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.0-green.svg)](https://swagger.io/specification/)

| Method | Endpoint    | Description                               | Status / Utility   |
| ------ | ----------- | ----------------------------------------- | ------------------ |
| POST   | /shorten    | Create a shortened URL from a long link   | Core API           |
| GET    | /r/{code}   | Redirect to the original destination      | Core API           |
| GET    | /metrics    | Prometheus metrics for monitoring         | DevOps             |
| GET    | /health     | Liveness probe (Server is alive)          | Kubernetes/Health  |
| GET    | /readyz     | Readiness probe (DB is connected)         | Kubernetes/Health  |

Example:xw
```bash
curl -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

Response:
```json
{"short":"http://localhost:8080/r/abc123"}
```

Full OpenAPI spec: `api/openapi.yaml`

## Observability (local)

Start the containerized service (URL Shortener + Prometheus + Grafana):
```bash
docker compose up
```

- App: http://localhost:8080 (metrics at http://localhost:8080/metrics)
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
