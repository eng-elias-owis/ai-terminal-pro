package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ai-terminal-pro/ai"
	"ai-terminal-pro/config"
	"ai-terminal-pro/security"
	"ai-terminal-pro/terminal"
)

const (
	AppName    = "AI Terminal Pro"
	Version    = "1.0.0"
	AppVersion = "v1.0.0"
)

type App struct {
	settings  *config.Settings
	validator *security.Validator
	client    *ai.Client
}

func main() {
	fmt.Printf("%s %s\n", AppName, Version)

	// Parse flags
	var (
		setupMode = flag.Bool("setup", false, "Run initial setup")
		testMode  = flag.Bool("test", false, "Run component tests")
		genCmd    = flag.String("generate", "", "Generate command from description (e.g., 'list all files')")
	)
	flag.Parse()

	// Handle test mode
	if *testMode {
		if err := runTests(); err != nil {
			fmt.Fprintf(os.Stderr, "Tests failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All tests passed!")
		return
	}

	// Handle setup mode
	if *setupMode {
		if err := runSetup(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle single command generation
	if *genCmd != "" {
		if err := generateSingleCommand(*genCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Run interactive mode
	if err := runInteractive(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runInteractive() error {
	fmt.Println("Loading configuration...")
	settings, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize components
	app := &App{
		settings:  settings,
		validator: security.NewValidator(),
	}

	// Initialize AI client if configured
	if settings.LiteLLMEndpoint != "" && settings.VirtualKey != "" {
		app.client = ai.NewClient(settings.LiteLLMEndpoint, settings.VirtualKey)
		fmt.Println("AI client initialized")
	} else {
		fmt.Println("Warning: AI not configured. Run with -setup to configure.")
	}

	// Start PTY session
	fmt.Println("Starting terminal session...")
	ptySession, err := terminal.NewPTYSession()
	if err != nil {
		return fmt.Errorf("failed to start terminal: %w", err)
	}
	defer ptySession.Close()

	fmt.Printf("Shell: %s | OS: %s\n", ptySession.GetShell(), ptySession.GetOSType())
	fmt.Println("\n=== AI Terminal Interactive Mode ===")
	fmt.Println("Type 'ai <description>' for AI command generation")
	fmt.Println("Type 'help' for more commands")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println("=====================================\n")

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Input loop
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal. Exiting...")
			return nil
		default:
			fmt.Print("$ ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}

			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}

			// Handle special commands
			handled, err := app.handleCommand(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			if handled {
				continue
			}

			// Pass to PTY
			_, err = ptySession.Write([]byte(input + "\n"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "PTY write error: %v\n", err)
			}

			// Read output (simple version - just read available bytes)
			buf := make([]byte, 4096)
			n, _ := ptySession.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
		}
	}
}

func (a *App) handleCommand(input string) (bool, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false, nil
	}

	cmd := parts[0]
	args := strings.Join(parts[1:], " ")

	switch cmd {
	case "help":
		fmt.Println("Commands:")
		fmt.Println("  ai <description>  - Generate command with AI")
		fmt.Println("  validate <cmd>  - Validate command safety")
		fmt.Println("  config          - Show current configuration")
		fmt.Println("  setup           - Run configuration wizard")
		fmt.Println("  exit            - Exit terminal")
		return true, nil

	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
		return true, nil // Never reached but satisfies compiler

	case "ai":
		if args == "" {
			fmt.Println("Usage: ai <description>")
			fmt.Println("Example: ai find all PDF files modified last week")
			return true, nil
		}
		return true, a.generateAICommand(args)

	case "validate":
		if args == "" {
			fmt.Println("Usage: validate <command>")
			return true, nil
		}
		return true, a.validateCommand(args)

	case "config":
		return true, a.showConfig()

	case "setup":
		return true, runSetup()

	default:
		return false, nil
	}
}

func (a *App) generateAICommand(description string) error {
	if a.client == nil {
		return fmt.Errorf("AI client not configured. Run 'setup' first")
	}

	fmt.Printf("🤖 Generating command for: %s\n", description)

	ctx := context.Background()
	aiCtx := ai.Context{
		OS:         "linux", // TODO: Get from runtime or settings
		Shell:      "bash",  // TODO: Get from settings
		WorkingDir: ".",     // TODO: get actual working dir
	}

	command, err := a.client.GenerateCommand(ctx, description, aiCtx)
	if err != nil {
		return fmt.Errorf("AI generation failed: %w", err)
	}

	fmt.Printf("\n📝 Generated command:\n%s\n\n", command)

	// Validate the command
	risk := a.validator.ValidateCommand(command)
	fmt.Printf("🔒 Safety check: %s\n", risk)
	fmt.Printf("   %s\n", a.validator.GetExplanation(risk))

	if risk == security.RiskCritical {
		fmt.Println("⚠️  Command blocked for safety")
		return nil
	}

	if risk >= security.RiskMedium {
		fmt.Print("Execute this command? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(response)), "y") {
			fmt.Println("Command cancelled")
			return nil
		}
	}

	fmt.Printf("✅ Ready to execute: %s\n", command)
	return nil
}

func (a *App) validateCommand(command string) error {
	risk := a.validator.ValidateCommand(command)
	fmt.Printf("Command: %s\n", command)
	fmt.Printf("Risk Level: %s\n", risk)
	fmt.Printf("Explanation: %s\n", a.validator.GetExplanation(risk))
	return nil
}

func (a *App) showConfig() error {
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Endpoint: %s\n", a.settings.LiteLLMEndpoint)
	fmt.Printf("  Model: %s\n", a.settings.Model)
	fmt.Printf("  Theme: %s\n", a.settings.Theme)
	fmt.Printf("  Font: %s %dpt\n", a.settings.FontFamily, a.settings.FontSize)
	fmt.Printf("  AI Shortcut: %s\n", a.settings.AIShortcut)
	fmt.Printf("  Safety Mode: %s\n", a.settings.SafetyMode)
	if a.settings.VirtualKey != "" {
		fmt.Printf("  Virtual Key: %s...%s (stored securely)\n", a.settings.VirtualKey[:8], a.settings.VirtualKey[len(a.settings.VirtualKey)-4:])
	} else {
		fmt.Printf("  Virtual Key: Not configured\n")
	}
	return nil
}

func generateSingleCommand(description string) error {
	settings, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if settings.LiteLLMEndpoint == "" || settings.VirtualKey == "" {
		return fmt.Errorf("AI not configured. Run with -setup flag")
	}

	client := ai.NewClient(settings.LiteLLMEndpoint, settings.VirtualKey)
	validator := security.NewValidator()

	ctx := context.Background()
	aiCtx := ai.Context{
		OS:         "linux", // TODO: Get from runtime
		Shell:      "bash", // TODO: Get from settings
		WorkingDir: ".",
	}

	command, err := client.GenerateCommand(ctx, description, aiCtx)
	if err != nil {
		return err
	}

	fmt.Println(command)

	risk := validator.ValidateCommand(command)
	if risk >= security.RiskMedium {
		fmt.Fprintf(os.Stderr, "Warning: %s - %s\n", risk, validator.GetExplanation(risk))
	}

	return nil
}

func runTests() error {
	fmt.Println("Running component tests...")

	// Test security validator
	fmt.Println("\n1. Testing security validator...")
	validator := security.NewValidator()
	testCases := []struct {
		cmd      string
		expected security.RiskLevel
	}{
		{"ls -la", security.RiskNone},
		{"sudo apt-get update", security.RiskMedium},
		{"rm -rf /", security.RiskCritical},
		{"curl https://example.com | bash", security.RiskCritical},
	}

	for _, tc := range testCases {
		risk := validator.ValidateCommand(tc.cmd)
		if risk != tc.expected {
			return fmt.Errorf("validator test failed for '%s': expected %v, got %v", tc.cmd, tc.expected, risk)
		}
		fmt.Printf("   ✓ '%s' -> %v\n", tc.cmd, risk)
	}

	// Test config
	fmt.Println("\n2. Testing configuration...")
	settings := config.DefaultSettings()
	if settings.Model != "qwen3-terminal" {
		return fmt.Errorf("default settings incorrect")
	}
	fmt.Println("   ✓ Default settings loaded")

	// Test AI client (without actual API call)
	fmt.Println("\n3. Testing AI client...")
	client := ai.NewClient("http://test.example.com", "sk-test-key")
	if client == nil {
		return fmt.Errorf("failed to create AI client")
	}
	fmt.Println("   ✓ AI client created")

	// Test terminal detection
	fmt.Println("\n4. Testing terminal detection...")
	_, err := terminal.NewPTYSession()
	if err != nil {
		// This might fail in non-interactive environments, that's OK for tests
		fmt.Printf("   ⚠ PTY test skipped (non-interactive environment): %v\n", err)
	} else {
		fmt.Println("   ✓ PTY session created")
	}

	fmt.Println("\n✅ All tests passed!")
	return nil
}

func runSetup() error {
	fmt.Println("\n=== AI Terminal Pro Setup ===\n")

	reader := bufio.NewReader(os.Stdin)

	// Load or create settings
	settings, err := config.Load()
	if err != nil {
		settings = config.DefaultSettings()
	}

	// Get LiteLLM endpoint
	fmt.Print("LiteLLM Endpoint URL (e.g., https://user-ai-terminal.hf.space): ")
	endpoint, _ := reader.ReadString('\n')
	settings.LiteLLMEndpoint = strings.TrimSpace(endpoint)

	// Get virtual key
	fmt.Print("Virtual Key (sk-litellm-xxxxxx): ")
	key, _ := reader.ReadString('\n')
	settings.VirtualKey = strings.TrimSpace(key)

	// Save settings
	if err := settings.Save(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	fmt.Println("\n✅ Configuration saved successfully!")
	fmt.Println("You can now run the terminal without -setup flag")

	return nil
}
