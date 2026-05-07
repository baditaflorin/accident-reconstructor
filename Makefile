SHELL := /usr/bin/env bash
GO_PACKAGES := ./cmd/... ./internal/... ./pkg/...
VERSION := $(shell node -p "require('./package.json').version")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
IMAGE := ghcr.io/baditaflorin/accident-reconstructor

.PHONY: help install-hooks dev build data test test-integration smoke lint fmt pages-preview docker-build docker-push release compose-up compose-down clean hooks-pre-commit hooks-pre-push hooks-post-checkout

help:
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-22s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-hooks: ## Wire local git hooks
	git config core.hooksPath .githooks
	chmod +x .githooks/*

dev: ## Run frontend dev server
	npm run dev

build: ## Build frontend into docs/ and backend binary
	npm run gen:api
	npm run build
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)" -o tmp/accident-server ./cmd/server
	test -f docs/index.html

data: ## Mode C has no static data pipeline
	@echo "Mode C runtime reconstruction uses no static data-generation pipeline."

test: ## Run unit tests
	go test $(GO_PACKAGES)
	npm test

test-integration: ## Run integration tests
	go test -tags=integration $(GO_PACKAGES)

smoke: ## Run local backend and Pages smoke tests
	bash scripts/smoke.sh

lint: ## Run linters
	npm run lint
	go vet $(GO_PACKAGES)
	@if command -v golangci-lint >/dev/null; then golangci-lint run $(GO_PACKAGES); else echo "golangci-lint not installed; skipping"; fi
	@if command -v govulncheck >/dev/null; then govulncheck $(GO_PACKAGES); else echo "govulncheck not installed; skipping"; fi
	npm audit --audit-level=high

fmt: ## Autoformat source
	npx prettier --write .
	gofmt -w cmd internal pkg
	@if command -v goimports >/dev/null; then goimports -w cmd internal pkg; fi

pages-preview: ## Serve the built Pages site locally
	npm run build
	rm -rf tmp/pages-preview
	mkdir -p tmp/pages-preview/accident-reconstructor
	cp -R docs/. tmp/pages-preview/accident-reconstructor/
	npx http-server tmp/pages-preview -a 127.0.0.1 -p 4173 -c-1

docker-build: ## Build amd64 backend image
	docker buildx build --platform linux/amd64 --load \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(IMAGE):latest \
		-t $(IMAGE):v$(VERSION) \
		-t $(IMAGE):$(COMMIT) .

docker-push: ## Push backend image tags to GHCR
	docker buildx build --platform linux/amd64 --push \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(IMAGE):latest \
		-t $(IMAGE):v$(VERSION) \
		-t $(IMAGE):$(COMMIT) .

release: test build smoke docker-push ## Tag and publish a release image
	git tag v$(VERSION)
	git push origin v$(VERSION)

compose-up: ## Run local compose stack
	docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.dev.yml up --build

compose-down: ## Stop local compose stack
	docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.dev.yml down

clean: ## Remove generated local artifacts
	rm -rf tmp coverage docs/assets docs/404.html

hooks-pre-commit:
	npx prettier --check .
	npm run lint
	npx tsc --noEmit
	go test $(GO_PACKAGES)
	go vet $(GO_PACKAGES)
	@if command -v golangci-lint >/dev/null; then golangci-lint run $(GO_PACKAGES); else echo "golangci-lint not installed; skipping"; fi
	@if command -v gitleaks >/dev/null; then gitleaks protect --staged --redact; else echo "gitleaks not installed; skipping"; fi

hooks-pre-push: test build smoke

hooks-post-checkout:
	npm run gen:api
