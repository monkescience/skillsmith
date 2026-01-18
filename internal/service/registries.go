package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/monke/skillsmith/internal/config"
)

// Registry management errors.
var (
	ErrPathNotExist        = errors.New("path does not exist")
	ErrPathNotDir          = errors.New("path is not a directory")
	ErrRegistryExists      = errors.New("registry with this name already exists")
	ErrCannotRemoveBuiltin = errors.New("cannot remove the builtin registry")
	ErrRegistryNotFound    = errors.New("registry not found")
)

// ListRegistries returns all configured registry sources.
func (s *Service) ListRegistries() ([]RegistryInfo, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Always include builtin as first entry
	registries := []RegistryInfo{
		{
			Name:    "builtin",
			Type:    "builtin",
			Enabled: true,
		},
	}

	for _, reg := range cfg.Registries {
		info := RegistryInfo{
			Name:    reg.Name,
			Path:    reg.Path,
			URL:     reg.URL,
			Enabled: reg.IsEnabled(),
		}

		switch {
		case reg.IsLocal():
			info.Type = "local"
		case reg.IsGit():
			info.Type = "git"
		default:
			info.Type = "unknown"
		}

		registries = append(registries, info)
	}

	return registries, nil
}

// AddRegistry adds a new local registry source.
func (s *Service) AddRegistry(input AddRegistryInput) error {
	path := input.Path

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}

		path = filepath.Join(homeDir, path[2:])
	}

	// Make path absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	// Verify path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrPathNotExist, absPath)
		}

		return fmt.Errorf("check path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrPathNotDir, absPath)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Check for duplicate name
	if input.Name == "builtin" {
		return fmt.Errorf("%w: %s", ErrRegistryExists, input.Name)
	}

	for _, reg := range cfg.Registries {
		if reg.Name == input.Name {
			return fmt.Errorf("%w: %s", ErrRegistryExists, input.Name)
		}
	}

	cfg.Registries = append(cfg.Registries, config.RegistrySource{
		Name: input.Name,
		Path: absPath,
	})

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// RemoveRegistry removes a registry by name.
func (s *Service) RemoveRegistry(name string) error {
	if name == "builtin" {
		return ErrCannotRemoveBuiltin
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	found := false
	newRegistries := make([]config.RegistrySource, 0, len(cfg.Registries))

	for _, reg := range cfg.Registries {
		if reg.Name == name {
			found = true

			continue
		}

		newRegistries = append(newRegistries, reg)
	}

	if !found {
		return fmt.Errorf("%w: %s", ErrRegistryNotFound, name)
	}

	cfg.Registries = newRegistries

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
