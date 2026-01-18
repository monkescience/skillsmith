// Package loader provides functions to load registries from configuration.
package loader

import (
	"fmt"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
)

// LoadFromConfig creates a MultiRegistry from the user's config.
// Always includes the builtin embedded source first (highest priority).
func LoadFromConfig() (*registry.MultiRegistry, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	multi := registry.NewMultiRegistry()

	// Always add builtin source first (highest priority)
	multi.AddSource(registry.NewEmbeddedSource("builtin"))

	// Add configured sources
	for _, src := range cfg.Registries {
		if !src.IsEnabled() {
			continue
		}

		switch {
		case src.IsLocal():
			multi.AddSource(registry.NewLocalSource(src.Name, src.Path))
		case src.IsGit():
			// Git source will be implemented in Phase 2
			// For now, skip (could log this)
			continue
		}
	}

	// Load all sources
	err = multi.Load()
	if err != nil {
		return nil, fmt.Errorf("load sources: %w", err)
	}

	return multi, nil
}

// LoadBuiltinOnly creates a MultiRegistry with only the builtin embedded source.
func LoadBuiltinOnly() (*registry.MultiRegistry, error) {
	multi, err := registry.LoadWithBuiltin()
	if err != nil {
		return nil, fmt.Errorf("load builtin: %w", err)
	}

	return multi, nil
}
