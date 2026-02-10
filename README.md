# AI Terminal Pro

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![React Version](https://img.shields.io/badge/react-18-blue.svg)](https://reactjs.org)
[![Wails](https://img.shields.io/badge/wails-v2.11-purple.svg)](https://wails.io)

**AI-Enhanced Cross-Platform Terminal Application**

Built with Go + Wails + React, powered by fine-tuned Qwen3-0.6B model via LiteLLM proxy.

## 🎯 Overview

AI Terminal Pro transforms traditional terminal usage by adding AI-powered command generation. Users can describe what they want to do in natural language, and the AI generates the appropriate command for their specific operating system and shell.

### Key Features

- 🤖 **AI Command Generation** - Natural language to terminal commands
- 🖥️ **Cross-Platform** - Windows, macOS, Linux support
- 🔒 **Security First** - Virtual key system, command validation
- 💰 **Cost Optimized** - Scale-to-zero, hybrid online/offline mode
- 📊 **Full Observability** - Cost tracking, metrics, alerting
- ⚡ **Production Ready** - Enterprise-grade architecture

## 🏗️ Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Go Terminal    │────▶│  LiteLLM Proxy  │────▶│  HF Dedicated   │
│  App (Wails)    │     │  (HF Spaces)    │     │  Endpoint       │
└─────────────────┘     └─────────────────┘     └─────────────────┘
       │                       │
       │              ┌────────▼────────┐
       │              │  Supabase DB    │
       │              │  (PostgreSQL)   │
       │              └─────────────────┘
       │
┌──────▼────────┐
│  Local GGUF    │
│  (Fallback)    │
└─────────────────┘
```

### Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Backend** | Go 1.21+ + Wails v2 | Native app, PTY management |
| **Frontend** | React 18 + TypeScript | UI components |
| **Terminal** | xterm.js | Terminal emulator |
| **Styling** | Tailwind CSS | Modern styling |
| **AI Proxy** | LiteLLM | Unified API gateway |
| **Database** | Supabase (PostgreSQL) | Cost tracking |
| **Model** | HuggingFace Dedicated | Fine-tuned Qwen3-0.6B |

## 🚀 Quick Start

### Prerequisites

- **Go 1.21** or later
- **Node.js 18** or later
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Linux Desktop Dependencies** (for Linux builds):
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev
  
  # Fedora
  sudo dnf install gtk3-devel webkit2gtk3-devel
  
  # Arch
  sudo pacman -S gtk3 webkit2gtk
  ```

### Installation

```bash
# Clone repository
git clone https://github.com/EngEliasOwis/ai-terminal-pro.git
cd ai-terminal-pro

# Install Go dependencies
go mod tidy

# Install Node dependencies
cd frontend && npm install
cd ..

# Run in development mode
wails dev

# Build for production
wails build
```

### Alternative: CLI-Only Version

For headless environments or testing without GUI:

```bash
# Build CLI version
go build -o ai-terminal-pro .

# Run tests
./ai-terminal-pro -test

# Setup configuration
./ai-terminal-pro -setup
```

## 📁 Project Structure

```
ai-terminal-pro/
├── go.mod              # Go module (root level)
├── go.sum              # Go dependencies
├── main.go             # Wails entry point
├── app.go              # Wails app lifecycle
├── wails.json          # Wails configuration
│
├── ai/                 # LiteLLM AI client
│   └── client.go
├── config/             # Settings management
│   └── settings.go
├── security/           # Command validation
│   └── validator.go
├── terminal/           # PTY management
│   └── pty.go
│
├── frontend/           # React + TypeScript
│   ├── index.html      # Vite entry point
│   ├── package.json    # NPM dependencies
│   ├── src/            # React components
│   │   ├── App.tsx
│   │   ├── components/
│   │   └── styles/
│   └── dist/           # Build output (embedded)
│
├── deploy/             # Deployment configs
│   └── litellm/        # HF Spaces setup
│       ├── config.yaml
│       ├── Dockerfile
│       └── README.md
│
├── docs/               # Documentation
│   ├── architecture.md
│   └── security.md
│
├── tests/              # Test suites
├── BUILD.md            # Detailed build instructions
└── README.md           # This file
```

## 🚀 Deployment

### 1. Deploy LiteLLM Proxy to HuggingFace Spaces

```bash
cd deploy/litellm
# Follow README.md for HF Spaces deployment
```

**Required Environment Variables:**
- `HF_TOKEN` - Your HuggingFace token
- `LITELLM_MASTER_KEY` - Admin key for dashboard
- `DATABASE_URL` - Supabase PostgreSQL connection

### 2. Create Supabase Database

1. Sign up at https://supabase.com
2. Create new project (free tier - 500MB storage)
3. Get connection string from Settings → Database
4. Add to HF Spaces secrets

### 3. Configure HuggingFace Dedicated Endpoint

1. Deploy your Qwen3-0.6B fine-tuned model
2. Enable scale-to-zero for cost savings
3. Copy endpoint URL to HF Spaces secrets

## 📖 Documentation

- **[BUILD.md](BUILD.md)** - Detailed build instructions and requirements
- **[docs/architecture.md](docs/architecture.md)** - System architecture overview
- **[docs/security.md](docs/security.md)** - Security model and access control
- **[deploy/litellm/README.md](deploy/litellm/README.md)** - Deployment guide

## 🔒 Security

### Virtual Key System

Unlike raw HF tokens, our virtual keys provide:
- ✅ **Revocable** - Disable instantly if compromised
- ✅ **Scoped** - Limited to specific models
- ✅ **Budget-limited** - Prevent overage
- ✅ **Rate-limited** - Prevent abuse
- ✅ **Trackable** - Per-key analytics

### Command Validation

All AI-generated commands pass through security validation:
- **Risk Classification** (None, Low, Medium, High, Critical)
- **Pattern Detection** - Blocks dangerous commands
- **User Confirmation** - Required for risky operations

## 💰 Cost Analysis

| Component | Service | Monthly Cost |
|-----------|---------|--------------|
| LiteLLM Proxy | HF Spaces | $0 |
| Database | Supabase | $0 |
| Model Hosting | HF Dedicated | $3-60* |
| Local Fallback | llama.cpp | $0 |
| **Total** | | **$3-60** |

*Cost varies by usage (scale-to-zero saves 50-90%)

## 🛠️ Development

### Build Commands

```bash
# Development (with hot reload)
wails dev

# Build for current platform
wails build

# Build for specific platform
wails build -platform linux/amd64
wails build -platform windows/amd64
wails build -platform darwin/amd64

# Build CLI version only
go build .
```

### Testing

```bash
# Run component tests
./ai-terminal-pro -test

# Test CLI mode
./ai-terminal-pro -generate "list all files"
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **Ready Tensor** - LLMED Certification Program
- **HuggingFace** - Model hosting and inference
- **LiteLLM** - AI proxy and cost management
- **Wails** - Go + Webview framework
- **xterm.js** - Terminal emulator

## 📧 Contact

- **Author:** Elias Owis
- **Email:** elias@engelias.website
- **Project:** https://github.com/eng-elias-owis/ai-terminal-pro

---

Built with ❤️ by Elias Owis (EClaw)

**Module 2 Project - LLM Engineering & Deployment Certification**
