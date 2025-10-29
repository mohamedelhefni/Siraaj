# Channel Analytics Feature

## Overview

The Channel Analytics feature automatically classifies incoming events based on their traffic source, providing insights into where your users are coming from. This feature is essential for understanding acquisition channels and optimizing marketing efforts.

## Channel Types

The system classifies traffic into five distinct channels:

### 1. 🔵 Direct
Traffic that comes directly to your site without a referrer.
- **Examples**: 
  - User types URL directly in browser
  - Clicks from bookmark
  - Email clients (some)
  - No referrer information
- **Detection Logic**: Empty or null referrer, or referrer matches current domain

### 2. 🟢 Organic
Traffic from organic search engine results.
- **Examples**:
  - Google search results
  - Bing search results
  - Yahoo search results
- **Supported Search Engines**: Google, Bing, Yahoo, DuckDuckGo, Baidu, Yandex, Ask, AOL, Ecosia
- **Detection Logic**: Referrer domain matches known search engines

### 3. 🟣 Social
Traffic from social media platforms.
- **Examples**:
  - Facebook posts/ads
  - Twitter/X links
  - LinkedIn shares
  - Instagram bio links
- **Supported Platforms**: Facebook, Twitter, LinkedIn, Instagram, Pinterest, Reddit, TikTok, YouTube, Snapchat, Tumblr
- **Detection Logic**: Referrer domain matches known social platforms

### 4. 🟠 Referral
Traffic from other websites that link to your content.
- **Examples**:
  - Blog mentions
  - News articles
  - Partner websites
  - Forums and communities
- **Detection Logic**: Any referrer that doesn't match Direct, Organic, Social, or Paid criteria

### 5. 🔴 Paid
Traffic from paid advertising campaigns.
- **Examples**:
  - Google Ads (PPC)
  - Facebook Ads
  - Microsoft Advertising
  - Display ads
- **Detection Logic**: URL contains:
  - UTM parameters: `utm_medium=cpc`, `utm_medium=ppc`, `utm_medium=paid`, `utm_medium=display`
  - Click IDs: `gclid` (Google), `fbclid` (Facebook), `msclkid` (Microsoft), `twclid` (Twitter)
- **Priority**: Paid detection takes precedence over all other channels

## Technical Implementation

### Database Schema

```sql
-- New column added to events table
ALTER TABLE events ADD COLUMN channel VARCHAR;

-- Indexes for performance
CREATE INDEX idx_channel ON events(channel);
CREATE INDEX idx_timestamp_channel ON events(timestamp DESC, channel);
```

### Channel Detection Algorithm

```go
// Priority order:
1. Check for Paid parameters (highest priority)
2. Check for empty/null referrer → Direct
3. Check if referrer is same domain → Direct
4. Check if referrer is social media → Social
5. Check if referrer is search engine → Organic
6. Default to Referral
```

### API Endpoint

**GET** `/api/channels`

Query Parameters:
- `start` (string): Start date in YYYY-MM-DD format
- `end` (string): End date in YYYY-MM-DD format
- `project` (string, optional): Filter by project ID
- `country` (string, optional): Filter by country
- `browser` (string, optional): Filter by browser
- `device` (string, optional): Filter by device
- `os` (string, optional): Filter by operating system
- `botFilter` (string, optional): Filter bot traffic ("human", "bot")

Response:
```json
[
  {
    "channel": "Direct",
    "total_events": 1500,
    "unique_users": 450,
    "total_visits": 600,
    "page_views": 1200,
    "conversion_rate": 2.0
  },
  {
    "channel": "Organic",
    "total_events": 2300,
    "unique_users": 890,
    "total_visits": 1100,
    "page_views": 2100,
    "conversion_rate": 1.91
  }
  // ... more channels
]
```

## Usage Examples

### Frontend Integration

