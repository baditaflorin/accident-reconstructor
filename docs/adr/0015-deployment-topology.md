# 0015: Deployment Topology

## Status

Accepted

## Context

Mode C requires a public static frontend and an independently deployed backend.

## Decision

Serve the frontend from GitHub Pages. Run the backend with Docker Compose behind nginx on host port `25342`, pulling images from GHCR. nginx terminates TLS, applies security headers, rate-limits `/api/`, allows the GitHub Pages origin via CORS, and blocks public `/metrics`.

## Consequences

The frontend remains cheap and static. The backend can be upgraded or rolled back independently.

## Alternatives Considered

Serving the frontend from the Go API was rejected because Pages is a first-class deliverable.
