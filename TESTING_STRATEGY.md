# Channel Analytics Feature - Testing Strategy

## Overview
This document outlines the testing strategy for the Channel Analytics feature, which classifies incoming events based on their traffic source.

## Components to Test

### 1. Channel Detector (`internal/channeldetector/`)

#### Unit Tests (✅ Already Implemented)
- **Test File**: `channeldetector_test.go`
- **Coverage**:
  - Direct traffic detection (empty referrer, same domain)
  - Organic search detection (Google, Bing, Yahoo, etc.)
  - Social media detection (Facebook, Twitter, LinkedIn, etc.)
  - Paid advertising detection (UTM parameters, click IDs)
  - Referral traffic detection
  - Edge cases (malformed URLs, null values, case sensitivity)

#### Run Tests:
```bash
cd /Users/hefni/projects/siraaj
go test ./internal/channeldetector/... -v
```

### 2. Database Migrations

#### Manual Testing Steps:
1. **Run migrations** to add the channel column:
```bash
# Start the server (migrations run automatically)
go run main.go
```

2. **Verify schema**:
```bash
# Connect to DuckDB and check schema
duckdb data/analytics.db
> DESCRIBE events;
> SELECT * FROM schema_migrations ORDER BY version;
```

Expected results:
- `channel` column should exist in events table
- Migrations 8, 9, and 10 should be applied
- Indexes should be created on channel column

### 3. Event Tracking with Channel Detection

#### Test Scenarios:

**Scenario 1: Direct Traffic**
```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "test-user-1",
    "session_id": "session-1",
    "url": "https://example.com/home",
    "referrer": "",
    "user_agent": "Mozilla/5.0"
  }'
```
Expected: `channel` = "Direct"

**Scenario 2: Organic Search**
```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "test-user-2",
    "session_id": "session-2",
    "url": "https://example.com/products",
    "referrer": "https://www.google.com/search?q=test",
    "user_agent": "Mozilla/5.0"
  }'
```
Expected: `channel` = "Organic"

**Scenario 3: Social Media**
```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "test-user-3",
    "session_id": "session-3",
    "url": "https://example.com/blog",
    "referrer": "https://www.facebook.com",
    "user_agent": "Mozilla/5.0"
  }'
```
Expected: `channel` = "Social"

**Scenario 4: Paid Advertising**
```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "test-user-4",
    "session_id": "session-4",
    "url": "https://example.com/landing?utm_source=google&utm_medium=cpc",
    "referrer": "https://www.google.com",
    "user_agent": "Mozilla/5.0"
  }'
```
Expected: `channel` = "Paid"

**Scenario 5: Referral Traffic**
```bash
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "test-user-5",
    "session_id": "session-5",
    "url": "https://example.com/article",
    "referrer": "https://news.ycombinator.com",
    "user_agent": "Mozilla/5.0"
  }'
```
Expected: `channel` = "Referral"

**Verify stored data**:
```bash
duckdb data/analytics.db
> SELECT event_name, url, referrer, channel FROM events ORDER BY timestamp DESC LIMIT 10;
```

### 4. API Endpoint Testing

#### Test GET /api/channels

**Basic Request**:
```bash
curl "http://localhost:8080/api/channels?start=2025-01-01&end=2025-12-31"
```

Expected Response:
```json
[
  {
    "channel": "Direct",
    "total_events": 150,
    "unique_users": 50,
    "total_visits": 75,
    "page_views": 120,
    "conversion_rate": 1.6
  },
  {
    "channel": "Organic",
    "total_events": 200,
    "unique_users": 80,
    "total_visits": 100,
    "page_views": 180,
    "conversion_rate": 1.8
  },
  ...
]
```

**With Filters**:
```bash
# Filter by project
curl "http://localhost:8080/api/channels?start=2025-01-01&end=2025-12-31&project=default"

# Filter by country
curl "http://localhost:8080/api/channels?start=2025-01-01&end=2025-12-31&country=US"

# Filter by bot traffic
curl "http://localhost:8080/api/channels?start=2025-01-01&end=2025-12-31&botFilter=human"
```

#### Batch Event Testing
```bash
curl -X POST http://localhost:8080/api/track/batch \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {
        "event_name": "page_view",
        "user_id": "batch-user-1",
        "url": "https://example.com/page1",
        "referrer": "https://www.google.com"
      },
      {
        "event_name": "page_view",
        "user_id": "batch-user-2",
        "url": "https://example.com/page2",
        "referrer": "https://www.facebook.com"
      },
      {
        "event_name": "page_view",
        "user_id": "batch-user-3",
        "url": "https://example.com/page3?utm_medium=cpc",
        "referrer": ""
      }
    ]
  }'
```

