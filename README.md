<div align="center">
  <img src="./logo.png" alt="Siraaj Logo" width="200"/>
  <h1>Siraaj Analytics</h1>
  <p><strong>Fast, Simple, Self-Hosted Analytics</strong></p>
  
  <p>
    <a href="#features">Features</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#installation">Installation</a> •
    <a href="#documentation">Documentation</a> •
    <a href="#api">API</a> •
    <a href="#contributing">Contributing</a>
  </p>

  <p>
    <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version"/>
    <img src="https://img.shields.io/badge/DuckDB-Powered-yellow?style=flat" alt="DuckDB"/>
    <img src="https://img.shields.io/badge/Svelte-5-FF3E00?style=flat&logo=svelte" alt="Svelte"/>
    <img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="License"/>
    <a href="https://github.com/mohamedelhefni/siraaj/actions/workflows/test.yml">
      <img src="https://github.com/mohamedelhefni/siraaj/actions/workflows/test.yml/badge.svg" alt="Tests"/>
    </a>
    <a href="https://github.com/mohamedelhefni/siraaj/actions/workflows/docker-build-push.yml">
      <img src="https://github.com/mohamedelhefni/siraaj/actions/workflows/docker-build-push.yml/badge.svg" alt="Docker Build"/>
    </a>
    <a href="https://hub.docker.com/r/mohamedelhefni/siraaj">
      <img src="https://img.shields.io/docker/pulls/mohamedelhefni/siraaj?style=flat&logo=docker" alt="Docker Pulls"/>
    </a>
  </p>
</div>

---

## 🚀 Overview

Siraaj is a lightweight, privacy-focused analytics platform that you can self-host. Built with Go and DuckDB, it provides real-time insights into your web traffic without compromising user privacy or relying on third-party services.

### Why Siraaj?

- **🔒 Privacy First**: No cookies, no tracking across sites, fully GDPR compliant
- **⚡ Blazing Fast**: Powered by DuckDB for lightning-fast analytics queries
- **📊 Beautiful Dashboard**: Modern, responsive UI built with SvelteKit
- **🎯 Simple Integration**: Drop-in JavaScript SDK, works anywhere
- **🌍 Multi-Project**: Track multiple websites from one instance
- **📈 Real-Time**: See your data as it happens
- **💾 Self-Hosted**: Your data stays on your infrastructure
- **🏗️ Clean Architecture**: Maintainable, testable, production-ready code

---

## ✨ Features

### Analytics Capabilities

- **📊 Core Metrics**
  - Total events, unique visitors, page views
  - Total visits (sessions), bounce rate
  - Real-time online users

- **📈 Trend Analysis**
  - Compare current vs previous period
  - Percentage change indicators
  - Dynamic time granularity (hourly/daily/monthly)

- **🌐 Geographic Insights**
  - Country-level visitor tracking
  - IP geolocation support (optional)

- **🖥️ Technical Breakdown**
  - Browser statistics
  - Operating system distribution
  - Device types (desktop/mobile/tablet)

- **🔍 Traffic Analysis**
  - Top pages and URLs
  - Referrer sources
  - Custom event tracking

### Dashboard Features

- **🎨 Interactive UI**
  - Real-time data updates (10s, 30s, 1min, 5min intervals)
  - Click-to-filter on all metrics
  - URL parameter persistence (shareable links)
  - Responsive design for all devices

- **📅 Flexible Time Ranges**
  - Today, Yesterday
  - Last 7/30 days
  - This/Last month
  - Last 3/6 months
  - This year
  - Custom date range

- **🔎 Advanced Filtering**
  - Filter by project, source, country, browser, event
  - Combine multiple filters
  - Metric-specific chart views

---

## 🏃 Quick Start

### Prerequisites

- Go 1.24 or higher
- Node.js 18+ (for building dashboard)
- pnpm (optional, or use npm)

### Run in 3 Steps

```bash
# 1. Clone the repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# 2. Build the project
go build -o siraaj

# 3. Run the server
./siraaj
```

