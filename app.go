package main

import (
	"context"
	"fmt"
	"runtime"

	"ai-terminal-pro/ai"
	"ai-terminal-pro/config"
	"ai-terminal-pro/security"
	"ai-terminal-pro/terminal"
)

// App struct
type App struct {
	ctx       context.Context
	settings  *config.Settings
	validator *security.Validator
	client    *ai.Client
	terminal  *terminal.PTYSession
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	settings, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		settings = config.DefaultSettings()
	}
	a.settings = settings

	// Initialize security validator
	a.validator = security.NewValidator()

	// Initialize AI client if configured
	if settings.LiteLLMEndpoint != "" && settings.VirtualKey != "" {
		a.client = ai.NewClient(settings.LiteLLMEndpoint, settings.VirtualKey)
	}

	// Start terminal session
	ptySession, err := terminal.NewPTYSession()
	if err != nil {
		fmt.Printf("Failed to start terminal: %v\n", err)
		return
	}
	a.terminal = ptySession
}

// OnDomReady is called after front-end resources have been loaded
func (a *App) OnDomReady(ctx context.Context) {
	// Frontend is ready
}

// OnBeforeClose is called when the application is about to quit
func (a *App) OnBeforeClose(ctx context.Context) bool {
	// Clean up terminal session
	if a.terminal != nil {
		a.terminal.Close()
	}
	return false
}

// OnShutdown is called at application termination
func (a *App) OnShutdown(ctx context.Context) {
	// Cleanup resources
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetOS returns the current operating system
func (a *App) GetOS() string {
	return runtime.GOOS
}

// GenerateCommand generates a terminal command using AI
func (a *App) GenerateCommand(description string) (map[string]interface{}, error) {
	if a.client == nil {
		return nil, fmt.Errorf("AI client not configured")
	}

	aiCtx := ai.Context{
		OS:         a.settings.GetOSType(),
		Shell:      a.settings.GetShell(),
		WorkingDir: ".", // TODO: Get actual working directory
	}

	command, err := a.client.GenerateCommand(a.ctx, description, aiCtx)
	if err != nil {
		return nil, err
	}

	// Validate the command
	risk := a.validator.ValidateCommand(command)
	explanation := a.validator.GetExplanation(risk)

	return map[string]interface{}{
		"command":     command,
		"risk":        risk.String(),
		"explanation": explanation,
		"blocked":     risk == security.RiskCritical,
	}, nil
}

// ValidateCommand checks a command's safety level
func (a *App) ValidateCommand(command string) map[string]interface{} {
	risk := a.validator.ValidateCommand(command)
	explanation := a.validator.GetExplanation(risk)

	return map[string]interface{}{
		"command":     command,
		"risk":        risk.String(),
		"explanation": explanation,
		"blocked":     risk == security.RiskCritical,
	}
}

// GetSettings returns the current application settings
func (a *App) GetSettings() *config.Settings {
	return a.settings
}

// SaveSettings saves the application settings
func (a *App) SaveSettings(settings *config.Settings) error {
	if err := settings.Save(); err != nil {
		return err
	}
	a.settings = settings
	return nil
}

// WriteToTerminal writes data to the terminal
func (a *App) WriteToTerminal(data string) error {
	if a.terminal == nil {
		return fmt.Errorf("terminal not initialized")
	}
	_, err := a.terminal.Write([]byte(data))
	return err
}

// ReadFromTerminal reads data from the terminal
func (a *App) ReadFromTerminal() (string, error) {
	if a.terminal == nil {
		return "", fmt.Errorf("terminal not initialized")
	}
	buf := make([]byte, 4096)
	n, err := a.terminal.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// ResizeTerminal resizes the terminal
func (a *App) ResizeTerminal(rows, cols int) error {
	if a.terminal == nil {
		return fmt.Errorf("terminal not initialized")
	}
	return a.terminal.Resize(rows, cols)
}

// GetShell returns the detected shell
func (a *App) GetShell() string {
	if a.terminal == nil {
		return "bash" // Default
	}
	return a.terminal.GetShell()
}
