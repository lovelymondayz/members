.PHONY: dev build up down logs clean deploy

# Start development environment
dev:
	cd backend && go run ./cmd/server &
	cd frontend && npm run dev
	@echo "Backend: http://localhost:8082 | Frontend: http://localhost:3003"

# Production build
build:
	cd frontend && npm ci && npm run build
	cd backend && go build -o members-api ./cmd/server
	@echo "Build complete"

# Docker operations
up:
	docker compose up -d --build
	@echo "Members running — FE: http://localhost:3003, BE: http://localhost:8082"

down:
	docker compose down

# Utility
logs:
	docker compose logs -f

clean:
	docker compose down -v
	rm -rf frontend/dist backend/members-api

deploy:
	./update.sh