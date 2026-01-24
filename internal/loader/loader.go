// Package loader provides functions to load registries from configuration.
package loader

import (
	"fmt"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/project"
	"github.com/monke/skillsmith/internal/registry"
)

// LoadFromConfig creates a MultiRegistry from the user's config.
// Sources are loaded in order with last source winning for duplicates:
// 1. Builtin (embedded) - lowest priority
// 2. Global registries from ~/.config/skillsmith/config.yaml
// 3. Project registries from ./.skillsmith.yaml - highest priority
func LoadFromConfig() (*registry.MultiRegistry, error) {
	return LoadFromConfigWithProject(true)
}

// LoadFromConfigWithProject creates a MultiRegistry with optional project config.
// If includeProject is true and a project config exists, project registries are included.
func LoadFromConfigWithProject(includeProject bool) (*registry.MultiRegistry, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	multi := registry.NewMultiRegistry()

	// 1. Add builtin source first (lowest priority, can be overridden)
	multi.AddSource(registry.NewEmbeddedSource("builtin"))

	// 2. Add global configured sources
	addRegistrySources(multi, cfg.Registries)

	// 3. Add project-specific sources (highest priority)
	if includeProject {
		projectCfg, _, err := project.Load()
		if err == nil && projectCfg != nil {
			addRegistrySources(multi, projectCfg.Registries)
		}
		// Ignore error - project config is optional
	}

	// Load all sources
	err = multi.Load()
	if err != nil {
		return nil, fmt.Errorf("load sources: %w", err)
	}

	return multi, nil
}

// LoadFromConfigOnly creates a MultiRegistry from global config only (no project).
func LoadFromConfigOnly() (*registry.MultiRegistry, error) {
	return LoadFromConfigWithProject(false)
}

// LoadBuiltinOnly creates a MultiRegistry with only the builtin embedded source.
func LoadBuiltinOnly() (*registry.MultiRegistry, error) {
	multi, err := registry.LoadWithBuiltin()
	if err != nil {
		return nil, fmt.Errorf("load builtin: %w", err)
	}

	return multi, nil
}

// addRegistrySources adds registry sources to a MultiRegistry.
func addRegistrySources(multi *registry.MultiRegistry, sources []config.RegistrySource) {
	for _, src := range sources {
		if !src.IsEnabled() {
			continue
		}

		switch {
		case src.IsLocal():
			multi.AddSource(registry.NewLocalSource(src.Name, src.Path))
		case src.IsGit():
			multi.AddSource(registry.NewGitSource(src.Name, src.URL))
		}
	}
}
