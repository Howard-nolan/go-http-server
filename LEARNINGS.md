
## 2025-10-26 - PR #9: Finalized GitHub Merge Workflow

**Change Summary:** Refined the GitHub workflow and PR template to simplify learning capture and align automation with new section names.

**How It Works:** Updated the PR template to use concise fields — *Change Summary*, *How It Works*, and *Additional Notes*.  
Modified the workflow to extract those sections and append them to LEARNINGS.md.

**Additional Notes:** Verified the new workflow syntax locally and ensured backward compatibility with existing entries.


## 2025-11-05 - PR #10: Added Repo Scaffolding and Config File

**Change Summary:** 
- Added initial project scaffold and /health endpoint.
- Introduced config loading (PORT) and basic logger.
- Implemented graceful shutdown with context and signal handling.

**How It Works:** 
- main.go wires config, logger, and routes through a local ServeMux.
- Server listens on configurable port and shuts down cleanly on SIGINT/SIGTERM.

**Additional Notes:** 
- Foundation for future router, structured logs, and error handling in Phase 2.


## 2025-11-10 - PR #12: Added Chi Routing, API Stubs

**Change Summary:** 
- Added Chi router setup to handle API routes.
- Created stub handlers for /r/{code} and /shorten endpoints to prepare for future URL redirection and shortening logic.
- This sets the foundation for versioned API routing and future OpenAPI documentation integration.

**How It Works:** 
- The server now initializes a Chi router in main.go, registering routes under /.
- GET /r/{code} will later handle redirecting a short code to its target URL.
- POST /shorten will accept a long URL in the request body and return a shortened code.
- Both endpoints currently return placeholder JSON responses to confirm routing is functional.

**Additional Notes:** 
- Future work: implement persistence layer for shortened URLs (e.g., MongoDB, Redis).
- Can easily extend versioning by adding additional route groups (e.g., /v2/).
- No functional logic yet—this MR only sets up routing and scaffolding.


## 2025-11-13 - PR #13: Added response.go, consistent error response

**Change Summary:** 
- Added a response.go helper to centralize JSON response formatting.
- Introduced a consistent JSON error shape used across all HTTP handlers.
- Updated handlers to use shared helpers instead of writing responses inline.

**How It Works:** 
- response.go has WriteJSON and WriteError functions

**Additional Notes:** 
- Existing error responses have changed shape; any clients depending on the old error body format may need to be updated.
- Follow-up: we can extend ErrorResponse with request IDs / correlation IDs or more structured error codes if needed.


## 2025-11-15 - PR #14: Minor GitHub Workflow Tweak

**Change Summary:**
- LEARNINGS.md formatting improvement; minor change to append-learnings.yml

**How It Works:**
- Add new line after headers in LEARNINGS.md


## 2025-11-17 - PR #15: added openapi spec

**Change Summary:**
- Added OpenAPI spec

**How It Works:**
- openapi.yaml contains spec for three apis.
- can be viewed locally with OpenAPI extension


## 2025-11-30 - PR #16: Wire in SQLite DB; Have routes write / read from the db

**Change Summary:**
- Wired DB into handlers; shortened URLs now persist to SQLite and redirects fetch from DB.
- Added random code generation with collision retries, URL validation/normalization, and real redirects.
- Fixed goose embed/dialect setup and used request-aware DB calls.

