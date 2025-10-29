package detector

import (
	"path/filepath"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewDetectionEngine(t *testing.T) {
	engine := NewDetectionEngine()
	if engine == nil {
		t.Fatal("NewDetectionEngine() returned nil")
	}

	if engine.providers == nil {
		t.Fatal("providers map is nil")
	}

	// Check that providers are initialized
	providers := engine.GetAvailableProviders()
	if len(providers) == 0 {
		t.Fatal("no providers initialized")
	}

	expectedProviders := []string{"node", "python", "java", "go", "php", "ruby", "rust", "deno", "shell", "staticfile"}
	for _, expected := range expectedProviders {
		found := false
		for _, provider := range providers {
			if provider == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected provider %s not found", expected)
		}
	}
}

func TestDetectProject_NodeExpress(t *testing.T) {
	engine := NewDetectionEngine()

	// Use testdata fixture
	projectPath := filepath.Join("../../testdata/node/express")
	files := []types.FileInfo{
		{
			Path:        "package.json",
			Name:        "package.json",
			IsDirectory: false,
			Size:        func() *int64 { v := int64(100); return &v }(),
			Extension:   func() *string { v := ".json"; return &v }(),
		},
	}

	options := &types.CLIOptions{
		Verbose: true,
	}

	gitHandler := git.NewGitHandler()
	results, err := engine.DetectProject(projectPath, files, gitHandler, options)
	if err != nil {
		t.Fatalf("DetectProject failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("no detection results returned")
	}

	// Should detect Node.js
	found := false
	for _, result := range results {
		if result.Language == "node" {
			found = true
			if result.Confidence <= 0 {
				t.Errorf("expected positive confidence, got %f", result.Confidence)
			}
			break
		}
	}
	if !found {
		t.Error("expected to detect Node.js project")
	}
}

func TestDetectProject_Python(t *testing.T) {
	engine := NewDetectionEngine()

	projectPath := filepath.Join("../../testdata/python/basic")
	files := []types.FileInfo{
		{
			Path:        "main.py",
			Name:        "main.py",
			IsDirectory: false,
			Size:        func() *int64 { v := int64(50); return &v }(),
			Extension:   func() *string { v := ".py"; return &v }(),
		},
		{
			Path:        "requirements.txt",
			Name:        "requirements.txt",
			IsDirectory: false,
			Size:        func() *int64 { v := int64(30); return &v }(),
			Extension:   func() *string { v := ".txt"; return &v }(),
		},
	}

	options := &types.CLIOptions{
		Verbose: true,
	}

	gitHandler := git.NewGitHandler()
	results, err := engine.DetectProject(projectPath, files, gitHandler, options)
	if err != nil {
		t.Fatalf("DetectProject failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("no detection results returned")
	}

	// Should detect Python
	found := false
	for _, result := range results {
		if result.Language == "python" {
			found = true
			if result.Confidence <= 0 {
				t.Errorf("expected positive confidence, got %f", result.Confidence)
			}
			break
		}
	}
	if !found {
		t.Error("expected to detect Python project")
	}
}

func TestDetectProject_Static(t *testing.T) {
	engine := NewDetectionEngine()

	projectPath := filepath.Join("../../testdata/static/basic")
	files := []types.FileInfo{
		{
			Path:        "index.html",
			Name:        "index.html",
			IsDirectory: false,
			Size:        func() *int64 { v := int64(100); return &v }(),
			Extension:   func() *string { v := ".html"; return &v }(),
		},
		{
			Path:        "style.css",
			Name:        "style.css",
			IsDirectory: false,
			Size:        func() *int64 { v := int64(50); return &v }(),
			Extension:   func() *string { v := ".css"; return &v }(),
		},
	}

	options := &types.CLIOptions{
		Verbose: true,
	}

	gitHandler := git.NewGitHandler()
	results, err := engine.DetectProject(projectPath, files, gitHandler, options)
	if err != nil {
		t.Fatalf("DetectProject failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("no detection results returned")
	}

	// Should detect static files
	found := false
	for _, result := range results {
		if result.Language == "staticfile" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to detect static file project")
	}
}

func TestGetBestResult(t *testing.T) {
	engine := NewDetectionEngine()

	results := []*types.DetectResult{
		{
			Language:   "node",
			Framework:  "express",
			Confidence: 0.8,
		},
		{
			Language:   "python",
			Framework:  "flask",
			Confidence: 0.6,
		},
		{
			Language:   "staticfile",
			Framework:  "",
			Confidence: 0.3,
		},
	}

	best := engine.GetBestResult(results)
	if best == nil {
		t.Fatal("GetBestResult returned nil")
	}

	if best.Language != "node" {
		t.Errorf("expected best result to be 'node', got '%s'", best.Language)
	}

	if best.Confidence != 0.8 {
		t.Errorf("expected confidence 0.8, got %f", best.Confidence)
	}
}

func TestFilterResults(t *testing.T) {
	engine := NewDetectionEngine()

	results := []*types.DetectResult{
		{
			Matched:    true,
			Language:   "node",
			Confidence: 0.8,
		},
		{
			Matched:    true,
			Language:   "python",
			Confidence: 0.6,
		},
		{
			Matched:    true,
			Language:   "staticfile",
			Confidence: 0.3,
		},
	}

	filtered := engine.FilterResults(results, 0.5)
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered results, got %d", len(filtered))
	}

	for _, result := range filtered {
		if result.Confidence < 0.5 {
			t.Errorf("filtered result has confidence %f, expected >= 0.5", result.Confidence)
		}
	}
}

func TestGetDetectionStats(t *testing.T) {
	engine := NewDetectionEngine()

	results := []*types.DetectResult{
		{
			Matched:    true,
			Language:   "node",
			Framework:  "express",
			Confidence: 0.8,
		},
		{
			Matched:    true,
			Language:   "python",
			Framework:  "flask",
			Confidence: 0.6,
		},
	}

	stats := engine.GetDetectionStats(results)
	if stats == nil {
		t.Fatal("GetDetectionStats returned nil")
	}

	if stats.Total != 2 {
		t.Errorf("expected total 2, got %d", stats.Total)
	}

	if stats.Detected != 2 {
		t.Errorf("expected detected 2, got %d", stats.Detected)
	}

	expectedAvg := (0.8 + 0.6) / 2
	if stats.AvgConfidence != expectedAvg {
		t.Errorf("expected average confidence %f, got %f", expectedAvg, stats.AvgConfidence)
	}

	if len(stats.Languages) != 2 {
		t.Errorf("expected 2 languages, got %d", len(stats.Languages))
	}

	if len(stats.Frameworks) != 2 {
		t.Errorf("expected 2 frameworks, got %d", len(stats.Frameworks))
	}
}
