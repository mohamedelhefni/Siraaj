# Siraaj Load Testing Tools 🚀

This directory contains **extreme-speed** load testing tools optimized for ClickHouse's legendary ingestion performance.

## ⚡ Performance Highlights

- **200K-500K events/second** on modern hardware
- **10M events in ~30 seconds** with async_insert
- **100M events in ~3-5 minutes** with optimal settings
- Leverages ClickHouse's async_insert, parallel workers, and batch optimization

**[See EXTREME_SPEED_GUIDE.md for detailed performance tuning →](./EXTREME_SPEED_GUIDE.md)**

---

## Tools Available

### 1. General Load Test (`main.go`) - EXTREME SPEED EDITION ⚡
General-purpose load tester with ClickHouse-optimized extreme performance.

**Features:**
- 🔥 **async_insert** enabled by default for 5-10x speed
- 🚀 **Parallel workers** (10 concurrent) for maximum throughput
- 💾 **Large batches** (50K default) for optimal ClickHouse performance
- 📊 **Real-time metrics** showing instant speed and total throughput
- DB mode: Direct ClickHouse insertion (fastest)
- HTTP mode: HTTP API load testing
- CSV mode: Generate CSV files for import

**Quick Start:**
```bash
# EXTREME SPEED - 10M events (default, ~30 seconds)
go run main.go

# Custom 100M events (~3-5 minutes)
go run main.go -events=100000000 -batch=100000 -users=1000000

# Small test - 1M events
go run main.go -events=1000000 -batch=10000 -users=50000
```

**Usage:**
```bash
# Database mode - Default extreme speed
go run main.go -mode=db -events=10000000 -batch=50000 -users=100000

# HTTP mode - 50k events with 100 workers
go run main.go -mode=http -events=50000 -workers=100 -users=5000

# Custom ClickHouse DSN
go run main.go -mode=db -events=20000000 -dsn="clickhouse://localhost:9000/siraaj?username=default&password="
```

**Options:**
- `-mode`: `db`, `http`, or `csv` (default: `db`)
- `-events`: Total number of events (default: `10000000` - 10M)
- `-batch`: Batch size for DB mode (default: `50000` - optimized for ClickHouse)
- `-workers`: Number of concurrent workers for HTTP mode (default: `50`)
- `-users`: Number of unique users (default: `100000`)
- `-project`: Project ID for events (default: `test_project`)
- `-dsn`: ClickHouse DSN (default: `clickhouse://localhost:9000/siraaj?username=default&password=`)
- `-endpoint`: API endpoint for HTTP mode (default: `http://localhost:8080/api/events`)
- `-csv`: CSV file path for CSV mode (default: `../data/loadtest.csv`)

**Performance Tips:**
- Larger batches = faster (try `-batch=100000` with sufficient RAM)
- More users = more realistic data distribution
- async_insert is **automatically enabled** for maximum speed
- 10 parallel workers run simultaneously

---

### 2. Funnel Data Generator (`funnel/main.go`)
Specialized load tester that generates realistic user journey data for funnel analysis.

**Features:**
- Generates sequential events following realistic user journeys
- Multiple funnel templates (e-commerce, SaaS, content, etc.)
- Realistic drop-off rates at each step
- Time-distributed data over configurable date ranges
- Perfect for testing funnel analysis features

**Funnel Templates:**
1. **E-commerce Purchase** (7 steps)
   - Home → Product → Add to Cart → Cart → Checkout → Payment → Purchase
   - Drop-off rates: 20% → 30% → 40% → 35% → 50% → 10% → 5%

2. **E-commerce Browse** (3 steps)
   - Home → Product 1 → Product 2
   - High drop-off: 30% → 70% → 80%

3. **SaaS Activation** (6 steps)
   - Landing → Pricing → CTA Click → Signup → Dashboard → Feature Use
   - Drop-off rates: 25% → 40% → 15% → 45% → 30% → 25%

4. **Content Newsletter** (4 steps)
   - Blog → Article → CTA Click → Form Submit
   - Drop-off rates: 15% → 40% → 20% → 30%

5. **Mobile App Install** (4 steps)
   - Landing → CTA Click → App Install → App Open
   - Drop-off rates: 30% → 25% → 60% → 35%

6. **Video Engagement** (4 steps)
   - Features → Video Play → CTA Click → Signup
   - Drop-off rates: 20% → 45% → 35% → 50%

7. **Support Journey** (4 steps)
   - Help → Search → Support → Contact Click
   - Drop-off rates: 25% → 40% → 50% → 30%

**Usage:**
```bash
# Generate 10,000 user journeys directly to ClickHouse
cd funnel
go run main.go -mode=db -users=10000

# Generate 5,000 journeys via HTTP API
go run main.go -mode=http -users=5000

# Generate data for last 7 days instead of 30
go run main.go -mode=db -users=10000 -days=7

# Custom ClickHouse DSN
go run main.go -mode=db -users=20000 -dsn="clickhouse://localhost:9000/siraaj?username=default&password="

# Custom project ID
go run main.go -mode=db -users=10000 -project=my_project
```

