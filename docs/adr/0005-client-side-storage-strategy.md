# 0005: Client-Side Storage Strategy

## Status

Accepted

## Context

Users may handle sensitive footage and need local drafts without accounts. The frontend also needs cached job summaries and API settings.

## Decision

Use IndexedDB through `idb` for case drafts and last successful API responses. Use `localStorage` only for small UI preferences and API base URL overrides. Avoid cross-device sync in v1.

## Consequences

User data stays in the browser unless explicitly uploaded to the backend. Clearing browser storage removes local drafts.

## Alternatives Considered

Server accounts and cloud sync were rejected because v1 should avoid auth and sensitive centralized storage.
