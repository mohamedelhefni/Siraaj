# Core SDK

The Siraaj Analytics core SDK is a lightweight, framework-agnostic JavaScript library.

## Installation

### CDN

```html
<script src="http://localhost:8080/sdk/analytics.js"></script>
```

### NPM

```bash
npm install @hefni101/siraaj
```

## Initialization

```javascript
import { analytics } from '@hefni101/siraaj';

analytics.init({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  autoTrack: true,
  debug: false
});
```

## Configuration

See [Configuration Guide](/sdk/configuration) for all options.

## Methods

### track()

Track custom events.

```javascript
analytics.track('event_name', { property: 'value' });
```

### pageView()

Track page views.

```javascript
analytics.pageView('/path');
```

### identify()

Identify users.

```javascript
analytics.identify('user-id', { email: 'user@example.com' });
```

## Next Steps

- [Configuration →](/sdk/configuration)
- [Framework Integration →](/sdk/overview)
