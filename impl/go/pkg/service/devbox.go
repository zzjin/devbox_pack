package service

import (
	"fmt"

	"github.com/labring/devbox-pack/pkg/detector"
	"github.com/labring/devbox-pack/pkg/formatters"
	"github.com/labring/devbox-pack/pkg/generators"
	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

// DevBoxPack core service class
type DevBoxPack struct {
	gitHandler      *git.GitHandler
	detectionEngine *detector.DetectionEngine
	planGenerator   *generators.ExecutionPlanGenerator
	outputUtils     *formatters.OutputUtils
}

// NewDevBoxPack creates a DevBox Pack instance
func NewDevBoxPack() *DevBoxPack {
	return &DevBoxPack{
		gitHandler:      git.NewGitHandler(),
		detectionEngine: detector.NewDetectionEngine(),
		planGenerator:   generators.NewExecutionPlanGenerator(),
		outputUtils:     formatters.NewOutputUtils(),
	}
}

// GeneratePlan generates execution plan
func (d *DevBoxPack) GeneratePlan(repoPath string, options *types.CLIOptions) (*types.ExecutionPlan, error) {
	// 1. Prepare project directory
	d.outputUtils.OutputInfo("Preparing project directory...", options)
	projectPath, err := d.gitHandler.PrepareProject(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare project: %w", err)
	}

	// 2. Scan project files
	d.outputUtils.OutputInfo("Scanning project files...", options)
	scanOptions := &types.ScanOptions{
		MaxDepth: 3,
		MaxFiles: 1000,
	}
	files, err := d.gitHandler.ScanProject(projectPath, scanOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	d.outputUtils.OutputDebug(fmt.Sprintf("Scanned %d files", len(files)), options)

	// 3. Detect language and framework
	d.outputUtils.OutputInfo("Detecting language and framework...", options)
	fileInfos := make([]types.FileInfo, len(files))
	for i, file := range files {
		fileInfos[i] = *file
	}
	detectResults, err := d.detectionEngine.DetectProject(projectPath, fileInfos, d.gitHandler, options)
	if err != nil {
		return nil, fmt.Errorf("failed to detect project: %w", err)
	}

	if len(detectResults) == 0 {
		return nil, fmt.Errorf("no supported language or framework detected in path: %s", projectPath)
	}

	// 4. Generate execution plan
	d.outputUtils.OutputInfo("Generating execution plan...", options)
	detectResultValues := make([]types.DetectResult, len(detectResults))
	for i, result := range detectResults {
		detectResultValues[i] = *result
	}
	plan, err := d.planGenerator.GeneratePlan(detectResultValues, *options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	d.outputUtils.OutputSuccess("Execution plan generated successfully", options)
	return plan, nil
}

// Run executes the complete workflow
func (d *DevBoxPack) Run(repoPath string, options *types.CLIOptions) error {
	plan, err := d.GeneratePlan(repoPath, options)
	if err != nil {
		d.outputUtils.OutputError(err, options)
		return err
	}

	return d.outputUtils.OutputPlan(plan, options)
}
