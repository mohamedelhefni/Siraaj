# React Integration

Complete guide for integrating Siraaj Analytics into your React application.

## Installation

```bash
npm install @hefni101/siraaj
# or
pnpm add @hefni101/siraaj
# or
yarn add @hefni101/siraaj
```

## Quick Start

### 1. Provider Setup

Wrap your app with the `AnalyticsProvider`:

```jsx
// App.jsx or App.tsx
import { AnalyticsProvider } from '@hefni101/siraaj/react';

function App() {
  return (
    <AnalyticsProvider 
      config={{
        apiUrl: 'http://localhost:8080',
        projectId: 'my-react-app',
        autoTrack: true,
        debug: process.env.NODE_ENV === 'development'
      }}
    >
      <YourApp />
    </AnalyticsProvider>
  );
}

export default App;
```

### 2. Use the Hook

```jsx
import { useAnalytics } from '@hefni101/siraaj/react';

function SignupButton() {
  const { track } = useAnalytics();
  
  const handleClick = () => {
    track('signup_clicked', {
      location: 'hero',
      plan: 'premium'
    });
  };
  
  return (
    <button onClick={handleClick}>
      Sign Up
    </button>
  );
}
```

## Hooks API

### `useAnalytics()`

Main hook for tracking events and identifying users.

```jsx
import { useAnalytics } from '@hefni101/siraaj/react';

function MyComponent() {
  const { track, identify, pageView, reset } = useAnalytics();
  
  // Track custom event
  const handlePurchase = (productId) => {
    track('purchase_completed', {
      product_id: productId,
      price: 99.99,
      currency: 'USD'
    });
  };
  
  // Identify user
  const handleLogin = (user) => {
    identify(user.id, {
      email: user.email,
      name: user.name,
      plan: user.plan
    });
  };
  
  // Manual page view tracking
  const handleNavigation = (path) => {
    pageView(path);
  };
  
  // Reset on logout
  const handleLogout = () => {
    reset();
  };
  
  return (
    <div>
      <button onClick={() => handlePurchase('prod-123')}>Buy Now</button>
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
}
```

### `usePageTracking()`

Auto-track page views (useful with React Router):

```jsx
import { usePageTracking } from '@hefni101/siraaj/react';

function App() {
  // Automatically tracks page views on route changes
  usePageTracking();
  
  return (
    <Router>
      {/* Your routes */}
    </Router>
  );
}
```

With custom properties:

```jsx
import { usePageTracking } from '@hefni101/siraaj/react';
import { useLocation } from 'react-router-dom';

function App() {
  const location = useLocation();
  
  usePageTracking({
    // Custom properties sent with each page view
    app_version: '1.0.0',
    environment: process.env.NODE_ENV
  });
  
  return (
    <Router>
      {/* Your routes */}
    </Router>
  );
}
```

### `useIdentify()`

Auto-identify users when they login:

```jsx
import { useIdentify } from '@hefni101/siraaj/react';

function App() {
  const { user } = useAuth();
  
  // Automatically identifies user when user object changes
  useIdentify(user?.id, {
    email: user?.email,
    plan: user?.plan,
    signup_date: user?.created_at
  });
  
  return <YourApp />;
}
```

## React Router Integration

### With React Router v6

```jsx
import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom';
import { AnalyticsProvider, usePageTracking } from '@hefni101/siraaj/react';

function AnalyticsWrapper({ children }) {
  usePageTracking();
  return children;
}

function App() {
  return (
    <AnalyticsProvider config={{ apiUrl: '...', projectId: '...' }}>
      <BrowserRouter>
        <AnalyticsWrapper>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/about" element={<About />} />
            <Route path="/products" element={<Products />} />
          </Routes>
        </AnalyticsWrapper>
      </BrowserRouter>
    </AnalyticsProvider>
  );
}
```

### Custom Route Tracking

```jsx
import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { useAnalytics } from '@hefni101/siraaj/react';

function RouteTracker() {
  const location = useLocation();
  const { pageView } = useAnalytics();
  
  useEffect(() => {
    pageView(location.pathname, {
      search: location.search,
      hash: location.hash,
      state: location.state
    });
  }, [location, pageView]);
  
  return null;
}
```

## TypeScript Support

Full type definitions included:

```tsx
import { useAnalytics } from '@hefni101/siraaj/react';
import type { AnalyticsConfig } from '@hefni101/siraaj';

const config: AnalyticsConfig = {
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  debug: true
};

function MyComponent() {
  const { track } = useAnalytics();
  
  track('event_name', {
    property1: 'value1',
    property2: 123,
    property3: true
  });
  
  return <div>...</div>;
}
```

## Advanced Patterns

### Track Form Submissions

```jsx
import { useAnalytics } from '@hefni101/siraaj/react';

function ContactForm() {
  const { trackForm } = useAnalytics();
  
  const handleSubmit = (e) => {
    e.preventDefault();
    
    trackForm('contact-form', {
      source: 'homepage',
      form_type: 'contact'
    });
    
    // Submit form...
  };
  
  return (
    <form onSubmit={handleSubmit} id="contact-form">
      {/* Form fields */}
      <button type="submit">Submit</button>
    </form>
  );
}
```

