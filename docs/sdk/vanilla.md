# Vanilla JavaScript

Use Siraaj Analytics with pure JavaScript - no framework required.

## CDN Installation

The easiest way to get started:

```html
<!DOCTYPE html>
<html>
<head>
  <title>My Website</title>
  
  <!-- Include Siraaj SDK -->
  <script src="http://localhost:8080/sdk/analytics.js"></script>
  
  <!-- Initialize -->
  <script>
    const analytics = new Analytics({
      apiUrl: 'http://localhost:8080',
      projectId: 'my-website',
      autoTrack: true,
      debug: false
    });
  </script>
</head>
<body>
  <h1>Welcome to My Website</h1>
  <button id="signup-btn">Sign Up</button>
  
  <script>
    // Track button clicks
    document.getElementById('signup-btn').addEventListener('click', function() {
      analytics.track('signup_clicked', {
        location: 'homepage',
        button_text: 'Sign Up'
      });
    });
  </script>
</body>
</html>
```

## NPM/Module Installation

```bash
npm install @hefni101/siraaj
```

### ES Modules

```javascript
import { analytics } from '@hefni101/siraaj';

analytics.init({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  autoTrack: true
});

// Track events
analytics.track('page_viewed', {
  path: window.location.pathname
});
```

### CommonJS

```javascript
const { analytics } = require('@hefni101/siraaj');

analytics.init({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app'
});
```

## Configuration

```javascript
const analytics = new Analytics({
  // Required
  apiUrl: 'http://localhost:8080',
  projectId: 'my-website',
  
  // Optional
  autoTrack: true,           // Auto-track page views, clicks, forms
  debug: false,              // Enable console logging
  bufferSize: 10,            // Events to buffer before sending
  flushInterval: 30000,      // Auto-flush interval (30 seconds)
  timeout: 10000,            // Request timeout (10 seconds)
  maxRetries: 3,             // Max retry attempts
  useBeacon: true,           // Use sendBeacon for reliability
  sampling: 1.0,             // Track 100% of users
  respectDoNotTrack: true,   // Honor DNT header
  enablePerformanceTracking: false // Track Web Vitals
});
```

## Tracking Events

### Page Views

```javascript
// Auto-tracked if autoTrack is enabled
// Manual tracking:
analytics.pageView();

// With custom URL
analytics.pageView('/custom-path');

// With properties
analytics.pageView('/products', {
  category: 'electronics',
  featured: true
});
```

### Custom Events

```javascript
analytics.track('button_clicked', {
  button_id: 'cta-signup',
  location: 'hero',
  text: 'Get Started'
});

analytics.track('video_played', {
  video_id: 'intro-video',
  duration: 120,
  quality: '1080p'
});

analytics.track('search_performed', {
  query: 'javascript analytics',
  results_count: 42
});
```

### Form Tracking

```javascript
// Auto-tracked if autoTrack is enabled
// Manual tracking:
const form = document.getElementById('contact-form');

form.addEventListener('submit', function(e) {
  analytics.trackForm('contact-form', {
    form_type: 'contact',
    source: 'footer'
  });
});
```

### Click Tracking

```javascript
// Auto-tracks all links if autoTrack is enabled
// Manual tracking:
document.getElementById('download-btn').addEventListener('click', function() {
  analytics.trackClick('download-btn', {
    file: 'whitepaper.pdf',
    size: '2.5MB'
  });
});
```

### Error Tracking

```javascript
// Auto-tracked if autoTrack is enabled
// Manual tracking:
try {
  riskyOperation();
} catch (error) {
  analytics.trackError(error, {
    context: 'checkout_process',
    step: 'payment'
  });
}

// Track custom errors
analytics.trackError('Payment gateway timeout', {
  gateway: 'stripe',
  attempt: 3
});
```

## User Identification

```javascript
// On login
analytics.identify('user-123', {
  email: 'user@example.com',
  name: 'John Doe',
  plan: 'premium',
  signup_date: '2024-01-15'
});

// Update user properties
analytics.setUserProperties({
  plan: 'enterprise',
  mrr: 299
});

// On logout
analytics.reset();
```

## Advanced Features

### Manual Flushing

```javascript
// Flush events immediately
await analytics.flush();

// Flush before page unload (auto-handled)
window.addEventListener('beforeunload', function() {
  analytics.flush(true); // Use sendBeacon
});
```

### Performance Tracking

Enable Web Vitals tracking:

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  enablePerformanceTracking: true
});

// Automatically tracks:
// - LCP (Largest Contentful Paint)
// - FID (First Input Delay)
// - CLS (Cumulative Layout Shift)
```

### Sampling

Track only a percentage of users:

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  sampling: 0.5 // Track 50% of users
});
```

### Debug Mode

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  debug: true // Logs all events to console
});

// Console output:
// [Siraaj] Initialized {apiUrl: '...', projectId: '...'}
// [Siraaj] Event: page_view
// [Siraaj] Flushing 5 events
```

## Single Page Applications

### Manual Page Tracking

```javascript
// Disable auto-tracking
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-spa',
  autoTrack: false
});

// Track page changes manually
function navigateTo(path) {
  // Update URL
  history.pushState({}, '', path);
  
  // Track page view
  analytics.pageView(path);
  
  // Render content
  renderPage(path);
}

// Listen to popstate
window.addEventListener('popstate', function() {
  analytics.pageView(window.location.pathname);
});
```

### Hash-based Routing

```javascript
window.addEventListener('hashchange', function() {
  analytics.pageView(window.location.hash);
});

