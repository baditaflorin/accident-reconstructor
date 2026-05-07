# 0004: API Contract

## Status

Accepted

## Context

Mode C requires a runtime API that the static frontend can call. The API must be stable enough for generated clients and smoke tests.

## Decision

Maintain an OpenAPI 3.1 contract in `api/openapi.yaml`. Generate TypeScript types with `openapi-typescript` and call the API through `openapi-fetch`.

## Consequences

The frontend does not hand-write response shapes. Breaking API changes require contract updates and regenerated client types.

## Alternatives Considered

GraphQL was rejected because the interaction pattern is simple job creation and artifact retrieval. Hand-written REST clients were rejected to avoid drift.
