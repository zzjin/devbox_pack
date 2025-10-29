package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExecutionPlan 执行计划结构
type ExecutionPlan struct {
	Provider string `json:"provider"`
	Base     struct {
		Name     string `json:"name"`
		Platform string `json:"platform"`
	} `json:"base"`
	Runtime struct {
		Language string                 `json:"language"`
		Version  string                 `json:"version"`
		Tools    interface{}            `json:"tools"`
		Env      map[string]interface{} `json:"env"`
	} `json:"runtime"`
	Apt      interface{} `json:"apt"`
	AptDeps  interface{} `json:"aptDeps"`
	Dev      *Command    `json:"dev"`
	Build    *Command    `json:"build"`
	Start    *Command    `json:"start"`
	Commands struct {
		Dev   []string `json:"dev"`
		Build []string `json:"build"`
		Start []string `json:"start"`
	} `json:"commands"`
	Evidence struct {
		Files  interface{} `json:"files"`
		Reason string      `json:"reason"`
	} `json:"evidence"`
}

// Command 命令结构
type Command struct {
	Cmd     string   `json:"cmd"`
	Caches  []string `json:"caches,omitempty"`
	PortEnv string   `json:"portEnv,omitempty"`
	Notes   []string `json:"notes,omitempty"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	TestCase string   `json:"testCase"`
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
	Plan     *ExecutionPlan `json:"plan,omitempty"`
}

// ValidationSummary 验证摘要
type ValidationSummary struct {
	TotalCases    int                `json:"totalCases"`
	ValidCases    int                `json:"validCases"`
	InvalidCases  int                `json:"invalidCases"`
	SuccessRate   float64            `json:"successRate"`
	Results       []ValidationResult `json:"results"`
	ProviderStats map[string]int     `json:"providerStats"`
}

// 验证器
type PlanValidator struct {
	knownProviders map[string]bool
	knownCommands  map[string][]string
}

// NewPlanValidator 创建新的验证器
func NewPlanValidator() *PlanValidator {
	return &PlanValidator{
		knownProviders: map[string]bool{
			"node":       true,
			"python":     true,
			"java":       true,
			"go":         true,
			"php":        true,
			"ruby":       true,
			"elixir":     true,
			"deno":       true,
			"rust":       true,
			"staticfile": true,
			"shell":      true,
		},
		knownCommands: map[string][]string{
			"node":   {"npm", "yarn", "pnpm", "bun"},
			"python": {"pip", "poetry", "pipenv", "uv", "pdm"},
			"java":   {"mvn", "gradle"},
			"go":     {"go"},
			"php":    {"composer"},
			"ruby":   {"bundle"},
			"elixir": {"mix"},
			"deno":   {"deno"},
			"rust":   {"cargo"},
		},
	}
}

// ValidatePlan 验证单个执行计划
func (v *PlanValidator) ValidatePlan(testCase string, plan *ExecutionPlan) ValidationResult {
	result := ValidationResult{
		TestCase: testCase,
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
		Plan:     plan,
	}

	// 验证 Provider
	if plan.Provider == "" {
		result.Errors = append(result.Errors, "provider 不能为空")
		result.Valid = false
	} else if !v.knownProviders[plan.Provider] {
		result.Errors = append(result.Errors, fmt.Sprintf("未知的 Provider: %s", plan.Provider))
		result.Valid = false
	}

	// 验证基础镜像
	if plan.Base.Name == "" {
		result.Errors = append(result.Errors, "基础镜像名称不能为空")
		result.Valid = false
	}

	// 验证运行时
	if plan.Runtime.Language == "" {
		result.Errors = append(result.Errors, "运行时语言不能为空")
		result.Valid = false
	}

	// 验证工具
	if knownTools, exists := v.knownCommands[plan.Provider]; exists {
		if tools, ok := plan.Runtime.Tools.([]interface{}); ok {
			for _, tool := range tools {
				if toolStr, ok := tool.(string); ok {
					if toolStr == "" {
						result.Warnings = append(result.Warnings, "检测到空的工具名称")
						continue
					}
					found := false
					for _, knownTool := range knownTools {
						if strings.Contains(toolStr, knownTool) {
							found = true
							break
						}
					}
					if !found {
						result.Warnings = append(result.Warnings, fmt.Sprintf("未知的工具: %s (Provider: %s)", toolStr, plan.Provider))
					}
				}
			}
		}
	}

	// 检查命令 - 优先检查 Commands 字段，如果为空再检查传统字段
	hasCommands := len(plan.Commands.Dev) > 0 || len(plan.Commands.Build) > 0 || len(plan.Commands.Start) > 0
	
	if !hasCommands {
		// 如果 Commands 字段为空，检查传统字段
		if plan.Dev == nil || plan.Dev.Cmd == "" {
			result.Errors = append(result.Errors, "dev 命令不能为空")
			result.Valid = false
		}

		if plan.Build == nil || plan.Build.Cmd == "" {
			result.Errors = append(result.Errors, "build 命令不能为空")
			result.Valid = false
		}

		if plan.Start == nil || plan.Start.Cmd == "" {
			result.Errors = append(result.Errors, "start 命令不能为空")
			result.Valid = false
		}
	} else {
		// 如果 Commands 字段有内容，验证其完整性
		if len(plan.Commands.Dev) == 0 {
			result.Warnings = append(result.Warnings, "dev 命令列表为空")
		}
		if len(plan.Commands.Build) == 0 {
			result.Warnings = append(result.Warnings, "build 命令列表为空")
		}
		if len(plan.Commands.Start) == 0 {
			result.Warnings = append(result.Warnings, "start 命令列表为空")
		}
	}

	// 检查端口配置
	hasPortConfig := false
	if hasCommands {
		// 检查 Commands 字段中的端口配置
		for _, cmd := range plan.Commands.Start {
			if strings.Contains(cmd, "PORT") || strings.Contains(cmd, "port") {
				hasPortConfig = true
				break
			}
		}
	} else if plan.Start != nil && plan.Start.Cmd != "" {
		// 检查传统字段中的端口配置
		hasPortConfig = strings.Contains(plan.Start.Cmd, "PORT") || strings.Contains(plan.Start.Cmd, "port")
	}
	
	if !hasPortConfig && (hasCommands && len(plan.Commands.Start) > 0 || !hasCommands && plan.Start != nil && plan.Start.Cmd != "") {
		result.Warnings = append(result.Warnings, "启动命令缺少端口环境变量")
	}

	// 验证证据文件
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "没有检测到证据文件")
	}

	// 验证特定 Provider 的规则
	v.validateProviderSpecific(plan, &result)

	return result
}

// validateCommand 验证命令
func (v *PlanValidator) validateCommand(cmdType string, cmd *Command, result *ValidationResult) {
	if cmd == nil {
		return
	}

	if cmd.Cmd == "" {
		result.Errors = append(result.Errors, fmt.Sprintf("%s 命令不能为空", cmdType))
		result.Valid = false
	}

	// 验证启动命令的端口环境变量
	if cmdType == "start" && cmd.PortEnv == "" {
		result.Warnings = append(result.Warnings, "启动命令缺少端口环境变量")
	}
}

// validateProviderSpecific 验证特定 Provider 的规则
func (v *PlanValidator) validateProviderSpecific(plan *ExecutionPlan, result *ValidationResult) {
	switch plan.Provider {
	case "node":
		v.validateNodeProject(plan, result)
	case "python":
		v.validatePythonProject(plan, result)
	case "java":
		v.validateJavaProject(plan, result)
	case "go":
		v.validateGoProject(plan, result)
	case "rust":
		v.validateRustProject(plan, result)
	}
}

// validateNodeProject 验证 Node.js 项目
func (v *PlanValidator) validateNodeProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Node.js 项目缺少 package.json 文件")
		result.Warnings = append(result.Warnings, "Node.js 项目缺少锁文件")
		return
	}

	// 检查是否有 package.json
	hasPackageJson := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok && strings.Contains(fileStr, "package.json") {
				hasPackageJson = true
				break
			}
		}
	}
	if !hasPackageJson {
		result.Warnings = append(result.Warnings, "Node.js 项目缺少 package.json 文件")
	}

	// 检查工具和锁文件的一致性
	hasLockFile := false
	lockFiles := []string{"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"}
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok {
				for _, lockFile := range lockFiles {
					if strings.Contains(fileStr, lockFile) {
						hasLockFile = true
						break
					}
				}
				if hasLockFile {
					break
				}
			}
		}
	}
	if !hasLockFile {
		result.Warnings = append(result.Warnings, "Node.js 项目缺少锁文件")
	}
}

// validatePythonProject 验证 Python 项目
func (v *PlanValidator) validatePythonProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Python 项目缺少依赖文件")
		return
	}

	// 检查依赖文件
	hasDepsFile := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok {
				if strings.Contains(fileStr, "requirements.txt") ||
				   strings.Contains(fileStr, "pyproject.toml") ||
				   strings.Contains(fileStr, "Pipfile") ||
				   strings.Contains(fileStr, "poetry.lock") {
					hasDepsFile = true
					break
				}
			}
		}
	}
	if !hasDepsFile {
		result.Warnings = append(result.Warnings, "Python 项目缺少依赖文件")
	}
}

// validateJavaProject 验证 Java 项目
func (v *PlanValidator) validateJavaProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Java 项目缺少构建文件 (pom.xml 或 build.gradle)")
		return
	}

	// 检查构建文件
	hasBuildFile := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok {
				if strings.Contains(fileStr, "pom.xml") || strings.Contains(fileStr, "build.gradle") {
					hasBuildFile = true
					break
				}
			}
		}
	}
	if !hasBuildFile {
		result.Warnings = append(result.Warnings, "Java 项目缺少构建文件 (pom.xml 或 build.gradle)")
	}
}

// validateGoProject 验证 Go 项目
func (v *PlanValidator) validateGoProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Go 项目缺少 go.mod 文件")
		return
	}

	// 检查 go.mod
	hasGoMod := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok && strings.Contains(fileStr, "go.mod") {
				hasGoMod = true
				break
			}
		}
	}
	if !hasGoMod {
		result.Warnings = append(result.Warnings, "Go 项目缺少 go.mod 文件")
	}
}

// validateRustProject 验证 Rust 项目
func (v *PlanValidator) validateRustProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Rust 项目缺少 Cargo.toml 文件")
		return
	}

	// 检查 Cargo.toml
	hasCargoToml := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok && strings.Contains(fileStr, "Cargo.toml") {
				hasCargoToml = true
				break
			}
		}
	}
	if !hasCargoToml {
		result.Warnings = append(result.Warnings, "Rust 项目缺少 Cargo.toml 文件")
	}
}

// 主函数
func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run validate_plans.go <plans_directory>")
		os.Exit(1)
	}

	plansDir := os.Args[1]
	validator := NewPlanValidator()
	
	// 获取所有计划文件
	planFiles, err := filepath.Glob(filepath.Join(plansDir, "*.json"))
	if err != nil {
		fmt.Printf("错误: 无法读取计划文件: %v\n", err)
		os.Exit(1)
	}

	if len(planFiles) == 0 {
		fmt.Printf("错误: 在目录 %s 中没有找到 JSON 文件\n", plansDir)
		os.Exit(1)
	}

	fmt.Printf("找到 %d 个执行计划文件\n", len(planFiles))

	var results []ValidationResult
	providerStats := make(map[string]int)

	// 验证每个计划文件
	for _, planFile := range planFiles {
		testCase := strings.TrimSuffix(filepath.Base(planFile), ".json")
		
		// 读取计划文件
		data, err := ioutil.ReadFile(planFile)
		if err != nil {
			result := ValidationResult{
				TestCase: testCase,
				Valid:    false,
				Errors:   []string{fmt.Sprintf("无法读取文件: %v", err)},
			}
			results = append(results, result)
			continue
		}

		// 解析 JSON
		var plan ExecutionPlan
		if err := json.Unmarshal(data, &plan); err != nil {
			result := ValidationResult{
				TestCase: testCase,
				Valid:    false,
				Errors:   []string{fmt.Sprintf("JSON 解析错误: %v", err)},
			}
			results = append(results, result)
			continue
		}

		// 验证计划
		result := validator.ValidatePlan(testCase, &plan)
		results = append(results, result)

		// 统计 Provider
		if result.Valid {
			providerStats[plan.Provider]++
		}
	}

	// 生成摘要
	validCases := 0
	for _, result := range results {
		if result.Valid {
			validCases++
		}
	}

	summary := ValidationSummary{
		TotalCases:    len(results),
		ValidCases:    validCases,
		InvalidCases:  len(results) - validCases,
		SuccessRate:   float64(validCases) / float64(len(results)) * 100,
		Results:       results,
		ProviderStats: providerStats,
	}

	// 输出结果
	fmt.Printf("\n=== 验证结果摘要 ===\n")
	fmt.Printf("总计划数: %d\n", summary.TotalCases)
	fmt.Printf("有效计划: %d\n", summary.ValidCases)
	fmt.Printf("无效计划: %d\n", summary.InvalidCases)
	fmt.Printf("成功率: %.2f%%\n", summary.SuccessRate)

	fmt.Printf("\n=== Provider 统计 ===\n")
	var providers []string
	for provider := range providerStats {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	
	for _, provider := range providers {
		fmt.Printf("- %s: %d 个有效计划\n", provider, providerStats[provider])
	}

	// 输出详细结果
	fmt.Printf("\n=== 详细验证结果 ===\n")
	for _, result := range results {
		status := "✓"
		if !result.Valid {
			status = "✗"
		}
		
		fmt.Printf("%s %s", status, result.TestCase)
		if result.Plan != nil {
			fmt.Printf(" (Provider: %s)", result.Plan.Provider)
		}
		fmt.Printf("\n")

		if len(result.Errors) > 0 {
			for _, err := range result.Errors {
				fmt.Printf("  错误: %s\n", err)
			}
		}

		if len(result.Warnings) > 0 {
			for _, warning := range result.Warnings {
				fmt.Printf("  警告: %s\n", warning)
			}
		}
	}

	// 保存详细报告
	outputDir := filepath.Dir(plansDir)
	reportFile := filepath.Join(outputDir, "validation_report.json")
	
	reportData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Printf("错误: 无法生成报告: %v\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(reportFile, reportData, 0644); err != nil {
		fmt.Printf("错误: 无法保存报告: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n详细报告已保存到: %s\n", reportFile)

	// 如果有无效计划，以非零状态退出
	if summary.InvalidCases > 0 {
		os.Exit(1)
	}
}