```javascript
import { fetchChannels } from '$lib/api';

// Fetch channel data
const channels = await fetchChannels('2025-01-01', '2025-01-31', {
  project: 'my-project',
  botFilter: 'human'
});

// Display in chart
<ChannelsChart data={channels} chartType="pie" />
```

### Direct API Calls

```bash
# Get all channels for last 30 days
curl "http://localhost:8080/api/channels?start=2025-01-01&end=2025-01-31"

# Get channels for specific project, human traffic only
curl "http://localhost:8080/api/channels?start=2025-01-01&end=2025-01-31&project=website&botFilter=human"
```

### Event Tracking with Channel Detection

Channels are automatically detected when you track events:

```javascript
// Using the SDK
analytics.track('page_view', {
  url: 'https://example.com/products',
  referrer: 'https://www.google.com/search?q=products'
  // Channel will be automatically detected as "Organic"
});
```

```bash
# Direct API call
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "user-123",
    "url": "https://example.com/landing?utm_medium=cpc",
    "referrer": "https://www.google.com"
  }'
# Channel will be detected as "Paid" (utm_medium=cpc)
```

## Dashboard Features

### Channel Overview Page

Located at `/dashboard/channels`, the page provides:

1. **Summary Cards**
   - Total Events
   - Unique Users
   - Total Visits
   - Page Views

2. **Visualization**
   - Pie Chart: Shows percentage distribution
   - Bar Chart: Shows absolute counts
   - Toggle between chart types

3. **Filters**
   - Date range selection (Today, Last 7 days, Last 30 days, Custom)
   - Project filter
   - Bot traffic filter
   - Real-time data refresh

4. **Detailed Table**
   - Channel name with color coding
   - Events, Users, Visits, Page Views
   - Views per Visit ratio
   - Percentage of total traffic

### Color Coding

- 🔵 **Direct**: Blue (`#3b82f6`)
- 🟢 **Organic**: Green (`#10b981`)
- 🟣 **Social**: Purple (`#8b5cf6`)
- 🟠 **Referral**: Orange (`#f59e0b`)
- 🔴 **Paid**: Red (`#ef4444`)
- ⚫ **Unknown**: Gray (`#6b7280`)

## Configuration

### Customizing Channel Detection

To add more search engines or social platforms, edit:

```go
// internal/channeldetector/channeldetector.go

var organicSearchEngines = []string{
  "google.com",
  "bing.com",
  // Add more...
}

var socialPlatforms = []string{
  "facebook.com",
  "twitter.com",
  // Add more...
}
```

### Adding Custom Paid Parameters

```go
var paidParameters = []string{
  "utm_medium=cpc",
  "utm_medium=ppc",
  "gclid=",
  // Add custom parameters...
  "mycustom_param=",
}
```

## Performance Considerations

- **Indexing**: Composite indexes on `(timestamp, channel)` for fast queries
- **Caching**: Channel detection logic is stateless and can be cached
- **Batch Processing**: Batch event tracking processes channels efficiently
- **Query Optimization**: Aggregations use DuckDB's columnar storage for speed

### Typical Performance

- Channel detection: < 1ms per event
- Channel query (30 days): < 100ms
- Dashboard load: < 500ms

## Migration Guide

### Existing Data

For events tracked before this feature was added:

1. **Channel column will be NULL** for old events
2. **Display as "Unknown"** in analytics
3. **Optional**: Backfill channels for historical data

```sql
-- Backfill example (run carefully on production)
UPDATE events 
SET channel = CASE
  WHEN referrer = '' OR referrer IS NULL THEN 'Direct'
  WHEN referrer LIKE '%google.com%' 
    OR referrer LIKE '%bing.com%' 
    OR referrer LIKE '%yahoo.com%' THEN 'Organic'
  WHEN referrer LIKE '%facebook.com%' 
    OR referrer LIKE '%twitter.com%' 
    OR referrer LIKE '%linkedin.com%' THEN 'Social'
  WHEN url LIKE '%gclid=%' 
    OR url LIKE '%utm_medium=cpc%' 
    OR url LIKE '%utm_medium=ppc%' THEN 'Paid'
  ELSE 'Referral'
END
WHERE channel IS NULL;
```

