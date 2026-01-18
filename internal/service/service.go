package service

import (
	"fmt"

	"github.com/monke/skillsmith/internal/loader"
	"github.com/monke/skillsmith/internal/registry"
)

// Service provides the application-level API for skillsmith.
// It coordinates between registry, installer, and config packages.
type Service struct {
	registry *registry.Registry
}

// New creates a new Service, loading registries from config.
func New() (*Service, error) {
	multi, err := loader.LoadFromConfig()
	if err != nil {
		return nil, fmt.Errorf("load registry: %w", err)
	}

	return &Service{
		registry: multi.Registry(),
	}, nil
}

// NewWithRegistry creates a Service with a pre-loaded registry.
// Useful for testing or when registry is already loaded.
func NewWithRegistry(reg *registry.Registry) *Service {
	return &Service{
		registry: reg,
	}
}

// Reload reloads the registry from config.
// Call this after adding or removing registry sources.
func (s *Service) Reload() error {
	multi, err := loader.LoadFromConfig()
	if err != nil {
		return fmt.Errorf("reload registry: %w", err)
	}

	s.registry = multi.Registry()

	return nil
}
