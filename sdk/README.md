# Siraaj Analytics SDK

Lightweight, framework-agnostic analytics SDK with first-class support for React, Vue, Svelte, Preact, Next.js, and Nuxt.

## Features

- 🪶 **Lightweight** - Core is < 3KB gzipped
- ⚡ **Fast** - Batching, buffering, and automatic retries
- 🎯 **Framework-specific** - Optimized hooks and components for React, Vue, Svelte, Preact, Next.js, and Nuxt
- 📦 **Tree-shakeable** - Only import what you need
- 🔒 **Type-safe** - Full TypeScript support
- 🚀 **Auto-tracking** - Page views, clicks, forms, and errors out of the box
- 🔄 **SSR Ready** - Works seamlessly with Next.js and Nuxt server-side rendering

## Installation

```bash
# Using pnpm
pnpm add @siraaj/analytics

# Using npm
npm install @siraaj/analytics

# Using yarn
yarn add @siraaj/analytics
```

## Quick Start

### Vanilla JavaScript / TypeScript

```javascript
import { analytics } from '@siraaj/analytics';

// Initialize
analytics.init({
  endpoint: 'https://your-analytics-server.com',
  apiKey: 'your-api-key',
  autoTrack: true,
});

// Track events
analytics.track('button_clicked', { 
  buttonId: 'signup',
  location: 'header' 
});

// Track page views
analytics.page();

// Identify users
analytics.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium',
});
```

### React

```jsx
import { AnalyticsProvider, useAnalytics, usePageTracking } from '@siraaj/analytics/react';

function App() {
  return (
    <AnalyticsProvider config={{
      endpoint: 'https://your-analytics-server.com',
      apiKey: 'your-api-key',
    }}>
      <YourApp />
    </AnalyticsProvider>
  );
}

function YourComponent() {
  const { track, identify } = useAnalytics();
  
  // Auto-track page views
  usePageTracking();
  
  const handleClick = () => {
    track('button_clicked', { button: 'signup' });
  };
  
  return <button onClick={handleClick}>Sign Up</button>;
}
```

### Vue 3

```vue
<script setup>
import { useAnalytics, usePageTracking } from '@siraaj/analytics/vue';

const { track, identify } = useAnalytics();

// Auto-track page view
usePageTracking();

const handleClick = () => {
  track('button_clicked', { button: 'signup' });
};
</script>

<template>
  <button @click="handleClick">Sign Up</button>
</template>
```

**Plugin registration:**

```javascript
import { createApp } from 'vue';
import { AnalyticsPlugin } from '@siraaj/analytics/vue';
import App from './App.vue';

const app = createApp(App);

app.use(AnalyticsPlugin, {
  endpoint: 'https://your-analytics-server.com',
  apiKey: 'your-api-key',
});

app.mount('#app');
```

### Svelte

```svelte
<script>
import { createAnalytics, usePageTracking } from '@siraaj/analytics/svelte';

const { track, userId } = createAnalytics();

// Auto-track page view
usePageTracking();

const handleClick = () => {
  track('button_clicked', { button: 'signup' });
};
</script>

<button on:click={handleClick}>Sign Up</button>
```

**Initialization:**

```javascript
// In your main file
import { initAnalytics } from '@siraaj/analytics/svelte';

initAnalytics({
  endpoint: 'https://your-analytics-server.com',
  apiKey: 'your-api-key',
});
```

### Next.js (App Router)

```tsx
// app/layout.tsx
import { AnalyticsProvider } from '@siraaj/analytics/next';
import { initNextAnalytics } from '@siraaj/analytics/next';

initNextAnalytics({
  endpoint: 'https://your-analytics-server.com',
  apiKey: 'your-api-key',
});

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <body>
        <AnalyticsProvider config={{
          endpoint: 'https://your-analytics-server.com',
          apiKey: 'your-api-key',
        }}>
          {children}
        </AnalyticsProvider>
      </body>
    </html>
  );
}

// app/page.tsx
'use client';
import { useNextAnalytics, useAnalytics } from '@siraaj/analytics/next';

export default function Page() {
  const { track } = useAnalytics();
  
  // Auto-track route changes
  useNextAnalytics();
  
  return (
    <button onClick={() => track('clicked')}>
      Click Me
    </button>
  );
}
```

### Next.js (Pages Router)

