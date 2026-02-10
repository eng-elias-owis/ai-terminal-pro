#!/bin/bash
# Build script for Linux/macOS

set -e

echo "Building AI Terminal Pro..."

# Build frontend
cd frontend
npm install
npm run build
cd ..

# Build Go application
wails build

echo "Build complete! Check build/bin/ directory."
