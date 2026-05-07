# 0003: Frontend Framework and Build Tooling

## Status

Accepted

## Context

The UI needs strict TypeScript, a fast local dev loop, GitHub Pages output, 3D visualization, and strong form/state ergonomics.

## Decision

Use React, TypeScript strict mode, Vite, Tailwind CSS, Zod, TanStack Query, openapi-fetch, Three.js, and lucide-react.

## Consequences

The app has a familiar ecosystem and fast builds. Three.js is lazy-loaded behind the reconstruction viewer to protect the first-load budget.

## Alternatives Considered

SvelteKit and Astro were considered. React/Vite was chosen for the broad library ecosystem around Three.js, uploads, and OpenAPI clients.