**Options:**
- `-mode`: `db` or `http` (default: `db`)
- `-users`: Number of user journeys to generate (default: `10000`)
- `-project`: Project ID (default: `funnel_test`)
- `-dsn`: ClickHouse DSN for DB mode (default: `clickhouse://localhost:9000/siraaj?username=default&password=`)
- `-endpoint`: API endpoint for HTTP mode (default: `http://localhost:8080/api/track`)
- `-days`: Generate data for the last N days (default: `30`)

---

### 3. HTTP Load Tester (`http.go`)
High-performance HTTP-only load tester optimized for maximum throughput.

**Features:**
- Optimized for 10,000+ requests/sec
- Pre-generates events to eliminate generation overhead
- Large connection pool (1000 idle connections)
- Real-time progress reporting with instantaneous RPS

**Usage:**
```bash
# 100k requests with 200 workers
go run http.go 100000 200

# 500k requests with 500 workers (stress test)
go run http.go 500000 500

# Custom users and duration
go run http.go 100000 200 5000 300

# Custom server URL
go run http.go 10000 50 1000 30 http://localhost:8080
```

**Arguments:**
1. Total requests (default: `100000`)
2. Concurrency/workers (default: `100`)
3. Number of users (default: `10000`)
4. Duration in seconds (default: `600`)
5. Server URL (default: `http://localhost:8080`)

---

## Quick Start Examples

### Testing Funnel Analysis

**Step 1: Generate funnel data**
```bash
cd funnel
go run main.go -mode=db -users=10000 -days=30
```

**Step 2: Start the server**
```bash
cd ..
make run
```

**Step 3: Test funnel analysis**
- Open browser: `http://localhost:8080/dashboard/funnel`
- Try these funnel configurations:

**E-commerce Funnel:**
```
Step 1: event_name = "page_view", url = "/"
Step 2: event_name = "page_view", url = "/product/123"
Step 3: event_name = "add_to_cart", url = "/product/123"
Step 4: event_name = "checkout_started", url = "/checkout"
Step 5: event_name = "purchase", url = "/confirmation"
```

**SaaS Funnel:**
```
Step 1: event_name = "page_view", url = "/"
Step 2: event_name = "page_view", url = "/pricing"
Step 3: event_name = "signup", url = "/signup"
Step 4: event_name = "page_view", url = "/dashboard"
Step 5: event_name = "feature_used", url = "/dashboard"
```

**Content Funnel:**
```
Step 1: event_name = "page_view", url = "/blog"
Step 2: event_name = "page_view", url = "/blog/article-1"
Step 3: event_name = "button_click", url = "/blog/article-1"
Step 4: event_name = "form_submit", url = "/blog/article-1"
```

### Performance Testing

**HTTP API stress test:**
```bash
# Start server
make run

# In another terminal
go run http.go 100000 200
```

**Database throughput test:**
```bash
go run main.go -mode=db -events=500000 -batch=5000
```

---

## Expected Results

### Funnel Data Generator
- **User Journeys:** 10,000 users
- **Total Events:** ~40,000 events (avg 4 events per user)
- **Time:** ~30-60 seconds for 10k users (DB mode)
- **Completion Rates:** 
  - E-commerce purchase: ~2-5% complete all steps
  - SaaS activation: ~5-10% complete all steps
  - Content newsletter: ~15-25% complete all steps

### HTTP Load Tester
- **With Appender API:** 10,000+ req/s on modern hardware
- **With regular INSERT:** 1,000-2,000 req/s
- **Latency:** Avg 5-20ms per request

### Database Load Tester
- **Throughput:** 100,000-500,000 events/sec with ClickHouse batching
- **Time:** 100k events in 1-2 seconds

---

## Tips for Best Results

1. **Funnel Testing:**
   - Use `funnel/main.go` for realistic user journey data
   - Generate at least 10k users for meaningful statistics
   - Use 30 days time range for better temporal distribution

2. **Performance Testing:**
   - Use `http.go` for maximum HTTP throughput
   - Start with 100 workers, increase to 200-500 for stress testing
   - Monitor server CPU/memory during tests

3. **Data Generation:**
   - Use `main.go` with DB mode for fastest bulk data insertion
   - Use HTTP mode to test the full API stack
   - Combine both tools: funnel data for analysis, general load for volume

4. **ClickHouse:**
   - Ensure ClickHouse has enough disk space
   - Monitor table sizes with `SELECT table, formatReadableSize(sum(bytes)) FROM system.parts GROUP BY table`
   - Consider optimizing tables with `OPTIMIZE TABLE events FINAL`

---

## Troubleshooting

**"Cannot connect to server"**
- Ensure ClickHouse server is running on correct port
- Check firewall settings
- Verify DSN connection string

**"Connection refused"**
- Ensure ClickHouse is running: `clickhouse-client --query "SELECT 1"`
- Check ClickHouse port (default: 9000 for native, 8123 for HTTP)
- Verify database exists: `clickhouse-client --query "SHOW DATABASES"`

**Low throughput**
- Increase batch size (DB mode)
- Increase workers (HTTP mode)
- Check server resource usage
- Ensure ClickHouse has enough I/O capacity

**Out of memory**
- Reduce batch size
- Reduce number of concurrent workers
- Generate data in smaller chunks
- Check ClickHouse memory settings
