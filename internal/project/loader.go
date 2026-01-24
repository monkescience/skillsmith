package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Common errors.
var (
	ErrNotFound      = errors.New("project config not found")
	ErrAlreadyExists = errors.New("project config already exists")
)

// filePermissions for project config files.
const filePermissions = 0o644

// Load loads the project configuration from the current directory.
// It searches up the directory tree to find a .skillsmith.yaml file.
func Load() (*Config, string, error) {
	return LoadFrom(".")
}

// LoadFrom loads the project configuration starting from the given directory.
// It searches up the directory tree to find a .skillsmith.yaml file.
// Returns the config and the directory where it was found.
func LoadFrom(startDir string) (*Config, string, error) {
	absPath, err := filepath.Abs(startDir)
	if err != nil {
		return nil, "", fmt.Errorf("resolve path: %w", err)
	}

	// Walk up the directory tree looking for the config file
	dir := absPath
	for {
		configPath := filepath.Join(dir, ConfigFileName)

		if fileExists(configPath) {
			cfg, err := loadFile(configPath)
			if err != nil {
				return nil, "", err
			}

			return cfg, dir, nil
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, config not found
			return nil, "", ErrNotFound
		}

		dir = parent
	}
}

// LoadFromDir loads the project configuration from a specific directory.
// Unlike LoadFrom, this does not search parent directories.
func LoadFromDir(dir string) (*Config, error) {
	configPath := filepath.Join(dir, ConfigFileName)

	if !fileExists(configPath) {
		return nil, ErrNotFound
	}

	return loadFile(configPath)
}

// Exists checks if a project config exists in the current directory
// or any parent directory.
func Exists() bool {
	_, _, err := Load()
	return err == nil
}

// ExistsInDir checks if a project config exists in the specific directory.
func ExistsInDir(dir string) bool {
	configPath := filepath.Join(dir, ConfigFileName)
	return fileExists(configPath)
}

// Save saves the project configuration to the specified directory.
func Save(cfg *Config, dir string) error {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, filePermissions)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// SaveWithHeader saves the project configuration with a helpful header comment.
func SaveWithHeader(cfg *Config, dir string) error {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	header := `# Skillsmith Project Configuration
# https://github.com/monke/skillsmith
#
# Limit to specific tools (optional):
#   tools: [claude, opencode]
#
# Add project-specific registries (optional):
#   registries:
#     - name: team-skills
#       url: https://github.com/team/skills.git

`

	content := header + string(data)

	err = os.WriteFile(configPath, []byte(content), filePermissions)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// Init creates a new project configuration in the specified directory.
// Returns an error if a config already exists.
func Init(dir string) (*Config, error) {
	if ExistsInDir(dir) {
		return nil, ErrAlreadyExists
	}

	cfg := &Config{
		Tools:  []string{},
		Skills: []string{},
		Agents: []string{},
	}

	err := SaveWithHeader(cfg, dir)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// InitWithConfig creates a new project configuration with the given config.
// Returns an error if a config already exists.
func InitWithConfig(dir string, cfg *Config) error {
	if ExistsInDir(dir) {
		return ErrAlreadyExists
	}

	return Save(cfg, dir)
}

// GetConfigPath returns the path to the project config file in the given directory.
func GetConfigPath(dir string) string {
	return filepath.Join(dir, ConfigFileName)
}

// loadFile loads and parses a project config file.
func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is validated by caller
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
