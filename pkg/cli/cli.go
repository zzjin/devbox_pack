// Package cli provides command-line interface functionality
// for the DevBox Pack execution plan generator.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/labring/devbox-pack/pkg/service"
	"github.com/labring/devbox-pack/pkg/types"
	"github.com/labring/devbox-pack/pkg/utils"
)

// CLIApp CLI application structure
type CLIApp struct {
	version string
}

// NewCLIApp creates a new CLI application instance
func NewCLIApp() *CLIApp {
	return &CLIApp{
		version: "1.0.0",
	}
}

// showHelp displays help information
func (c *CLIApp) showHelp() {
	fmt.Print(`
DevBox Pack Execution Plan Generator

Usage:
  devbox-pack <repository> [options]

Arguments:
  repository               Git repository URL or local path

Options:
  -h, --help              Show help information
  -v, --version           Show version information
  --ref <ref>             Git branch or tag (default: main)
  --subdir <path>         Subdirectory path
  --provider <name>       Force use of specified Provider
  --format <format>      Output format (pretty|json, default: pretty)
  --verbose               Show detailed information
  --offline               Offline mode, do not clone repository
  --platform <arch>       Target platform (e.g.: linux/amd64)
  --base <name>           Specify base image

Examples:
  devbox-pack https://github.com/user/repo
  devbox-pack . --offline --verbose
  devbox-pack /path/to/project --format json
  devbox-pack https://github.com/user/repo --ref develop --subdir backend

Supported Providers:
  node, python, java, go, php, ruby, deno, rust, staticfile, shell

Output Formats:
  pretty    - Human readable format (default)
  json      - JSON format
`)
}

// showVersion displays version information
func (c *CLIApp) showVersion() {
	fmt.Println(c.version)
}

// parseArgs parses command line arguments
func (c *CLIApp) parseArgs(args []string) (string, map[string]interface{}, error) {
	if len(args) < 2 {
		return "", nil, fmt.Errorf("please provide repository path or URL")
	}

	options := make(map[string]interface{})
	var repo string

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if arg == "--help" || arg == "-h" {
			c.showHelp()
			os.Exit(0)
		}

		if arg == "--version" || arg == "-v" {
			c.showVersion()
			os.Exit(0)
		}

		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")

			if key == "verbose" || key == "offline" || key == "quiet" {
				options[key] = true
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				options[key] = args[i+1]
				i++ // Skip next argument
			} else {
				return "", nil, types.NewDevBoxPackError(
					fmt.Sprintf("option --%s requires a value", key),
					types.ErrorCodeInvalidArgument,
					nil,
				)
			}
		} else if repo == "" {
			repo = arg
		} else {
			return "", nil, types.NewDevBoxPackError(
				fmt.Sprintf("unknown argument: %s", arg),
				types.ErrorCodeInvalidArgument,
				nil,
			)
		}
	}

	return repo, options, nil
}

// validateOptions validates and converts CLI options
func (c *CLIApp) validateOptions(rawOptions map[string]interface{}) (*types.CLIOptions, error) {
	options := &types.CLIOptions{
		Format:  "pretty",
		Verbose: false,
		Offline: false,
	}

	// Set option values
	if ref, ok := rawOptions["ref"].(string); ok {
		options.Ref = &ref
	}
	if subdir, ok := rawOptions["subdir"].(string); ok {
		options.Subdir = &subdir
	}
	if provider, ok := rawOptions["provider"].(string); ok {
		options.Provider = &provider
	}
	if format, ok := rawOptions["format"].(string); ok {
		options.Format = format
	}
	if verbose, ok := rawOptions["verbose"].(bool); ok {
		options.Verbose = verbose
	}
	if offline, ok := rawOptions["offline"].(bool); ok {
		options.Offline = offline
	}
	if quiet, ok := rawOptions["quiet"].(bool); ok {
		options.Quiet = quiet
	}
	if platform, ok := rawOptions["platform"].(string); ok {
		options.Platform = &platform
	}
	if base, ok := rawOptions["base"].(string); ok {
		options.Base = &base
	}

	// Validate output format
	if options.Format != string(types.OutputFormatJSON) && options.Format != string(types.OutputFormatPretty) {
		return nil, types.NewDevBoxPackError(
			fmt.Sprintf("unsupported output format: %s", options.Format),
			types.ErrorCodeInvalidFormat,
			map[string]interface{}{
				"format":    options.Format,
				"supported": []string{"json", "pretty"},
			},
		)
	}

	// Validate Provider
	if options.Provider != nil {
		supportedProviders := []string{
			"node", "python", "java", "go", "php", "ruby",
			"deno", "rust", "staticfile", "shell",
		}
		found := false
		for _, p := range supportedProviders {
			if *options.Provider == p {
				found = true
				break
			}
		}
		if !found {
			return nil, types.NewDevBoxPackError(
				fmt.Sprintf("unsupported Provider: %s", *options.Provider),
				types.ErrorCodeInvalidProvider,
				nil,
			)
		}
	}

	// Validate platform
	if options.Platform != nil {
		supportedPlatforms := []string{
			"linux/amd64", "linux/arm64",
			"darwin/amd64", "darwin/arm64",
		}
		found := false
		for _, p := range supportedPlatforms {
			if *options.Platform == p {
				found = true
				break
			}
		}
		if !found {
			return nil, types.NewDevBoxPackError(
				fmt.Sprintf("unsupported platform: %s", *options.Platform),
				types.ErrorCodeInvalidPlatform,
				nil,
			)
		}
	}

	return options, nil
}

