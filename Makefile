.PHONY: help setup dev dev-backend dev-frontend build build-backend build-frontend clean backup docker-build docker-run test

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Install all dependencies
	@echo "TODO: install dependencies"

dev: ## Run backend and frontend in development mode
	@echo "TODO: run dev servers"

dev-backend: ## Run Go backend with hot reload
	@echo "TODO: run backend dev"

dev-frontend: ## Run Vite dev server
	@echo "TODO: run frontend dev"

build: build-frontend build-backend ## Build everything

build-frontend: ## Build React frontend
	@echo "TODO: build frontend"

build-backend: ## Build Go binary with embedded static files
	@echo "TODO: build backend"

test: ## Run all tests
	@echo "TODO: run tests"

clean: ## Remove build artifacts
	@echo "TODO: clean"

backup: ## Trigger a manual database backup to Google Drive
	@echo "TODO: trigger backup"

docker-build: ## Build Docker image
	@echo "TODO: docker build"

docker-run: ## Run with docker-compose
	@echo "TODO: docker run"
