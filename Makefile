.PHONY: help setup dev dev-backend dev-frontend build build-backend build-frontend clean docker-build docker-run test test-backend android-build

BINARY_NAME := expense-tracker
BIN_DIR     := bin
DOCKER_IMG  := expense-tracker

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ── Setup ──────────────────────────────────────────────────────────────────────

setup: ## Install all dependencies (Go modules + npm)
	cd backend && go mod download
	cd frontend && npm install

# ── Development ────────────────────────────────────────────────────────────────

dev-backend: ## Run Go backend (port 8080)
	cd backend && go run ./cmd/server/

dev-frontend: ## Run Vite dev server (port 5173, proxies /api to :8080)
	cd frontend && npm run dev

dev: ## Run backend + frontend in parallel
	@$(MAKE) -j2 dev-backend dev-frontend

# ── Build ──────────────────────────────────────────────────────────────────────

build-frontend: ## Build React app into backend/cmd/server/static/
	cd frontend && npm run build

build-backend: ## Build Go binary (requires frontend built first)
	mkdir -p $(BIN_DIR)
	cd backend && CGO_ENABLED=1 go build -ldflags="-s -w" -o ../$(BIN_DIR)/$(BINARY_NAME) ./cmd/server/

build: build-frontend build-backend ## Build everything into a single binary

# ── Test ───────────────────────────────────────────────────────────────────────

test-backend: ## Run Go tests
	cd backend && go test ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test -- --run

test: test-backend test-frontend ## Run all tests

# ── Docker ─────────────────────────────────────────────────────────────────────

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMG) .

docker-run: ## Run app with docker-compose
	docker-compose up

docker-up: ## Run app in background with docker-compose
	docker-compose up -d --build

docker-down: ## Stop docker-compose
	docker-compose down

# ── Android ────────────────────────────────────────────────────────────────────

android-build: ## Build Android debug APK (requires Android SDK + ANDROID_HOME set)
	cd mobile/android && ./gradlew assembleDebug
	@echo "APK: mobile/android/app/build/outputs/apk/debug/app-debug.apk"

android-install: ## Install APK on connected device/emulator
	cd mobile/android && ./gradlew installDebug

# ── Misc ───────────────────────────────────────────────────────────────────────

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -rf backend/cmd/server/static/assets
	rm -f  backend/cmd/server/static/index.html \
	       backend/cmd/server/static/favicon.svg \
	       backend/cmd/server/static/icons.svg
	rm -rf frontend/dist
	rm -rf mobile/android/app/build mobile/android/build mobile/android/.gradle

.DEFAULT_GOAL := help
