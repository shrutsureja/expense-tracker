# ── Stage 1: Build React frontend ─────────────────────────────────────────────
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ── Stage 2: Build Go backend ─────────────────────────────────────────────────
FROM golang:1.24-alpine AS backend-builder
# gcc + musl-dev required for CGO (sqlite3)
RUN apk add --no-cache gcc musl-dev
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Copy built frontend into the embed directory
COPY --from=frontend-builder /app/backend/cmd/server/static ./cmd/server/static/
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /expense-tracker ./cmd/server/

# ── Stage 3: Minimal runtime ──────────────────────────────────────────────────
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /expense-tracker .
RUN mkdir -p /app/data
VOLUME ["/app/data"]
EXPOSE 8080
ENV PORT=8080
ENV DB_PATH=/app/data/expense-tracker.db
CMD ["./expense-tracker"]
