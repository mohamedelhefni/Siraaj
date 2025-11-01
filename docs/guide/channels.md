# Channel Analytics

Understand where your traffic comes from with automatic channel classification.

## What are Channels?

Channels categorize your traffic based on the source and how users arrived at your website. Siraaj automatically detects and classifies traffic into six channels:

- **🔵 Direct** - Users who typed your URL or used a bookmark
- **🟢 Organic** - Traffic from search engines (Google, Bing, etc.)
- **🟣 Social** - Visitors from social media platforms
- **🟠 Referral** - Traffic from other websites
- **🔴 Paid** - Paid advertising campaigns
- **⚫ Unknown** - Unclassified traffic

## How Channel Detection Works

Channels are automatically detected server-side based on:

1. **Referrer URL** - Where the user came from
2. **URL Parameters** - UTM parameters and click IDs
3. **Domain matching** - Against known search engines and social platforms

### Detection Priority

The detection follows this priority order:

1. **Paid** (highest priority) - Checked first
2. **Direct** - No referrer or same domain
3. **Social** - Known social media platforms
4. **Organic** - Known search engines
5. **Referral** - All other external sources
6. **Unknown** - When classification fails

## Channel Definitions

### Direct Traffic 🔵

**Detected when:**
- No referrer information
- Referrer is empty or null
- Referrer matches your own domain

**Examples:**
- User types URL directly
- Click from bookmark
- Email client links (some)
- QR codes
- Mobile apps

### Organic Search 🟢

**Detected when referrer is from:**
- Google (`google.com`)
- Bing (`bing.com`)
- Yahoo (`yahoo.com`)
- DuckDuckGo (`duckduckgo.com`)
- Baidu (`baidu.com`)
- Yandex (`yandex.ru`)
- Ask.com (`ask.com`)
- AOL (`aol.com`)
- Ecosia (`ecosia.org`)

### Social Media 🟣

**Detected when referrer is from:**
- Facebook (`facebook.com`)
- Twitter/X (`twitter.com`, `x.com`)
- LinkedIn (`linkedin.com`)
- Instagram (`instagram.com`)
- Pinterest (`pinterest.com`)
- Reddit (`reddit.com`)
- TikTok (`tiktok.com`)
- YouTube (`youtube.com`)
- Snapchat (`snapchat.com`)
- Tumblr (`tumblr.com`)

### Paid Advertising 🔴

**Detected when URL contains:**

**UTM Parameters:**
- `utm_medium=cpc`
- `utm_medium=ppc`
- `utm_medium=paid`
- `utm_medium=display`

**Click IDs:**
- `gclid=` (Google Ads)
- `fbclid=` (Facebook Ads)
- `msclkid=` (Microsoft Advertising)
- `twclid=` (Twitter Ads)

**Example:**
```
https://example.com/landing?utm_source=google&utm_medium=cpc&utm_campaign=summer_sale
```

### Referral Traffic 🟠

**Detected when:**
- Traffic comes from an external website
- Referrer doesn't match Direct, Organic, Social, or Paid criteria

**Examples:**
- Blog mentions
- News articles
- Partner websites
- Forum links
- Hacker News, Product Hunt, etc.

### Unknown ⚫

**When classification fails:**
- Malformed referrer URLs
- Missing or corrupted data
- Very old events before channel detection was added

## API Endpoint

```http
GET /api/channels?start=2024-01-01&end=2024-01-31
```

### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| start | string | Start date (YYYY-MM-DD) |
| end | string | End date (YYYY-MM-DD) |
| project | string | Filter by project ID |
| country | string | Filter by country |
| browser | string | Filter by browser |
| device | string | Filter by device type |
| os | string | Filter by operating system |
| botFilter | string | Filter bot traffic (human/bot) |

### Response

```json
[
  {
    "channel": "Organic",
    "total_events": 2300,
    "unique_users": 890,
    "total_visits": 1100,
    "page_views": 2100,
    "conversion_rate": 1.91
  },
  {
    "channel": "Direct",
    "total_events": 1500,
    "unique_users": 450,
    "total_visits": 600,
    "page_views": 1200,
    "conversion_rate": 2.0
  },
  {
    "channel": "Social",
    "total_events": 800,
    "unique_users": 320,
    "total_visits": 400,
    "page_views": 750,
    "conversion_rate": 1.88
  }
]
```

## Fetching Channel Data

### JavaScript

```javascript
async function getChannels() {
  const response = await fetch(
    'http://localhost:8080/api/channels?start=2024-01-01&end=2024-01-31&project=my-website&botFilter=human'
  );
  
  const channels = await response.json();
  
  channels.forEach(channel => {
    console.log(`${channel.channel}: ${channel.total_events} events, ${channel.unique_users} users`);
  });
}
```

### cURL

```bash
curl "http://localhost:8080/api/channels?start=2024-01-01&end=2024-01-31&project=my-website"
```

### Python

```python
import requests

response = requests.get('http://localhost:8080/api/channels', params={
    'start': '2024-01-01',
    'end': '2024-01-31',
    'project': 'my-website',
    'botFilter': 'human'
})

channels = response.json()
for channel in channels:
    print(f"{channel['channel']}: {channel['total_events']} events")
```

## Tagging Campaigns for Paid Traffic

