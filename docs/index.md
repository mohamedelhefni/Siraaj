---
layout: home

hero:
  name: Siraaj Analytics
  text: Fast, Simple, Self-Hosted Analytics
  tagline: Privacy-focused web analytics powered by Go and DuckDB
  image:
    src: /hero-image.svg
    alt: Siraaj Analytics
  actions:
    - theme: brand
      text: Get Started
      link: /guide/quick-start
    - theme: alt
      text: View on GitHub
      link: https://github.com/mohamedelhefni/siraaj

features:
  - icon: 🔒
    title: Privacy First
    details: No cookies, no tracking across sites. Fully GDPR compliant. Your users' privacy is protected.
  
  - icon: ⚡
    title: Blazing Fast
    details: Powered by DuckDB for lightning-fast analytics queries. Get insights in milliseconds.
  
  - icon: 📊
    title: Beautiful Dashboard
    details: Modern, responsive UI built with SvelteKit. See your data in real-time with interactive charts.
  
  - icon: 🎯
    title: Simple Integration
    details: Drop-in JavaScript SDK that works anywhere. Framework-specific integrations for React, Vue, Svelte, and more.
  
  - icon: 🌍
    title: Multi-Project Support
    details: Track multiple websites from one instance. Perfect for agencies and multi-site businesses.
  
  - icon: 📈
    title: Real-Time Analytics
    details: See your data as it happens. Live visitor counts, real-time events, and instant insights.
  
  - icon: 💾
    title: Self-Hosted
    details: Your data stays on your infrastructure. Complete control over your analytics data.
  
  - icon: 🏗️
    title: Production Ready
    details: Clean architecture, fully tested, and battle-tested. Ready for production deployment.
  
  - icon: 🔍
    title: Advanced Tracking
    details: Custom events, user identification, funnel analysis, and channel attribution.
---

<style>
:root {
  --vp-home-hero-name-color: transparent;
  --vp-home-hero-name-background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  --vp-home-hero-image-background-image: linear-gradient(135deg, #3b82f620 0%, #8b5cf620 100%);
  --vp-home-hero-image-filter: blur(44px);
}

@media (min-width: 640px) {
  :root {
    --vp-home-hero-image-filter: blur(56px);
  }
}

@media (min-width: 960px) {
  :root {
    --vp-home-hero-image-filter: blur(68px);
  }
}
</style>

## Quick Example

```javascript
// Initialize Siraaj Analytics
const analytics = new Analytics({
  apiUrl: 'http://localhost:8080',
  projectId: 'my-website',
  autoTrack: true
});

// Track custom events
analytics.track('button_clicked', {
  button_id: 'signup',
  location: 'hero'
});

// Identify users
analytics.identify('user-123', {
  email: 'user@example.com',
  plan: 'premium'
});
```

## Why Siraaj?

Siraaj is built for developers who need a **fast**, **privacy-focused**, and **self-hosted** analytics solution. Unlike traditional analytics platforms:

- ✅ **No external dependencies** - All data stays on your server
- ✅ **Lightweight** - SDK is < 5KB gzipped
- ✅ **Fast queries** - DuckDB columnar storage for instant insights
- ✅ **Open source** - MIT licensed, contribute and customize
- ✅ **Developer friendly** - Clean APIs, TypeScript support, comprehensive docs

## Trusted By Developers

<div class="stats-grid">
  <div class="stat-item">
    <div class="stat-value">< 5KB</div>
    <div class="stat-label">SDK Size</div>
  </div>
  <div class="stat-item">
    <div class="stat-value">< 50ms</div>
    <div class="stat-label">Query Speed</div>
  </div>
  <div class="stat-item">
    <div class="stat-value">100%</div>
    <div class="stat-label">GDPR Compliant</div>
  </div>
  <div class="stat-item">
    <div class="stat-value">MIT</div>
    <div class="stat-label">Open Source</div>
  </div>
</div>

<style>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 2rem;
  margin: 3rem 0;
  padding: 2rem;
  background: var(--vp-c-bg-soft);
  border-radius: 12px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--vp-c-text-2);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
</style>
