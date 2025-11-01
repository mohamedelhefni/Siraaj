import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Siraaj Analytics',
  description: 'Fast, Simple, Self-Hosted Analytics',
  
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3b82f6' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'en' }],
    ['meta', { name: 'og:site_name', content: 'Siraaj Analytics' }],
  ],

  themeConfig: {
    logo: '/logo.svg',
    
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/introduction' },
      { text: 'SDK', link: '/sdk/overview' },
      { text: 'API', link: '/api/overview' },
      {
        text: 'v1.0.0',
        items: [
          { text: 'Changelog', link: '/changelog' },
          { text: 'Contributing', link: '/contributing' }
        ]
      }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Introduction', link: '/guide/introduction' },
            { text: 'Quick Start', link: '/guide/quick-start' },
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Configuration', link: '/guide/configuration' },
          ]
        },
        {
          text: 'Core Concepts',
          items: [
            { text: 'Architecture', link: '/guide/architecture' },
            { text: 'Events & Tracking', link: '/guide/events-tracking' },
            { text: 'Sessions', link: '/guide/sessions' },
            { text: 'Privacy', link: '/guide/privacy' },
          ]
        },
        {
          text: 'Features',
          items: [
            { text: 'Dashboard', link: '/guide/dashboard' },
            { text: 'Analytics', link: '/guide/analytics' },
            { text: 'Funnels', link: '/guide/funnels' },
            { text: 'Channels', link: '/guide/channels' },
          ]
        },
        {
          text: 'Deployment',
          items: [
            { text: 'Docker', link: '/guide/docker' },
            { text: 'Production', link: '/guide/production' },
            { text: 'Scaling', link: '/guide/scaling' },
          ]
        }
      ],
      
      '/sdk/': [
        {
          text: 'Overview',
          items: [
            { text: 'Introduction', link: '/sdk/overview' },
            { text: 'Core SDK', link: '/sdk/core' },
            { text: 'Configuration', link: '/sdk/configuration' },
          ]
        },
        {
          text: 'Frameworks',
          items: [
            { text: 'Vanilla JavaScript', link: '/sdk/vanilla' },
            { text: 'React', link: '/sdk/react' },
            { text: 'Vue', link: '/sdk/vue' },
            { text: 'Svelte', link: '/sdk/svelte' },
            { text: 'Next.js', link: '/sdk/nextjs' },
            { text: 'Nuxt', link: '/sdk/nuxt' },
            { text: 'Preact', link: '/sdk/preact' },
          ]
        },
        {
          text: 'Advanced',
          items: [
            { text: 'Custom Events', link: '/sdk/custom-events' },
            { text: 'User Identification', link: '/sdk/user-identification' },
            { text: 'Auto-Tracking', link: '/sdk/auto-tracking' },
            { text: 'Performance', link: '/sdk/performance' },
          ]
        }
      ],
      
      '/api/': [
        {
          text: 'API Reference',
          items: [
            { text: 'Overview', link: '/api/overview' },
            { text: 'Authentication', link: '/api/authentication' },
            { text: 'Track Events', link: '/api/track-events' },
            { text: 'Query Analytics', link: '/api/query-analytics' },
            { text: 'Error Handling', link: '/api/error-handling' },
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/mohamedelhefni/siraaj' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025 Mohamed Elhefni'
    },

    search: {
      provider: 'local'
    },

    editLink: {
      pattern: 'https://github.com/mohamedelhefni/siraaj/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    },

    lastUpdated: {
      text: 'Updated at',
      formatOptions: {
        dateStyle: 'full',
        timeStyle: 'medium'
      }
    }
  }
})
