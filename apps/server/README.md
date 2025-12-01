# Zik AI Server

OpenAI-compatible AI server powered by GLM-4.6 via Z.AI reverse-engineered API.

## Features

- **OpenAI-Compatible API** - Drop-in replacement for OpenAI chat completions
- **Streaming Support** - Real-time SSE streaming responses
- **Free GLM-4.6 Access** - Unlimited free access to GLM-4.6 model
- **Structured Logging** - Production-ready logging with zerolog
- **Clean Architecture** - Maintainable and testable codebase
- **Token Counting** - Accurate usage tracking with tiktoken
- **YAML Configuration** - Flexible configuration with ENV overrides

## Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/zarazaex/zik.git
cd zik/apps/server

# Install dependencies
go mod download

# Build
make build
```

### Configuration

Create a `.env` file or use environment variables:

```bash
# Server configuration
PORT=8080
DEBUG=false

# Z.AI API (optional, defaults to anonymous mode)
ZAI_TOKEN=your_token_here

# Model configuration
MODEL=GLM-4-6-API-V1
THINK_MODE=reasoning  # Options: reasoning, think, strip, details
```

Or use a YAML config file:

```bash
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml
export CONFIG_PATH=configs/config.yaml
```

### Running

```bash
# Development mode
make run

# Production mode
./bin/zik-server
```

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### List Models
```bash
curl http://localhost:8080/v1/models
```

### Chat Completions

**Non-streaming:**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": false
  }'
```

**Streaming:**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [{"role": "user", "content": "Count to 5"}],
    "stream": true
  }'
```

## Development

### Project Structure

```
apps/server/
├── cmd/
│   └── zik-server/          # Main entry point
├── internal/
│   ├── api/                 # HTTP handlers & routing
│   │   ├── handlers/        # Request handlers
│   │   └── middleware/      # HTTP middleware
│   ├── service/             # Business logic
│   │   ├── ai/              # Z.AI client
│   │   └── auth/            # Authentication
│   ├── domain/              # Domain models
│   ├── config/              # Configuration
│   └── pkg/                 # Internal utilities
│       ├── crypto/          # Signature generation
│       ├── logger/          # Structured logging
│       ├── utils/           # Helpers
│       └── validator/       # Request validation
├── configs/                 # Configuration files
├── test/                    # Integration tests
├── Makefile                 # Build commands
└── README.md
```

### Commands

```bash
make help           # Show all available commands
make build          # Build binary
make run            # Run server
make test           # Run all tests
make cover          # Generate coverage report
make lint           # Run linters
make clean          # Clean build artifacts
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make cover

# Run linters
make lint
```

## Architecture

The server follows **Clean Architecture** principles:

- **Domain Layer**: Core business models and interfaces
- **Service Layer**: Business logic and external integrations
- **API Layer**: HTTP handlers and routing
- **Infrastructure**: Configuration, logging, utilities

## License

BSD License - See LICENSE file for details

## Author

**zarazaex** - [GitHub](https://github.com/zarazaex)

## Acknowledgments

- Powered by [Z.AI](https://chat.z.ai) reverse-engineered API
- GLM-4.6 model by Zhipu AI
