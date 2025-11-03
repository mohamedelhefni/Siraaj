# Svelte Integration

Integrate Siraaj Analytics into your Svelte application.

## Installation

```bash
npm install @hefni101/siraaj
```

## Setup

```javascript
// main.js or +layout.js
import { initAnalytics } from '@hefni101/siraaj/svelte';

initAnalytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-svelte-app',
  autoTrack: true
});
```

## Usage

```svelte
<script>
import { createAnalytics, usePageTracking } from '@hefni101/siraaj/svelte';

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
import { initAnalytics } from '@hefni101/siraaj/svelte';

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
