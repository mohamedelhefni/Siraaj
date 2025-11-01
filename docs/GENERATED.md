# Siraaj Analytics Documentation

Complete VitePress documentation site generated from actual codebase analysis.

## 🎉 What Has Been Created

A comprehensive documentation website using VitePress with the default theme, including:

### 📁 Directory Structure
```
docs/
├── .vitepress/
│   ├── config.mjs          # VitePress configuration
│   └── cache/              # Build cache
├── index.md                # Landing page with hero
├── README.md               # Documentation README
├── package.json            # Dependencies & scripts
├── pnpm-lock.yaml          # Lock file
├── public/                 # Static assets
│   ├── logo.svg           # Siraaj logo
│   └── hero-image.svg     # Hero section image
├── guide/                  # User guides
│   ├── introduction.md
│   ├── quick-start.md
│   ├── installation.md
│   └── configuration.md
├── sdk/                    # SDK documentation
│   ├── overview.md
│   ├── core.md
│   ├── configuration.md
│   ├── vanilla.md
│   ├── react.md
│   ├── vue.md
│   ├── svelte.md
│   ├── nextjs.md
│   ├── nuxt.md
│   ├── preact.md
│   ├── custom-events.md
│   ├── user-identification.md
│   ├── auto-tracking.md
│   └── performance.md
└── api/                    # API reference
    └── overview.md
```

## 🔍 Technical Details from Codebase Analysis

### Architecture (from `main.go`)

**Tech Stack:**
- **Backend**: Go 1.24+ with official DuckDB driver (`github.com/duckdb/duckdb-go/v2`)
- **Database**: DuckDB (OLAP optimized)
- **Storage**: Parquet files with ZSTD compression
- **Frontend**: SvelteKit 2.0 (embedded)
- **Geolocation**: MaxMind DB-IP (optional)

**Clean Architecture Layers:**
1. **Domain Layer** (`internal/domain/`) - Business entities
2. **Repository Layer** (`internal/repository/`) - Data access with Parquet storage
3. **Service Layer** (`internal/service/`) - Business logic
4. **Handler Layer** (`internal/handler/`) - HTTP endpoints
5. **Middleware Layer** (`internal/middleware/`) - CORS & logging

### API Endpoints (from `main.go` & `event_handler.go`)

#### Event Tracking
- `POST /api/track` - Single event tracking
- `POST /api/track/batch` - Batch event tracking (max 100 events)

#### Analytics
- `GET /api/stats` - Comprehensive statistics
- `GET /api/stats/overview` - Top-level metrics
- `GET /api/stats/timeline` - Time-series data
- `GET /api/stats/pages` - Page analytics
- `GET /api/stats/pages/entry-exit` - Entry/exit pages
- `GET /api/stats/countries` - Geographic data
- `GET /api/stats/sources` - Traffic sources
- `GET /api/stats/events` - Event analytics
- `GET /api/stats/devices` - Browser/Device/OS breakdown

#### Other Endpoints
- `GET /api/events` - Raw events list
- `GET /api/online` - Real-time online users
- `GET /api/projects` - Project list
- `POST /api/funnel` - Funnel analysis
- `GET /api/channels` - Traffic channel breakdown
- `GET /api/health` - Health check
- `GET /api/geo` - Geolocation test

#### Debug Endpoints
- `GET /api/debug/events` - Last 50 events
- `GET /api/debug/storage` - Storage statistics

### Event Structure (from `domain/event.go`)

```go
type Event struct {
    ID              uint64    // Auto-incremented
    Timestamp       time.Time // ISO 8601
    EventName       string    // e.g., "page_view", "button_clicked"
    UserID          string    // Persistent user identifier
    SessionID       string    // Session identifier
    SessionDuration int       // Duration in seconds
    URL             string    // Page URL
    Referrer        string    // HTTP referrer
    UserAgent       string    // Browser user agent
    IP              string    // Client IP (from X-Forwarded-For or X-Real-IP)
    Country         string    // Detected country (via GeoIP)
    Browser         string    // Detected browser
    OS              string    // Detected operating system
    Device          string    // Desktop/Mobile/Tablet
    IsBot           bool      // Bot detection flag
    ProjectID       string    // Project identifier (default: "default")
    Channel         string    // Direct/Organic/Referral/Social/Paid
}
```

### Bot Detection (from `botdetector/botdetector.go`)

**Detected Bots:**
- **Search Engines**: Googlebot, Bingbot, Yahoo, DuckDuckBot, BaiduSpider, Yandex
- **Social Media**: Facebook, Twitter, LinkedIn, WhatsApp, Telegram, Discord, Slack
- **SEO Tools**: AhrefsBot, SEMrushBot, MJ12bot, DotBot, RogerBot, Screaming Frog
- **Monitoring**: Pingdom, UptimeRobot, NewRelic, StatusCake
- **HTTP Libraries**: cURL, Wget, Python Requests, Go HTTP Client, Axios
- **Headless Browsers**: Headless Chrome, PhantomJS, Selenium, Puppeteer

### Channel Detection (from `channeldetector/channeldetector.go`)

**Priority Order:**
1. **Paid** - UTM parameters (utm_medium=cpc/ppc/paid) or click IDs (gclid, fbclid, msclkid)
2. **Direct** - Empty referrer or same domain
3. **Social** - Facebook, Twitter, LinkedIn, Instagram, Pinterest, Reddit, TikTok, YouTube, etc.
4. **Organic** - Google, Bing, Yahoo, DuckDuckGo, Baidu, Yandex, etc.
5. **Referral** - All other external sources

### Parquet Storage (from `storage/parquet_storage.go`)

