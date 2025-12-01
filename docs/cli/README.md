# ZIK CLI - AI Tools for Developers

Command-line interface for ZIK AI assistant. Generate commit messages, chat with AI, and get instant coding help directly from your terminal.

## Features

- **Smart Commit Messages** - Generate conventional commit messages from git diff
- **Quick Questions** - Ask one-off questions without context
- **Interactive Chat** - Multi-turn conversations with context (coming soon)
- **Code Analysis** - Review and explain code (coming soon)

## Installation

### Quick Install

```bash
curl -fsSL https://zik.zarazaex.xyz/install | bash
```

### From Source

```bash
cd apps/cli
make build
sudo make install
```

### Manual Installation

Download the latest release from [GitHub Releases](https://github.com/zarazaex69/zik/releases) and place it in your PATH.

## Quick Start

### Generate Commit Message

```bash
# Stage your changes
git add .

# Generate commit message
zik commit

# Auto-apply commit
zik commit --apply

# Prefer specific commit type
zik commit --type feat
```

Interactive workflow:
- `Y` - Accept and commit
- `N` - Cancel
- `R` - Regenerate message

### Ask Questions

```bash
# Quick question
zik ask "How do I reverse a string in Go?"

# Streaming response (default)
zik ask "Explain async/await in JavaScript"

# Non-streaming
zik ask --stream=false "What is a closure?"
```

### Configuration

```bash
# View current config
zik config list

# Edit config file
vim ~/.config/zik/config.yaml
```

## Configuration

Config file location: `~/.config/zik/config.yaml`

```yaml
model: GLM-4-6-API-V1
temperature: 0.7
max_tokens: 2000
streaming: true

commit:
  conventional_commits: true
  preferred_type: feat
  auto_stage: false

chat:
  save_history: true
  history_limit: 100
  timeout: 30s
```

## Environment Variables

- `ZIK_API_URL` - Override API endpoint (for development only)

Example:
```bash
ZIK_API_URL=http://localhost:8802 zik ask "test"
```

## Commands

### `zik commit`

Generate conventional commit messages from git diff.

**Flags:**
- `-s, --staged` - Analyze staged changes only (default)
- `-a, --all` - Analyze all changes (staged + unstaged)
- `-y, --apply` - Auto-apply without confirmation
- `-t, --type` - Preferred commit type (feat, fix, docs, etc.)

**Examples:**
```bash
zik commit                    # Interactive mode
zik commit --apply            # Auto-apply
zik commit --type fix         # Prefer 'fix' type
zik commit --all              # Include unstaged changes
```

### `zik ask`

Ask a quick question to AI.

**Flags:**
- `-s, --stream` - Stream response in real-time (default: true)

**Examples:**
```bash
zik ask "What is the difference between let and const?"
zik ask "How do I reverse a string in Go?"
zik ask --stream=false "Explain closures"
```

### `zik chat`

Start an interactive chat session (coming soon).

### `zik code`

Code analysis commands (coming soon).

**Subcommands:**
- `zik code review` - Review code changes
- `zik code explain <file>` - Explain code in a file

### `zik config`

Manage configuration.

**Subcommands:**
- `zik config list` - Show current configuration
- `zik config get <key>` - Get a config value
- `zik config set <key> <value>` - Set a config value

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Install Locally

```bash
make install
```

### Clean

```bash
make clean
```

## Architecture

```
apps/cli/
├── cmd/zik/              # Command implementations
│   ├── main.go           # Entry point
│   ├── commit.go         # Commit command
│   ├── ask.go            # Ask command
│   ├── chat.go           # Chat command
│   ├── code.go           # Code commands
│   └── config.go         # Config command
├── internal/
│   ├── ai/               # AI client
│   │   ├── client.go     # HTTP client
│   │   └── stream.go     # SSE streaming
│   ├── git/              # Git operations
│   │   └── client.go     # Git commands
│   ├── config/           # Configuration
│   │   ├── config.go     # Config management
│   │   └── constants.go  # Hardcoded constants
│   └── prompt/           # Prompt templates
│       └── commit.go     # Commit prompts
└── Makefile
```

## Security

The API endpoint (`https://zik-api.zarazaex.xyz`) is hardcoded at compile time to prevent unauthorized API usage. For development, use the `ZIK_API_URL` environment variable.

## License

BSD License - See LICENSE file for details

## Author

**zarazaex** - [GitHub](https://github.com/zarazaex)