## Analytics Insights

### Key Metrics to Monitor

1. **Channel Distribution**: Which channels drive the most traffic?
2. **Conversion by Channel**: Which channels have highest conversion rates?
3. **User Quality**: Do certain channels bring more engaged users?
4. **Cost Efficiency**: For Paid channels, track ROI

### Sample Queries

```sql
-- Channel distribution for last 30 days
SELECT 
  channel,
  COUNT(*) as events,
  COUNT(DISTINCT user_id) as users,
  COUNT(DISTINCT session_id) as sessions
FROM events
WHERE timestamp >= CURRENT_DATE - INTERVAL 30 DAY
GROUP BY channel
ORDER BY events DESC;

-- Channel performance by day
SELECT 
  DATE_TRUNC('day', timestamp) as date,
  channel,
  COUNT(*) as events
FROM events
WHERE timestamp >= CURRENT_DATE - INTERVAL 7 DAY
GROUP BY date, channel
ORDER BY date, events DESC;

-- Top converting channels
SELECT 
  channel,
  COUNT(DISTINCT user_id) as users,
  COUNT(CASE WHEN event_name = 'purchase' THEN 1 END) as conversions,
  ROUND(COUNT(CASE WHEN event_name = 'purchase' THEN 1 END)::FLOAT / COUNT(DISTINCT user_id) * 100, 2) as conversion_rate
FROM events
WHERE timestamp >= CURRENT_DATE - INTERVAL 30 DAY
GROUP BY channel
ORDER BY conversion_rate DESC;
```

## Troubleshooting

### Issue: All events showing as "Unknown"

**Cause**: Referrer or URL not being captured
**Solution**: Ensure SDK is sending referrer data:
```javascript
analytics.track('page_view', {
  url: window.location.href,
  referrer: document.referrer
});
```

### Issue: Paid traffic showing as Organic

**Cause**: UTM parameters not present or misspelled
**Solution**: Verify campaign URLs include `utm_medium=cpc` or `utm_medium=ppc`

### Issue: Internal navigation showing as Referral

**Cause**: Subdomain differences or protocol mismatch
**Solution**: Ensure `currentDomain` parameter is set correctly in channel detection

## Best Practices

1. **Always include referrer data** in event tracking
2. **Use UTM parameters** consistently for all paid campaigns
3. **Tag social media posts** with `utm_source=social_platform`
4. **Monitor "Unknown" channel** for classification issues
5. **Review channel definitions** regularly as marketing channels evolve

## API Client Examples

### JavaScript/TypeScript
```typescript
const response = await fetch('/api/channels?start=2025-01-01&end=2025-01-31');
const channels = await response.json();

channels.forEach(channel => {
  console.log(`${channel.channel}: ${channel.total_events} events`);
});
```

### Python
```python
import requests

response = requests.get('http://localhost:8080/api/channels', params={
    'start': '2025-01-01',
    'end': '2025-01-31',
    'project': 'website'
})

channels = response.json()
for channel in channels:
    print(f"{channel['channel']}: {channel['total_events']} events")
```

### Go
```go
resp, err := http.Get("http://localhost:8080/api/channels?start=2025-01-01&end=2025-01-31")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

var channels []map[string]interface{}
json.NewDecoder(resp.Body).Decode(&channels)
```

## Contributing

To extend channel detection:

1. **Add new search engine/social platform**: Update arrays in `channeldetector.go`
2. **Add new paid parameter**: Update `paidParameters` array
3. **Test thoroughly**: Add test cases in `channeldetector_test.go`
4. **Update documentation**: Reflect changes in this file

## License

See LICENSE file in the repository root.