### 5. Frontend Testing

#### Dashboard Integration Tests

1. **Navigate to Channels Page**:
   - Visit: `http://localhost:3000/dashboard/channels` (or appropriate port)
   - Verify page loads without errors

2. **Test Chart Visualization**:
   - Verify pie chart displays correctly
   - Switch to bar chart view
   - Verify data matches API response

3. **Test Date Range Selection**:
   - Select "Today"
   - Select "Last 7 days"
   - Select "Last 30 days"
   - Select "Custom range" and pick dates
   - Verify data updates correctly

4. **Test Filters**:
   - Select a project from dropdown
   - Apply bot filter (Human/Bots)
   - Verify chart and table update

5. **Test Table Display**:
   - Verify all columns show correct data
   - Check percentage calculations
   - Verify conversion rates

6. **Test Auto-refresh**:
   - Enable auto-refresh
   - Set interval to 10 seconds
   - Verify data updates automatically

### 6. Performance Testing

#### Load Testing with k6 or Apache Bench

**Test 1: Single Event Tracking**
```bash
# Send 1000 requests
ab -n 1000 -c 10 -p event.json -T application/json http://localhost:8080/api/track
```

**Test 2: Batch Event Tracking**
```bash
# Send 100 batch requests with 10 events each
ab -n 100 -c 5 -p batch_events.json -T application/json http://localhost:8080/api/track/batch
```

**Test 3: Channel Analytics Query**
```bash
# Query channels endpoint
ab -n 500 -c 10 http://localhost:8080/api/channels?start=2025-01-01&end=2025-12-31
```

Expected Performance:
- Event tracking: < 50ms p95
- Batch tracking: < 100ms p95
- Channel query: < 200ms p95

#### Database Query Performance
```sql
-- Check query execution time
.timer on
SELECT channel, COUNT(*) as total_events
FROM events
WHERE timestamp BETWEEN '2025-01-01' AND '2025-12-31'
GROUP BY channel;

-- Check index usage
EXPLAIN SELECT channel, COUNT(*) as total_events
FROM events
WHERE timestamp BETWEEN '2025-01-01' AND '2025-12-31'
GROUP BY channel;
```

### 7. Integration Testing

#### End-to-End Test Scenario

1. **Setup**: Clean database
```bash
rm data/analytics.db
go run main.go # Starts server and runs migrations
```

2. **Track events** from all channels (Direct, Organic, Social, Referral, Paid)

3. **Query API** and verify correct counts

4. **View Dashboard** and verify charts match API data

5. **Apply filters** and verify filtered results

### 8. Edge Cases and Error Handling

#### Test Cases:

1. **Empty/Null Values**:
   - Event with no referrer
   - Event with null URL
   - Event with empty channel (should be "Unknown")

2. **Malformed URLs**:
   - Invalid URL format
   - URL without protocol
   - URL with special characters

3. **Large Datasets**:
   - Query with 1M+ events
   - Channel distribution with all 5 channels

4. **Concurrent Requests**:
   - Multiple simultaneous event tracking requests
   - Simultaneous queries from multiple users

## Automated Testing

### Running All Go Tests
```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### CI/CD Integration

Add to `.github/workflows/test.yml`:
```yaml
- name: Run Channel Detector Tests
  run: go test ./internal/channeldetector/... -v

- name: Verify Migrations
  run: |
    go run main.go &
    sleep 5
    # Check if channel column exists
    duckdb data/analytics.db "DESCRIBE events;" | grep channel
```

## Acceptance Criteria

✅ All unit tests pass
✅ Database migrations apply successfully
✅ Channel detection works correctly for all 5 channel types
✅ API endpoint returns correct data with filters
✅ Frontend displays charts and tables correctly
✅ Performance meets requirements (< 200ms for queries)
✅ No memory leaks or resource exhaustion
✅ Edge cases handled gracefully

## Known Limitations

1. **Domain Detection**: Currently requires exact domain match for "same domain" detection
2. **Channel Classification**: Priority order is fixed (Paid > Direct > Social > Organic > Referral)
3. **URL Parsing**: Malformed URLs may result in "Unknown" channel

## Future Enhancements

- Add support for custom channel definitions
- Implement channel attribution models (first-click, last-click, multi-touch)
- Add channel transition/flow analysis
- Support for campaign tracking parameters beyond UTM