```tsx
// pages/_app.tsx
import { AnalyticsProvider } from '@siraaj/analytics/next';
import { useNextPagesAnalytics } from '@siraaj/analytics/next';

export default function App({ Component, pageProps }) {
  // Auto-track route changes
  useNextPagesAnalytics();
  
  return (
    <AnalyticsProvider config={{
      endpoint: 'https://your-analytics-server.com',
      apiKey: 'your-api-key',
    }}>
      <Component {...pageProps} />
    </AnalyticsProvider>
  );
}
```

### Nuxt 3

```typescript
// plugins/analytics.client.ts
import { initNuxtAnalytics } from '@siraaj/analytics/nuxt';

export default defineNuxtPlugin(() => {
  initNuxtAnalytics({
    endpoint: 'https://your-analytics-server.com',
    apiKey: 'your-api-key',
  });
});

// composables/useTracking.ts
import { useNuxtAnalytics } from '@siraaj/analytics/nuxt';

export function useTracking() {
  // Auto-track route changes
  useNuxtAnalytics();
}

// pages/index.vue
<script setup>
import { useNuxtApp } from '#app';

const { $analytics } = useNuxtApp();

const handleClick = () => {
  $analytics.track('button_clicked', { button: 'signup' });
};
</script>

<template>
  <button @click="handleClick">Sign Up</button>
</template>
```

### Preact

```jsx
import { AnalyticsProvider, useAnalytics } from '@siraaj/analytics/preact';

export function App() {
  return (
    <AnalyticsProvider config={{
      endpoint: 'https://your-analytics-server.com',
      apiKey: 'your-api-key',
    }}>
      <YourApp />
    </AnalyticsProvider>
  );
}

function YourComponent() {
  const { track } = useAnalytics();
  
  return (
    <button onClick={() => track('clicked')}>
      Click Me
    </button>
  );
}
```

## Configuration

```typescript
interface AnalyticsConfig {
  endpoint: string;        // Your analytics server endpoint
  apiKey: string;          // API key for authentication
  autoTrack?: boolean;     // Enable auto-tracking (default: true)
  debug?: boolean;         // Enable debug logs (default: false)
  batchSize?: number;      // Events per batch (default: 10)
  flushInterval?: number;  // Flush interval in ms (default: 5000)
  maxRetries?: number;     // Max retry attempts (default: 3)
  timeout?: number;        // Request timeout in ms (default: 10000)
}
```

## API Reference

### Core Methods

#### `analytics.init(config)`
Initialize the analytics SDK with configuration.

#### `analytics.track(event, properties?)`
Track a custom event with optional properties.

```javascript
analytics.track('purchase_completed', {
  productId: 'abc123',
  price: 29.99,
  currency: 'USD',
});
```

#### `analytics.page(properties?)`
Track a page view with optional properties.

```javascript
analytics.page({
  category: 'documentation',
  section: 'getting-started',
});
```

#### `analytics.identify(userId, properties?)`
Identify a user with optional properties.

```javascript
analytics.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium',
  signupDate: '2024-01-01',
});
```

#### `analytics.flush()`
Manually flush the event queue.

```javascript
await analytics.flush();
```

#### `analytics.reset()`
Reset the session and user ID.

```javascript
analytics.reset();
```

## Auto-Tracking

When `autoTrack: true`, the SDK automatically tracks:

- **Page views** - On initialization and navigation
- **Clicks** - On links and buttons
- **Form submissions** - On form submit events
- **Errors** - JavaScript errors and exceptions

You can disable auto-tracking and manually track events as needed.

## Bundle Sizes

All sizes are **gzipped**:

| Package | Size | Description |
|---------|------|-------------|
| Core | < 3 KB | Vanilla JS/TS |
| React | < 4 KB | React hooks + core |
| Vue | < 4 KB | Vue composables + core |
| Svelte | < 4 KB | Svelte stores + core |
| Preact | < 4 KB | Preact hooks + core |
| Next.js | < 4 KB | Next.js integration + core |
| Nuxt | < 4 KB | Nuxt integration + core |

## Browser Support

- Chrome/Edge ≥ 90
- Firefox ≥ 88
- Safari ≥ 14
- All modern browsers with ES2020 support

## Development

```bash
# Install dependencies
pnpm install

# Build all packages
pnpm run build

# Watch mode
pnpm run dev

# Check bundle sizes
pnpm run size

# Analyze bundle
pnpm run analyze
```

## License

MIT © Mohamed Elhefni
