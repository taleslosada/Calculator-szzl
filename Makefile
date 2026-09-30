.PHONY: help backend frontend test test-backend test-frontend coverage docker

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

backend: ## Run the Go API on :8080
	cd backend && go run ./cmd/server

frontend: ## Run the React dev server on :5173 (proxies /api to :8080)
	cd frontend && npm install && npm run dev

test: test-backend test-frontend ## Run all unit tests

test-backend:
	cd backend && go test ./... -cover

test-frontend:
	cd frontend && npm test

coverage: ## Generate HTML coverage reports for both layers
	cd backend && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html && go tool cover -func=coverage.out
	cd frontend && npm run coverage

docker: ## Build and run the full stack in one container on :8080
	docker compose up --build
