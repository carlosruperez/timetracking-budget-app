# TimeTracking Budget App — Task Runner
# Requires: just (https://github.com/casey/just)

# Show available commands
default:
    @just --list

# ─── Docker Compose ──────────────────────────────────────────────────────────

# Start all services (postgres + backend + web) via Docker
up:
    docker compose up --build

# Start all services in background
up-d:
    docker compose up -d --build

# Stop all services
down:
    docker compose down

# Stop and remove volumes (resets DB)
down-v:
    docker compose down -v

# Only start the database
db:
    docker compose up postgres

# Only start the database in background
db-d:
    docker compose up -d postgres

# Rebuild and start all services
rebuild:
    docker compose up --build

# Show logs
logs service="":
    docker compose logs -f {{service}}

# ─── Backend (Go) ────────────────────────────────────────────────────────────

# Run backend in dev mode (requires postgres running)
backend:
    cd backend && go run ./cmd/server/main.go

# Build backend binary
backend-build:
    cd backend && go build -o bin/server ./cmd/server/main.go

# Run backend tests
backend-test:
    cd backend && go test ./...

# Run backend vet
backend-vet:
    cd backend && go vet ./...

# ─── Web (Next.js) ───────────────────────────────────────────────────────────

# Run web dev server
web:
    cd apps/web && pnpm dev

# Build web for production
web-build:
    cd apps/web && pnpm build

# Type-check web
web-check:
    cd apps/web && pnpm exec tsc --noEmit

# ─── Mobile (Expo) ───────────────────────────────────────────────────────────

# Start Expo dev server
mobile:
    cd apps/mobile && npx expo start

# Start Expo for Android
mobile-android:
    cd apps/mobile && npx expo start --android

# Start Expo for iOS
mobile-ios:
    cd apps/mobile && npx expo start --ios

# ─── Shared Packages ─────────────────────────────────────────────────────────

# Build @timetracking/core
core-build:
    cd packages/core && pnpm build

# ─── Dev (local, without Docker) ─────────────────────────────────────────────

# Start DB + backend + web locally (3 processes)
dev: db-d
    #!/usr/bin/env bash
    trap 'kill 0' SIGINT
    (cd backend && go run ./cmd/server/main.go) &
    (cd apps/web && pnpm dev) &
    wait

# ─── Utilities ───────────────────────────────────────────────────────────────

# Install all JS dependencies
install:
    pnpm install

# Check everything (backend vet + tests, web type-check)
check: backend-vet backend-test web-check
    @echo "All checks passed."
