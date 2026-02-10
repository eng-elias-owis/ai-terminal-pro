.PHONY: all frontend build dev clean test

all: frontend build

frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

build: frontend
	@echo "Building application..."
	wails build

dev: frontend
	@echo "Starting development server..."
	wails dev

test:
	@echo "Running tests..."
	go build -o ai-terminal-pro . && ./ai-terminal-pro -test

clean:
	@echo "Cleaning build files..."
	rm -rf build/bin frontend/dist ai-terminal-pro

setup:
	@echo "Setting up development environment..."
	cd frontend && npm install
	go mod tidy
