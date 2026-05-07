# 0009: Configuration and Secrets Management

## Status

Accepted

## Context

The frontend is public and must never contain secrets. The backend may need runtime configuration for origins, storage, and optional LLM endpoints.

## Decision

Configure the frontend through build-time `VITE_*` values only for public settings. Configure the backend exclusively through environment variables documented in `.env.example`. Keep real `.env` files gitignored. Run gitleaks in local hooks.

## Consequences

Deployments can be reproduced without secret sprawl. Public Pages builds remain safe to inspect.

## Alternatives Considered

Encrypted frontend secrets were rejected because public client-side secrets are still secrets exposed to users.
