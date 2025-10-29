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
	var lines []string

	// Title
	lines = append(lines, "🚀 DevBox Pack Execution Plan")
	lines = append(lines, strings.Repeat("═", 50))
	lines = append(lines, "")

	// Runtime information
	lines = append(lines, "📋 Runtime Configuration")
	lines = append(lines, strings.Repeat("─", 20))
	lines = append(lines, fmt.Sprintf("Language: %s", plan.Runtime.Language))

	// Handle version display, avoid pointer address display
	version := "unknown"
	if plan.Runtime.Version != nil {
		version = *plan.Runtime.Version
	}
	lines = append(lines, fmt.Sprintf("Version: %s", version))

	if len(plan.Runtime.Tools) > 0 {
		lines = append(lines, fmt.Sprintf("Tools: %s", strings.Join(plan.Runtime.Tools, ", ")))
	}

	if len(plan.Runtime.Environment) > 0 {
		lines = append(lines, "Environment Variables:")
		for key, value := range plan.Runtime.Environment {
			lines = append(lines, fmt.Sprintf("  %s=%s", key, value))
		}
	}
	lines = append(lines, "")

	// Base image
	lines = append(lines, "🐳 Base Image")
	lines = append(lines, strings.Repeat("─", 20))
	lines = append(lines, fmt.Sprintf("Image: %s", plan.Base))
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
		"build": "🏗️  Build Commands",
		"start": "▶️  Start Commands",
	}

	commands := []struct {
		key      string
		commands []string
	}{
		{"dev", plan.Commands.Dev},
		{"build", plan.Commands.Build},
		{"start", plan.Commands.Start},
	}

	for _, cmd := range commands {
		if len(cmd.commands) > 0 {
			lines = append(lines, commandLabels[cmd.key])
			lines = append(lines, strings.Repeat("─", 20))
			lines = append(lines, cmd.commands...)
			lines = append(lines, "")
		}
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

// Global instances
var (
	DefaultFormatterFactory = NewFormatterFactory()
	DefaultOutputUtils      = NewOutputUtils()
)

// Convenience functions
func FormatPlan(plan *types.ExecutionPlan, format string, options *types.CLIOptions) (string, error) {
	return DefaultFormatterFactory.Format(plan, format)
}

func OutputPlan(plan *types.ExecutionPlan, options *types.CLIOptions) error {
	return DefaultOutputUtils.OutputPlan(plan, options)
}

func OutputError(err error, options *types.CLIOptions) {
	DefaultOutputUtils.OutputError(err, options)
}

func OutputDebug(message string, options *types.CLIOptions) {
	DefaultOutputUtils.OutputDebug(message, options)
}

func OutputInfo(message string, options *types.CLIOptions) {
	DefaultOutputUtils.OutputInfo(message, options)
}

func OutputSuccess(message string, options *types.CLIOptions) {
	DefaultOutputUtils.OutputSuccess(message, options)
}

func OutputWarning(message string, options *types.CLIOptions) {
	DefaultOutputUtils.OutputWarning(message, options)
}
