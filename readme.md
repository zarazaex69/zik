# ZIK - AI Tools For Developers

[![Release](https://img.shields.io/github/v/release/zarazaex69/zik?style=flat-square&logo=github&color=blue)](https://github.com/zarazaex69/zik/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/zarazaex69/zik/release.yml?style=flat-square&logo=github)](https://github.com/zarazaex69/zik/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zarazaex69/zik?filename=apps%2Fserver%2Fgo.mod&style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-BSD-green?style=flat-square)](LICENSE)

Powerful AI-powered tools for developers: CLI assistant, commit message generator, and OpenAI-compatible API.

![Logo-by-littlelynxuwu-forked-zarazaex](apps/web/assets/board.png)

## Quick Start

### Install CLI

```bash
curl -fsSL zik.zarazaex.xyz/install | bash
```

## Overview

ZIK is a comprehensive AI toolkit designed for developers, providing seamless access to GLM-4.6 model through multiple interfaces. Built with modern technologies and clean architecture principles, it offers production-ready solutions for AI-powered development workflows.

## Key Features

- **CLI Tool** - Generate commit messages, ask questions, analyze code
- **API Server** - OpenAI-compatible API powered by GLM-4.6
- **Web UI** - Interactive interface for API exploration
- **Free & Unlimited** - No API keys required
- **Streaming Support** - Real-time SSE responses
- **Cross-Platform** - Linux, macOS, Windows support

## Technology Stack

### Backend
- **Go 1.23** - High-performance server and CLI
- **Chi Router** - Lightweight HTTP routing
- **Zerolog** - Structured logging
- **Tiktoken** - Token counting

### Frontend
- **React 19** - Modern UI framework
- **Bun** - Fast JavaScript runtime
- **TailwindCSS 4** - Utility-first styling
- **Radix UI** - Accessible components

### Infrastructure
- **Docker** - Containerization
- **GitHub Actions** - CI/CD pipeline
- **YAML** - Configuration management

### Use CLI

```bash
# Generate commit message
zik commit

# Ask a question
zik ask "How do I reverse a string in Go?"

# Interactive chat (coming soon)
zik chat
```

### Run API Server

```bash
# Using Docker
docker pull ghcr.io/zarazaex69/zik-server:latest
docker run -p 8080:8080 ghcr.io/zarazaex69/zik-server:latest

# From source
cd apps/server
make build && ./bin/zik-server
```

### API Usage

```bash
curl -X POST https://zik-api.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": false
  }'
```

## Documentation

- [CLI Documentation](docs/cli/README.md)
- [API Documentation](docs/server/README.md)
- [Deployment Guide](deployment/README.md)

## Project Structure

```
zik/
├── apps/
│   ├── cli/          # Command-line interface (Go)
│   ├── server/       # API server (Go)
│   └── web/          # Web UI (React + Bun)
├── packages/
│   ├── prompts/      # Shared prompt templates
│   └── sdk/          # Client SDKs
├── docs/             # Documentation
├── deployment/       # Docker deployment configs
└── .github/          # CI/CD workflows
```

## Development

### Prerequisites

- Go 1.23+
- Bun 1.0+
- Docker (optional)
- Make

### CLI Development

```bash
cd apps/cli
make build
./bin/zik --help
```

### Server Development

```bash
cd apps/server
cp .env.example .env
make run
```

### Web Development

```bash
cd apps/web
bun install
bun dev
```

### Running Tests

```bash
# Server tests
cd apps/server
make test
make cover

# CLI tests
cd apps/cli
make test
```

## Deployment

### Docker Compose

```bash
cd deployment
cp .env.example .env
docker compose up -d
```

### Manual Deployment

See [Deployment Guide](deployment/README.md) for detailed instructions.

## API Compatibility

ZIK API is fully compatible with OpenAI's chat completions API. You can use it as a drop-in replacement:

```python
from openai import OpenAI

client = OpenAI(
    base_url="https://zik-api.zarazaex.xyz/v1",
    api_key="not-needed"
)

response = client.chat.completions.create(
    model="GLM-4-6-API-V1",
    messages=[{"role": "user", "content": "Hello!"}]
)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

BSD License - See LICENSE file for details

## Author

**zarazaex** - [GitHub](https://github.com/zarazaex69)

## Links

- Website: [zik.zarazaex.xyz](https://zik.zarazaex.xyz)
- API: [zik-api.zarazaex.xyz](https://zik-api.zarazaex.xyz)
- GitHub: [github.com/zarazaex69/zik](https://github.com/zarazaex69/zik)

## Acknowledgments

- Powered by [Z.AI](https://chat.z.ai) reverse-engineered API
- GLM-4.6 model by Zhipu AI
- Logo by littlelynxuwu
