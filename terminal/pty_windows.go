//go:build windows

package terminal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// PTYSession represents a terminal session on Windows
type PTYSession struct {
	cmd    *exec.Cmd
	shell  string
	osType string
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
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

// start initializes the shell process
func (s *PTYSession) start() error {
	var cmd *exec.Cmd

	switch s.shell {
	case "pwsh":
		cmd = exec.Command("pwsh", "-NoExit", "-Command", "$OutputEncoding = [Console]::OutputEncoding = [System.Text.Encoding]::UTF8")
	case "powershell":
		cmd = exec.Command("powershell", "-NoExit", "-Command", "$OutputEncoding = [Console]::OutputEncoding = [System.Text.Encoding]::UTF8")
	default:
		cmd = exec.Command("cmd.exe")
	}

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

	cmd.Env = os.Environ()

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
