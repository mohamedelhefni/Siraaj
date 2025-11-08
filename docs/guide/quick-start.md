# Quick Start

Get Siraaj Analytics up and running in 3 minutes.

::: info Installation
Currently, Docker is the recommended installation method. Pre-built binaries are coming soon!
:::

## Docker Installation

**Prerequisites:** Docker installed

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name siraaj \
  mohamedelhefni/siraaj:latest
```

That's it! Siraaj is now running at `http://localhost:8080`

### Access the Dashboard

Open your browser:
```
http://localhost:8080/dashboard/
```

---

## Alternative: Build from Source

**Prerequisites:** Go 1.24+

```bash
# Clone repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# Build
go build -o siraaj

# Run
./siraaj
```

---

## Add Tracking to Your Website

### Step 1: Include the SDK

Add to your HTML `<head>`:

```html
<script>
  (function() {
    var script = document.createElement('script');
    script.src = 'http://localhost:8080/sdk/analytics.js';
    script.defer = true;
    document.head.appendChild(script);
    
    script.onload = function() {
      window.siraaj = new Analytics({
        apiUrl: 'http://localhost:8080',
        projectId: 'my-website',
        autoTrack: true  // Automatically track page views
      });
    };
  })();
</script>
```

### Step 2: Track Custom Events (Optional)

```javascript
// Track button clicks
siraaj.track('button_clicked', {
  button_id: 'signup',
  location: 'hero'
});

// Track form submissions
siraaj.track('form_submitted', {
  form_name: 'contact',
  page: 'landing'
});

// Identify users
siraaj.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium'
});
```

### Step 3: View Your Analytics

Open the dashboard:
```
http://localhost:8080/dashboard/
```

You should see:
- ✅ Real-time page views
- ✅ Unique visitors
- ✅ Geographic distribution
- ✅ Traffic channels
- ✅ Browser/OS/Device stats
- ✅ Custom events

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

### Test Event Tracking

```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "test_event",
    "user_id": "test-user",
    "project_id": "my-website",
    "url": "https://example.com/test",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
  }'
```

Expected response:
```json
{
  "status": "ok"
}
```

---

## Docker Compose (Recommended for Production)

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
      - PARQUET_FILE=/data/events
      - DUCKDB_MEMORY_LIMIT=4GB
      - DUCKDB_THREADS=4
      # - CORS=https://example.com,https://app.example.com
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

Run:
```bash
docker-compose up -d
```

View logs:
```bash
docker-compose logs -f siraaj
```

---

## Configuration (Optional)

### Environment Variables

```bash
# Server
export PORT=8080
export DB_PATH=data/analytics.db
export PARQUET_FILE=data/events

# DuckDB Performance
export DUCKDB_MEMORY_LIMIT=4GB
export DUCKDB_THREADS=4

# CORS (for production)
export CORS=https://example.com,https://app.example.com

# Run
./siraaj
```

### Using Different Port

```bash
PORT=3000 ./siraaj
```

Or with Docker:
```bash
docker run -d -p 3000:3000 -e PORT=3000 -v $(pwd)/data:/data mohamedelhefni/siraaj:latest
```

---

## Next Steps

<div class="tip custom-block">
  <p class="custom-block-title">📦 Installation Guide</p>
  <p>More installation options and detailed setup</p>
  <a href="/guide/installation">Installation Guide →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">⚙️ Configuration</p>
  <p>Customize Siraaj for your needs</p>
  <a href="/guide/configuration">Configuration Guide →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">🎯 SDK Integration</p>
  <p>Framework-specific guides (React, Vue, Svelte, Next.js)</p>
  <a href="/sdk/overview">SDK Documentation →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">📡 API Reference</p>
  <p>Complete API documentation</p>
  <a href="/api/overview">API Documentation →</a>
</div>

---

## Troubleshooting

### Port Already in Use

```bash
# Use a different port
PORT=3000 ./siraaj
```

### Permission Denied

```bash
chmod +x siraaj
./siraaj
```

### CORS Errors

Add your domain to CORS:

```bash
export CORS=https://your-domain.com
./siraaj
```

Or with Docker:
```bash
docker run -d \
  -p 8080:8080 \
  -e CORS=https://your-domain.com \
  -v $(pwd)/data:/data \
  mohamedelhefni/siraaj:latest
```

### No Data Showing

1. **Check if events are being tracked:**
   ```bash
   curl http://localhost:8080/api/debug/events
   ```

2. **Check SDK is loaded:**
   Open browser console and check for `siraaj` object

3. **Verify API endpoint:**
   ```bash
   curl http://localhost:8080/api/health
   ```

---

## Getting Help

- 📖 [Full Documentation](/)
- 💬 [GitHub Issues](https://github.com/mohamedelhefni/siraaj/issues)
- 📧 [Community Support](https://github.com/mohamedelhefni/siraaj/discussions)

---

**That's it!** You now have a fully functional analytics platform. Start tracking your website and gain valuable insights! 🎉
