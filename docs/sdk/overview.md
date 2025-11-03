# SDK Overview

Siraaj provides a comprehensive JavaScript SDK with first-class support for multiple frameworks and environments.

## Installation

### CDN (Browser)

The simplest way to get started:

```html
<script src="http://your-server:8080/sdk/analytics.js"></script>
<script>
  const analytics = new Analytics({
    apiUrl: 'http://your-server:8080',
    projectId: 'my-website'
  });
</script>
```

### NPM Package (Coming Soon)

```bash
npm install @hefni101/siraaj
# or
pnpm add @hefni101/siraaj
# or
yarn add @hefni101/siraaj
```

## Framework Support

Siraaj SDK works everywhere JavaScript runs:

### Vanilla JavaScript
```javascript
import { analytics } from '@hefni101/siraaj';

analytics.init({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app'
});
```

[Learn more →](/sdk/vanilla)

### React
```jsx
import { AnalyticsProvider, useAnalytics } from '@hefni101/siraaj/react';

function App() {
  return (
    <AnalyticsProvider config={{ apiUrl: '...', projectId: '...' }}>
      <YourApp />
    </AnalyticsProvider>
  );
}
```

[Learn more →](/sdk/react)

### Vue 3
```vue
<script setup>
import { useAnalytics } from '@hefni101/siraaj/vue';

const { track } = useAnalytics();
</script>
```

[Learn more →](/sdk/vue)

### Svelte
```svelte
<script>
import { createAnalytics } from '@hefni101/siraaj/svelte';

const { track } = createAnalytics();
</script>
```

[Learn more →](/sdk/svelte)

### Next.js
```tsx
import { AnalyticsProvider } from '@hefni101/siraaj/next';
```

[Learn more →](/sdk/nextjs)

### Nuxt 3
```typescript
import { useNuxtAnalytics } from '@hefni101/siraaj/nuxt';
```

[Learn more →](/sdk/nuxt)

## Core Features

### Auto-Tracking

Enable automatic tracking of:
- Page views
- Link clicks
- Form submissions
- JavaScript errors

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-website',
  autoTrack: true  // Enable auto-tracking
});
```

### Custom Events

Track any custom event:

```javascript
analytics.track('button_clicked', {
  button_id: 'signup',
  location: 'hero',
  variant: 'primary'
});
```

### User Identification

Identify users for cohort analysis:

```javascript
analytics.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium',
  signup_date: '2024-01-01'
});
```

### Page Views

Track page navigation:

```javascript
# or
npm install @hefni101/siraaj
analytics.pageView('/products');
pnpm add @hefni101/siraaj
// Or with custom properties
analytics.pageView('/products', {
yarn add @hefni101/siraaj
  category: 'electronics',
  subcategory: 'phones'
});
```
import { analytics } from '@hefni101/siraaj';
## Configuration Options

```typescript
interface AnalyticsConfig {
  // Required
  apiUrl: string;              // Your Siraaj server URL
  projectId: string;           // Project identifier
  
