# API Overview

RESTful API reference for Siraaj Analytics.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Currently, Siraaj doesn't require authentication. API key authentication is coming soon.

## Endpoints

### Health Check

Check server status.

```http
GET /api/health
```

**Response**

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": 3600
}
```

### Track Event

Send a single event.

```http
POST /api/track
Content-Type: application/json
```

**Request Body**

```json
{
  "event_name": "page_view",
  "user_id": "user-123",
  "session_id": "session-456",
  "url": "https://example.com/products",
  "referrer": "https://google.com",
  "user_agent": "Mozilla/5.0...",
  "timestamp": "2024-01-15T10:30:00Z",
  "browser": "Chrome",
  "os": "MacOS",
  "device": "Desktop",
  "project_id": "my-website",
  "channel": "Organic",
  "custom_property": "value"
}
```

**Response**

```json
{
  "success": true,
  "event_id": "evt-789"
}
```

### Track Batch

Send multiple events at once (recommended).

```http
POST /api/track/batch
Content-Type: application/json
```

**Request Body**

```json
{
  "events": [
    {
      "event_name": "page_view",
      "user_id": "user-123",
      ...
    },
    {
      "event_name": "button_clicked",
      "user_id": "user-123",
      ...
    }
  ]
}
```

**Response**

```json
{
  "success": true,
  "events_received": 2
}
```

### Query Analytics

Get analytics data (Coming Soon).

```http
GET /api/analytics
```

**Query Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| project_id | string | Project identifier |
| start_date | string | Start date (ISO 8601) |
| end_date | string | End date (ISO 8601) |
| metric | string | Metric to query |

**Example**

```bash
curl "http://localhost:8080/api/analytics?\
project_id=my-website&\
start_date=2024-01-01T00:00:00Z&\
end_date=2024-01-31T23:59:59Z&\
metric=page_views"
```

**Response**

```json
{
  "metric": "page_views",
  "value": 10543,
  "change": 15.3,
  "data": [
    {
      "date": "2024-01-01",
      "value": 340
    },
    ...
  ]
}
```

## Error Responses

### 400 Bad Request

```json
{
  "error": "Invalid request body",
  "details": "Missing required field: event_name"
}
```

### 429 Too Many Requests

```json
{
  "error": "Rate limit exceeded",
  "retry_after": 60
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal server error",
  "request_id": "req-123"
}
```

## Rate Limits

Default rate limits:

- **Track endpoint**: 1000 requests/minute per IP
- **Batch endpoint**: 100 requests/minute per IP
- **Query endpoint**: 60 requests/minute per IP

## CORS

Configure allowed origins:

```bash
CORS=https://example.com,https://app.example.com ./siraaj
```

## Examples

### cURL

```bash
# Track single event
curl -X POST http://localhost:8080/api/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_name": "page_view",
    "user_id": "user-123",
    "project_id": "my-website",
    "url": "https://example.com",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
  }'
```

### JavaScript

```javascript
// Track event
fetch('http://localhost:8080/api/track', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    event_name: 'button_clicked',
    user_id: 'user-123',
    project_id: 'my-app',
    timestamp: new Date().toISOString()
  })
});
```

### Python

```python
import requests
from datetime import datetime

response = requests.post('http://localhost:8080/api/track', json={
    'event_name': 'purchase_completed',
    'user_id': 'user-123',
    'project_id': 'my-store',
    'timestamp': datetime.utcnow().isoformat() + 'Z',
    'amount': 99.99,
    'currency': 'USD'
})

print(response.json())
```

## Next Steps

- [Track Events Guide →](/api/track-events)
- [Query Analytics →](/api/query-analytics)
- [Error Handling →](/api/error-handling)
