# OpenCode Integration Guide

## Overview

OpenCode is an AI coding agent that can help develop the AI Terminal Pro project. This guide explains how to set it up and use it effectively.

## Installation

### Option 1: Using go install
```bash
go install github.com/opencode-ai/opencode@latest
```

### Option 2: Pre-built binaries
Download from: https://github.com/opencode-ai/opencode/releases

### Option 3: Homebrew (macOS/Linux)
```bash
brew install opencode-ai/tap/opencode
```

## Verification

Check if OpenCode is installed:
```bash
opencode --version
# Should show: 1.1.51 or later
```

## Configuration

The project includes `.opencode.yaml` with:
- Project name and description
- Model configuration (kimi-k2.5)
- Custom prompt for AI Terminal Pro context
- Build/test commands
- File patterns to include/exclude

## Usage

### Interactive TUI Mode
```bash
# Navigate to project
cd ai-terminal-pro

# Start OpenCode interactive mode
opencode

# Or with specific model
opencode --model kimi-k2.5
```

### One-shot Commands
```bash
# Run a specific task
opencode run "fix the terminal I/O issue"

# Run with specific model
opencode run "add shell selection feature" --model kimi-k2.5

# Continue last session
opencode run "continue from where we left off" --continue
```

### Web Interface
```bash
# Start web server and open browser
opencode web

# Then access at http://localhost:8080
```

### Server Mode (Headless)
```bash
# Start server
opencode serve --port 8080

# Attach from another terminal
opencode attach http://localhost:8080
```

## Available Commands in Config

The `.opencode.yaml` defines these commands:

```bash
# Build the project
opencode run "/build"

# Run tests
opencode run "/test"

# Development mode
opencode run "/dev"

# Production build
opencode run "/build-wails"
```

## Best Practices

### 1. Start with Context
```bash
opencode run "review the terminal/pty_windows.go file and suggest improvements"
```

### 2. Incremental Changes
```bash
opencode run "fix the PowerShell execution issue - the command runs but doesn't show output"
```

### 3. Testing Changes
```bash
opencode run "build the project and verify it compiles without errors"
```

### 4. Cross-Platform Awareness
```bash
opencode run "ensure the shell detection works on both Windows and Linux"
```

## Troubleshooting

### Issue: "Failed to change directory"
**Cause:** Running opencode with wrong arguments
**Fix:** Use `opencode run "message"` not `opencode version`

### Issue: Slow or hanging
**Cause:** TUI mode waiting for input
**Fix:** Use `opencode run "message"` for one-shot commands

### Issue: Model not found
**Cause:** No API key configured
**Fix:** 
```bash
opencode auth
# Follow prompts to add API key
```

### Issue: Build failures
**Cause:** Missing dependencies
**Fix:**
```bash
# Install Go dependencies
go mod tidy

# Install Node dependencies
cd frontend && npm install
```

## API Key Setup

### OpenAI
```bash
export OPENAI_API_KEY="your-key-here"
```

### Anthropic (Claude)
```bash
export ANTHROPIC_API_KEY="your-key-here"
```

### Other Providers
See: https://docs.opencode.ai/authentication

## Session Management

### List Sessions
```bash
opencode session list
```

### Continue Session
```bash
opencode --continue
# or
opencode --session <session-id>
```

### Export Session
```bash
opencode export <session-id> > session-backup.json
```

## Integration with This Project

The `.opencode.yaml` is configured with:
- **Project Context:** Wails v2, Go + React, terminal emulation
- **Model:** kimi-k2.5 (change with `--model` flag)
- **Build Commands:** Pre-configured for this project structure
- **File Patterns:** Watches .go, .tsx, .ts, .json, .yaml, .md
- **Ignore:** node_modules, dist, build, binaries

## Example Workflows

### Add a New Feature
```bash
opencode run "add a settings panel to change terminal font size in the frontend"
```

### Fix a Bug
```bash
opencode run "fix issue where terminal doesn't scroll to bottom on new output"
```

### Refactor Code
```bash
opencode run "refactor the terminal I/O handling to use a cleaner interface"
```

### Cross-Platform Fix
```bash
opencode run "ensure the shell detection works correctly on Windows PowerShell, CMD, and WSL"
```

## Advanced Usage

### Custom Prompts
Edit `.opencode.yaml` to modify the system prompt for different contexts.

### Multiple Models
```bash
# Use Claude for complex reasoning
opencode run "design the architecture" --model claude-3.5-sonnet

# Use GPT-4 for code generation
opencode run "implement the feature" --model gpt-4

# Use local model
opencode run "quick edit" --model ollama/llama3
```

### GitHub Integration
```bash
# Checkout and work on PR
opencode pr 123

# Use GitHub agent
opencode github
```

## Resources

- **Documentation:** https://docs.opencode.ai
- **GitHub:** https://github.com/opencode-ai/opencode
- **Discord:** https://discord.gg/opencode

## Current Status

✅ **OpenCode is installed** (v1.1.51)
✅ **Configuration file created** (.opencode.yaml)
✅ **Ready to use** for future development

**Next Steps:**
1. Configure API key: `opencode auth`
2. Start using: `opencode run "your task"`
3. Or use TUI: `opencode` (interactive mode)
