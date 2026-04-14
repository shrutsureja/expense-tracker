# Family Wallet — Expense Tracker

A family expense tracking app with web and Android support. Designed for ultra-low-friction expense entry — even elder family members can use it comfortably.

## Features

- **2-tap quick add**: tap a frequent category → enter amount → done
- **29 built-in categories** with emoji icons (Groceries, Petrol, Medicine, etc.)
- **Smart suggestions**: most-used tags shown first based on your history
- **Family roles**: super admin → family owner → family members (PIN login)
- **Analytics dashboard**: per-person, per-category, daily trend, monthly comparison
- **Mobile-first**: large touch targets, bottom navigation, elder-friendly design
- **Android APK**: WebView wrapper — runs the same web app on Android

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go + chi router |
| Database | SQLite (embedded in binary) |
| Frontend | React + TypeScript + Vite + Tailwind CSS |
| Charts | Recharts |
| Mobile | Android WebView wrapper |
| Deployment | Docker (single 9MB binary + SQLite file) |

## Quick Start

### Development

```bash
# 1. Set up dependencies
make setup

# 2. Create .env file
cp .env.example .env
# Edit .env — set SUPER_ADMIN_PASSWORD and JWT_SECRET

# 3. Run backend + frontend dev servers (in parallel)
make dev
# Backend: http://localhost:8080
# Frontend: http://localhost:5173 (proxies /api to backend)
```

### Production (single binary)

```bash
make build
SUPER_ADMIN_PASSWORD=yourpassword JWT_SECRET=yoursecret ./bin/expense-tracker
# App running at http://localhost:8080
```

### Docker

```bash
# Copy and edit env file
cp .env.example .env

docker-compose up --build
# App running at http://localhost:8080
```

## Auth Setup

1. Log in as super admin: username `shrutsureja`, password from `SUPER_ADMIN_PASSWORD` env var
2. Go to Admin → Create a family (e.g. "Sureja Family") with an owner account
3. Log in as the family owner → Family → Add Members with userid + PIN
4. Share the userid and PIN verbally with family members

## Android APK

Change `SERVER_URL` in `mobile/android/app/build.gradle` to your server's IP, then:

```bash
# Requires Android SDK and ANDROID_HOME set
make android-build
# APK: mobile/android/app/build/outputs/apk/debug/app-debug.apk

make android-install   # Install on connected device/emulator
```

For local network, set `SERVER_URL` to `http://192.168.1.YOUR_IP:8080`.

## All Make Commands

```
make help              Show all commands
make setup             Install Go + npm dependencies
make dev               Run both dev servers in parallel
make dev-backend       Run Go backend only (port 8080)
make dev-frontend      Run Vite dev server only (port 5173)
make build             Build everything (frontend → embed → Go binary)
make build-frontend    Build React app only
make build-backend     Build Go binary only
make docker-build      Build Docker image
make docker-up         Start with docker-compose (background)
make docker-down       Stop docker-compose
make test              Run all tests
make test-backend      Run Go tests
make android-build     Build Android APK
make android-install   Install APK on device/emulator
make clean             Remove all build artifacts
```

## Project Structure

```
expense-tracker/
├── backend/
│   ├── cmd/server/         # Entry point + embedded static files
│   └── internal/
│       ├── auth/           # JWT utilities
│       ├── config/         # Env-based config
│       ├── database/       # SQLite + migrations + seed
│       ├── handler/        # HTTP handlers + middleware
│       ├── models/         # Data types
│       ├── repository/     # DB queries
│       └── service/        # Business logic
├── frontend/
│   └── src/
│       ├── api/            # Typed API client modules
│       ├── components/     # UI, layout, expense, dashboard, family
│       ├── context/        # AuthContext
│       ├── pages/          # Login, Home, Add, Dashboard, Family, Admin
│       ├── types/          # TypeScript interfaces
│       └── utils/          # Formatters, constants
├── mobile/android/         # Android WebView wrapper
├── .env.example
├── Makefile
├── Dockerfile
└── docker-compose.yml
```
