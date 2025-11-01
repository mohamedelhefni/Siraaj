# Svelte Integration

Integrate Siraaj Analytics into your Svelte application.

## Installation

```bash
npm install @siraaj/analytics
```

## Setup

```javascript
// main.js or +layout.js
import { initAnalytics } from '@siraaj/analytics/svelte';

initAnalytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-svelte-app',
  autoTrack: true
});
```

## Usage

```svelte
<script>
import { createAnalytics, usePageTracking } from '@siraaj/analytics/svelte';

const { track, identify } = createAnalytics();

// Auto-track page view
usePageTracking();

function handleClick() {
  track('button_clicked', {
    button: 'signup'
  });
}
</script>

<button on:click={handleClick}>Sign Up</button>
```

## SvelteKit Integration

```javascript
// src/routes/+layout.js
import { initAnalytics } from '@siraaj/analytics/svelte';

export const load = async () => {
  if (typeof window !== 'undefined') {
    initAnalytics({
      apiUrl: 'http://localhost:8080',
      projectId: 'my-sveltekit-app'
    });
  }
};
```

## Next Steps

- [Custom Events →](/sdk/custom-events)
- [User Identification →](/sdk/user-identification)
