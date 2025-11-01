# Performance Tracking

Track Web Vitals and performance metrics.

## Enable Performance Tracking

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  enablePerformanceTracking: true
});
```

## Tracked Metrics

### Core Web Vitals

- **LCP** - Largest Contentful Paint
- **FID** - First Input Delay
- **CLS** - Cumulative Layout Shift

These metrics are automatically tracked when enabled.

## View Performance Data

Performance metrics appear in your analytics dashboard under the "Performance" section.

## Custom Performance Tracking

```javascript
// Track custom performance metrics
const startTime = performance.now();

// ... operation ...

const endTime = performance.now();
analytics.track('operation_performance', {
  operation: 'data_fetch',
  duration: endTime - startTime
});
```

## Next Steps

- [Overview →](/sdk/overview)
- [Configuration →](/sdk/configuration)
