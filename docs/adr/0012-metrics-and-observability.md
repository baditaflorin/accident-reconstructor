# 0012: Metrics and Observability

## Status

Accepted

## Context

The backend should be scrape-ready. The frontend should avoid analytics by default.

## Decision

Expose Prometheus metrics at `/metrics`, including Go runtime metrics, HTTP request counts/durations, reconstruction job counts, upload byte totals, and pipeline duration histograms. Block `/metrics` publicly at nginx. Do not add client analytics in v1.

## Consequences

Operators can observe backend health without collecting user behavior from the static site.

## Alternatives Considered

Plausible analytics was deferred because v1 does not need product analytics.
