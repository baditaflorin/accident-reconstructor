# Deployment

Live frontend:

https://baditaflorin.github.io/accident-reconstructor/

Backend image:

ghcr.io/baditaflorin/accident-reconstructor:latest

## Server Prerequisites

- Docker Engine with Compose.
- DNS A record for the API host.
- TLS certificates mounted into the `letsencrypt` volume or bind-mounted at `/etc/letsencrypt`.
- Optional local Ollama endpoint if LLM narrative summaries are desired.

## First Deploy

```sh
cd deploy
cp .env.example .env
docker compose pull
docker compose up -d
```

The public HTTPS port is `25342`, mapped to nginx `443`. The app service listens only on the internal Docker network at `:8080`.

## Rollback

```sh
docker compose pull app
docker compose up -d app
```

For a specific version:

```sh
docker pull ghcr.io/baditaflorin/accident-reconstructor:v0.1.0
```

Then edit `deploy/docker-compose.yml` to pin that tag and run:

```sh
docker compose up -d app
```

## Logs

```sh
docker compose logs -f app
docker compose logs -f nginx
```

## Backups

Case artifacts live in the `cases` named volume:

```sh
docker run --rm -v accident-reconstructor_cases:/data -v "$PWD":/backup alpine tar czf /backup/cases-backup.tgz /data
```

## Notes

The distroless production image runs the Go API and deterministic fallback pipeline. For native COLMAP/GDAL/OpenCV execution, build an operator-specific image or host wrapper with those tools available and keep the same API contract.
