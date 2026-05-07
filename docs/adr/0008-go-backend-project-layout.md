# 0008: Go Backend Project Layout

## Status

Accepted

## Context

The backend must stay modular while following familiar Go conventions.

## Decision

Use the `golang-standards/project-layout` shape: `cmd/`, `internal/`, `pkg/`, `api/`, `configs/`, `scripts/`, and `test/`.

## Consequences

Runtime-only internals stay private, shared schemas can live in `pkg/`, and deploy scripts remain separate from application code.

## Alternatives Considered

A flat Go package was rejected because reconstruction, HTTP, config, and job lifecycle concerns would blur quickly.
