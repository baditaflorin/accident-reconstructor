# 0014: Error Handling Conventions

## Status

Accepted

## Context

Video processing fails often due to codecs, corrupted files, missing tools, or low visual overlap. Failures must be explicit and recoverable.

## Decision

Return structured JSON errors from the API with `code`, `message`, and optional `details`. In Go, wrap errors with `%w`, avoid panics, and provide `internal/utils.HandleErrorOrLogWithMessages(err, errMsg, successMsg)` for the requested convention.

## Consequences

The frontend can show clear messages and users can export reports that include processing limitations.

## Alternatives Considered

Plain text errors were rejected because they are hard for generated clients and UI states to consume.
