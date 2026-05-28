# syntax=docker/dockerfile:1.7

# --- Stage 1: Build Frontend (SPA) ---
FROM oven/bun:alpine AS frontend-builder
WORKDIR /app/web

# Install dependencies
COPY web/package.json web/bun.lock ./
RUN --mount=type=cache,target=/root/.bun bun install --frozen-lockfile

# Build the frontend as static SPA
COPY web/ ./
RUN bun run build

# --- Stage 2: Build Backend ---
FROM golang:1.26.1-alpine AS backend-builder
WORKDIR /app

# Install deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy only the source required for the backend build
COPY cmd ./cmd
COPY config ./config
COPY internal ./internal
COPY pkg ./pkg
COPY server ./server
COPY sqlc ./sqlc
COPY docs ./docs

# Copy frontend dist from stage 1 into the backend build context
COPY --from=frontend-builder /app/web/dist ./web/dist

# Build backend binary
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pos-server ./cmd/app

# --- Stage 3: Final Image (minimal, no Node.js) ---
FROM alpine:3.21
WORKDIR /app

# Add ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy backend binary
COPY --from=backend-builder /app/pos-server ./pos-server

# Copy frontend static files
COPY --from=frontend-builder /app/web/dist ./web/dist

# Copy swagger docs
COPY --from=backend-builder /app/docs ./docs

# Expose only backend port
EXPOSE 8080

CMD ["./pos-server"]
