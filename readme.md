
# ZIK - AI Tools For Developers

Powerful AI-powered tools for developers: CLI assistant, commit message generator, and OpenAI-compatible API.

![Logo-by-littlelynxuwu-forked-zarazaex](apps/web/assets/board.png)


## Features

- **CLI Tool** - Generate commit messages, ask questions, analyze code
- **API Server** - OpenAI-compatible API powered by GLM-4.6
- **Web UI** - Interactive interface for API exploration
- **Free & Unlimited** - No API keys required

## Quick Start

### Install CLI

```bash
curl -fsSL https://zik.zarazaex.xyz/install | bash
```

### Use CLI

```bash
# Generate commit message
zik commit

# Ask a question
zik ask "How do I reverse a string in Go?"

# Interactive chat (coming soon)
zik chat
```

## Documentation

- [CLI Documentation](docs/cli/README.md)
- [API Documentation](docs/server/README.md)
- [Deployment Guide](deployment/README.md)

## Project Structure

```
zik/
├── apps/
│   ├── cli/          # Command-line interface
│   ├── server/       # API server
│   └── web/          # Web UI
├── docs/             # Documentation
├── deployment/       # Docker deployment configs
└── .github/          # CI/CD workflows
```

## Development

### CLI

```bash
cd apps/cli
make build
./bin/zik --help
```

### Server

```bash
cd apps/server
make build
./bin/zik-server
```

### Web

```bash
cd apps/web
bun install
bun dev
```


## License

BSD License - See LICENSE file for details

## Author

**zarazaex** - [GitHub](https://github.com/zarazaex69)

## Links

- Website: [zik.zarazaex.xyz](https://zik.zarazaex.xyz)
- API: [zik-api.zarazaex.xyz](https://zik-api.zarazaex.xyz)
- GitHub: [github.com/zarazaex69/zik](https://github.com/zarazaex69/zik)
