# Funnel Analysis

Track multi-step conversion funnels to understand where users drop off in your conversion process.

## What are Funnels?

Funnels help you analyze the conversion path users take through specific steps on your website. For example:

- **Sign-up funnel**: Landing → Sign-up form → Email verification → Welcome page
- **Purchase funnel**: Product page → Add to cart → Checkout → Payment → Confirmation
- **Onboarding funnel**: First login → Profile setup → Tutorial completion

## API Endpoint

```http
POST /api/funnel
Content-Type: application/json
```

## Creating a Funnel

### Basic Example

```javascript
const funnelRequest = {
  steps: [
    {
      name: "Landing Page",
      event_name: "page_view",
      url: "/landing"
    },
    {
      name: "Sign Up Form",
      event_name: "page_view",
      url: "/signup"
    },
    {
      name: "Sign Up Completed",
      event_name: "signup_completed"
    }
  ],
  start_date: "2024-01-01",
  end_date: "2024-01-31",
  filters: {
    project: "my-website"
  }
};

const response = await fetch('http://localhost:8080/api/funnel', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(funnelRequest)
});

const result = await response.json();
```

### Response Structure

```json
{
  "steps": [
    {
      "step": {
        "name": "Landing Page",
        "event_name": "page_view",
        "url": "/landing"
      },
      "user_count": 1000,
      "session_count": 1200,
      "event_count": 1500,
      "conversion_rate": 100.0,
      "overall_rate": 100.0,
      "dropoff_rate": 0.0,
      "avg_time_to_next": 45.2,
      "median_time_to_next": 30.0
    },
    {
      "step": {
        "name": "Sign Up Form",
        "event_name": "page_view",
        "url": "/signup"
      },
      "user_count": 400,
      "session_count": 450,
      "event_count": 500,
      "conversion_rate": 40.0,
      "overall_rate": 40.0,
      "dropoff_rate": 60.0,
      "avg_time_to_next": 120.5,
      "median_time_to_next": 90.0
    },
    {
      "step": {
        "name": "Sign Up Completed",
        "event_name": "signup_completed"
      },
      "user_count": 200,
      "session_count": 210,
      "event_count": 220,
      "conversion_rate": 50.0,
      "overall_rate": 20.0,
      "dropoff_rate": 50.0,
      "avg_time_to_next": 0,
      "median_time_to_next": 0
    }
  ],
  "total_users": 1000,
  "completed_users": 200,
  "completion_rate": 20.0,
  "avg_completion": 165.7,
  "time_range": "2024-01-01 to 2024-01-31"
}
```

## Funnel Step Configuration

### Step Properties

Each funnel step can include:

- **name** (string, required): Display name for the step
- **event_name** (string, required): Event name to match
- **url** (string, optional): URL pattern to match
- **filters** (object, optional): Additional filters

### URL Matching

```javascript
{
  name: "Product Pages",
  event_name: "page_view",
  url: "/products/"  // Matches any URL starting with /products/
}
```

### Event-Only Steps

```javascript
{
  name: "Add to Cart",
  event_name: "add_to_cart"  // Matches this event regardless of URL
}
```

## Advanced Filtering

### Global Filters

Apply filters to all funnel steps:

```javascript
{
  steps: [...],
  start_date: "2024-01-01",
  end_date: "2024-01-31",
  filters: {
    project: "my-website",
    country: "US",
    device: "Desktop"
  }
}
```

### Per-Step Filters

Apply specific filters to individual steps:

```javascript
{
  steps: [
    {
      name: "Landing Page",
      event_name: "page_view",
      url: "/landing",
      filters: {
        source: "google"  // Only track users from Google
      }
    },
    {
      name: "Sign Up",
      event_name: "signup_completed"
    }
  ]
}
```

## Understanding Metrics

### User Count
Number of unique users who completed this step.

### Session Count
Number of unique sessions that reached this step.

### Event Count
Total number of times this event occurred.

### Conversion Rate
Percentage of users from the previous step who completed this step.

### Overall Rate
Percentage of users from the first step who completed this step.

### Dropoff Rate
Percentage of users who left after the previous step.

### Time to Next
Average/median time (in seconds) users took to reach the next step.

## Common Use Cases

### E-commerce Funnel

```javascript
{
  steps: [
    { name: "Product View", event_name: "page_view", url: "/products/" },
    { name: "Add to Cart", event_name: "add_to_cart" },
    { name: "Checkout", event_name: "page_view", url: "/checkout" },
    { name: "Payment", event_name: "payment_initiated" },
    { name: "Purchase", event_name: "purchase_completed" }
  ]
}
```

