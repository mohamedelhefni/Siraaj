# Custom Events

Track custom events with Siraaj Analytics.

## Basic Usage

```javascript
analytics.track('event_name', {
  property1: 'value1',
  property2: 123,
  property3: true
});
```

## Event Naming

Use descriptive, action-based names:

```javascript
// ✅ Good
analytics.track('signup_completed');
analytics.track('video_played');
analytics.track('purchase_initiated');

// ❌ Bad
analytics.track('event1');
analytics.track('click');
analytics.track('btn');
```

## Event Properties

```javascript
analytics.track('purchase_completed', {
  order_id: 'order-123',
  total: 99.99,
  currency: 'USD',
  items: [
    {
      product_id: 'prod-456',
      name: 'Product Name',
      price: 99.99,
      quantity: 1
    }
  ],
  payment_method: 'credit_card',
  shipping_method: 'express'
});
```

## Common Events

### E-commerce

```javascript
// Product viewed
analytics.track('product_viewed', {
  product_id: 'prod-123',
  name: 'Product Name',
  price: 99.99,
  category: 'Electronics'
});

// Add to cart
analytics.track('add_to_cart', {
  product_id: 'prod-123',
  quantity: 1,
  price: 99.99
});

// Checkout started
analytics.track('checkout_started', {
  cart_total: 199.98,
  item_count: 2
});

// Purchase completed
analytics.track('purchase_completed', {
  order_id: 'order-123',
  total: 199.98,
  items: [...]
});
```

### User Actions

```javascript
// Signup
analytics.track('signup_completed', {
  method: 'email',
  plan: 'premium'
});

// Login
analytics.track('login_completed', {
  method: 'google'
});

// Search
analytics.track('search_performed', {
  query: 'analytics tool',
  results_count: 42
});
```

## Best Practices

1. **Be consistent** - Use snake_case for event names
2. **Add context** - Include relevant properties
3. **Track outcomes** - Focus on user actions and results
4. **Keep it simple** - Don't track everything

## Next Steps

- [User Identification →](/sdk/user-identification)
- [Auto-Tracking →](/sdk/auto-tracking)