// Track on load
analytics.pageView(window.location.hash);
```

## E-commerce Tracking

```javascript
// Product viewed
analytics.track('product_viewed', {
  product_id: 'prod-123',
  product_name: 'Wireless Headphones',
  price: 99.99,
  currency: 'USD',
  category: 'Electronics'
});

// Add to cart
analytics.track('add_to_cart', {
  product_id: 'prod-123',
  quantity: 1,
  price: 99.99
});

// Purchase completed
analytics.track('purchase_completed', {
  order_id: 'order-456',
  total: 109.99,
  currency: 'USD',
  items: [
    {
      product_id: 'prod-123',
      quantity: 1,
      price: 99.99
    }
  ],
  shipping: 10.00,
  tax: 0
});
```

## A/B Testing

```javascript
// Track which variant user sees
const variant = Math.random() < 0.5 ? 'A' : 'B';

analytics.track('ab_test_assigned', {
  test_name: 'homepage_hero',
  variant: variant
});

// Show variant
showHeroVariant(variant);

// Track conversion
document.getElementById('cta-btn').addEventListener('click', function() {
  analytics.track('ab_test_conversion', {
    test_name: 'homepage_hero',
    variant: variant
  });
});
```

## Examples

### Basic Website

```html
<!DOCTYPE html>
<html>
<head>
  <title>My Blog</title>
  <script src="http://localhost:8080/sdk/analytics.js"></script>
  <script>
    const analytics = new Analytics({
      apiUrl: 'http://localhost:8080',
      projectId: 'my-blog',
      autoTrack: true
    });
  </script>
</head>
<body>
  <article>
    <h1>Blog Post Title</h1>
    <button id="share-btn">Share</button>
  </article>
  
  <script>
    // Track share clicks
    document.getElementById('share-btn').addEventListener('click', function() {
      analytics.track('article_shared', {
        title: document.querySelector('h1').textContent,
        platform: 'twitter'
      });
    });
    
    // Track reading time
    let startTime = Date.now();
    window.addEventListener('beforeunload', function() {
      const readingTime = Math.floor((Date.now() - startTime) / 1000);
      analytics.track('article_read', {
        title: document.querySelector('h1').textContent,
        reading_time: readingTime
      });
    });
  </script>
</body>
</html>
```

### Landing Page

```html
<!DOCTYPE html>
<html>
<head>
  <title>Product Launch</title>
  <script src="http://localhost:8080/sdk/analytics.js"></script>
  <script>
    const analytics = new Analytics({
      apiUrl: 'http://localhost:8080',
      projectId: 'product-launch',
      autoTrack: true,
      enablePerformanceTracking: true
    });
  </script>
</head>
<body>
  <section id="hero">
    <h1>Amazing Product</h1>
    <button id="cta">Get Early Access</button>
  </section>
  
  <form id="signup-form">
    <input type="email" name="email" required>
    <button type="submit">Sign Up</button>
  </form>
  
  <script>
    // Track CTA clicks
    document.getElementById('cta').addEventListener('click', function() {
      analytics.track('cta_clicked', {
        location: 'hero',
        text: this.textContent
      });
      
      // Scroll to form
      document.getElementById('signup-form').scrollIntoView({
        behavior: 'smooth'
      });
    });
    
    // Track form submissions
    document.getElementById('signup-form').addEventListener('submit', function(e) {
      e.preventDefault();
      
      const email = this.email.value;
      
      analytics.track('signup_submitted', {
        source: 'landing_page'
      });
      
      // Submit to backend
      fetch('/api/signup', {
        method: 'POST',
        body: JSON.stringify({ email })
      });
    });
  </script>
</body>
</html>
```

## TypeScript

```typescript
import { analytics } from '@hefni101/siraaj';
import type { AnalyticsConfig } from '@hefni101/siraaj';

const config: AnalyticsConfig = {
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  autoTrack: true,
  debug: false
};

analytics.init(config);

// Type-safe event tracking
analytics.track('button_clicked', {
  button_id: 'signup',
  location: 'hero'
});
```

## Best Practices

### 1. Initialize Early

```javascript
// ✅ Good - in <head>
<head>
  <script src="analytics.js"></script>
  <script>
    const analytics = new Analytics({...});
  </script>
</head>

// ❌ Bad - at end of <body>
```

### 2. Use Descriptive Event Names

```javascript
// ✅ Good
analytics.track('checkout_completed');
analytics.track('video_played');

// ❌ Bad
analytics.track('event1');
analytics.track('click');
```

### 3. Add Context

```javascript
// ✅ Good
analytics.track('search_performed', {
  query: 'javascript',
  results_count: 42,
  filters_applied: ['tutorial', 'beginner']
});

// ❌ Bad
analytics.track('search_performed');
```

### 4. Handle Errors

```javascript
// Wrap in try-catch for critical operations
try {
  analytics.track('payment_initiated', {
    amount: 99.99
  });
} catch (error) {
  console.error('Analytics error:', error);
  // Continue with payment anyway
}
```

## Browser Support

- ✅ Chrome 90+
- ✅ Firefox 88+
- ✅ Safari 14+
- ✅ Edge 90+
- ✅ Opera 76+

## Next Steps

- [Configuration Options →](/sdk/configuration)
- [Custom Events →](/sdk/custom-events)
- [User Identification →](/sdk/user-identification)
