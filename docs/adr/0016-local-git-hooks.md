# 0016: Local Git Hooks

## Status

Accepted

## Context

The project forbids GitHub Actions but still needs repeatable checks.

## Decision

Use plain `.githooks/` scripts wired by `make install-hooks`. Hooks run formatting, linting, TypeScript checks, Go checks, gitleaks, Conventional Commits validation, build, tests, and smoke tests.

## Consequences

Checks run locally before commits and pushes. Contributors must install required tools listed in the README.

## Alternatives Considered

lefthook was considered but plain hooks keep the repo dependency-light.
