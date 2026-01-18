package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// RegistrySource represents a configured registry source.
type RegistrySource struct {
	// Name is the unique identifier for this source.
	Name string `yaml:"name"`

	// Type is the source type (local, git). Defaults to "local" if Path is set, "git" if URL is set.
	Type string `yaml:"type,omitempty"`

	// Path is the local filesystem path (for local sources).
	Path string `yaml:"path,omitempty"`

	// URL is the Git repository URL (for git sources).
	URL string `yaml:"url,omitempty"`

	// Enabled controls whether this source is active. Defaults to true.
	Enabled *bool `yaml:"enabled,omitempty"`
}

// IsEnabled returns whether this source is enabled (defaults to true).
func (s *RegistrySource) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}

	return *s.Enabled
}

// IsLocal returns true if this is a local filesystem source.
func (s *RegistrySource) IsLocal() bool {
	if s.Type == "local" {
		return true
	}

	if s.Type == "" && s.Path != "" {
		return true
	}

	return false
}

// IsGit returns true if this is a Git repository source.
func (s *RegistrySource) IsGit() bool {
	if s.Type == "git" {
		return true
	}

	if s.Type == "" && s.URL != "" {
		return true
	}

	return false
}

// SkillsmithConfig represents the user's skillsmith configuration.
type SkillsmithConfig struct {
	// Registries is the list of configured registry sources.
	Registries []RegistrySource `yaml:"registries"`
}

// DefaultConfig returns the default configuration with only the builtin registry.
func DefaultConfig() *SkillsmithConfig {
	return &SkillsmithConfig{
		Registries: []RegistrySource{},
	}
}

// GetConfigPath returns the path to the skillsmith config file.
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "skillsmith", "config.yaml"), nil
}

// LoadConfig loads the skillsmith configuration from disk.
// Returns default config if file doesn't exist.
func LoadConfig() (*SkillsmithConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path) //nolint:gosec // path is constructed internally
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}

		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg SkillsmithConfig

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig writes the configuration to disk.
func SaveConfig(cfg *SkillsmithConfig) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	err = EnsureDir(path)
	if err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	err = os.WriteFile(path, data, filePermissions)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// filePermissions for config files.
const filePermissions = 0o600
