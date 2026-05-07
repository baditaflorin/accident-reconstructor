# Runbook

Frontend:

https://baditaflorin.github.io/accident-reconstructor/

Backend image:

ghcr.io/baditaflorin/accident-reconstructor:latest

## Local Debugging

```sh
make dev
go run ./cmd/server
```

## Production Logs

```sh
cd deploy
docker compose logs -f app
docker compose logs -f nginx
```

## Health Checks

```sh
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/metrics
```

## Resource Sizing

Minimum: 2 CPU cores, 4 GB RAM, 20 GB disk.

Recommended for native photogrammetry workloads: 8 CPU cores, 16 GB RAM, 100 GB disk.

## Common Failures

- Backend unavailable from Pages: check `VITE_API_BASE_URL` in the UI and CORS `PAGES_ORIGIN`.
- Case stuck processing: inspect `docker compose logs app`.
- Native reconstruction unavailable: call `/api/v1/tools` and verify COLMAP/GDAL/OpenCV adapters on the server image.
- Upload rejected: check `MAX_UPLOAD_BYTES` and nginx `client_max_body_size`.
