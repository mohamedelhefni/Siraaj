# Installation

Multiple ways to install and run Siraaj Analytics.

## Method 1: Pre-built Binary

Download and run the latest release:

```bash
# Linux (amd64)
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-linux-amd64 -o siraaj
chmod +x siraaj
./siraaj

# macOS (amd64)
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-darwin-amd64 -o siraaj
chmod +x siraaj
./siraaj

# macOS (arm64 / M1/M2/M3)
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-darwin-arm64 -o siraaj
chmod +x siraaj
./siraaj

# Windows
# Download from releases page and run
```

## Method 2: Docker (Recommended)

### Quick Start

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name siraaj \
  mohamedelhefni/siraaj:latest
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    container_name: siraaj
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - PORT=8080
      - DB_PATH=/data/analytics.db
      - GEODB_PATH=/data/geodb/dbip-country.mmdb
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Run:

```bash
docker-compose up -d
```

### Custom Build

```bash
# Clone repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# Build Docker image
docker build -t siraaj:custom .

# Run
docker run -d -p 8080:8080 -v $(pwd)/data:/data siraaj:custom
```

## Method 3: Build from Source

### Prerequisites

- Go 1.24 or higher
- Node.js 18+ and pnpm (for building dashboard)
- Git

### Build Steps

```bash
# 1. Clone repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# 2. Build dashboard
cd dashboard
pnpm install
pnpm build
cd ..

# 3. Build Go binary
go build -o siraaj

# 4. Run
./siraaj
```

### Development Build

```bash
# Install dependencies
go mod download
cd dashboard && pnpm install && cd ..

# Run backend (with hot reload using air)
air

# In another terminal, run frontend dev server
cd dashboard
pnpm dev
```

## Method 4: Go Install

```bash
go install github.com/mohamedelhefni/siraaj@latest
```

## Verify Installation

Check that Siraaj is running:

```bash
# Health check
curl http://localhost:8080/api/health

# Expected response
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": 3600
}
```

## Next Steps

- [Configuration →](/guide/configuration)
- [SDK Integration →](/sdk/overview)
- [Production Deployment →](/guide/production)
