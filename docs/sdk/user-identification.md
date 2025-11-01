# User Identification

Identify and track users across sessions.

## Basic Usage

```javascript
analytics.identify('user-123', {
  email: 'user@example.com',
  name: 'John Doe',
  plan: 'premium'
});
```

## When to Identify

### On Login

```javascript
function handleLogin(user) {
  analytics.identify(user.id, {
    email: user.email,
    name: user.name,
    plan: user.subscription.plan,
    signup_date: user.created_at
  });
}
```

### On Signup

```javascript
function handleSignup(user) {
  analytics.identify(user.id, {
    email: user.email,
    signup_method: 'email',
    plan: 'free',
    referral_source: getReferralSource()
  });
  
  analytics.track('signup_completed', {
    method: 'email'
  });
}
```

### On Logout

```javascript
function handleLogout() {
  analytics.reset();
  // User gets a new anonymous ID
}
```

## User Properties

```javascript
// Set user properties
analytics.identify('user-123', {
  // Demographics
  name: 'John Doe',
  email: 'john@example.com',
  company: 'Acme Inc',
  
  // Subscription
  plan: 'enterprise',
  mrr: 299,
  
  // Attributes
  signup_date: '2024-01-15',
  last_login: '2024-01-20',
  feature_flags: ['beta_feature'],
  
  // Custom
  role: 'admin',
  team_size: 10
});
```

## Update Properties

```javascript
// Update user properties
analytics.setUserProperties({
  plan: 'enterprise', // Upgraded
  mrr: 499,
  last_purchase_date: '2024-01-20'
});
```

## Best Practices

1. **Don't store PII** unless necessary
2. **Use consistent IDs** across platforms
3. **Set on every login**
4. **Reset on logout**

## Next Steps

- [Auto-Tracking →](/sdk/auto-tracking)
- [Performance →](/sdk/performance)
