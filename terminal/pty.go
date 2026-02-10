package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/creack/pty"
)

// PTYSession represents a pseudo-terminal session
type PTYSession struct {
	PTY    *os.File
	cmd    *exec.Cmd
	shell  string
	osType string
}

// NewPTYSession creates a new PTY session with appropriate shell
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

// detectShell determines the appropriate shell for the current OS
func detectShell() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Check for PowerShell 7, then 5, then cmd
		if _, err := exec.LookPath("pwsh"); err == nil {
			return "pwsh", nil
		}
		if _, err := exec.LookPath("powershell"); err == nil {
			return "powershell", nil
		}
		return "cmd", nil

	case "darwin", "linux":
		// Check SHELL environment variable
		if shell := os.Getenv("SHELL"); shell != "" {
			return shell, nil
		}
		// Try common shells
		for _, sh := range []string{"/bin/zsh", "/bin/bash", "/bin/sh"} {
			if _, err := os.Stat(sh); err == nil {
				return sh, nil
			}
		}
		return "", fmt.Errorf("no shell found")

	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// start initializes the PTY and shell process
func (s *PTYSession) start() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// On Windows, use cmd.exe with special flags for better compatibility
		if s.shell == "cmd" {
			cmd = exec.Command("cmd.exe", "/Q", "/K", "prompt $P$G")
		} else {
			cmd = exec.Command(s.shell, "-NoExit", "-Command", "[Console]::OutputEncoding = [System.Text.Encoding]::UTF8")
		}
	default:
		cmd = exec.Command(s.shell, "-l") // Login shell
	}

	// Set up environment
	cmd.Env = os.Environ()

	// Create PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}

	s.PTY = ptmx
	s.cmd = cmd

	return nil
}

// Write sends input to the PTY
func (s *PTYSession) Write(data []byte) (int, error) {
	return s.PTY.Write(data)
}

// Read reads output from the PTY
func (s *PTYSession) Read(p []byte) (int, error) {
	return s.PTY.Read(p)
}

// Resize updates the terminal size
func (s *PTYSession) Resize(rows, cols int) error {
	return pty.Setsize(s.PTY, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

// Close terminates the PTY session
func (s *PTYSession) Close() error {
	if err := s.PTY.Close(); err != nil {
		return err
	}
	return s.cmd.Wait()
}

// GetShell returns the detected shell name
func (s *PTYSession) GetShell() string {
	return s.shell
}

// GetOSType returns the operating system
func (s *PTYSession) GetOSType() string {
	return s.osType
}

// IsWindows returns true if running on Windows
func (s *PTYSession) IsWindows() bool {
	return s.osType == "windows"
}
