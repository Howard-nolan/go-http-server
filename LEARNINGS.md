
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
- Created stub handlers for /v1/r/{code} and /v1/shorten endpoints to prepare for future URL redirection and shortening logic.
- This sets the foundation for versioned API routing and future OpenAPI documentation integration.

**How It Works:** 
- The server now initializes a Chi router in main.go, registering routes under /v1/.
- GET /v1/r/{code} will later handle redirecting a short code to its target URL.
- POST /v1/shorten will accept a long URL in the request body and return a shortened code.
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

