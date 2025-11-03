# Siraaj Analytics SDK Examples

This directory contains framework-specific examples demonstrating how to use the Siraaj Analytics SDK.

## Examples

- **react-example.html** - React hooks and provider pattern
- **vue-example.html** - Vue 3 composables and plugin
- **svelte-example.html** - Svelte stores and lifecycle
- **preact-example.html** - Preact hooks (similar to React)
- **vanilla-example.html** - Pure JavaScript usage

## Running Examples

Simply open any `.html` file in your browser. No build step required!

Each example demonstrates:
- SDK initialization
- Event tracking
- User identification
- Real-time event display

## Production Usage

In production, install the SDK via npm/pnpm:

```bash
pnpm add @hefni101/siraaj
```

Then import the framework-specific package:

```javascript
// React
import { useAnalytics } from '@hefni101/siraaj/react';

// Vue
import { useAnalytics } from '@hefni101/siraaj/vue';

// Svelte
import { createAnalytics } from '@hefni101/siraaj/svelte';

// Preact
import { useAnalytics } from '@hefni101/siraaj/preact';

// Vanilla JS
import { analytics } from '@hefni101/siraaj';
```
