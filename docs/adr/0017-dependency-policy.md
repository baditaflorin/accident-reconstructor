# 0017: Dependency Policy

## Status

Accepted

## Context

Accident reconstruction is safety-sensitive, so custom low-level implementations should be minimized.

## Decision

Prefer production-ready libraries: `chi`, `validator`, `prometheus/client_golang`, `slog`, `openapi-fetch`, `zod`, `tanstack-query`, `three`, `idb`, Vite, and Playwright. External forensic tools are invoked through adapters with explicit availability reporting.

## Consequences

The project stands on maintained components and can disclose when a native tool is missing or a fallback estimator was used.

## Alternatives Considered

Hand-rolled routing, metrics, OpenAPI clients, and 3D rendering were rejected.
