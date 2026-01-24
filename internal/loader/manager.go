package loader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/installer"
	"github.com/monke/skillsmith/internal/registry"
)

// Manager errors.
var (
	ErrItemNotFound        = errors.New("item not found")
	ErrItemNotCompatible   = errors.New("item not compatible with tool")
	ErrPathNotExist        = errors.New("path does not exist")
	ErrPathNotDir          = errors.New("path is not a directory")
	ErrRegistryExists      = errors.New("registry with this name already exists")
	ErrCannotRemoveBuiltin = errors.New("cannot remove the builtin registry")
	ErrRegistryNotFound    = errors.New("registry not found")
)

// Manager provides the main API for working with the registry.
// It coordinates loading, installation, and configuration.
type Manager struct {
	registry *registry.Registry
}

// NewManager creates a new Manager, loading registries from config.
func NewManager() (*Manager, error) {
	multi, err := LoadFromConfig()
	if err != nil {
		return nil, fmt.Errorf("load registry: %w", err)
	}

	return &Manager{
		registry: multi.Registry(),
	}, nil
}

// NewManagerWithRegistry creates a Manager with a pre-loaded registry.
// Useful for testing or when registry is already loaded.
func NewManagerWithRegistry(reg *registry.Registry) *Manager {
	return &Manager{
		registry: reg,
	}
}

// Registry returns the underlying registry.
func (m *Manager) Registry() *registry.Registry {
	return m.registry
}

// Reload reloads the registry from config.
// Call this after adding or removing registry sources.
func (m *Manager) Reload() error {
	multi, err := LoadFromConfig()
	if err != nil {
		return fmt.Errorf("reload registry: %w", err)
	}

	m.registry = multi.Registry()

	return nil
}

// GetItem returns a single item by name.
func (m *Manager) GetItem(name string) (*registry.Item, error) {
	for i := range m.registry.Items {
		if m.registry.Items[i].Name == name {
			return &m.registry.Items[i], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrItemNotFound, name)
}

// ListItems returns items filtered by tool, with optional type filter.
func (m *Manager) ListItems(tool registry.Tool, itemType registry.ItemType) []registry.Item {
	if itemType != "" {
		return m.registry.ByToolAndType(tool, itemType)
	}

	return m.registry.ByTool(tool)
}

// ListItemsWithState returns items with their installation state.
func (m *Manager) ListItemsWithState(
	tool registry.Tool, scope config.Scope, itemType registry.ItemType,
) []installer.ItemWithState {
	items := m.ListItems(tool, itemType)
	result := make([]installer.ItemWithState, 0, len(items))

	for _, item := range items {
		state, path, _ := installer.GetItemState(item, tool, scope)
		result = append(result, installer.ItemWithState{
			Item:        item,
			State:       state,
			InstallPath: path,
		})
	}

	return result
}

// Install installs an item.
func (m *Manager) Install(
	itemName string, tool registry.Tool, scope config.Scope, force bool,
) (*installer.Result, string, error) {
	item, err := m.GetItem(itemName)
	if err != nil {
		return nil, "", err
	}

	// Check compatibility
	if !item.IsCompatibleWith(tool) {
		return &installer.Result{Success: false}, "", fmt.Errorf("%w: %s", ErrItemNotCompatible, tool)
	}

	// Get path for result
	path, err := config.GetInstallPath(item.Name, item.Type, tool, scope)
	if err != nil {
		return nil, "", fmt.Errorf("get install path: %w", err)
	}

	// Install
	result, err := installer.Install(*item, tool, scope, force)
	if err != nil {
		return nil, path, fmt.Errorf("install: %w", err)
	}

	return result, path, nil
}

// Uninstall removes an installed item.
func (m *Manager) Uninstall(
	itemName string, tool registry.Tool, scope config.Scope,
) (*installer.Result, string, error) {
	item, err := m.GetItem(itemName)
	if err != nil {
		return nil, "", err
	}

	// Get path for result
	path, err := config.GetInstallPath(item.Name, item.Type, tool, scope)
	if err != nil {
		return nil, "", fmt.Errorf("get install path: %w", err)
	}

	// Uninstall
	result, err := installer.Uninstall(*item, tool, scope)
	if err != nil {
		return nil, path, fmt.Errorf("uninstall: %w", err)
	}

	return result, path, nil
}

// RegistryInfo represents a configured registry source.
type RegistryInfo struct {
	Name    string
	Type    string // "builtin", "local", "git"
	Path    string
	URL     string
	Enabled bool
}

// ListRegistries returns all configured registry sources.
func (m *Manager) ListRegistries() ([]RegistryInfo, error) {
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
func (m *Manager) AddRegistry(name, path string) error {
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
	if name == "builtin" {
		return fmt.Errorf("%w: %s", ErrRegistryExists, name)
	}

	for _, reg := range cfg.Registries {
		if reg.Name == name {
			return fmt.Errorf("%w: %s", ErrRegistryExists, name)
		}
	}

	cfg.Registries = append(cfg.Registries, config.RegistrySource{
		Name: name,
		Path: absPath,
	})

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// RemoveRegistry removes a registry by name.
func (m *Manager) RemoveRegistry(name string) error {
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