### Track Clicks

```jsx
import { useAnalytics } from '@hefni101/siraaj/react';

function NavBar() {
  const { trackClick } = useAnalytics();
  
  return (
    <nav>
      <button 
        id="signup-btn"
        onClick={() => trackClick('signup-btn', { 
          location: 'navbar',
          variant: 'primary'
        })}
      >
        Sign Up
      </button>
    </nav>
  );
}
```

### Track Errors

```jsx
import { useAnalytics } from '@hefni101/siraaj/react';
import { useEffect } from 'react';

function ErrorBoundary({ error, children }) {
  const { trackError } = useAnalytics();
  
  useEffect(() => {
    if (error) {
      trackError(error, {
        component: 'ErrorBoundary',
        fatal: true
      });
    }
  }, [error, trackError]);
  
  if (error) {
    return <div>Something went wrong</div>;
  }
  
  return children;
}
```

### E-commerce Tracking

```jsx
import { useAnalytics } from '@hefni101/siraaj/react';

function ProductPage({ product }) {
  const { track } = useAnalytics();
  
  const handleAddToCart = () => {
    track('add_to_cart', {
      product_id: product.id,
      product_name: product.name,
      price: product.price,
      currency: 'USD',
      category: product.category
    });
  };
  
  const handlePurchase = (orderId, total) => {
    track('purchase_completed', {
      order_id: orderId,
      total: total,
      currency: 'USD',
      items: [
        {
          product_id: product.id,
          quantity: 1,
          price: product.price
        }
      ]
    });
  };
  
  return (
    <div>
      <h1>{product.name}</h1>
      <button onClick={handleAddToCart}>Add to Cart</button>
      <button onClick={() => handlePurchase('ord-123', 99.99)}>
        Buy Now
      </button>
    </div>
  );
}
```

## Context API

Access the analytics instance directly:

```jsx
import { useContext } from 'react';
import { AnalyticsContext } from '@hefni101/siraaj/react';

function MyComponent() {
  const analytics = useContext(AnalyticsContext);
  
  if (!analytics) {
    throw new Error('Analytics not initialized');
  }
  
  // Use analytics instance directly
  analytics.track('custom_event');
  
  return <div>...</div>;
}
```

## SSR / Next.js

For Next.js, use the dedicated Next.js integration:

```jsx
// See /sdk/nextjs for details
import { AnalyticsProvider } from '@hefni101/siraaj/next';
```

[Next.js Integration Guide →](/sdk/nextjs)

## Best Practices

### 1. Initialize Early

Place the provider as high as possible in your component tree:

```jsx
// ✅ Good
ReactDOM.render(
  <AnalyticsProvider config={...}>
    <App />
  </AnalyticsProvider>,
  document.getElementById('root')
);

// ❌ Bad - too deep in the tree
function App() {
  return (
    <div>
      <AnalyticsProvider config={...}>
        <Component />
      </AnalyticsProvider>
    </div>
  );
}
```

### 2. Memoize Event Handlers

```jsx
import { useCallback } from 'react';
import { useAnalytics } from '@hefni101/siraaj/react';

function MyComponent() {
  const { track } = useAnalytics();
  
  const handleClick = useCallback(() => {
    track('button_clicked');
  }, [track]);
  
  return <button onClick={handleClick}>Click Me</button>;
}
```

### 3. Conditional Tracking

```jsx
const { track } = useAnalytics();

// Only track in production
if (process.env.NODE_ENV === 'production') {
  track('event');
}
```

### 4. Cleanup on Unmount

```jsx
import { useEffect } from 'react';
import { useAnalytics } from '@hefni101/siraaj/react';

function MyComponent() {
  const { track } = useAnalytics();
  
  useEffect(() => {
    track('component_mounted');
    
    return () => {
      track('component_unmounted');
    };
  }, [track]);
  
  return <div>...</div>;
}
```

## Troubleshooting

### Provider Not Found

```
Error: useAnalytics must be used within AnalyticsProvider
```

**Solution**: Ensure your component is wrapped with `AnalyticsProvider`:

```jsx
<AnalyticsProvider config={...}>
  <YourComponent />
</AnalyticsProvider>
```

### Events Not Tracking

1. Check the browser console for errors
2. Enable debug mode:
   ```jsx
   <AnalyticsProvider config={{ debug: true, ... }}>
   ```
3. Verify the server URL is correct
4. Check CORS settings

### TypeScript Errors

Make sure you have the types installed:

```bash
npm install --save-dev @types/react
```

## Examples

Check out complete examples:

- [Basic React App](https://github.com/mohamedelhefni/siraaj/tree/main/sdk/examples/react-basic)
- [React Router](https://github.com/mohamedelhefni/siraaj/tree/main/sdk/examples/react-router)
- [E-commerce App](https://github.com/mohamedelhefni/siraaj/tree/main/sdk/examples/react-ecommerce)

## Next Steps

- [Track Custom Events →](/sdk/custom-events)
- [User Identification →](/sdk/user-identification)