The server will start on `http://localhost:8080`

- **Dashboard**: http://localhost:8080/dashboard/

---

## 📦 Installation

### Method 2: Docker (Recommended)

```bash
# Pull and run the latest version
docker pull mohamedelhefni/siraaj:latest

docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name siraaj \
  mohamedelhefni/siraaj:latest
```

**With Docker Compose:**

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

Save as `docker-compose.yml` and run:
```bash
docker-compose up -d
```

### Method 3: Pre-built Binary

```bash
# Download the latest release
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-linux-amd64 -o siraaj
chmod +x siraaj
./siraaj
```

### Method 4: Build from Source

```bash
# Clone the repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# Build the dashboard (optional, already embedded)
cd dashboard
pnpm install
pnpm build
cd ..

# Build the Go binary
go build -o siraaj

# Run
./siraaj
```

---

## 📚 Documentation

### JavaScript SDK Integration

Add the SDK to your website:

```html
<!-- Include the SDK -->
<script src="http://your-server:8080/sdk/analytics.js"></script>

<!-- Initialize -->
<script>
  const analytics = new Analytics({
    apiUrl: 'http://your-server:8080',
    projectId: 'my-website',
    autoTrack: true, // Automatically track page views
    debug: false
  });
</script>
```

### Track Custom Events

```javascript
// Track a custom event
analytics.track('button_clicked', {
  button_id: 'signup',
  location: 'hero'
});

// Track page views manually
analytics.trackPageView();

// Track with custom properties
analytics.track('purchase', {
  product: 'Premium Plan',
  price: 99.99,
  currency: 'USD'
});

// Identify users (optional)
analytics.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium'
});
```

### SDK Configuration Options

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',    // Your analytics server URL
  projectId: 'default',               // Project identifier
  autoTrack: true,                    // Auto-track page views
  bufferSize: 10,                     // Events to buffer before sending
  flushInterval: 30000,               // Auto-flush interval (ms)
  debug: false,                       // Enable debug logging
  respectDoNotTrack: true             // Honor DNT header
});
```

## 🏗️ Architecture

Siraaj follows Clean Architecture principles:

```
siraaj/
├── main.go                 # Application entry point
├── internal/
│   ├── domain/            # Business entities
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   └── migrations/        # Database migrations
├── sdk/
│   └── analytics.js       # JavaScript SDK
├── dashboard/             # SvelteKit dashboard
│   ├── src/
│   │   ├── routes/        # Pages
│   │   └── lib/           # Components
│   └── package.json
├── geolocation/           # GeoIP service
└── ui/                    # Static assets
```

### Technology Stack

**Backend:**
- Go 1.24+ (HTTP server, business logic)
- DuckDB (Analytics database, OLAP queries)
- MaxMind DB-IP (Geolocation)

**Frontend:**
- SvelteKit 2.0 (UI framework)
- Chart.js 4.4 (Data visualization)
- Tailwind CSS (Styling)
- shadcn-svelte (UI components)

**SDK:**
- Vanilla JavaScript (No dependencies)
- UMD format (Works everywhere)

---

## ⚙️ Configuration

### Environment Variables

```bash
# Server configuration
PORT=8080
DB_PATH=analytics.db

# CORS (optional)
CORS=https://mysite.com,https://app.mysite.com
```

---

### Data Retention

Configure retention policies by manually cleaning old data:

```sql
-- Delete events older than 90 days
DELETE FROM events WHERE timestamp < NOW() - INTERVAL 90 DAYS;
```

---

## 🛠️ Development

### Development Workflow

```bash
# 1. Clone and setup
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# 2. Install dependencies
go mod download
cd dashboard && pnpm install && cd ..

# 3. Run with hot reload
air  # Backend with auto-reload

# In another terminal
cd dashboard
pnpm dev  # Frontend dev server on :5173
```

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/mohamedelhefni">Mohamed Elhefni</a></p>
  <p>⭐ Star this repo if you find it useful!</p>
</div>
