# Configuration

Complete guide to configuring Siraaj Analytics.

## Environment Variables

Configure Siraaj using environment variables:

```bash
# Server
PORT=8080                           # Server port (default: 8080)
DB_PATH=data/analytics.db           # DuckDB database path
PARQUET_FILE=data/events            # Parquet storage directory

# DuckDB Performance
DUCKDB_MEMORY_LIMIT=4GB             # Memory limit (default: 4GB)
DUCKDB_THREADS=4                    # Number of threads (default: 4)

# CORS
CORS=https://example.com,https://app.example.com  # Allowed origins (comma-separated)
```

### Load from File

Create `.env` file:

```bash
PORT=8080
DB_PATH=./data/analytics.db
PARQUET_FILE=./data/events
DUCKDB_MEMORY_LIMIT=4GB
DUCKDB_THREADS=4
CORS=https://example.com,https://app.example.com
```

Load and run:

```bash
export $(cat .env | xargs)
./siraaj
```

---

## Server Configuration

### Port

```bash
# Custom port
PORT=3000 ./siraaj

# Or with Docker
docker run -d -p 3000:3000 -e PORT=3000 mohamedelhefni/siraaj:latest
```

### Database Path

```bash
# Absolute path
DB_PATH=/var/lib/siraaj/analytics.db ./siraaj

# Relative path
DB_PATH=./custom/analytics.db ./siraaj
```

**Note:** DuckDB creates the database file automatically if it doesn't exist.

---

## Storage Configuration

### Parquet Storage

```bash
# Set storage directory
PARQUET_FILE=data/events ./siraaj

# Docker
docker run -d \
  -e PARQUET_FILE=/data/events \
  -v $(pwd)/data:/data \
  mohamedelhefni/siraaj:latest
```

**Storage Structure:**
- Events are buffered (10,000 events default)
- Flushed every 30 seconds
- Stored as compressed Parquet files
- Automatic file merging when > 100 files

---

## DuckDB Performance Tuning

### Memory Limit

```bash
# Increase for high-traffic sites
DUCKDB_MEMORY_LIMIT=8GB ./siraaj

# Decrease for resource-constrained environments
DUCKDB_MEMORY_LIMIT=1GB ./siraaj
```

### Thread Count

```bash
# Match your CPU cores
DUCKDB_THREADS=8 ./siraaj

# For containers, set to container CPU limit
DUCKDB_THREADS=2 ./siraaj
```

**Recommendations:**
- **Low traffic** (< 10k events/day): 1-2 threads, 1GB memory
- **Medium traffic** (10k-100k events/day): 2-4 threads, 4GB memory
- **High traffic** (> 100k events/day): 4-8 threads, 8GB+ memory

---

## CORS Configuration

### Allow Specific Domains

```bash
CORS=https://example.com,https://app.example.com ./siraaj
```

### Allow All Domains (Development Only)

```bash
CORS=* ./siraaj
```

:::warning Security
Never use `CORS=*` in production! Always specify exact domains.
:::

### Docker CORS

```yaml
services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    environment:
      - CORS=https://example.com,https://app.example.com
```

---

## Docker Configuration

### Docker Compose

Full production configuration:

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
      # Server
      - PORT=8080
      - DB_PATH=/data/analytics.db
      - PARQUET_FILE=/data/events
      
      # Performance
      - DUCKDB_MEMORY_LIMIT=4GB
      - DUCKDB_THREADS=4
      
      # Security
      - CORS=https://example.com,https://app.example.com
      
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

---

## Data Retention

DuckDB and Parquet storage grow over time. Clean old data periodically:

### Manual Cleanup

```sql
-- Connect to database
duckdb data/analytics.db

-- Delete events older than 90 days
DELETE FROM events WHERE timestamp < NOW() - INTERVAL 90 DAYS;

-- Compact database
VACUUM;
```

### Automated Cleanup Script

```bash
#!/bin/bash
# cleanup.sh

DB_PATH="data/analytics.db"
RETENTION_DAYS=90

duckdb "$DB_PATH" <<SQL
DELETE FROM events 
WHERE timestamp < CURRENT_TIMESTAMP - INTERVAL '$RETENTION_DAYS days';
VACUUM;
SQL

echo "Cleanup complete. Deleted events older than $RETENTION_DAYS days"
```

Run with cron:

```bash
# Run daily at 2 AM
0 2 * * * /path/to/cleanup.sh
```

---

## Logging

### Log Level

Siraaj logs to stdout. Control verbosity:

```bash
# Production: info level (default)
./siraaj 2>&1 | tee siraaj.log

# Development: verbose logging
./siraaj
```

### Docker Logs

```bash
# View logs
docker logs siraaj

# Follow logs
docker logs -f siraaj

# Last 100 lines
docker logs --tail 100 siraaj
```

---

## Health Checks

### Health Endpoint

```bash
curl http://localhost:8080/api/health
```

Response:

```json
{
  "status": "healthy",
  "database": "connected"
}
```

### Docker Healthcheck

Already configured in recommended docker-compose.yml:

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:8080/api/health"]
  interval: 30s
  timeout: 3s
  retries: 3
  start_period: 5s
```

---

## Production Checklist

✅ Set specific CORS origins (not `*`)  
✅ Configure adequate memory limit  
✅ Set appropriate thread count  
✅ Mount persistent volume for data  
✅ Enable healthchecks  
✅ Configure restart policy  
✅ Set up log rotation  
✅ Plan data retention policy  
✅ Monitor disk space  
✅ Test backup/restore  

---

## Environment-Specific Configs

### Development

```yaml
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
      - CORS=*  # Allow all in dev
```

### Staging

```yaml
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
      - DUCKDB_MEMORY_LIMIT=2GB
      - CORS=https://staging.example.com
    restart: unless-stopped
