#!/bin/bash
# OpenCode wrapper script for ai-terminal-pro

set -e

PROJECT_DIR="/root/.openclaw/workspace/ai-terminal-pro"

echo "🦅 AI Terminal Pro - OpenCode Integration"
echo "=========================================="

# Check if opencode is installed
if ! command -v opencode &> /dev/null; then
    echo "❌ OpenCode not found. Installing..."
    go install github.com/opencode-ai/opencode@latest
fi

# Navigate to project directory
cd "$PROJECT_DIR"

# Check for config
if [ ! -f .opencode.yaml ]; then
    echo "❌ .opencode.yaml not found in $PROJECT_DIR"
    exit 1
fi

echo "✅ OpenCode is ready to use!"
echo ""
echo "Usage examples:"
echo "  opencode run \"fix terminal I/O issue\""
echo "  opencode run \"add new feature\" --model kimi-k2.5"
echo "  opencode web  # Start web interface"
echo ""
echo "Available commands in .opencode.yaml:"
echo "  - build: Build frontend and Go backend"
echo "  - test: Run component tests"
echo "  - dev: Run wails dev"
echo "  - build-wails: Build production binary"
