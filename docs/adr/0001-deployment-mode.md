# 0001: Deployment Mode

## Status

Accepted

## Context

The product must accept large videos, run photogrammetry and computer vision tooling, estimate speeds, and generate court/insurance-oriented reconstruction artifacts. GitHub Pages is preferred for public delivery, but Pages cannot run native COLMAP/GDAL workloads, persist large case artifacts server-side, or provide reliable cross-origin isolation headers for production-grade threaded WASM.

## Decision

Use Mode C: GitHub Pages frontend plus Docker backend. The browser app is static and hosted from GitHub Pages. A Dockerized Go API performs heavy reconstruction work by orchestrating COLMAP, OpenCV-compatible adapters, GDAL, FFmpeg, and an optional local LLM.

## Consequences

The public surface remains a static site. Operators need to run a backend for real processing. The frontend can still demo and inspect existing case artifacts without secrets.

## Alternatives Considered

- Mode A: Rejected for v1 because COLMAP and video-scale processing are too brittle and heavy in the browser.
- Mode B: Rejected because users need private, per-case uploads and runtime processing rather than a public precomputed dataset.
