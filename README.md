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
- **API**: http://localhost:8080/api/

---

## 📦 Installation

### Method 1: Pre-built Binary

```bash
# Download the latest release
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-linux-amd64 -o siraaj
chmod +x siraaj
./siraaj
```

### Method 2: Build from Source

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

### Method 3: Docker (Coming Soon)

```bash
docker run -p 8080:8080 -v ./data:/data mohamedelhefni/siraaj
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

---

## 🔌 API Reference

### Track Event

```bash
POST /api/track
Content-Type: application/json

{
  "event_name": "page_view",
  "user_id": "user-123",
  "session_id": "session-456",
  "url": "https://example.com/page",
  "referrer": "https://google.com",
  "user_agent": "Mozilla/5.0...",
  "project_id": "my-website",
  "properties": "{\"custom\": \"data\"}"
}

Response: {"status": "ok"}
```

### Get Statistics

```bash
GET /api/stats?start=2024-01-01&end=2024-01-31&project=my-website&limit=50

Response:
{
  "total_events": 1500,
  "unique_users": 450,
  "total_visits": 1200,
  "page_views": 2000,
  "bounce_rate": 45.5,
  "events_change": 25.5,
  "users_change": 15.2,
  "timeline": [...],
  "timeline_format": "day",
  "top_events": [...],
  "top_pages": [...],
  "browsers": [...],
  "top_countries": [...],
  "top_sources": [...]
}
```

### Query Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `start` | Start date (YYYY-MM-DD) | `2024-01-01` |
| `end` | End date (YYYY-MM-DD) | `2024-01-31` |
| `project` | Filter by project ID | `my-website` |
| `source` | Filter by referrer | `google.com` |
| `country` | Filter by country | `United States` |
| `browser` | Filter by browser | `Chrome` |
| `event` | Filter by event name | `page_view` |
| `limit` | Result limit (max 1000) | `50` |

### Get Events

```bash
GET /api/events?start=2024-01-01&end=2024-01-31&limit=100&offset=0

Response:
{
  "events": [...],
  "total": 5000,
  "limit": 100,
  "offset": 0
}
```

### Online Users

```bash
GET /api/online?window=5

Response:
{
  "online_users": 12,
  "active_sessions": 18,
  "time_window_mins": 5,
  "cutoff_time": "2024-10-24T10:30:00Z"
}
```

### List Projects

```bash
GET /api/projects

Response: ["default", "website-1", "website-2"]
```

### Get Top Properties

```bash
GET /api/properties?start=2024-01-01&end=2024-01-31&limit=20

Response:
[
  {
    "key": "button_id",
    "value": "signup",
    "count": 1250,
    "event_types": 3
  },
  {
    "key": "plan",
    "value": "premium",
    "count": 890,
    "event_types": 5
  }
]
```

### Health Check

```bash
GET /api/health

Response:
{
  "status": "ok",
  "database": "duckdb",
  "version": "1.0.0",
  "geolocation": true
}
```

---

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

# Geolocation (optional)
GEOIP_DB_PATH=/path/to/dbip-city.mmdb

# CORS (optional)
CORS_ORIGINS=https://mysite.com,https://app.mysite.com
```

### Database

Siraaj uses DuckDB, a fast in-process analytical database. The database file (`analytics.db`) is created automatically on first run.

**Schema:**
- Auto-migrated on startup
- Indexed for fast queries
- Optimized for analytical workloads

---

## 🔒 Privacy & GDPR Compliance

Siraaj is designed with privacy in mind:

- ✅ No cookies used
- ✅ No cross-site tracking
- ✅ No PII collected by default
- ✅ IP addresses can be anonymized
- ✅ Respects Do Not Track (DNT)
- ✅ Self-hosted (you control the data)
- ✅ Easy data export/deletion

### Data Retention

Configure retention policies by manually cleaning old data:

```sql
-- Delete events older than 90 days
DELETE FROM events WHERE timestamp < NOW() - INTERVAL 90 DAYS;
```

---

## 📊 Performance

Siraaj is built for speed:

- **Query Performance**: Sub-second queries on millions of events
- **Ingestion**: 10,000+ events/second (single instance)
- **Memory**: Efficient memory usage with DuckDB
- **Scalability**: Vertical scaling (add more CPU/RAM)

### Benchmarks

```
Hardware: MacBook Pro M1, 16GB RAM
Dataset: 1M events

Query              | Time
-------------------|-------
Total stats        | 45ms
Timeline (30 days) | 67ms
Top pages          | 23ms
Country breakdown  | 31ms
```

---

## 🛠️ Development

### Prerequisites

```bash
# Install Go
brew install go  # macOS
# or download from https://golang.org/

# Install Node.js and pnpm
brew install node pnpm

# Install Air (hot reload for Go)
go install github.com/air-verse/air@latest
```

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

### Project Structure

```
.
├── cmd/                   # Command-line tools
├── internal/              # Private application code
│   ├── domain/           # Core business models
│   ├── repository/       # Database operations
│   ├── service/          # Business logic
│   ├── handler/          # HTTP handlers
│   └── middleware/       # HTTP middleware
├── dashboard/             # Frontend application
│   ├── src/
│   │   ├── routes/       # SvelteKit routes
│   │   └── lib/          # Shared components
│   └── static/           # Static assets
├── sdk/                   # JavaScript SDK
├── geolocation/          # GeoIP functionality
└── main.go               # Application entry
```

### Running Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/repository

# Integration tests
go test -tags=integration ./...
```

### Building for Production

```bash
# Build dashboard
cd dashboard
pnpm build
cd ..

# Build binary
go build -ldflags="-s -w" -o siraaj

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o siraaj-linux-amd64
GOOS=windows GOARCH=amd64 go build -o siraaj-windows-amd64.exe
GOOS=darwin GOARCH=arm64 go build -o siraaj-darwin-arm64
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'Add amazing feature'`)
4. **Push to the branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

### Code Style

- Go: Follow standard Go conventions (`gofmt`, `golint`)
- JavaScript: Prettier + ESLint
- Svelte: Standard Svelte formatting

### Commit Messages

Use conventional commits:
```
feat: add user authentication
fix: resolve timezone issue in charts
docs: update API documentation
refactor: improve query performance
test: add repository tests
```

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [DuckDB](https://duckdb.org/) - Amazing analytical database
- [SvelteKit](https://kit.svelte.dev/) - Fantastic web framework
- [Chart.js](https://www.chartjs.org/) - Beautiful charts
- [MaxMind](https://www.maxmind.com/) - GeoIP data
- [DB-IP](https://db-ip.com/) - Free GeoIP alternative

---

## 📮 Contact & Support

- **Issues**: [GitHub Issues](https://github.com/mohamedelhefni/siraaj/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mohamedelhefni/siraaj/discussions)
- **Author**: [@mohamedelhefni](https://github.com/mohamedelhefni)

---

## 🗺️ Roadmap

- [ ] Docker support with Docker Compose
- [ ] Kubernetes deployment manifests
- [ ] Email reports (daily/weekly/monthly)
- [ ] Webhook integrations
- [ ] Custom dashboards
- [ ] A/B testing support
- [ ] Funnel analysis
- [ ] User session recordings (optional)
- [ ] API rate limiting
- [ ] Multi-user authentication
- [ ] Database backup automation
- [ ] Data export (CSV, JSON)

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/mohamedelhefni">Mohamed Elhefni</a></p>
  <p>⭐ Star this repo if you find it useful!</p>
</div>
