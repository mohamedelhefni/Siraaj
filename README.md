<div align="center">
  <img src="./logo.png" alt="Siraaj Logo" width="200"/>
  <h1>Siraaj Analytics</h1>
  <p><strong>Privacy-First, Self-Hosted Web Analytics</strong></p>
  
  <p>
    <a href="#-features">Features</a> •
    <a href="#-quick-start">Quick Start</a> •
    <a href="#-installation">Installation</a> •
    <a href="#-documentation">Documentation</a>
  </p>

  <p>
    <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version"/>
    <img src="https://img.shields.io/badge/DuckDB-Powered-yellow?style=flat" alt="DuckDB"/>
    <img src="https://img.shields.io/badge/Svelte-5-FF3E00?style=flat&logo=svelte" alt="Svelte"/>
    <img src="https://img.shields.io/badge/License-AGPL--3.0-blue?style=flat" alt="License"/>
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

## 🚀 What is Siraaj?

Siraaj is a **lightweight, privacy-focused analytics platform** built with Go and DuckDB. It provides real-time insights into your web traffic without cookies, third-party tracking, or compromising user privacy.

**Key Highlights:**
- 🔒 **Privacy-First**: No cookies, GDPR compliant, your data stays on your server
- ⚡ **Lightning Fast**: DuckDB columnar storage delivers sub-50ms query performance
- 📊 **Beautiful Dashboard**: Modern SvelteKit 5 UI with real-time updates
- 🎯 **Simple Integration**: Drop-in JavaScript SDK (< 5KB gzipped)
- 🌍 **Multi-Project**: Track unlimited websites from one instance
- 🏗️ **Production-Ready**: Clean architecture, well-tested, Docker-ready

---

## ✨ Features

- **Core Analytics**: Page views, unique visitors, sessions, bounce rate, online users
- **Traffic Channels**: Automatic classification (Direct, Organic, Social, Referral, Paid)
- **Geographic Data**: Country-level tracking with optional MaxMind geolocation
- **Technical Insights**: Browser, OS, device type distribution
- **Custom Events**: Track any event with custom properties
- **Funnel Analysis**: Measure conversion through multi-step funnels
- **Real-Time Dashboard**: Live updates with configurable refresh intervals
- **Advanced Filtering**: Filter by project, date, country, browser, OS, device, source
- **Entry/Exit Pages**: Track where users enter and leave your site
- **Bot Detection**: Automatic bot filtering with 50+ known bots
- **Parquet Storage**: Efficient columnar storage with automatic file management

---

## 🏃 Quick Start

### Option 1: Docker (Recommended)

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name siraaj \
  mohamedelhefni/siraaj:latest
```

### Option 2: Binary

```bash
# Download and run
curl -L https://github.com/mohamedelhefni/siraaj/releases/latest/download/siraaj-linux-amd64 -o siraaj
chmod +x siraaj
./siraaj
```

### Option 3: From Source

```bash
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj
go build -o siraaj
./siraaj
```

**Access the dashboard:** http://localhost:8080/dashboard/

---

## 📦 Installation

### Docker Compose

Create `docker-compose.yml`:

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
      - DUCKDB_MEMORY_LIMIT=4GB
      - DUCKDB_THREADS=4
    restart: unless-stopped
```

Run: `docker-compose up -d`

### Build from Source

**Prerequisites**: Go 1.24+, Node.js 18+, pnpm

```bash
# Clone repository
git clone https://github.com/mohamedelhefni/siraaj.git
cd siraaj

# Build dashboard (optional - already embedded)
cd dashboard
pnpm install
pnpm build
cd ..

# Build Go binary
go build -o siraaj

# Run
./siraaj
```

---

## 📚 Usage

### Add to Your Website

Include the SDK in your HTML:

```html
<script>
  (function() {
    var script = document.createElement('script');
    script.src = 'http://your-server:8080/sdk/analytics.js';
    script.defer = true;
    document.head.appendChild(script);
    
    script.onload = function() {
      window.siraaj = new Analytics({
        apiUrl: 'http://your-server:8080',
        projectId: 'my-website',
        autoTrack: true
      });
    };
  })();
</script>
```

### Track Custom Events

```javascript
// Track button clicks
siraaj.track('button_clicked', {
  button_id: 'signup',
  location: 'hero'
});

// Track conversions
siraaj.track('purchase', {
  product: 'Premium Plan',
  price: 99.99,
  currency: 'USD'
});

// Identify users (optional)
siraaj.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium'
});
```

### Framework Integrations

**React/Next.js:**
```bash
npm install @hefni101/siraaj
```

See [SDK Documentation](./sdk/README.md) for React, Vue, Svelte, Next.js, and Nuxt integrations.

---

## ⚙️ Configuration

### Environment Variables

