# Nuxt Integration

Integrate Siraaj Analytics into your Nuxt 3 application.

## Installation

```bash
npm install @siraaj/analytics
```

## Setup

### Plugin

```typescript
// plugins/analytics.client.ts
import { initNuxtAnalytics } from '@siraaj/analytics/nuxt';

export default defineNuxtPlugin(() => {
  initNuxtAnalytics({
    apiUrl: useRuntimeConfig().public.analyticsUrl,
    projectId: 'my-nuxt-app',
    autoTrack: true
  });
});
```

### Runtime Config

```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      analyticsUrl: process.env.NUXT_PUBLIC_ANALYTICS_URL || 'http://localhost:8080'
    }
  }
});
```

## Usage

```vue
<script setup>
import { useNuxtAnalytics } from '@siraaj/analytics/nuxt';

// Auto-track route changes
useNuxtAnalytics();

const { $analytics } = useNuxtApp();

function handleClick() {
  $analytics.track('button_clicked', {
    button: 'signup'
  });
}
</script>

<template>
  <button @click="handleClick">Sign Up</button>
</template>
```

## Environment Variables

```env
# .env
NUXT_PUBLIC_ANALYTICS_URL=http://localhost:8080
```

## Next Steps

- [Vue Integration →](/sdk/vue)
- [Custom Events →](/sdk/custom-events)