### Sign-up Funnel

```javascript
{
  steps: [
    { name: "Landing", event_name: "page_view", url: "/" },
    { name: "Pricing Page", event_name: "page_view", url: "/pricing" },
    { name: "Sign Up Form", event_name: "page_view", url: "/signup" },
    { name: "Account Created", event_name: "signup_completed" },
    { name: "Email Verified", event_name: "email_verified" }
  ]
}
```

### Content Engagement Funnel

```javascript
{
  steps: [
    { name: "Blog Home", event_name: "page_view", url: "/blog" },
    { name: "Article View", event_name: "page_view", url: "/blog/" },
    { name: "Scroll Depth 50%", event_name: "scroll_depth", filters: { depth: "50" } },
    { name: "Scroll Depth 100%", event_name: "scroll_depth", filters: { depth: "100" } },
    { name: "Subscribe", event_name: "newsletter_subscribe" }
  ]
}
```

## Best Practices

### 1. Define Clear Steps

Make sure each step is clearly defined and represents a meaningful action in your conversion process.

### 2. Track Custom Events

Use custom events for important actions:

```javascript
// On your website
analytics.track('add_to_cart', { product_id: '123' });
analytics.track('payment_initiated', { amount: 99.99 });
analytics.track('purchase_completed', { order_id: 'ord_456' });
```

### 3. Analyze Time Metrics

Pay attention to `avg_time_to_next` and `median_time_to_next` to identify:
- Steps where users get stuck
- Optimal timing for interventions
- Fast vs. slow conversion paths

### 4. Segment Your Funnels

Create separate funnels for different user segments:

```javascript
// Mobile users
{ filters: { device: "Mobile" } }

// Organic traffic
{ filters: { channel: "Organic" } }

// Specific countries
{ filters: { country: "US" } }
```

### 5. Monitor Dropoff Rates

High dropoff rates indicate:
- Confusing UI/UX
- Technical issues
- Missing information
- Price resistance

## Tracking Events for Funnels

### JavaScript SDK

```javascript
// Initialize
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-website',
  autoTrack: true
});

// Track custom funnel events
analytics.track('add_to_cart', {
  product_id: 'prod_123',
  product_name: 'Widget',
  price: 29.99
});

analytics.track('checkout_started', {
  cart_total: 89.97,
  item_count: 3
});

analytics.track('purchase_completed', {
  order_id: 'ord_456',
  total: 89.97,
  payment_method: 'credit_card'
});
```

## Example: Full Funnel Analysis

```javascript
async function analyzePurchaseFunnel() {
  const response = await fetch('http://localhost:8080/api/funnel', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      steps: [
        { name: "Homepage", event_name: "page_view", url: "/" },
        { name: "Product Page", event_name: "page_view", url: "/products/" },
        { name: "Add to Cart", event_name: "add_to_cart" },
        { name: "Checkout", event_name: "checkout_started" },
        { name: "Purchase", event_name: "purchase_completed" }
      ],
      start_date: "2024-01-01",
      end_date: "2024-01-31",
      filters: {
        project: "my-store"
      }
    })
  });
  
  const funnel = await response.json();
  
  console.log(`Total users who entered: ${funnel.total_users}`);
  console.log(`Users who completed: ${funnel.completed_users}`);
  console.log(`Completion rate: ${funnel.completion_rate}%`);
  
  funnel.steps.forEach((step, index) => {
    console.log(`\nStep ${index + 1}: ${step.step.name}`);
    console.log(`  Users: ${step.user_count}`);
    console.log(`  Conversion rate: ${step.conversion_rate}%`);
    console.log(`  Dropoff rate: ${step.dropoff_rate}%`);
  });
}
```

## Troubleshooting

### No Users in Funnel

- Verify events are being tracked correctly
- Check date range includes relevant data
- Ensure event names match exactly
- Verify URL patterns are correct

### Low Conversion Rates

- Review UX between steps
- Check for technical errors
- Analyze time metrics for bottlenecks
- Segment by traffic source to identify issues

### Missing Time Metrics

- Ensure events have correct timestamps
- Verify events are tracked in sequence
- Check that users complete steps in order

## Next Steps

- [Channels →](/guide/channels)
- [Dashboard →](/guide/dashboard)
- [Custom Events →](/sdk/custom-events)
- [API Reference →](/api/overview)
