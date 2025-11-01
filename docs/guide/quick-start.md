# Quick Start

Get Siraaj Analytics up and running in just a few minutes.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.24+** installed ([Download](https://golang.org/dl/))
- **Git** for cloning the repository
- **Basic terminal knowledge**

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj
```

### 2. Build the Project

```bash
go build -o siraaj
```

This creates a binary named `siraaj` in your current directory.

### 3. Run the Server

```bash
./siraaj
```

The server will start on `http://localhost:8080`

:::tip Success!
You should see output like:
```
🚀 Siraaj Analytics Server
📊 Dashboard: http://localhost:8080/dashboard/
🔌 API: http://localhost:8080/api
```
:::

## Verify Installation

### Access the Dashboard

Open your browser and navigate to:

```
http://localhost:8080/dashboard/
```

You should see the Siraaj Analytics dashboard.

### Test the API

```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

## Add Tracking to Your Website

### 1. Include the SDK

Add this script to your website's `<head>` section:

```html
<!DOCTYPE html>
<html>
<head>
  <!-- Your other head tags -->
  
  <!-- Siraaj Analytics SDK -->
  <script src="http://localhost:8080/sdk/analytics.js"></script>
  <script>
    const analytics = new Analytics({
      apiUrl: 'http://localhost:8080',
      projectId: 'my-website',
      autoTrack: true
    });
  </script>
</head>
<body>
  <!-- Your website content -->
</body>
</html>
```

### 2. Track Custom Events (Optional)

```javascript
// Track button clicks
document.getElementById('signup-btn').addEventListener('click', () => {
  analytics.track('signup_clicked', {
    location: 'hero',
    plan: 'premium'
  });
});

// Track form submissions
analytics.trackForm('contact-form', {
  source: 'landing-page'
});
```

### 3. View Your Data

Open the dashboard at `http://localhost:8080/dashboard/` and you should see:
- Real-time page views
- Visitor statistics
- Event tracking data

## Using Docker (Alternative)

Prefer Docker? Here's the fastest way:

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name siraaj \
  mohamedelhefni/siraaj:latest
```

### Docker Compose

Create a `docker-compose.yml`:

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
    restart: unless-stopped
```

Run it:

```bash
docker-compose up -d
```

## Configuration

### Environment Variables

Create a `.env` file:

```bash
# Server Configuration
PORT=8080
DB_PATH=analytics.db

# CORS (optional)
CORS=https://mysite.com,https://app.mysite.com

# Geolocation (optional)
GEODB_PATH=data/geodb/dbip-country.mmdb
```

Load it:

```bash
export $(cat .env | xargs) && ./siraaj
```

### Command-line Flags

```bash
./siraaj --port 3000 --db analytics.db
```

Available flags:
- `--port`: Server port (default: 8080)
- `--db`: Database path (default: analytics.db)
- `--cors`: Allowed origins (comma-separated)

## Next Steps

Congratulations! 🎉 Siraaj is now running. Here's what to explore next:

### Learn More

<div class="tip custom-block">
  <p class="custom-block-title">📚 Configuration</p>
  <p>Customize Siraaj for your needs</p>
  <a href="/guide/configuration">Configuration Guide →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">🎯 SDK Integration</p>
  <p>Framework-specific guides</p>
  <a href="/sdk/overview">SDK Documentation →</a>
</div>

<div class="tip custom-block">
  <p class="custom-block-title">🚀 Production Deployment</p>
  <p>Deploy Siraaj to production</p>
  <a href="/guide/production">Production Guide →</a>
</div>

## Troubleshooting

### Port Already in Use

If port 8080 is already in use:

```bash
./siraaj --port 3000
```

Or set the environment variable:

```bash
PORT=3000 ./siraaj
```

### Permission Denied

If you get a permission error:

```bash
chmod +x siraaj
./siraaj
```

### Database Issues

If you encounter database errors, try deleting the database file:

```bash
rm analytics.db
./siraaj
```

:::warning
This will delete all your analytics data!
:::

### CORS Errors

If you're getting CORS errors in the browser:

```bash
CORS=https://your-domain.com ./siraaj
```

Or add to `.env`:
```
CORS=https://your-domain.com,http://localhost:3000
```

## Getting Help

- 📖 Check the [Configuration Guide](/guide/configuration)
- 💬 Open an issue on [GitHub](https://github.com/mohamedelhefni/siraaj/issues)
- 📧 Email: [your-email@example.com]

---

Ready to dive deeper? Continue to the [Installation Guide](/guide/installation) for more deployment options.