```bash
# Server
PORT=8080                          # Server port
DB_PATH=data/analytics.db          # DuckDB database path
PARQUET_FILE=data/events           # Parquet storage directory

# DuckDB Performance
DUCKDB_MEMORY_LIMIT=4GB            # Memory limit for DuckDB
DUCKDB_THREADS=4                   # Number of threads

# CORS (optional)
CORS=https://example.com,https://app.example.com
```

### API Endpoints

- `POST /api/track` - Track single event
- `POST /api/track/batch` - Track multiple events
- `GET /api/stats` - Get analytics statistics
- `GET /api/stats/overview` - Dashboard overview
- `GET /api/stats/timeline` - Timeline data
- `GET /api/stats/pages` - Top pages
- `GET /api/stats/countries` - Country distribution
- `GET /api/stats/sources` - Traffic sources
- `GET /api/stats/devices` - Browser/OS/Device stats
- `GET /api/channels` - Channel analytics
- `GET /api/online` - Real-time online users
- `GET /api/projects` - List projects
- `POST /api/funnel` - Funnel analysis
- `GET /api/health` - Health check

See [API Documentation](./docs/api/overview.md) for details.

---

## 🏗️ Architecture

Siraaj follows **Clean Architecture** principles for maintainability and testability:

```
siraaj/
├── main.go                    # Application entry point
├── internal/
│   ├── domain/               # Business entities (Event, Stats)
│   ├── repository/           # Data access layer (DuckDB/Parquet)
│   ├── service/              # Business logic
│   ├── handler/              # HTTP handlers
│   ├── middleware/           # CORS, logging
│   ├── migrations/           # Database migrations
│   ├── storage/              # Parquet storage engine
│   ├── botdetector/          # Bot detection (50+ known bots)
│   └── channeldetector/      # Traffic channel classification
├── sdk/
│   ├── analytics.js          # Vanilla JS SDK
│   └── dist/                 # Framework integrations (React, Vue, Svelte, etc.)
├── dashboard/                # SvelteKit 5 dashboard
├── geolocation/              # GeoIP service (MaxMind)
└── docs/                     # VitePress documentation

```

**Tech Stack:**
- **Backend**: Go 1.24, DuckDB (OLAP queries), Parquet (columnar storage)
- **Frontend**: SvelteKit 5, Chart.js, Tailwind CSS, shadcn-svelte
- **SDK**: Vanilla JS/TypeScript (UMD), framework adapters
- **Storage**: DuckDB + Parquet files with automatic merging
- **Deployment**: Docker, binary releases for Linux/macOS/Windows

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific tests
make test-unit          # Unit tests only
make test-integration   # Integration tests
make test-race          # Race condition detection

# Run benchmarks
make bench
```

**Test Coverage**: ~85% with unit, integration, and benchmark tests

---

## 📈 Performance

- **Query Speed**: < 50ms (p95) for analytics queries
- **Event Tracking**: < 20ms (p95) for batch inserts
- **Storage**: Columnar Parquet format with compression
- **Concurrency**: Handles thousands of concurrent events
- **Memory**: ~100MB base, configurable DuckDB memory limit

**Optimizations:**
- Parquet columnar storage with buffering (10k events)
- Automatic file merging (< 100 files maintained)
- DuckDB parallel query execution
- Indexed queries on timestamp, project, channel
- Bot traffic filtering

---

## 🔒 Privacy & Compliance

- **No cookies**: Uses anonymous user/session IDs
- **No cross-site tracking**: Data scoped per domain
- **GDPR compliant**: No personal data collected by default
- **Self-hosted**: Complete data ownership and control
- **Bot filtering**: Excludes known bots from analytics
- **IP anonymization**: Optional geolocation only stores country

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Development Setup:**
```bash
# Install dependencies
go mod download
cd dashboard && pnpm install && cd ..

# Run tests
make test

# Run with hot reload (requires air)
air
```

---

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

**Key Points:**
- ✅ Free to use, modify, and distribute
- ✅ Source code must be made available when distributed
- ✅ Network use is distribution (if you modify and host publicly, share the code)
- ✅ Commercial use allowed
- ⚠️ Changes must be documented
- ⚠️ Same license for derivatives

See [LICENSE](LICENSE) file for full details.

---

## 🙏 Acknowledgments

- **DuckDB** - Fast OLAP database engine
- **SvelteKit** - Modern web framework
- **MaxMind** - GeoIP database
- All contributors and users of Siraaj

---

<div align="center">
  <p>Built with ❤️ by <a href="https://github.com/mohamedelhefni">Mohamed Elhefni</a></p>
  <p>
    <a href="https://github.com/mohamedelhefni/siraaj">GitHub</a> •
    <a href="https://github.com/mohamedelhefni/siraaj/issues">Issues</a> •
    <a href="./docs">Documentation</a>
  </p>
  <p>⭐ Star this repo if you find it useful!</p>
</div>
