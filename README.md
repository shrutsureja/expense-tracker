# Expense Tracker

A family expense tracking application with web and Android support. Designed for ultra-low-friction expense entry that even elder family members can use comfortably.

## Features

- Quick 2-tap expense entry with smart tag suggestions
- Family-based expense management with role-based access
- Analytics dashboard with charts and breakdowns
- Excel export for detailed reports
- Weekly automatic backups to Google Drive
- Mobile-friendly responsive design
- Android WebView app

## Tech Stack

- **Backend**: Go (chi router) + SQLite
- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Charts**: Recharts
- **Mobile**: Android WebView wrapper
- **Deployment**: Docker (single container)

## Quick Start

1. Copy the example config:
   ```bash
   cp config.example.yaml config.yaml
   ```

2. Edit `config.yaml` with your settings (super admin credentials, JWT secret, Google Drive backup config)

3. Run in development mode:
   ```bash
   make setup
   make dev
   ```

4. Build for production:
   ```bash
   make build
   ./bin/expense-tracker
   ```

5. Or use Docker:
   ```bash
   docker-compose up --build
   ```

## Configuration

All configuration is in `config.yaml`. See `config.example.yaml` for all available options.

### Google Drive Backup Setup

1. Create a Google Cloud project and enable the Google Drive API
2. Create a service account and download the credentials JSON
3. Place the JSON file at the path specified in `config.yaml` (default: `./credentials/gdrive-service-account.json`)
4. Share your target Google Drive folder with the service account email
5. Set the folder ID in `config.yaml`

Backups run automatically on the configured interval (default: weekly). Backups older than 2 months are auto-deleted from both local storage and Google Drive.

## Available Commands

```bash
make help              # Show all available commands
make setup             # Install all dependencies
make dev               # Run backend + frontend dev servers
make build             # Build everything for production
make backup            # Trigger a manual backup
make docker-build      # Build Docker image
make docker-run        # Run with docker-compose
make test              # Run all tests
make clean             # Remove build artifacts
```

## Project Structure

```
expense-tracker/
├── backend/           # Go backend (API + static file server)
│   ├── cmd/server/    # Entry point
│   └── internal/      # Application code
├── frontend/          # React frontend
│   └── src/
├── mobile/            # Android WebView wrapper
│   └── android/
├── config.example.yaml
├── Makefile
├── Dockerfile
└── docker-compose.yml
```
