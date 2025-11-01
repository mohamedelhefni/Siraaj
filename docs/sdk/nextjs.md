# Next.js Integration

Complete guide for integrating Siraaj Analytics into Next.js applications (App Router and Pages Router).

## Installation

```bash
npm install @siraaj/analytics
```

## App Router (Next.js 13+)

### Root Layout

```tsx
// app/layout.tsx
import { AnalyticsProvider } from '@siraaj/analytics/next';

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AnalyticsProvider 
          config={{
            apiUrl: process.env.NEXT_PUBLIC_ANALYTICS_URL,
            projectId: 'my-nextjs-app',
            autoTrack: true
          }}
        >
          {children}
        </AnalyticsProvider>
      </body>
    </html>
  );
}
```

### Client Components

```tsx
'use client';

import { useAnalytics, useNextAnalytics } from '@siraaj/analytics/next';

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

## Pages Router

```tsx
// pages/_app.tsx
import { AnalyticsProvider } from '@siraaj/analytics/next';
import { useNextPagesAnalytics } from '@siraaj/analytics/next';

export default function App({ Component, pageProps }) {
  // Auto-track route changes
  useNextPagesAnalytics();
  
  return (
    <AnalyticsProvider config={{
      apiUrl: process.env.NEXT_PUBLIC_ANALYTICS_URL,
      projectId: 'my-app'
    }}>
      <Component {...pageProps} />
    </AnalyticsProvider>
  );
}
```

## Environment Variables

```env
# .env.local
NEXT_PUBLIC_ANALYTICS_URL=http://localhost:8080
```

## Next Steps

- [Nuxt Integration →](/sdk/nuxt)
- [Custom Events →](/sdk/custom-events)
