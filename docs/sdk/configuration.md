# SDK Configuration

Complete configuration reference for Siraaj Analytics SDK.

## Configuration Options

```typescript
interface AnalyticsConfig {
  // Required
  apiUrl: string;              // Your Siraaj server URL
  projectId: string;           // Project identifier
  
  // Tracking
  autoTrack?: boolean;         // Auto-track page views, clicks, forms (default: true)
  
  // Debug
  debug?: boolean;             // Enable debug logging (default: false)
  
  // Performance
  bufferSize?: number;         // Events per batch (default: 10, max: 100)
  flushInterval?: number;      // Auto-flush interval in ms (default: 30000)
  timeout?: number;            // Request timeout in ms (default: 10000)
  
  // Reliability
  maxRetries?: number;         // Max retry attempts (default: 3)
  useBeacon?: boolean;         // Use sendBeacon API (default: true)
  maxQueueSize?: number;       // Max failed events to queue (default: 50)
  
  // Privacy
  respectDoNotTrack?: boolean; // Honor DNT header (default: true)
  sampling?: number;           // Sample rate 0.0-1.0 (default: 1.0)
  
  // Advanced
  enablePerformanceTracking?: boolean; // Track Web Vitals (default: false)
}
```

## Examples

### Development

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app-dev',
  autoTrack: true,
  debug: true,
  bufferSize: 5,
  flushInterval: 5000
});
```

### Production

```javascript
const analytics = new Analytics({
  apiUrl: 'https://analytics.example.com',
  projectId: 'my-app',
  autoTrack: true,
  debug: false,
  bufferSize: 20,
  flushInterval: 30000,
  respectDoNotTrack: true,
  enablePerformanceTracking: true
});
```

## Next Steps

- [Auto-Tracking →](/sdk/auto-tracking)
- [Custom Events →](/sdk/custom-events)