**Configuration:**
- Default buffer size: 10,000 events
- Default flush interval: 30 seconds
- Compression: ZSTD
- Row group size: 100,000
- Max files before merge: 100
- Merge check interval: 5 minutes

**Features:**
- Append-only partitioned files
- Automatic background flushing
- Automatic file merging when threshold reached
- Crash-safe with graceful shutdown

**Parquet Columns:**
```
id, timestamp, date_hour, date_day, date_month,
event_name, user_id, session_id, session_duration,
url, referrer, user_agent, ip, country,
browser, os, device, is_bot, project_id, channel
```

### DuckDB Optimizations (from `main.go`)

**Settings:**
- Memory limit: 4GB (configurable via `DUCKDB_MEMORY_LIMIT`)
- Threads: 4 (configurable via `DUCKDB_THREADS`)
- Connection pool: Max 5 open, 2 idle
- Enabled features:
  - Object cache
  - Force parallelism
  - HTTP metadata cache
  - Experimental parallel CSV

### Database Indexes (from `migrations/migrations.go`)

**Single Column Indexes:**
- `idx_timestamp` (DESC)
- `idx_event_name`
- `idx_user_id`
- `idx_session_id`
- `idx_country`
- `idx_referrer`
- `idx_project_id`
- `idx_is_bot`
- `idx_channel`

**Composite Indexes:**
- `idx_timestamp_event_name`
- `idx_timestamp_project`
- `idx_session_timestamp`
- `idx_user_timestamp`
- `idx_timestamp_url`
- `idx_timestamp_country`
- `idx_timestamp_is_bot`
- `idx_timestamp_browser`
- `idx_timestamp_device`
- `idx_timestamp_os`
- `idx_timestamp_referrer`
- `idx_timestamp_channel`
- `idx_pageview_covering` (covering index)

### Middleware (from `middleware/middleware.go`)

**CORS:**
- Configurable via `CORS` environment variable
- Default: `*` (allow all)
- Allowed methods: GET, POST, PUT, DELETE, OPTIONS
- Allowed headers: Content-Type, Authorization

**Logging:**
- Format: `{METHOD} {URI} {DURATION}`
- All requests logged to stdout

### Environment Variables

**Required:**
- None (all have defaults)

**Optional:**
- `PORT` - Server port (default: 8080)
- `DB_PATH` - Database path (default: data/analytics.db)
- `PARQUET_FILE` - Parquet directory (default: data/events)
- `CORS` - Allowed origins (default: *)
- `GEODB_PATH` - GeoIP database path
- `DUCKDB_MEMORY_LIMIT` - Memory limit (default: 4GB)
- `DUCKDB_THREADS` - Thread count (default: 4)

### SDK Event Tracking (from `sdk/src/core/analytics.ts`)

**Configuration Options:**
```typescript
{
  apiUrl: string              // Required
  projectId: string           // Required
  autoTrack?: boolean         // default: true
  debug?: boolean             // default: false
  bufferSize?: number         // default: 10, max: 100
  flushInterval?: number      // default: 30000ms
  timeout?: number            // default: 10000ms
  maxRetries?: number         // default: 3
  useBeacon?: boolean         // default: true
  sampling?: number           // default: 1.0 (100%)
  maxQueueSize?: number       // default: 50
  respectDoNotTrack?: boolean // default: true
  enablePerformanceTracking?: boolean  // default: false
}
```

**Auto-tracked Events:**
- Page views (on load & navigation)
- Link clicks
- Form submissions  
- JavaScript errors
- Unhandled promise rejections
- Page visibility changes
- Web Vitals (if enabled): LCP, FID, CLS

### Docker (from `Dockerfile`)

**Multi-stage Build:**
1. **dashboard-builder** - Node 24 Alpine, builds SvelteKit dashboard
2. **go-builder** - Go 1.24 Bookworm, builds Go binary with CGO
3. **Final** - Debian Bookworm Slim, runtime only

**Runtime:**
- Non-root user (UID: 1000)
- Health check: `/api/health` every 30s
- Data volume: `/data`
- Exposed port: 8080

## 📝 Running the Documentation

```bash
cd docs

# Install dependencies
pnpm install

# Start dev server
pnpm dev

# Build for production
pnpm build

# Preview production build
pnpm preview
```

The docs will be available at `http://localhost:5173`

## 🎨 Features

- ✅ VitePress default theme
- ✅ Responsive landing page with hero section
- ✅ Full navigation structure
- ✅ Search functionality
- ✅ Dark mode support
- ✅ Mobile-friendly
- ✅ Custom SVG logos and graphics
- ✅ Syntax highlighting for code blocks
- ✅ Auto-generated sidebar
- ✅ GitHub links
- ✅ Edit on GitHub links

## 📚 Documentation Sections

### Guide
- Introduction to Siraaj
- Quick Start (3-minute setup)
- Installation methods (binary, Docker, source)
- Configuration guide

### SDK
- Overview of all SDKs
- Core SDK documentation  
- Framework guides (React, Vue, Svelte, Next.js, Nuxt, Preact)
- Vanilla JavaScript guide
- Configuration options
- Custom events
- User identification
- Auto-tracking features
- Performance tracking

### API
- API overview
- All endpoints documented based on actual code
- Request/response examples
- Error handling

## 🔗 Next Steps

1. Start the dev server: `cd docs && pnpm dev`
2. Add remaining guide pages (architecture, dashboard, etc.)
3. Add more API endpoint documentation
4. Add code examples from actual usage
5. Add deployment guides
6. Customize theme colors if needed

All documentation is generated based on actual code analysis from the Siraaj codebase!
