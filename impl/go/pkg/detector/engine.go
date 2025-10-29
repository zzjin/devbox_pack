/**
 * DevBox Pack Execution Plan Generator - Detection Engine
 */

package detector

import (
	"fmt"
	"sort"

	"github.com/labring/devbox-pack/pkg/providers"
	"github.com/labring/devbox-pack/pkg/types"
)

// DetectionEngine detection engine
// Responsible for coordinating all Providers for project detection
type DetectionEngine struct {
	providers map[string]Provider
}

// NewDetectionEngine creates a new detection engine
func NewDetectionEngine() *DetectionEngine {
	engine := &DetectionEngine{
		providers: make(map[string]Provider),
	}
	engine.initializeProviders()
	return engine
}

// initializeProviders initializes all Providers
func (e *DetectionEngine) initializeProviders() {
	providerList := []Provider{
		providers.NewStaticFileProvider(), // Move StaticFileProvider to last, lower priority
		providers.NewNodeProvider(),
		providers.NewPythonProvider(),
		providers.NewJavaProvider(),
		providers.NewGoProvider(),
		providers.NewPHPProvider(),
		providers.NewRubyProvider(),
		providers.NewDenoProvider(),
		providers.NewRustProvider(),
		providers.NewShellProvider(),
	}

	for _, provider := range providerList {
		e.providers[provider.GetName()] = provider
	}
}

// GetAvailableProviders gets all available Providers
func (e *DetectionEngine) GetAvailableProviders() []string {
	var names []string
	for name := range e.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetProvider gets the specified Provider
func (e *DetectionEngine) GetProvider(name string) (Provider, bool) {
	provider, exists := e.providers[name]
	return provider, exists
}

// DetectProject detects project language and framework
func (e *DetectionEngine) DetectProject(
	projectPath string,
	files []types.FileInfo,
	gitHandler interface{},
	options *types.CLIOptions,
) ([]*types.DetectResult, error) {
	var results []*types.DetectResult

	if options != nil && options.Provider != nil && *options.Provider != "" {
		// Use specified Provider
		provider, exists := e.providers[*options.Provider]
		if !exists {
			return nil, fmt.Errorf("unknown Provider: %s, available Providers: %v",
				*options.Provider, e.GetAvailableProviders())
		}

		result, err := e.runProvider(provider, projectPath, files, gitHandler)
		if err != nil {
			return nil, err
		}
		if result != nil {
			results = append(results, result)
		}
	} else {
		// Run all Providers
		providers := e.getProvidersByPriority()

		for _, provider := range providers {
			result, err := e.runProvider(provider, projectPath, files, gitHandler)
			if err != nil {
				fmt.Printf("Provider %s detection failed: %v\n", provider.GetName(), err)
				// Continue with other Providers
				continue
			}
			if result != nil && result.Matched {
				results = append(results, result)
			}
		}
	}

	// Filter valid results
	var validResults []*types.DetectResult
	for _, result := range results {
		if e.isValidDetectResult(result) {
			validResults = append(validResults, result)
		}
	}

	return validResults, nil
}

// runProvider runs a single Provider
func (e *DetectionEngine) runProvider(
	provider Provider,
	projectPath string,
	files []types.FileInfo,
	gitHandler interface{},
) (*types.DetectResult, error) {
	result, err := provider.Detect(projectPath, files, interface{}(gitHandler))
	if err != nil {
		return nil, fmt.Errorf("Provider %s execution failed: %w", provider.GetName(), err)
	}

	// Validate result
	if !e.isValidDetectResult(result) {
		fmt.Printf("Provider %s returned invalid detection result\n", provider.GetName())
		return nil, nil
	}

	return result, nil
}

// getProvidersByPriority gets Provider list sorted by priority
func (e *DetectionEngine) getProvidersByPriority() []Provider {
	var providers []Provider
	for _, provider := range e.providers {
		providers = append(providers, provider)
	}

	// Sort by priority (smaller number means higher priority)
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].GetPriority() < providers[j].GetPriority()
	})

	return providers
}

// isValidDetectResult validates the validity of detection results
func (e *DetectionEngine) isValidDetectResult(result *types.DetectResult) bool {
	if result == nil {
		return false
	}

	// Check confidence range
	if result.Confidence < 0 || result.Confidence > 1 {
		return false
	}

	// If detected, should have language information
	if result.Matched && result.Language == "" {
		return false
	}

	return true
}

// GetBestResult gets the best detection result
func (e *DetectionEngine) GetBestResult(results []*types.DetectResult) *types.DetectResult {
	if len(results) == 0 {
		return nil
	}

	// Return result with highest confidence
	best := results[0]
	for _, result := range results[1:] {
		if result.Confidence > best.Confidence {
			best = result
		}
	}

	return best
}

// FilterResults filters detection results
func (e *DetectionEngine) FilterResults(results []*types.DetectResult, minConfidence float64) []*types.DetectResult {
	if minConfidence == 0 {
		minConfidence = 0.5 // Default minimum confidence
	}

	var filtered []*types.DetectResult
	for _, result := range results {
		if result.Matched && result.Confidence >= minConfidence {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// MergeResults merges similar detection results
func (e *DetectionEngine) MergeResults(results []*types.DetectResult) []*types.DetectResult {
	var merged []*types.DetectResult
	processed := make(map[string]bool)

	for _, result := range results {
		if result.Language == "" || processed[result.Language] {
			continue
		}

		// Find other results with the same language
		var sameLanguageResults []*types.DetectResult
		for _, r := range results {
			if r.Language == result.Language && r.Matched {
				sameLanguageResults = append(sameLanguageResults, r)
			}
		}

		if len(sameLanguageResults) == 1 {
			merged = append(merged, result)
		} else {
			// Merge results of the same language, choose the one with highest confidence
			bestResult := sameLanguageResults[0]
			for _, r := range sameLanguageResults[1:] {
				if r.Confidence > bestResult.Confidence {
					bestResult = r
				}
			}
			merged = append(merged, bestResult)
		}

		processed[result.Language] = true
	}

	return merged
}

// DetectionStats detection statistics
type DetectionStats struct {
	Total         int      `json:"total"`
	Detected      int      `json:"detected"`
	Languages     []string `json:"languages"`
	Frameworks    []string `json:"frameworks"`
	AvgConfidence float64  `json:"avgConfidence"`
}

// GetDetectionStats gets detection statistics
func (e *DetectionEngine) GetDetectionStats(results []*types.DetectResult) *DetectionStats {
	var detected []*types.DetectResult
	for _, r := range results {
		if r.Matched {
			detected = append(detected, r)
		}
	}

	// Collect languages and frameworks
	languageSet := make(map[string]bool)
	frameworkSet := make(map[string]bool)
	var totalConfidence float64

	for _, r := range detected {
		if r.Language != "" {
			languageSet[r.Language] = true
		}
		if r.Framework != "" {
			frameworkSet[r.Framework] = true
		}
		totalConfidence += r.Confidence
	}

	// Convert to slices
	var languages []string
	for lang := range languageSet {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	var frameworks []string
	for framework := range frameworkSet {
		frameworks = append(frameworks, framework)
	}
	sort.Strings(frameworks)

	// Calculate average confidence
	var avgConfidence float64
	if len(detected) > 0 {
		avgConfidence = totalConfidence / float64(len(detected))
	}

	return &DetectionStats{
		Total:         len(results),
		Detected:      len(detected),
		Languages:     languages,
		Frameworks:    frameworks,
		AvgConfidence: avgConfidence,
	}
}