To ensure your paid campaigns are properly tracked, use UTM parameters:

### Basic Campaign Tagging

```
https://example.com/?utm_source=google&utm_medium=cpc&utm_campaign=summer_sale
```

### Campaign URL Builder

```javascript
function buildCampaignUrl(baseUrl, params) {
  const url = new URL(baseUrl);
  
  if (params.source) url.searchParams.set('utm_source', params.source);
  if (params.medium) url.searchParams.set('utm_medium', params.medium);
  if (params.campaign) url.searchParams.set('utm_campaign', params.campaign);
  if (params.content) url.searchParams.set('utm_content', params.content);
  if (params.term) url.searchParams.set('utm_term', params.term);
  
  return url.toString();
}

// Example usage
const campaignUrl = buildCampaignUrl('https://example.com/landing', {
  source: 'facebook',
  medium: 'cpc',
  campaign: 'spring_promo',
  content: 'ad_variant_a'
});
```

### UTM Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| utm_source | Where traffic comes from | google, facebook, newsletter |
| utm_medium | Marketing medium | cpc, email, social |
| utm_campaign | Campaign name | summer_sale, product_launch |
| utm_content | Ad variant | banner_a, text_ad |
| utm_term | Keyword (for paid search) | running_shoes |

## Analyzing Channel Performance

### Key Metrics

**Total Events**
- All events from this channel
- Includes page views and custom events

**Unique Users**
- Number of distinct users
- Counted by user ID

**Total Visits**
- Number of sessions
- Counted by session ID

**Page Views**
- Only page_view events
- Subset of total events

**Conversion Rate**
- Views per visit ratio
- Higher = more engaged users

### Comparing Channels

```javascript
async function compareChannels() {
  const response = await fetch('http://localhost:8080/api/channels?start=2024-01-01&end=2024-01-31');
  const channels = await response.json();
  
  // Sort by total events
  channels.sort((a, b) => b.total_events - a.total_events);
  
  console.log('Top Channels by Traffic:');
  channels.forEach((channel, index) => {
    const percentage = (channel.total_events / channels.reduce((sum, c) => sum + c.total_events, 0) * 100).toFixed(1);
    console.log(`${index + 1}. ${channel.channel}: ${percentage}%`);
  });
  
  // Find best conversion rate
  const bestChannel = channels.reduce((best, current) => 
    current.conversion_rate > best.conversion_rate ? current : best
  );
  
  console.log(`\nBest converting channel: ${bestChannel.channel} (${bestChannel.conversion_rate.toFixed(2)})`);
}
```

## Dashboard Integration

The channel data appears in the Siraaj dashboard at `/dashboard/channels`.

**Features:**
- Pie chart visualization
- Bar chart view
- Detailed table with all metrics
- Date range selection
- Project filtering
- Real-time updates

## Use Cases

### Marketing Attribution

Understand which channels drive the most valuable traffic:

```javascript
// Fetch channels with conversion data
const channels = await fetch('/api/channels?start=2024-01-01&end=2024-01-31')
  .then(r => r.json());

// Calculate ROI if you track revenue
channels.forEach(channel => {
  const revenue = calculateRevenue(channel); // Your function
  const cost = getCost(channel); // Your function
  channel.roi = (revenue - cost) / cost * 100;
});
```

### Content Strategy

Identify where your content performs best:

```javascript
// Organic: Focus on SEO
// Social: Create shareable content
// Direct: Build brand loyalty
// Referral: Build partnerships
```

### Budget Allocation

Compare paid vs organic performance:

```javascript
const paid = channels.find(c => c.channel === 'Paid');
const organic = channels.find(c => c.channel === 'Organic');

console.log('Paid Cost Per User:', paidBudget / paid.unique_users);
console.log('Organic Cost Per User: $0');
```

## Troubleshooting

### All Traffic Shows as Direct

**Possible causes:**
- Referrer not being sent
- HTTPS → HTTP downgrade strips referrer
- SDK not capturing referrer

**Solution:**
```javascript
// Ensure SDK captures referrer
analytics.track('page_view', {
  url: window.location.href,
  referrer: document.referrer  // Explicitly send referrer
});
```

### Paid Traffic Shows as Organic

**Cause:** Missing or incorrect UTM parameters

**Solution:** Always include `utm_medium=cpc` or `utm_medium=ppc`:
```
https://example.com/landing?utm_source=google&utm_medium=cpc
```

### Social Traffic Shows as Referral

**Cause:** Social platform not in detection list

**Solution:** Check `internal/channeldetector/channeldetector.go` and add missing platforms.

## Best Practices

### 1. Consistent UTM Tagging

Always use lowercase and consistent naming:
- ✅ `utm_medium=cpc`
- ❌ `utm_medium=CPC` or `utm_medium=cost-per-click`

### 2. Track All Campaigns

Tag every paid campaign URL, even for testing.

### 3. Monitor Unknown Channel

If "Unknown" channel has significant traffic, investigate and fix classification issues.

### 4. Segment Analysis

Compare channels across:
- Time periods
- Geographic locations
- Device types
- User segments

### 5. Attribution Windows

Consider that users may visit from multiple channels before converting.

## Next Steps

- [Funnels →](/guide/funnels)
- [Custom Events →](/sdk/custom-events)
- [API Reference →](/api/overview)
