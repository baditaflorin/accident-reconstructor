# Postmortem

## What Was Built

Accident Reconstructor now has a public GitHub Pages frontend, a Go reconstruction API, OpenAPI contract, Docker deployment assets, local hooks, tests, smoke coverage, and documentation.

Live site:

https://baditaflorin.github.io/accident-reconstructor/

Repository:

https://github.com/baditaflorin/accident-reconstructor

## Was Mode C Correct?

Yes for v1. GitHub Pages is right for the public UI, but browser-only COLMAP/GDAL/OpenCV processing remains too fragile for large private videos. A runtime backend is justified.

## What Worked

- Pages was enabled from the first commit.
- The UI can run as a static app and still show a sample reconstruction when no backend is configured.
- The backend produces structured artifacts, reports, metrics, and case bundles.

## What Did Not Work

- A courtroom-quality native COLMAP pipeline cannot honestly be completed as a pure browser task.
- The distroless image is excellent for the Go API, but native photogrammetry tools need a heavier operator image or worker layer.

## Surprises

- The latest Vite scaffold pulled TypeScript 6 preview, while OpenAPI tooling still expects TypeScript 5.x. The repo pins TypeScript 5.9.

## Accepted Tech Debt

- The native COLMAP/GDAL/OpenCV execution path is represented by adapters and tool discovery, with deterministic fallback geometry in the base image.
- Case state is in memory, while artifacts persist on disk. A restart keeps files but not the live job list.
- Browser UI uses a direct backend URL setting rather than a production API discovery mechanism.

## Next Improvements

1. Add a toolchain image or worker process with COLMAP, GDAL, Python OpenCV, and controlled frame extraction.
2. Persist case metadata in SQLite so restarts recover historical cases.
3. Add calibration UI for lane width, skid marks, and known road reference points.

## Time

Estimate: 2-3 focused days for a production-grade native reconstruction path.

Actual bootstrap: one implementation pass for a deployable v1 scaffold and deterministic reconstruction workflow.
