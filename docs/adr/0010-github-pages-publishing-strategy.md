# 0010: GitHub Pages Publishing Strategy

## Status

Accepted

## Context

The live site must work from day one, with no GitHub Actions. The repo also needs durable documentation under `docs/`.

## Decision

Publish GitHub Pages from the `main` branch `/docs` folder at `https://baditaflorin.github.io/accident-reconstructor/`. Configure Vite with `base: "/accident-reconstructor/"`, hashed assets, `outDir: "docs"`, and `emptyOutDir: false` so ADRs and documentation remain committed. Generate `docs/404.html` as an SPA fallback.

## Consequences

The Pages output is committed. Stale assets may need `make clean` before a rebuild, but documentation is not destroyed by Vite.

## Alternatives Considered

A `gh-pages` branch was rejected because it adds branch choreography. Publishing from repository root was rejected because source files would be publicly served as the site.
