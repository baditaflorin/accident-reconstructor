# 0013: Testing Strategy

## Status

Accepted

## Context

The repo needs local checks rather than GitHub Actions. Tests must remain fast enough for pre-push hooks.

## Decision

Use Go unit tests for `internal/`, Vitest for frontend logic, and Playwright smoke tests against a built Pages preview. `make test` runs Go and frontend unit tests. `make smoke` builds, serves `docs/`, and verifies one happy path.

## Consequences

Local contributors can run the same checks as hooks. Browser smoke coverage catches broken Pages builds.

## Alternatives Considered

GitHub Actions were rejected by project constraint. Full visual regression was deferred to keep smoke checks under one minute.
