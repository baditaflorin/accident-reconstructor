# 0007: Data Generation Pipeline

## Status

Accepted

## Context

Mode B would require a static data generation pipeline. This project uses Mode C because every reconstruction is private, user-supplied, and processed at runtime.

## Decision

Do not include a Mode B data generation pipeline in v1. Static sample fixtures may live under `public/samples/` or `test/fixtures/`, but they are not authoritative generated artifacts.

## Consequences

`make data` is intentionally a no-op with an explanatory message.

## Alternatives Considered

Publishing reconstructed public sample cases was rejected because it does not solve private crash reconstruction for users.
