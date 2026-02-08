# AI Terminal Pro

AI-Enhanced Cross-Platform Terminal Application

## Overview

AI Terminal Pro is a production-grade terminal application that combines traditional terminal functionality with AI-powered command generation using a fine-tuned Qwen3-0.6B model.

## Architecture

```
┌─────────────────┐
│   Go + Wails    │
│  Terminal App   │
└────────┬────────┘
         │ Virtual Key
         ▼
┌─────────────────┐
│  LiteLLM Proxy  │
│   (HF Spaces)   │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌──────┐  ┌──────┐
│  HF  │  │ Local│
│Endpoint│ │ GGUF │
└──────┘  └──────┘
```

## Tech Stack

- **Backend:** Go + Wails v2
- **Frontend:** React + TypeScript + xterm.js
- **AI Proxy:** LiteLLM (HuggingFace Spaces)
- **Database:** Supabase (PostgreSQL)
- **Model Hosting:** HuggingFace Dedicated Endpoint

## Development

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/ai-terminal-pro.git
cd ai-terminal-pro

# Install dependencies
cd backend && go mod tidy
cd ../frontend && npm install

# Run in development mode
wails dev

# Build for production
wails build
```

## Project Structure

- `backend/` - Go backend with PTY management, AI client, security
- `frontend/` - React frontend with xterm.js terminal
- `deploy/` - Deployment configurations (LiteLLM, model)
- `docs/` - Documentation
- `tests/` - Test suites

## License

MIT License
