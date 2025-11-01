# Auto-Tracking

Siraaj automatically tracks common user interactions.

## Enabled by Default

When `autoTrack: true`, Siraaj automatically tracks:

- Page views
- Link clicks
- Form submissions
- JavaScript errors
- Page visibility changes

## Page Views

```javascript
// Tracked automatically on page load
// and navigation in SPAs
```

## Link Clicks

```javascript
// Automatically tracks clicks on <a> tags
// Captures: URL, link text, external/internal
```

## Form Submissions

```javascript
// Automatically tracks form submissions
// Captures: form ID, action URL
```

## Errors

```javascript
// Automatically tracks JavaScript errors
// Captures: error message, stack trace, filename
```

## Disable Auto-Tracking

```javascript
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-app',
  autoTrack: false // Disable all auto-tracking
});
```

## Next Steps

- [Custom Events →](/sdk/custom-events)
