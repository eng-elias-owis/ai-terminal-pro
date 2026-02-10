# Troubleshooting Guide

## Common Issues

### 1. "no index.html could be found in your Assets fs.FS"

**Cause:** The frontend hasn't been built before running `wails dev` or `wails build`.

**Solution:**

```bash
# Step 1: Build the frontend first
cd frontend
npm install
npm run build
cd ..

# Step 2: Now run wails
wails dev
# or
wails build
```

**Alternative - Use build scripts:**

**Windows:**
```cmd
build-windows.bat
```

**Linux/macOS:**
```bash
make all
# or
./build.sh
```

---

### 2. "can't be reached" when opening browser link

**Cause:** The Wails dev server requires the app to be running in the background.

**Solution:**
- Keep `wails dev` running in the terminal
- The browser link (http://localhost:34115) only works while the terminal window is open
- For testing, open the link in a browser while `wails dev` is still running

**Note:** Wails creates a desktop application, not just a web server. The dev server is for development convenience only.

---

### 3. "Nothing happens" when opening the .exe

**Cause:** On Windows, the application might be missing WebView2 runtime.

**Solution:**

1. **Install WebView2 Runtime:**
   Download and install from: https://developer.microsoft.com/en-us/microsoft-edge/webview2/

2. **Check Windows Defender:**
   - The .exe might be blocked by Windows Defender
   - Right-click the .exe → Properties → Check "Unblock" if present

3. **Run from Command Line:**
   ```cmd
   cd build\bin
   ai-terminal-pro.exe
   ```
   This will show any error messages.

---

### 4. Build fails on Windows with Cygwin/Git Bash

**Cause:** Path resolution issues between Windows and Unix-style paths.

**Solution:**

**Option A: Use Command Prompt or PowerShell**
```cmd
cd ai-terminal-pro
build-windows.bat
```

**Option B: Use WSL (Windows Subsystem for Linux)**
```bash
cd /mnt/e/path/to/ai-terminal-pro
make all
```

---

### 5. "go.mod file not found" error

**Cause:** Running `wails dev` from the wrong directory.

**Solution:**
```bash
# Must be in the project root where go.mod is located
cd ai-terminal-pro
ls go.mod  # Should show the file
wails dev
```

---

### 6. Frontend changes not reflecting

**Cause:** Vite dev server and Wails dev server are separate.

**Solution:**

**For development with hot reload:**
```bash
# Terminal 1: Start Go backend
wails dev

# Terminal 2: Start frontend dev server
cd frontend
npm run dev
```

Or use the combined approach:
```bash
# Build frontend first, then run wails
cd frontend && npm run build && cd ..
wails dev
```

---

### 7. "Embedded file not found" on build

**Cause:** The embed directive can't find the files.

**Solution:**

Ensure frontend/dist exists with files:
```bash
ls frontend/dist/
# Should show: index.html and assets/
```

If empty, rebuild:
```bash
cd frontend
npm run build
cd ..
```

Then rebuild the Go binary:
```bash
go build .
```

---

## Platform-Specific Issues

### Windows

**Requirements:**
- Go 1.21+
- Node.js 18+
- WebView2 Runtime (usually pre-installed on Windows 11)
- Visual Studio Build Tools or MinGW-w64 (for CGO)

**Common Errors:**

1. **"gcc not found"**
   - Install MinGW-w64: https://www.mingw-w64.org/
   - Or install Visual Studio Build Tools

2. **"webview2loader.dll not found"**
   - Install WebView2 Runtime
   - Or add `-tags native_webview2loader` to wails build

### Linux

**Requirements:**
- Go 1.21+
- Node.js 18+
- GTK3 development libraries
- WebKit2GTK development libraries

**Install dependencies:**
```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk3-devel

# Arch
sudo pacman -S gtk3 webkit2gtk
```

### macOS

**Requirements:**
- Go 1.21+
- Node.js 18+
- Xcode Command Line Tools

**Install dependencies:**
```bash
xcode-select --install
```

---

## Still Having Issues?

1. **Clean build:**
   ```bash
   make clean
   make all
   ```

2. **Check versions:**
   ```bash
   go version
   node --version
   wails version
   ```

3. **Verify project structure:**
   ```bash
   ls -la
   # Should see: go.mod, main.go, wails.json, frontend/
   ```

4. **Run in verbose mode:**
   ```bash
   wails build -v 2
   ```

5. **Check GitHub Issues:**
   https://github.com/wailsapp/wails/issues

---

## Quick Reference: Correct Build Order

```bash
# 1. Clone and enter directory
git clone https://github.com/eng-elias-owis/ai-terminal-pro.git
cd ai-terminal-pro

# 2. Install dependencies
cd frontend && npm install && cd ..
go mod tidy

# 3. Build frontend
cd frontend && npm run build && cd ..

# 4. Build application
wails build

# 5. Run the binary
./build/bin/ai-terminal-pro  # Linux/macOS
.\build\bin\ai-terminal-pro.exe  # Windows
```

Or use the provided scripts:
- **Windows:** `build-windows.bat`
- **Linux/macOS:** `make all` or `./build.sh`
