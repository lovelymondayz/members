.PHONY: help build run test clean

APP_NAME=members
GO_CMD=go
DOCKER_CMD=docker compose

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## ─── Backend ────────────────────────────────────────────

build: ## Build Go binary
	CGO_ENABLED=0 $(GO_CMD) build -o bin/$(APP_NAME) ./backend/cmd/server

run: ## Run backend (requires .env)
	./bin/$(APP_NAME)

dev: ## Run backend with hot reload (requires air)
	air

test: ## Run Go tests
	$(GO_CMD) test ./... -v -count=1

lint: ## Run golangci-lint
	golangci-lint run ./...

clean: ## Remove build artifacts
	rm -rf bin/

## ─── Docker ────────────────────────────────────────────

docker-up: ## Start all containers
	$(DOCKER_CMD) up -d

docker-down: ## Stop all containers
	$(DOCKER_CMD) down

docker-logs: ## Tail container logs
	$(DOCKER_CMD) logs -f

docker-rebuild: ## Rebuild and restart
	$(DOCKER_CMD) down && $(DOCKER_CMD) build --no-cache && $(DOCKER_CMD) up -d

## ─── Frontend ──────────────────────────────────────────

fe-install: ## Install frontend dependencies
	cd frontend && npm install

fe-dev: ## Start frontend dev server
	cd frontend && npm run dev

fe-build: ## Build frontend for production
	cd frontend && npm run build

fe-preview: ## Preview production build
	cd frontend && npm run preview

## ─── Deploy ────────────────────────────────────────────

deploy: build fe-build docker-rebuild ## Full deploy
