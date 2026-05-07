# 0011: Logging Strategy

## Status

Accepted

## Context

Mode C needs server logs suitable for production operations. The frontend should not leak sensitive case details through console output.

## Decision

Use Go `slog` JSON logs to stdout with `trace_id`, method, path, status, and duration. Frontend production builds keep console output minimal and route user-visible failures through UI toasts.

## Consequences

Backend logs work with Docker and log collectors. Browser logs remain quiet for public users.

## Alternatives Considered

Text logs were rejected for production because JSON is easier to query and correlate.
