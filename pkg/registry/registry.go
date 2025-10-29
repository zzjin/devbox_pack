/**
 * DevBox Pack Execution Plan Generator - Provider Registry
 */

package registry

import (
	"github.com/labring/devbox-pack/pkg/detector"
	"github.com/labring/devbox-pack/pkg/providers"
)

// ProviderRegistry Provider registry
type ProviderRegistry struct {
	providers map[string]detector.Provider
}

// NewProviderRegistry creates a new Provider registry
func NewProviderRegistry() *ProviderRegistry {
	registry := &ProviderRegistry{
		providers: make(map[string]detector.Provider),
	}
	registry.initializeProviders()
	return registry
}

// initializeProviders initializes all Providers
func (r *ProviderRegistry) initializeProviders() {
	// Register providers for various languages
	r.providers["node"] = providers.NewNodeProvider()
	r.providers["python"] = providers.NewPythonProvider()
	r.providers["java"] = providers.NewJavaProvider()
	r.providers["go"] = providers.NewGoProvider()
	r.providers["php"] = providers.NewPHPProvider()
	r.providers["ruby"] = providers.NewRubyProvider()
	r.providers["rust"] = providers.NewRustProvider()
	r.providers["deno"] = providers.NewDenoProvider()
	r.providers["shell"] = providers.NewShellProvider()
	r.providers["staticfile"] = providers.NewStaticFileProvider()
}

// GetProvider gets Provider by name
func (r *ProviderRegistry) GetProvider(name string) detector.Provider {
	return r.providers[name]
}

// GetAllProviders gets all Providers
func (r *ProviderRegistry) GetAllProviders() map[string]detector.Provider {
	return r.providers
}
