# Preact Integration

Integrate Siraaj Analytics into your Preact application.

## Installation

```bash
npm install @hefni101/siraaj
```

## Setup

```jsx
import { render } from 'preact';
import { AnalyticsProvider } from '@hefni101/siraaj/preact';
import App from './App';

render(
  <AnalyticsProvider config={{
    apiUrl: 'http://localhost:8080',
    projectId: 'my-preact-app',
    autoTrack: true
  }}>
    <App />
  </AnalyticsProvider>,
  document.getElementById('app')
);
```

## Usage

```jsx
import { useAnalytics } from '@hefni101/siraaj/preact';

export function MyComponent() {
  const { track } = useAnalytics();
  
  return (
    <button onClick={() => track('clicked')}>
      Click Me
    </button>
  );
}
```

## Next Steps

- [React Integration →](/sdk/react)
- [Custom Events →](/sdk/custom-events)