// parseGitRepository parses Git repository information
func (c *CLIApp) parseGitRepository(repo string, options *types.CLIOptions) (*types.GitRepository, error) {
	// Check if it's a local path
	if options.Offline || (!strings.HasPrefix(repo, "http") && !strings.HasPrefix(repo, "git@")) {
		return &types.GitRepository{
			URL:    repo,
			Ref:    options.Ref,
			Subdir: options.Subdir,
		}, nil
	}

	// Validate Git URL format
	if !strings.HasPrefix(repo, "https://") &&
		!strings.HasPrefix(repo, "http://") &&
		!strings.HasPrefix(repo, "git@") &&
		!strings.HasPrefix(repo, "ssh://") {
		return nil, types.NewDevBoxPackError(
			fmt.Sprintf("invalid Git repository URL: %s", repo),
			types.ErrorCodeInvalidGitURL,
			nil,
		)
	}

	return &types.GitRepository{
		URL:    repo,
		Ref:    options.Ref,
		Subdir: options.Subdir,
	}, nil
}

// handleAnalyze handles analyze command
func (c *CLIApp) handleAnalyze(repo string, rawOptions map[string]interface{}) error {
	// Validate and convert options
	options, err := c.validateOptions(rawOptions)
	if err != nil {
		return err
	}

	if options.Verbose {
		fmt.Println(utils.Blue("🔍 DevBox Pack Execution Plan Generator"))
		fmt.Println(utils.Gray(fmt.Sprintf("Analysis target: %s", repo)))
		optionsJSON, _ := json.MarshalIndent(options, "", "  ")
		fmt.Println(utils.Gray(fmt.Sprintf("Options: %s", string(optionsJSON))))
		fmt.Println()
	}

	// Parse Git repository information
	gitRepo, err := c.parseGitRepository(repo, options)
	if err != nil {
		return err
	}

	if options.Verbose {
		fmt.Println(utils.Blue("📋 Repository Information:"))
		fmt.Println(utils.Gray(fmt.Sprintf("  URL: %s", gitRepo.URL)))
		if gitRepo.Ref != nil {
			fmt.Println(utils.Gray(fmt.Sprintf("  Ref: %s", *gitRepo.Ref)))
		}
		if gitRepo.Subdir != nil {
			fmt.Println(utils.Gray(fmt.Sprintf("  Subdirectory: %s", *gitRepo.Subdir)))
		}
		fmt.Println()
	}

	// Execute main logic
	if gitRepo.URL == "" {
		return types.NewDevBoxPackError(
			"please provide repository path or URL",
			types.ErrorCodeInvalidInput,
			nil,
		)
	}

	// Create DevBoxPack instance and run
	devBoxPack := service.NewDevBoxPack()
	return devBoxPack.Run(gitRepo.URL, options)
}

// handleError handles errors
func (c *CLIApp) handleError(err error) {
	if devBoxErr, ok := err.(*types.DevBoxPackError); ok {
		fmt.Fprintf(os.Stderr, "%s\n", utils.Red(fmt.Sprintf("❌ Error [%s]: %s", devBoxErr.Code, devBoxErr.Message)))
		if devBoxErr.Details != nil {
			detailsJSON, _ := json.MarshalIndent(devBoxErr.Details, "", "  ")
			fmt.Fprintf(os.Stderr, "%s\n", utils.Gray(fmt.Sprintf("Details: %s", string(detailsJSON))))
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", utils.Red(fmt.Sprintf("❌ Unknown error: %s", err.Error())))
	}
}

// Run runs the CLI application
func (c *CLIApp) Run(args []string) error {
	repo, options, err := c.parseArgs(args)
	if err != nil {
		c.handleError(err)
		return err
	}

	if repo == "" {
		fmt.Fprintf(os.Stderr, "%s\n", utils.Red("❌ Error: please provide repository path or URL"))
		fmt.Println()
		c.showHelp()
		return fmt.Errorf("missing repository argument")
	}

	err = c.handleAnalyze(repo, options)
	if err != nil {
		c.handleError(err)
		return err
	}

	return nil
}
