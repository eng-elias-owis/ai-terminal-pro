//go:build windows

package terminal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// ShellInfo represents a detected shell
type ShellInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"` // cmd, powershell, pwsh, wsl, gitbash
	Description string `json:"description"`
}

// PTYSession represents a terminal session on Windows
type PTYSession struct {
	cmd         *exec.Cmd
	shell       string
	osType      string
	stdin       io.WriteCloser
	stdout      io.ReadCloser
	stderr      io.ReadCloser
	outputMutex sync.Mutex
	outputBuf   []byte
}

// NewPTYSession creates a new terminal session
func NewPTYSession() (*PTYSession, error) {
	shell, err := detectShell()
	if err != nil {
		return nil, fmt.Errorf("failed to detect shell: %w", err)
	}

	session := &PTYSession{
		shell:  shell,
		osType: runtime.GOOS,
	}

	if err := session.start(); err != nil {
		return nil, err
	}

	return session, nil
}

// NewPTYSessionWithShell creates a session with a specific shell
func NewPTYSessionWithShell(shellPath string) (*PTYSession, error) {
	session := &PTYSession{
		shell:  filepath.Base(shellPath),
		osType: runtime.GOOS,
	}

	if err := session.startWithPath(shellPath); err != nil {
		return nil, err
	}

	return session, nil
}

// detectShell determines the appropriate shell
func detectShell() (string, error) {
	// Check for PowerShell 7, then 5, then cmd
	if _, err := exec.LookPath("pwsh"); err == nil {
		return "pwsh", nil
	}
	if _, err := exec.LookPath("powershell"); err == nil {
		return "powershell", nil
	}
	return "cmd", nil
}

// GetAvailableShells returns all installed shells on the system
func GetAvailableShells() []ShellInfo {
	var shells []ShellInfo

	// Check for PowerShell 7
	if path, err := exec.LookPath("pwsh"); err == nil {
		shells = append(shells, ShellInfo{
			Name:        "PowerShell 7",
			Path:        path,
			Type:        "pwsh",
			Description: "Cross-platform PowerShell",
		})
	}

	// Check for Windows PowerShell (5.1)
	if path, err := exec.LookPath("powershell"); err == nil {
		shells = append(shells, ShellInfo{
			Name:        "Windows PowerShell",
			Path:        path,
			Type:        "powershell",
			Description: "Built-in Windows PowerShell",
		})
	}

	// Check for Command Prompt
	if path, err := exec.LookPath("cmd"); err == nil {
		shells = append(shells, ShellInfo{
			Name:        "Command Prompt",
			Path:        path,
			Type:        "cmd",
			Description: "Classic Windows CMD",
		})
	}

	// Check for WSL
	if path, err := exec.LookPath("wsl"); err == nil {
		shells = append(shells, ShellInfo{
			Name:        "WSL (Ubuntu)",
			Path:        path,
			Type:        "wsl",
			Description: "Windows Subsystem for Linux",
		})
	}

	// Check for Git Bash (common on developer machines)
	gitPaths := []string{
		`C:\Program Files\Git\bin\bash.exe`,
		`C:\Program Files (x86)\Git\bin\bash.exe`,
	}
	for _, path := range gitPaths {
		if _, err := os.Stat(path); err == nil {
			shells = append(shells, ShellInfo{
				Name:        "Git Bash",
				Path:        path,
				Type:        "gitbash",
				Description: "Git for Windows Bash",
			})
			break
		}
	}

	return shells
}

// start initializes the shell process with default detection
func (s *PTYSession) start() error {
	switch s.shell {
	case "pwsh":
		return s.startWithPath("pwsh")
	case "powershell":
		return s.startWithPath("powershell")
	default:
		return s.startWithPath("cmd")
	}
}

// startWithPath initializes the shell process with specific path
func (s *PTYSession) startWithPath(shellPath string) error {
	var cmd *exec.Cmd
	shellType := strings.ToLower(filepath.Base(shellPath))

	switch {
	case strings.Contains(shellType, "pwsh"):
		// PowerShell 7 - use interactive mode
		cmd = exec.Command(shellPath, "-NoLogo", "-NoExit")
	case strings.Contains(shellType, "powershell"):
		// Windows PowerShell 5.1 - use interactive mode
		cmd = exec.Command(shellPath, "-NoLogo", "-NoExit")
	case strings.Contains(shellType, "wsl"):
		// WSL - start bash
		cmd = exec.Command(shellPath, "bash", "-l")
	case strings.Contains(shellType, "bash") || strings.Contains(shellType, "git"):
		// Git Bash or other bash
		cmd = exec.Command(shellPath, "-l", "-i")
	default:
		// Command Prompt - simple invocation
		cmd = exec.Command(shellPath)
	}

	// Set up environment
	cmd.Env = os.Environ()

	// Create pipes for stdin/stdout/stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	s.cmd = cmd
	s.stdin = stdin
	s.stdout = stdout
	s.stderr = stderr

	return nil
}

// Write sends input to the shell
func (s *PTYSession) Write(data []byte) (int, error) {
	return s.stdin.Write(data)
}

// Read reads output from the shell
func (s *PTYSession) Read(p []byte) (int, error) {
	return s.stdout.Read(p)
}

// Resize is a no-op on Windows (would require ConPTY)
func (s *PTYSession) Resize(rows, cols int) error {
	// Windows resize not implemented - would need ConPTY
	return nil
}

// Close terminates the shell
func (s *PTYSession) Close() error {
	if s.stdin != nil {
		s.stdin.Close()
	}
	if s.stdout != nil {
		s.stdout.Close()
	}
	if s.stderr != nil {
		s.stderr.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	return nil
}

// GetShell returns the detected shell name
func (s *PTYSession) GetShell() string {
	return s.shell
}

// GetOSType returns the operating system
func (s *PTYSession) GetOSType() string {
	return s.osType
}

// IsWindows returns true
func (s *PTYSession) IsWindows() bool {
	return true
}
