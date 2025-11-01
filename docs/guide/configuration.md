# Configuration

Complete guide to configuring Siraaj Analytics.

## Environment Variables

Configure Siraaj using environment variables:

```bash
# Server
PORT=8080                    # Server port
DB_PATH=analytics.db         # Database file path

# CORS
CORS=https://example.com,https://app.example.com

# Geolocation (optional)
GEODB_PATH=data/geodb/dbip-country.mmdb

# Performance
MAX_WORKERS=10               # Number of worker goroutines
BUFFER_SIZE=1000            # Event buffer size
FLUSH_INTERVAL=10s          # Auto-flush interval
```

### Load from .env file

```bash
# Create .env file
cat > .env << EOF
PORT=8080
DB_PATH=./data/analytics.db
CORS=https://example.com
EOF

# Load and run
export $(cat .env | xargs)
./siraaj
```

## Command-line Flags

Override environment variables with flags:

```bash
./siraaj \
  --port 3000 \
  --db ./custom/analytics.db \
  --cors "https://example.com,https://app.example.com"
```

Available flags:

| Flag | Description | Default |
|------|-------------|---------|
| `--port` | Server port | 8080 |
| `--db` | Database path | analytics.db |
| `--cors` | Allowed origins | * |
| `--geodb` | GeoIP database path | data/geodb/dbip-country.mmdb |

## Configuration File

Create `config.yaml`:

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
