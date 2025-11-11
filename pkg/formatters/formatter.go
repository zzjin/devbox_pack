/**
 * DevBox Pack Execution Plan Generator - Output Formatter
 */

package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

// Formatter output formatter interface
type Formatter interface {
	Format(plan *types.ExecutionPlan, options *types.CLIOptions) (string, error)
}

// JSONFormatter JSON formatter
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format formats execution plan as JSON
func (f *JSONFormatter) Format(plan *types.ExecutionPlan, options *types.CLIOptions) (string, error) {
	if plan == nil {
		return "", fmt.Errorf("execution plan cannot be nil")
	}

	var data []byte
	var err error

	if options != nil && options.Pretty {
		data, err = json.MarshalIndent(plan, "", "  ")
	} else {
		data, err = json.Marshal(plan)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal plan to JSON: %w", err)
	}

	return string(data), nil
}

// PrettyFormatter Pretty formatter (human-readable format)
type PrettyFormatter struct{}

// NewPrettyFormatter creates a new Pretty formatter
func NewPrettyFormatter() *PrettyFormatter {
	return &PrettyFormatter{}
}

// Format formats execution plan as human-readable format
func (f *PrettyFormatter) Format(plan *types.ExecutionPlan, options *types.CLIOptions) (string, error) {
	if plan == nil {
		return "", fmt.Errorf("execution plan cannot be nil")
	}

	var lines []string

	// Title
	lines = append(lines, "🚀 DevBox Pack Execution Plan")
	lines = append(lines, strings.Repeat("═", 50))
	lines = append(lines, "")

	// Runtime information
	lines = append(lines, "📋 Runtime Configuration")
	lines = append(lines, strings.Repeat("─", 20))
	lines = append(lines, fmt.Sprintf("Provider: %s", plan.Provider))

	// Framework information if available
	if plan.Runtime.Framework != nil {
		lines = append(lines, fmt.Sprintf("Framework: %s", *plan.Runtime.Framework))
	}

	// Environment variables are now at the top level
	if len(plan.Environment) > 0 {
		lines = append(lines, "Environment Variables:")
		for key, value := range plan.Environment {
			lines = append(lines, fmt.Sprintf("  %s=%s", key, value))
		}
	}
	lines = append(lines, "")

	// Base image
	lines = append(lines, "🐳 Base Image")
	lines = append(lines, strings.Repeat("─", 20))
	lines = append(lines, fmt.Sprintf("Image: %s", plan.Runtime.Image))
	lines = append(lines, "")

	// APT dependencies
	if len(plan.Apt) > 0 {
		lines = append(lines, "📦 System Dependencies")
		lines = append(lines, strings.Repeat("─", 20))
		for _, pkg := range plan.Apt {
			lines = append(lines, fmt.Sprintf("• %s", pkg))
		}
		lines = append(lines, "")
	}

	// Commands
	commandLabels := map[string]string{
		"dev":   "🔧 Development Commands",
		"setup": "🏗️  Setup Commands",
		"build": "Build Commands",
		"run":   "Production Commands",
	}

	commands := []struct {
		key      string
		commands []string
	}{
		{"dev", plan.Commands.Dev},
		{"setup", plan.Commands.Setup},
		{"build", plan.Commands.Build},
		{"run", plan.Commands.Run},
	}

	for _, cmd := range commands {
		if len(cmd.commands) > 0 {
			lines = append(lines, commandLabels[cmd.key])
			lines = append(lines, strings.Repeat("─", 20))
			lines = append(lines, cmd.commands...)
			lines = append(lines, "")
		}
	}

	// Port information
	if plan.Port > 0 {
		lines = append(lines, fmt.Sprintf("Port: %d", plan.Port))
		lines = append(lines, "")
	}

	// Detection evidence
	if len(plan.Evidence.Files) > 0 || plan.Evidence.Reason != "" {
		lines = append(lines, "🔍 Detection Evidence")
		lines = append(lines, strings.Repeat("─", 20))

		if len(plan.Evidence.Files) > 0 {
			lines = append(lines, "Files:")
			for _, file := range plan.Evidence.Files {
				lines = append(lines, fmt.Sprintf("  • %s", file))
			}
		}

		if plan.Evidence.Reason != "" {
			if len(plan.Evidence.Files) > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, fmt.Sprintf("Reason: %s", plan.Evidence.Reason))
		}
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n"), nil
}

// FormatterFactory formatter factory
type FormatterFactory struct {
	formatters map[string]Formatter
}

// NewFormatterFactory creates a new formatter factory
func NewFormatterFactory() *FormatterFactory {
	factory := &FormatterFactory{
		formatters: make(map[string]Formatter),
	}

	// Register default formatters
	factory.formatters["json"] = NewJSONFormatter()
	factory.formatters["pretty"] = NewPrettyFormatter()

	return factory
}

// GetFormatter gets formatter
func (f *FormatterFactory) GetFormatter(format string) (Formatter, error) {
	formatter, exists := f.formatters[format]
	if !exists {
		return nil, fmt.Errorf("unsupported output format: %s, supported formats: %v",
			format, f.GetSupportedFormats())
	}
	return formatter, nil
}

// Format formats execution plan
func (f *FormatterFactory) Format(plan *types.ExecutionPlan, format string) (string, error) {
	formatter, err := f.GetFormatter(format)
	if err != nil {
		return "", err
	}
	return formatter.Format(plan, nil)
}

// GetSupportedFormats gets supported format list
func (f *FormatterFactory) GetSupportedFormats() []string {
	formats := make([]string, 0, len(f.formatters))
	for format := range f.formatters {
		formats = append(formats, format)
	}
	return formats
}

// RegisterFormatter registers custom formatter
func (f *FormatterFactory) RegisterFormatter(format string, formatter Formatter) {
	f.formatters[format] = formatter
}

// OutputUtils output utility class
type OutputUtils struct {
	factory *FormatterFactory
}

// NewOutputUtils creates new output utility
func NewOutputUtils() *OutputUtils {
	return &OutputUtils{
		factory: NewFormatterFactory(),
	}
}

// OutputPlan outputs execution plan
func (u *OutputUtils) OutputPlan(plan *types.ExecutionPlan, options *types.CLIOptions) error {
	output, err := u.factory.Format(plan, string(options.Format))
	if err != nil {
		return fmt.Errorf("failed to format plan: %w", err)
	}
	fmt.Println(output)
	return nil
}

// OutputError outputs error information
func (u *OutputUtils) OutputError(err error, options *types.CLIOptions) {
	if options != nil && options.Verbose {
		fmt.Printf("❌ Error: %s\n", err.Error())
	} else {
		fmt.Printf("❌ %s\n", err.Error())
	}
}

// OutputDebug outputs debug information
func (u *OutputUtils) OutputDebug(message string, options *types.CLIOptions) {
	if options != nil && options.Verbose {
		fmt.Printf("🔍 %s\n", message)
	}
}

// OutputInfo outputs information
func (u *OutputUtils) OutputInfo(message string, options *types.CLIOptions) {
	if options == nil || !options.Quiet {
		fmt.Printf("ℹ️  %s\n", message)
	}
}

// OutputSuccess outputs success information
func (u *OutputUtils) OutputSuccess(message string, options *types.CLIOptions) {
	if options == nil || !options.Quiet {
		fmt.Printf("✅ %s\n", message)
	}
}

// OutputWarning outputs warning information
func (u *OutputUtils) OutputWarning(message string, options *types.CLIOptions) {
	if options == nil || !options.Quiet {
		fmt.Printf("⚠️  %s\n", message)
	}
}
