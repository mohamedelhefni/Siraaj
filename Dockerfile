# Build stage for dashboard
FROM node:24-alpine AS dashboard-builder

WORKDIR /app/dashboard

# Copy dashboard files
COPY dashboard/package.json dashboard/pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm && \
    pnpm install --frozen-lockfile

# Copy dashboard source
COPY dashboard/ ./

# Accept PUBLIC_API_URL as build argument (defaults to empty for relative path)
ARG PUBLIC_API_URL=
ENV PUBLIC_API_URL=${PUBLIC_API_URL}

# Build dashboard
RUN pnpm build

# Build stage for Go application
FROM golang:1.24-bookworm AS go-builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    g++ \
    make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy built dashboard from previous stage
# SvelteKit builds to ../ui/dashboard, so we copy from the build location
COPY --from=dashboard-builder /app/ui/dashboard ./ui/dashboard

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o siraaj .

# Final stage - minimal runtime image
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1000 siraaj && \
    useradd -r -u 1000 -g siraaj siraaj

# Create data directory
RUN mkdir -p /data && \
    chown -R siraaj:siraaj /data

WORKDIR /app

# Copy binary from builder
COPY --from=go-builder /app/siraaj .

# Copy geolocation database directory if needed
RUN mkdir -p /app/geolocation

# Switch to non-root user
USER siraaj

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080 \
    DB_PATH=/data/analytics.db

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./siraaj"]
