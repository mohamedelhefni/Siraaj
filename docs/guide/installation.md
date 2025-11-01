# Installation

Complete installation guide for Siraaj Analytics with multiple deployment options.

## Method 1: Docker (Recommended)

Docker provides the easiest and most reliable deployment method.

### Quick Start

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name siraaj \
  mohamedelhefni/siraaj:latest
```

### Docker Compose (Production Recommended)

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
      # Server Configuration
      - PORT=8080
      - DB_PATH=/data/analytics.db
      - PARQUET_FILE=/data/events
      
      # DuckDB Performance
      - DUCKDB_MEMORY_LIMIT=4GB
      - DUCKDB_THREADS=4
      
      # CORS - Add your domains
      - CORS=https://example.com,https://app.example.com
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

Run:

```bash
docker-compose up -d
```

View logs:

```bash
docker-compose logs -f siraaj
```

Stop:

```bash
docker-compose down
```

### Custom Docker Build

Build from source:

```bash
# Clone repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# Build image
docker build -t siraaj:custom .

# Run
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  siraaj:custom
```

---

## Method 2: Pre-built Binary

Download and run standalone binaries (no dependencies required).

### Linux (amd64)

```bash
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-linux-amd64 -o siraaj
chmod +x siraaj
./siraaj
```

### Linux (arm64)

```bash
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-linux-arm64 -o siraaj
chmod +x siraaj
./siraaj
```

### macOS (Intel)

```bash
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-darwin-amd64 -o siraaj
chmod +x siraaj
./siraaj
```

### macOS (Apple Silicon M1/M2/M3)

```bash
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-darwin-arm64 -o siraaj
chmod +x siraaj
./siraaj
```

### Windows

Download from [GitHub Releases](https://github.com/mohamedelhefni/siraaj/releases/latest) and run `siraaj.exe`.

---

## Method 3: Build from Source

Build Siraaj yourself for development or customization.

### Prerequisites

- **Go 1.24+** ([Download](https://golang.org/dl/))
- **Node.js 18+** and **pnpm** (for building dashboard)
- **Git**

### Build Steps

```bash
# 1. Clone repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# 2. Download Go dependencies
go mod download

# 3. Build dashboard (optional - already embedded in releases)
cd dashboard
pnpm install
pnpm build
cd ..

# 4. Build Go binary
go build -o siraaj

# 5. Run
./siraaj
```

### Development Build with Hot Reload

```bash
# Install dependencies
go mod download
cd dashboard && pnpm install && cd ..

# Terminal 1: Run backend with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Terminal 2: Run frontend dev server
cd dashboard
pnpm dev
```

Frontend dev server runs on `http://localhost:5173`  
Backend API runs on `http://localhost:8080`

---

## Method 4: Go Install

Install directly using Go:

```bash
go install github.com/mohamedelhefni/siraaj@latest

# Run (make sure $GOPATH/bin is in your PATH)
siraaj
```

---

## Verify Installation

### Health Check

```bash
curl http://localhost:8080/api/health
```

Expected response:

```json
{
  "status": "healthy",
  "database": "connected"
}
```

### Access Dashboard

Open browser: `http://localhost:8080/dashboard/`

### Test Event Tracking

```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "test_event",
    "user_id": "test-user",
    "project_id": "default",
    "url": "https://example.com",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
  }'
```

---

## Systemd Service (Linux)

Create a systemd service for automatic startup.

### Create Service File

```bash
sudo nano /etc/systemd/system/siraaj.service
```

Add:

```ini
[Unit]
Description=Siraaj Analytics
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/siraaj
ExecStart=/opt/siraaj/siraaj
Restart=on-failure
RestartSec=5s

# Environment variables
Environment="PORT=8080"
Environment="DB_PATH=/var/lib/siraaj/analytics.db"
Environment="PARQUET_FILE=/var/lib/siraaj/events"
Environment="DUCKDB_MEMORY_LIMIT=4GB"
Environment="DUCKDB_THREADS=4"

[Install]
WantedBy=multi-user.target
```

### Enable and Start

```bash
# Create directories
sudo mkdir -p /opt/siraaj /var/lib/siraaj
sudo chown www-data:www-data /var/lib/siraaj

# Copy binary
sudo cp siraaj /opt/siraaj/
sudo chown www-data:www-data /opt/siraaj/siraaj

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable siraaj
sudo systemctl start siraaj

# Check status
sudo systemctl status siraaj

# View logs
sudo journalctl -u siraaj -f
```

---

## Reverse Proxy Setup

### Nginx

```nginx
server {
    listen 80;
    server_name analytics.example.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name analytics.example.com;
    
    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/analytics.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/analytics.example.com/privkey.pem;
    
    # Proxy to Siraaj
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### Caddy

```caddy
analytics.example.com {
    reverse_proxy localhost:8080
}
```

---

## Environment-Specific Configurations

### Development

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - PORT=8080
      - DB_PATH=/data/analytics.db
      - CORS=*  # Allow all origins in dev
```

### Production

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - PORT=8080
      - DB_PATH=/data/analytics.db
      - PARQUET_FILE=/data/events
      - DUCKDB_MEMORY_LIMIT=8GB
      - DUCKDB_THREADS=8
      - CORS=https://example.com,https://app.example.com
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

---

## Data Persistence

### Docker Volumes

Named volume (recommended):

```yaml
services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    volumes:
      - siraaj-data:/data

volumes:
  siraaj-data:
```

Bind mount (easier backup):

```yaml
services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    volumes:
      - ./data:/data
```

### Backup Data

```bash
# Stop Siraaj
docker-compose down

# Backup
tar -czf siraaj-backup-$(date +%Y%m%d).tar.gz data/

# Restart
docker-compose up -d
```

### Restore Data

```bash
# Stop Siraaj
docker-compose down

# Restore
tar -xzf siraaj-backup-20241101.tar.gz

# Restart
docker-compose up -d
```

---

## Resource Requirements

### Minimum

- **CPU**: 1 core
- **RAM**: 512MB
- **Disk**: 1GB (+ data storage)
- **Network**: 1Mbps

### Recommended

- **CPU**: 2+ cores
- **RAM**: 2GB (4GB for DuckDB memory limit)
- **Disk**: 10GB+ SSD
- **Network**: 10Mbps+

### High Traffic (100k+ events/day)

- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Disk**: 50GB+ SSD
- **Network**: 100Mbps+

---

## Next Steps

<div class="tip custom-block">
  <p class="custom-block-title">⚙️ Configuration</p>
  <p>Configure Siraaj for your environment</p>
  <a href="/guide/configuration">Configuration Guide →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">🎯 SDK Integration</p>
  <p>Add tracking to your website</p>
  <a href="/sdk/overview">SDK Documentation →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">🔒 Security</p>
  <p>Secure your installation</p>
  <a href="/guide/security">Security Guide →</a>
</div>

---

## Troubleshooting

### Docker Issues

```bash
# View logs
docker logs siraaj

# Restart container
docker restart siraaj

# Remove and recreate
docker-compose down
docker-compose up -d
```

### Permission Issues

```bash
# Fix data directory permissions
sudo chown -R 1000:1000 data/
```

### Build Errors

```bash
# Clear Go cache
go clean -cache -modcache

# Re-download dependencies
go mod download
go build -o siraaj
```

---

Need help? [Open an issue](https://github.com/mohamedelhefni/siraaj/issues) or check the [FAQ](/guide/faq).