```

### Production

```yaml
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

## Troubleshooting

### High Memory Usage

```bash
# Reduce memory limit
DUCKDB_MEMORY_LIMIT=2GB ./siraaj
```

### Slow Queries

```bash
# Increase threads
DUCKDB_THREADS=8 ./siraaj
```

### CORS Errors

```bash
# Check CORS configuration
echo $CORS

# Add your domain
export CORS=https://example.com
./siraaj
```

### Disk Space Issues

```bash
# Check database size
du -h data/analytics.db

# Check Parquet files
du -sh data/events/

# Clean old data (see Data Retention section)
```

---

## Next Steps

- [Funnels →](/guide/funnels)
- [Channels →](/guide/channels)
- [SDK Integration →](/sdk/overview)
- [API Reference →](/api/overview)

```yaml
server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 30s
  write_timeout: 30s
  
database:
  path: ./data/analytics.db
  max_connections: 10
  
cors:
  allowed_origins:
    - https://example.com
    - https://app.example.com
  allowed_methods:
    - GET
    - POST
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization
    
geolocation:
  enabled: true
  database_path: ./data/geodb/dbip-country.mmdb
  
performance:
  max_workers: 10
  buffer_size: 1000
  flush_interval: 10s
  
logging:
  level: info
  format: json
```

Load configuration:

```bash
./siraaj --config config.yaml
```

## CORS Configuration

### Allow Specific Domains

```bash
CORS=https://example.com,https://app.example.com ./siraaj
```

### Allow All Domains (Development Only)

```bash
CORS=* ./siraaj
```

:::warning Security
Never use `CORS=*` in production!
:::

### Docker CORS

```yaml
services:
  siraaj:
    image: mohamedelhefni/siraaj:latest
    environment:
      - CORS=https://example.com,https://app.example.com
```

## Database Configuration

### SQLite (Default)

```bash
DB_PATH=./data/analytics.db
```

### In-Memory (Testing)

```bash
DB_PATH=:memory:
```

:::warning Data Loss
In-memory databases are cleared on restart!
:::

### Custom Path

```bash
# Absolute path
DB_PATH=/var/lib/siraaj/analytics.db

# Relative path
DB_PATH=./custom/location/analytics.db
```

## Geolocation

Enable country-level geolocation:

### Download Database

```bash
# Create directory
mkdir -p data/geodb

# Download DB-IP database (free)
curl -L https://download.db-ip.com/free/dbip-country-lite-2024-01.mmdb.gz \
  | gunzip > data/geodb/dbip-country.mmdb
```

### Configure Path

```bash
GEODB_PATH=data/geodb/dbip-country.mmdb ./siraaj
```

### Disable Geolocation

```bash
# Don't set GEODB_PATH or set to empty
GEODB_PATH= ./siraaj
```

## Performance Tuning

### Worker Pool

```bash
# Increase for high traffic
MAX_WORKERS=20
```

### Buffer Size

```bash
# Larger buffer for batch processing
BUFFER_SIZE=5000
```

### Flush Interval

```bash
# Flush more frequently
FLUSH_INTERVAL=5s

# Flush less frequently (better performance)
FLUSH_INTERVAL=30s
```

## Logging

### Log Level

```bash
LOG_LEVEL=debug  # debug, info, warn, error
```

### Log Format

```bash
LOG_FORMAT=json  # json, text
```

### Log to File

```bash
./siraaj 2>&1 | tee siraaj.log
```

## Production Configuration

Example production setup:

```yaml
# config.production.yaml
server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 30s
  write_timeout: 30s
  
database:
  path: /var/lib/siraaj/analytics.db
  max_connections: 50
  
cors:
  allowed_origins:
    - https://example.com
    - https://app.example.com
    
geolocation:
  enabled: true
  database_path: /var/lib/siraaj/geodb/dbip-country.mmdb
  
performance:
  max_workers: 20
  buffer_size: 5000
  flush_interval: 10s
  
logging:
  level: info
  format: json
  
security:
  rate_limit: 1000  # requests per minute
  max_event_size: 1048576  # 1MB
```

## Security

### Rate Limiting

```bash
RATE_LIMIT=1000  # requests per minute
```

### Max Event Size

```bash
MAX_EVENT_SIZE=1048576  # 1MB in bytes
```

### Authentication (Coming Soon)

```bash
API_KEY=your-secret-key
```

## Data Retention

Configure automatic data cleanup:

```yaml
retention:
  enabled: true
  max_age_days: 90  # Delete events older than 90 days
  cleanup_interval: 24h  # Run cleanup daily
```

Or manually clean old data:

```sql
-- Delete events older than 90 days
DELETE FROM events WHERE timestamp < datetime('now', '-90 days');

-- Vacuum to reclaim space
VACUUM;
```

## Health Checks

Configure health check endpoint:

```bash
# Access health status
curl http://localhost:8080/api/health

# Response
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": 3600,
  "database": "ok",
  "workers": 10
}
```

## Monitoring

### Prometheus Metrics (Coming Soon)

```bash
METRICS_ENABLED=true
METRICS_PORT=9090
```

Access metrics:

```bash
curl http://localhost:9090/metrics
```

## Troubleshooting

### Check Configuration

```bash
./siraaj --version
./siraaj --help
```

### Verify Database

```bash
# Check database file
ls -lh analytics.db

# Check database size
du -h analytics.db
```

### Test CORS

```bash
curl -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -X OPTIONS \
  http://localhost:8080/api/track
```

## Next Steps

- [Production Deployment →](/guide/production)
- [Scaling Guide →](/guide/scaling)
- [Dashboard Guide →](/guide/dashboard)
