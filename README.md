# AI Terminal Pro

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![React Version](https://img.shields.io/badge/react-18-blue.svg)](https://reactjs.org)

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

- Go 1.21 or later
- Node.js 18 or later
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Installation

```bash
# Clone repository
git clone https://github.com/EngEliasOwis/ai-terminal-pro.git
cd ai-terminal-pro

# Install Go dependencies
cd backend && go mod tidy
cd ..

# Install Node dependencies
cd frontend && npm install
cd ..

# Run in development mode
wails dev

# Build for production
wails build
```

### Deployment

1. **Deploy LiteLLM Proxy to HuggingFace Spaces:**
   ```bash
   cd deploy/litellm
   # Follow README.md for HF Spaces deployment
   ```

2. **Create Supabase Database:**
   - Sign up at https://supabase.com
   - Create new project (free tier)
   - Copy connection string to HF Spaces secrets

3. **Configure HuggingFace Dedicated Endpoint:**
   - Deploy your Qwen3-0.6B fine-tuned model
   - Get endpoint URL
   - Add to HF Spaces secrets

## 📖 Documentation

- [Architecture Overview](docs/architecture.md)
- [Security Model](docs/security.md)
- [Deployment Guide](deploy/litellm/README.md)

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

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **Ready Tensor** - LLMED Certification Program
- **HuggingFace** - Model hosting and inference
- **LiteLLM** - AI proxy and cost management
- **Wails** - Go + Webview framework
- **xterm.js** - Terminal emulator

---

Built with ❤️ by Elias Owis (EClaw)

Module 2 Project - LLM Engineering & Deployment Certification
