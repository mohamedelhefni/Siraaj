# Vue Integration

Integrate Siraaj Analytics into your Vue 3 application.

## Installation

```bash
npm install @hefni101/siraaj
```

## Setup

### Plugin Registration

```javascript
// main.js
import { createApp } from 'vue';
import { AnalyticsPlugin } from '@hefni101/siraaj/vue';
import App from './App.vue';

const app = createApp(App);

app.use(AnalyticsPlugin, {
  apiUrl: 'http://localhost:8080',
  projectId: 'my-vue-app',
  autoTrack: true
});

app.mount('#app');
```

## Usage

### Composition API

```vue
<script setup>
import { useAnalytics, usePageTracking } from '@hefni101/siraaj/vue';

const { track, identify } = useAnalytics();

// Auto-track page views
usePageTracking();

const handleClick = () => {
  track('button_clicked', {
    button: 'signup'
  });
};
</script>

<template>
  <button @click="handleClick">Sign Up</button>
</template>
```

### Options API

```vue
<script>
export default {
  methods: {
    handleClick() {
      this.$analytics.track('button_clicked', {
        button: 'signup'
      });
    }
  }
}
</script>

<template>
  <button @click="handleClick">Sign Up</button>
</template>
```

## Router Integration

```javascript
import { createRouter } from 'vue-router';
import { useAnalytics } from '@hefni101/siraaj/vue';

const router = createRouter({...});

router.afterEach((to) => {
  const { pageView } = useAnalytics();
  pageView(to.path);
});
```

## Next Steps

- [Nuxt Integration →](/sdk/nuxt)
- [Custom Events →](/sdk/custom-events)
