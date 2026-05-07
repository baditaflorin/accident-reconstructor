# 0006: WASM Modules

## Status

Accepted

## Context

The original concept mentioned browser-portable photogrammetry and computer vision. GitHub Pages cannot set arbitrary COOP/COEP headers, and production COLMAP/GDAL workloads remain more reliable as native tools.

## Decision

Do not depend on heavyweight WASM modules in v1. Use Three.js in the browser for visualization and run COLMAP/GDAL/OpenCV-compatible processing through the Docker backend. Future WASM modules must be lazy-loaded and documented with cross-origin isolation requirements.

## Consequences

The v1 frontend remains lightweight and deployable on Pages. Heavy processing requires the backend.

## Alternatives Considered

COLMAP-WASM and OpenCV.js were rejected for the v1 critical path due to load size, threading, browser support, and operational risk.