**How It Works:**
- cmd/server/main.go opens ./data/app.db, runs migrations, builds a handlers.Handler with the shared *sql.DB, and registers routes.
- ShortenHandler validates/normalizes the incoming URL (adds https:// if missing), generates a 6-char code, inserts (code, url) into links via ExecContext with retries on collisions, and returns the short link.
- RedirectHandler looks up the URL by code using QueryRowContext; 404s if missing, otherwise issues an HTTP 302 to the stored URL.

**Additional Notes:**
- Bare hostnames like www.google.com are accepted; malformed URLs still 400.


## 2025-12-10 - PR #17: Feature/observability

**Change Summary:**
- Observability: Integrated Zap for structured, level-based JSON logging to improve debugging capabilities.
- Monitoring: Added a Prometheus middleware to track key Golden
Signals (request latency and traffic volume) and exposed them via a new /metrics endpoint.
- Local Dev Stack: Created a docker-compose.yml file to spin up a local monitoring stack (Prometheus and Grafana) alongside the application.
- Performance: Implemented a caching layer to reduce database load for frequently accessed URL redirects.

**How It Works:**
- Metrics Middleware: A new MetricsMiddleware wraps the HTTP router. It captures the start time of every request and records:
- http_requests_total (Counter): Labeled by method, route, and status.
- http_request_duration_seconds (Histogram): Captures latency buckets for P99/P95 calculations.

- Route Normalization: The middleware uses parameterized route names (e.g., /r/{code}) rather than raw paths to prevent high-cardinality issues in Prometheus.

- Docker Stack:
- Prometheus: Configured to scrape host.docker.internal:8080/metrics every 5 seconds.
- Grafana: Pre-provisioned to visualize the Prometheus data on port 3000.
- Logging: Replaced standard fmt print statements with a global Zap logger configuration.

**Additional Notes:**
- Access: Grafana is available at http://localhost:3000 (default creds: admin/admin) when running docker compose up.

- Trade-offs: Deferred OpenTelemetry (tracing) implementation as the current architecture is single-service and basic structured logging/metrics provide sufficient visibility for now.

- Verification: Verified that the /metrics endpoint is correctly outputting Prometheus-formatted text with appropriate buckets.


## 2025-12-18 - PR #18: Feature/tests and quality gates

**Change Summary:**
- Expanded coverage with table-driven handler unit tests (sqlmock) and an integration-tagged black-box test that exercises real HTTP on an in-memory DB.
- Wired server to new middleware/logger/cache-aware handler signature and recorded supporting dependencies (prometheus, hashicorp/golang-lru, sqlmock).

**How It Works:**
- Unit tests in internal/http/handlers/handlers_test.go drive handlers through table cases with sqlmock; integration_test.go (build tag integration) spins up httptest.NewServer with in-memory SQLite to hit /shorten then /r/{code} end-to-end.

**Additional Notes:**
- Run integration suite with go test -tags=integration ./...; it’s skipped in default runs.


## 2025-12-19 - PR #19: Better Timeout Logic; Added Idempotency Key to DB

**Change Summary:**
- Added request timeouts, readiness endpoint, and cache-aware handlers with idempotent shorten inserts; wrapped the server in http.TimeoutHandler and handled DB timeouts explicitly.
- Introduced idempotency key support (schema + insert logic) and documented readiness in OpenAPI; added ready route, updated README with test commands.
- Expanded tests: handler unit tests now cover timeouts/idempotency and redirect timeouts; integration test runs against in-memory DB.

**How It Works:**
- ShortenHandler reads the request ID as an idempotency key, uses ExecContext with an upsert-by-key, returns 408 on context cancel, and reuses the existing code on conflict.
- RedirectHandler checks cache, then QueryRowContext; returns 408 on context cancellation, 404 on miss.
- /readyz pings the DB with PingContext; server wraps the router with http.TimeoutHandler. Schema includes a unique idempotency_key on links.

**Additional Notes:**
- Integration tests run with go test -tags=integration ./...; unit tests via go test ./....


## 2025-12-26 - PR #20: Updated Readme

**Change Summary:**
- Updated README to reflect latest status of the project

**How It Works:**
- LEARNINGS.md still linked. Additional commands added.

## 2025-12-28 - PR #21: Feature/deploy

**Change Summary:**
- Added DB_PATH and BASE_URL config plumbing; handlers now build short URLs from BASE_URL and routes drop /v1 in favor of /shorten and /r/{code}.
- Swapped SQLite driver to modernc.org/sqlite for CGO‑free static builds; added multi‑stage distroless Dockerfile, .dockerignore, and an app service + volume in docker-compose.yml.
- Updated tests, OpenAPI, README, Prometheus config, and LEARNINGS to match the new paths and metrics endpoint.

**How It Works:**
- config.Load() reads PORT, DB_PATH, and BASE_URL; main opens SQLite with db.OpenAndMigrate(cfg.DbPath) and passes cfg.BaseURL into the handler.
- Handler trims the base URL; ShortenHandler inserts a link and returns /r/{code}, while /r/{code} redirects to the original URL; /metrics exposes Prometheus text format.
- Docker build compiles a static binary with CGO_ENABLED=0 and runs it on distroless; compose wires env vars and mounts a named volume for persistent SQLite.

**Additional Notes:**
- Default DB_PATH is app.db; compose overrides it to app.db for persistence—Cloud Run will need an explicit DB_PATH or a durable datastore.
- OpenAPI/README examples now use /shorten, /r/{code}, and /metrics.