  // Optional
import { AnalyticsProvider, useAnalytics } from '@hefni101/siraaj/react';
  debug?: boolean;             // Enable debug logs (default: false)
  bufferSize?: number;         // Events per batch (default: 10)
  flushInterval?: number;      // Auto-flush interval in ms (default: 30000)
  timeout?: number;            // Request timeout in ms (default: 10000)
  maxRetries?: number;         // Max retry attempts (default: 3)
  useBeacon?: boolean;         // Use sendBeacon API (default: true)
  sampling?: number;           // Sample rate 0-1 (default: 1.0)
  respectDoNotTrack?: boolean; // Honor DNT header (default: true)
  enablePerformanceTracking?: boolean; // Track Web Vitals (default: false)
}
```

## API Methods
import { useAnalytics } from '@hefni101/siraaj/vue';
### `analytics.track(eventName, properties?)`

Track a custom event.

```javascript
analytics.track('purchase_completed', {
  product_id: 'abc123',
  price: 29.99,
import { createAnalytics } from '@hefni101/siraaj/svelte';
});
```

### `analytics.pageView(url?, properties?)`

Track a page view.

import { AnalyticsProvider } from '@hefni101/siraaj/next';
analytics.pageView('/products/123', {
  product_name: 'Widget',
  category: 'electronics'
});
import { useNuxtAnalytics } from '@hefni101/siraaj/nuxt';

### `analytics.identify(userId, traits?)`

Identify a user.

```javascript
analytics.identify('user-456', {
  name: 'John Doe',
  email: 'john@example.com',
  plan: 'enterprise'
});
```

### `analytics.trackClick(elementId, properties?)`

Track a click event.

```javascript
analytics.trackClick('signup-button', {
  location: 'navbar',
  variant: 'cta'
});
```

### `analytics.trackForm(formId, properties?)`

Track a form submission.

```javascript
analytics.trackForm('contact-form', {
  form_type: 'contact',
  source: 'landing-page'
});
```

### `analytics.trackError(error, context?)`

Track an error.

```javascript
try {
  // risky operation
} catch (error) {
  analytics.trackError(error, {
    component: 'checkout',
    action: 'process_payment'
  });
}
```

### `analytics.flush()`

Manually flush the event buffer.

```javascript
await analytics.flush();
```

### `analytics.reset()`

Reset session and user ID.

```javascript
// Call on logout
analytics.reset();
```

## Bundle Sizes

All sizes are **gzipped**:

| Package | Size | Description |
|---------|------|-------------|
| Core | < 3 KB | Vanilla JS/TS |
| React | < 4 KB | React hooks + core |
| Vue | < 4 KB | Vue composables + core |
| Svelte | < 4 KB | Svelte stores + core |
| Next.js | < 4 KB | Next.js integration + core |
| Nuxt | < 4 KB | Nuxt integration + core |
| Preact | < 4 KB | Preact hooks + core |

## Browser Support

- Chrome/Edge ≥ 90
- Firefox ≥ 88
- Safari ≥ 14
- All modern browsers with ES2020 support

## TypeScript Support

Full TypeScript definitions included:

```typescript
import type { AnalyticsConfig, EventData } from '@hefni101/siraaj';

const config: AnalyticsConfig = {
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  debug: true
};
```

## Privacy Features

### Respect Do Not Track

Automatically respects the DNT header:

```javascript
const analytics = new Analytics({
  apiUrl: '...',
  projectId: '...',
  respectDoNotTrack: true  // default: true
});
```

### Sampling

Reduce tracking volume with sampling:

```javascript
const analytics = new Analytics({
  apiUrl: '...',
  projectId: '...',
  sampling: 0.5  // Track 50% of users
});
```

### No Cookies

Siraaj uses `sessionStorage` and `localStorage` instead of cookies:
- No cookie banners needed
- No third-party tracking
- Session-based analytics

## Performance

### Batching

Events are automatically batched:

```javascript
const analytics = new Analytics({
  bufferSize: 20,      // Send after 20 events
  flushInterval: 30000 // Or every 30 seconds
});
```

### Retry Logic

Failed requests are automatically retried with exponential backoff:

```javascript
const analytics = new Analytics({
  maxRetries: 3,       // Retry up to 3 times
  timeout: 10000       // 10 second timeout
});
```

### SendBeacon API

Uses `navigator.sendBeacon` for reliability during page unload:

```javascript
const analytics = new Analytics({
  useBeacon: true  // default: true
});
```

## Next Steps

Choose your framework:

<div class="grid-container">
  <a href="/sdk/vanilla" class="grid-item">
    <h3>📦 Vanilla JS</h3>
    <p>Pure JavaScript integration</p>
  </a>
  
  <a href="/sdk/react" class="grid-item">
    <h3>⚛️ React</h3>
    <p>Hooks and providers</p>
  </a>
  
  <a href="/sdk/vue" class="grid-item">
    <h3>💚 Vue</h3>
    <p>Composables and plugins</p>
  </a>
  
  <a href="/sdk/svelte" class="grid-item">
    <h3>🧡 Svelte</h3>
    <p>Stores and actions</p>
  </a>
  
  <a href="/sdk/nextjs" class="grid-item">
    <h3>▲ Next.js</h3>
    <p>App & Pages Router</p>
  </a>
  
  <a href="/sdk/nuxt" class="grid-item">
    <h3>💚 Nuxt</h3>
    <p>Plugins and composables</p>
  </a>
</div>

<style>
.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin: 2rem 0;
}

.grid-item {
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.2s;
}

.grid-item:hover {
  border-color: var(--vp-c-brand);
  box-shadow: 0 2px 12px rgba(59, 130, 246, 0.1);
}

.grid-item h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
}

.grid-item p {
  margin: 0;
  color: var(--vp-c-text-2);
  font-size: 0.9rem;
}
</style>
