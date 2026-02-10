package security

import (
	"regexp"
	"strings"
)

// RiskLevel represents the danger level of a command
type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskLow
	RiskMedium
	RiskHigh
	RiskCritical
)

// String returns the string representation of risk level
func (r RiskLevel) String() string {
	switch r {
	case RiskNone:
		return "None"
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	case RiskCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// Validator handles command validation
type Validator struct {
	blockedPatterns []*regexp.Regexp
	warningPatterns []*regexp.Regexp
}

// NewValidator creates a new command validator
func NewValidator() *Validator {
	v := &Validator{}
	v.compilePatterns()
	return v
}

// compilePatterns initializes regex patterns for dangerous commands
func (v *Validator) compilePatterns() {
	// Critical - Always blocked
	criticalPatterns := []string{
		`(?i)rm\s+-rf\s*/`,                    // Delete root
		`(?i)rm\s+-rf\s+/\.`,                 // Delete root variations
		`(?i):\(\)\s*\{\s*:\|:&\s*\};`,       // Fork bomb
		`(?i)mkfs\.`,                         // Format filesystem
		`(?i)dd\s+if=.*of=/dev/sd`,          // Direct disk write
		`(?i)>\s*/dev/sd`,                   // Write to disk device
		`(?i)mv\s+/\s+/dev/null`,            // Move root to null
	}

	// High risk - Require confirmation
	highPatterns := []string{
		`(?i)curl.*\|\s*(bash|sh|zsh)`,      // Pipe curl to shell
		`(?i)wget.*\|\s*(bash|sh|zsh)`,      // Pipe wget to shell
		`(?i)eval\s*\(`,                      // Eval with user input
		`(?i)exec\s*\(`,                      // Exec with user input
		`(?i)python.*-c.*\b(import\s+os|import\s+subprocess|exec|eval)\b`, // Python code execution
	}

	// Warning patterns - Notify but allow
	warningPatterns := []string{
		`(?i)sudo`,                           // Elevated privileges
		`(?i)chmod\s+777`,                    // World-writable
		`(?i)chown\s+-R`,                     // Recursive ownership change
		`(?i)rm\s+-rf`,                      // Recursive delete (with caution)
		`(?i)>.*/etc/`,                      // Write to /etc
	}

	for _, p := range criticalPatterns {
		if re, err := regexp.Compile(p); err == nil {
			v.blockedPatterns = append(v.blockedPatterns, re)
		}
	}

	for _, p := range highPatterns {
		if re, err := regexp.Compile(p); err == nil {
			v.blockedPatterns = append(v.blockedPatterns, re)
		}
	}

	for _, p := range warningPatterns {
		if re, err := regexp.Compile(p); err == nil {
			v.warningPatterns = append(v.warningPatterns, re)
		}
	}
}

// ValidateCommand checks a command and returns its risk level
func (v *Validator) ValidateCommand(command string) RiskLevel {
	command = strings.TrimSpace(command)
	
	if command == "" {
		return RiskNone
	}

	// Check critical patterns (hard block)
	for _, pattern := range v.blockedPatterns {
		if pattern.MatchString(command) {
			return RiskCritical
		}
	}

	// Check warning patterns
	for _, pattern := range v.warningPatterns {
		if pattern.MatchString(command) {
			return RiskMedium
		}
	}

	// Check for low-risk patterns
	lowRiskPatterns := []string{
		`(?i)sudo`,
		`(?i)apt-get`,
		`(?i)yum`,
		`(?i)brew install`,
	}

	for _, p := range lowRiskPatterns {
		if matched, _ := regexp.MatchString(p, command); matched {
			return RiskLow
		}
	}

	return RiskNone
}

// IsBlocked returns true if command should be completely blocked
func (v *Validator) IsBlocked(command string) bool {
	return v.ValidateCommand(command) == RiskCritical
}

// GetExplanation returns a human-readable explanation for the risk
func (v *Validator) GetExplanation(level RiskLevel) string {
	switch level {
	case RiskNone:
		return "This command appears safe."
	case RiskLow:
		return "This command requires elevated privileges or makes system changes."
	case RiskMedium:
		return "This command could impact system files or security settings. Please verify before executing."
	case RiskHigh:
		return "WARNING: This command downloads and executes remote code. Only execute if you trust the source."
	case RiskCritical:
		return "CRITICAL: This command is blocked as it could cause severe system damage or security compromise."
	default:
		return "Unknown risk level."
	}
}
