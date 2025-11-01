# Siraaj Documentation

Official documentation for Siraaj Analytics - Fast, Simple, Self-Hosted Analytics.

## Getting Started

### Development

```bash
# Install dependencies
pnpm install

# Start dev server
pnpm dev

# Visit http://localhost:5173
```

### Build

```bash
# Build for production
pnpm build

# Preview production build
pnpm preview
```

## Structure

```
docs/
├── .vitepress/
│   └── config.js          # VitePress configuration
├── index.md               # Landing page
├── guide/
│   ├── introduction.md
│   ├── quick-start.md
│   ├── installation.md
│   ├── configuration.md
│   ├── architecture.md
│   ├── dashboard.md
│   └── ...
├── sdk/
│   ├── overview.md
│   ├── vanilla.md
│   ├── react.md
│   ├── vue.md
│   ├── svelte.md
│   ├── nextjs.md
│   └── ...
└── api/
    ├── overview.md
    └── ...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally with `pnpm dev`
5. Submit a pull request

## License

MIT © Mohamed Elhefni
