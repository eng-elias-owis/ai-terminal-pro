# AI Terminal Pro - Build Instructions

## Prerequisites

### For Development (All Platforms)
- Go 1.21+
- Node.js 18+
- Wails CLI v2.11.0+

### For Linux Desktop Build
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk3-devel

# Arch
sudo pacman -S gtk3 webkit2gtk
```

## Build Steps

### 1. Clone and Setup
```bash
git clone https://github.com/eng-elias-owis/ai-terminal-pro.git
cd ai-terminal-pro

# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd frontend && npm install
```

### 2. Development Mode
```bash
# Run development server with hot reload
wails dev
```

### 3. Production Build
```bash
# Build for current platform
wails build

# Build for specific platform
wails build -platform linux/amd64
wails build -platform windows/amd64
wails build -platform darwin/amd64
```

### 4. Build Output
- Binary location: `build/bin/ai-terminal-pro`
- Size: ~20-30MB (includes embedded frontend)

## Project Structure

```
ai-terminal-pro/
├── go.mod              # Go dependencies
├── main.go             # Wails entry point
├── app.go              # Wails app with lifecycle hooks
├── wails.json          # Wails configuration
│
├── ai/                 # LiteLLM AI client
├── config/             # Settings management
├── security/           # Command validation
├── terminal/           # PTY management
│
├── frontend/           # React + TypeScript
│   ├── index.html      # HTML entry point
│   ├── src/            # React components
│   └── dist/           # Build output (embedded)
│
└── deploy/             # Deployment configs
    └── litellm/        # HF Spaces setup
```

## Current Status

✅ **Working:**
- Go backend compiles successfully
- Frontend builds with Vite
- TypeScript type checking passes
- Wails bindings generate correctly

⚠️ **Requires Desktop Dependencies:**
- On headless servers, install `libwebkit2gtk-4.0-dev` for builds
- GUI requires display server (X11/Wayland)

## Testing CLI Version

The project also includes a CLI-only version for testing:

```bash
# Build CLI version
go build -o ai-terminal-cli .

# Run tests
./ai-terminal-cli -test

# Setup configuration
./ai-terminal-cli -setup
```
