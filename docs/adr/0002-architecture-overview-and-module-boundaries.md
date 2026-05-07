# 0002: Architecture Overview and Module Boundaries

## Status

Accepted

## Context

The project needs a clear split between static UI, API contracts, reconstruction orchestration, and deploy assets.

## Decision

Use these boundaries:

- `src/`: Vite React frontend organized by feature.
- `api/`: OpenAPI contract.
- `cmd/server/`: runtime Go API entrypoint.
- `internal/httpapi/`: routing, middleware, and handlers.
- `internal/jobs/`: in-memory job state and artifact lifecycle.
- `internal/reconstruction/`: tool discovery, video probing, COLMAP/GDAL/LLM adapters, and deterministic fallback.
- `deploy/`: production Docker Compose and nginx.
- `docs/`: GitHub Pages output plus project documentation.

## Consequences

The frontend and backend can evolve independently while sharing a stable OpenAPI contract.

## Alternatives Considered

A monolithic Node app was rejected because the backend must orchestrate native forensic tooling and ship as a small production Docker image.
