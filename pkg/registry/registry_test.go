package registry

import (
	"testing"

	"github.com/labring/devbox-pack/pkg/detector"
)

func TestNewProviderRegistry(t *testing.T) {
	registry := NewProviderRegistry()

	if registry == nil {
		t.Fatal("NewProviderRegistry() returned nil")
	}

	if registry.providers == nil {
		t.Fatal("providers map should be initialized")
	}

	// Check that providers are initialized
	if len(registry.providers) == 0 {
		t.Error("expected providers to be initialized")
	}

	// Check for expected providers
	expectedProviders := []string{
		"node", "python", "java", "go", "php",
		"ruby", "rust", "deno", "shell", "staticfile",
	}

	for _, providerName := range expectedProviders {
		if _, exists := registry.providers[providerName]; !exists {
			t.Errorf("expected provider '%s' to be registered", providerName)
		}
	}
}

func TestProviderRegistry_GetProvider_ValidProvider(t *testing.T) {
	registry := NewProviderRegistry()

	tests := []struct {
		name         string
		providerName string
	}{
		{"node provider", "node"},
		{"python provider", "python"},
		{"java provider", "java"},
		{"go provider", "go"},
		{"php provider", "php"},
		{"ruby provider", "ruby"},
		{"rust provider", "rust"},
		{"deno provider", "deno"},
		{"shell provider", "shell"},
		{"staticfile provider", "staticfile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := registry.GetProvider(tt.providerName)

			if provider == nil {
				t.Fatalf("GetProvider(%s) returned nil provider", tt.providerName)
			}

			if provider.GetName() != tt.providerName {
				t.Errorf("expected provider name '%s', got '%s'", tt.providerName, provider.GetName())
			}
		})
	}
}

func TestProviderRegistry_GetProvider_InvalidProvider(t *testing.T) {
	registry := NewProviderRegistry()

	invalidProviders := []string{
		"nonexistent",
		"invalid",
		"",
		"PYTHON", // case sensitive
		"Node",   // case sensitive
	}

	for _, providerName := range invalidProviders {
		t.Run("invalid_"+providerName, func(t *testing.T) {
			provider := registry.GetProvider(providerName)

			if provider != nil {
				t.Errorf("expected nil provider for invalid provider '%s'", providerName)
			}
		})
	}
}

func TestProviderRegistry_GetAllProviders(t *testing.T) {
	registry := NewProviderRegistry()

	providersMap := registry.GetAllProviders()

	if len(providersMap) == 0 {
		t.Fatal("GetAllProviders() returned empty map")
	}

	// Check that we have the expected number of providers
	expectedCount := 10 // node, python, java, go, php, ruby, rust, deno, shell, staticfile
	if len(providersMap) != expectedCount {
		t.Errorf("expected %d providers, got %d", expectedCount, len(providersMap))
	}

	// Check that all providers are valid
	providerNames := make(map[string]bool)
	for name, provider := range providersMap {
		if provider == nil {
			t.Error("found nil provider in GetAllProviders() result")
			continue
		}

		if provider.GetName() != name {
			t.Errorf("provider name mismatch: key '%s', provider name '%s'", name, provider.GetName())
			continue
		}

		if providerNames[name] {
			t.Errorf("found duplicate provider name: %s", name)
		}
		providerNames[name] = true

		// Check that provider has valid priority
		priority := provider.GetPriority()
		if priority <= 0 {
			t.Errorf("provider %s has invalid priority: %d", name, priority)
		}
	}

	// Check for expected provider names
	expectedProviders := []string{
		"node", "python", "java", "go", "php",
		"ruby", "rust", "deno", "shell", "staticfile",
	}

	for _, expectedName := range expectedProviders {
		if !providerNames[expectedName] {
			t.Errorf("expected provider '%s' not found in GetAllProviders() result", expectedName)
		}
	}
}

func TestProviderRegistry_GetAllProviders_Sorted(t *testing.T) {
	registry := NewProviderRegistry()

	providersMap := registry.GetAllProviders()

	// Test that all providers in the map have valid priorities
	for name, provider := range providersMap {
		priority := provider.GetPriority()
		if priority <= 0 {
			t.Errorf("provider %s has invalid priority: %d", name, priority)
		}
	}
}

func TestProviderRegistry_InitializeProviders(t *testing.T) {
	registry := &ProviderRegistry{
		providers: make(map[string]detector.Provider),
	}

	// Call initializeProviders
	registry.initializeProviders()

	// Check that providers are initialized
	if len(registry.providers) == 0 {
		t.Error("initializeProviders() did not initialize any providers")
	}

	// Check that each provider is properly initialized
	for name, provider := range registry.providers {
		if provider == nil {
			t.Errorf("provider '%s' is nil after initialization", name)
			continue
		}

		if provider.GetName() != name {
			t.Errorf("provider name mismatch: key '%s', provider name '%s'", name, provider.GetName())
		}

		if provider.GetPriority() <= 0 {
			t.Errorf("provider '%s' has invalid priority: %d", name, provider.GetPriority())
		}
	}
}

func TestProviderRegistry_Concurrent(t *testing.T) {
	registry := NewProviderRegistry()

	// Test concurrent access to GetProvider
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Try to get different providers concurrently
			providers := []string{"node", "python", "java", "go", "php"}
			for _, name := range providers {
				provider := registry.GetProvider(name)
				if provider == nil {
					t.Errorf("concurrent GetProvider(%s) returned nil", name)
					return
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestProviderRegistry_GetAllProviders_Concurrent(t *testing.T) {
	registry := NewProviderRegistry()

	// Test concurrent access to GetAllProviders
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			defer func() { done <- true }()

			providersMap := registry.GetAllProviders()
			if len(providersMap) == 0 {
				t.Error("concurrent GetAllProviders() returned empty map")
				return
			}

			// Verify providers are valid
			for _, provider := range providersMap {
				if provider == nil {
					t.Error("concurrent GetAllProviders() returned nil provider")
					return
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}
}
