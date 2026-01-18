package installer

import (
	"fmt"
	"os"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
	"github.com/monke/skillsmith/internal/transformer"
)

// filePermissions is the default permission for created files.
const filePermissions = 0o600

// Result represents the outcome of an installation.
type Result struct {
	Success bool
	Path    string
	Message string
	Existed bool
}

// Install installs an item for a specific tool to the specified scope.
func Install(item registry.Item, tool registry.Tool, scope config.Scope, force bool) (*Result, error) {
	// Check compatibility
	if !item.IsCompatibleWith(tool) {
		return &Result{
			Success: false,
			Message: fmt.Sprintf("item not compatible with %s", tool),
		}, nil
	}

	path, err := config.GetInstallPath(item, tool, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get install path: %w", err)
	}

	// Check if file already exists
	if config.Exists(path) && !force {
		return &Result{
			Success: false,
			Path:    path,
			Message: "file already exists (use force to overwrite)",
			Existed: true,
		}, nil
	}

	// Transform content for the target tool
	content, err := transformer.Transform(item, tool)
	if err != nil {
		return nil, fmt.Errorf("failed to transform content: %w", err)
	}

	// Ensure parent directory exists
	err = config.EnsureDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the content
	err = os.WriteFile(path, []byte(content), filePermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &Result{
		Success: true,
		Path:    path,
		Message: "installed successfully",
		Existed: false,
	}, nil
}

// Uninstall removes an installed item for a specific tool.
func Uninstall(item registry.Item, tool registry.Tool, scope config.Scope) (*Result, error) {
	path, err := config.GetInstallPath(item, tool, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get install path: %w", err)
	}

	if !config.Exists(path) {
		return &Result{
			Success: false,
			Path:    path,
			Message: "file does not exist",
			Existed: false,
		}, nil
	}

	err = os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("failed to remove file: %w", err)
	}

	return &Result{
		Success: true,
		Path:    path,
		Message: "uninstalled successfully",
		Existed: true,
	}, nil
}

// IsInstalled checks if an item is already installed for a specific tool at the given scope.
func IsInstalled(item registry.Item, tool registry.Tool, scope config.Scope) (bool, string, error) {
	path, err := config.GetInstallPath(item, tool, scope)
	if err != nil {
		return false, "", fmt.Errorf("get install path: %w", err)
	}

	return config.Exists(path), path, nil
}

// InstallStatus represents installation status for both scopes.
type InstallStatus struct {
	LocalInstalled  bool
	LocalPath       string
	GlobalInstalled bool
	GlobalPath      string
}

// GetInstallStatus returns installation status for both scopes for a specific tool.
func GetInstallStatus(item registry.Item, tool registry.Tool) (*InstallStatus, error) {
	localInstalled, localPath, err := IsInstalled(item, tool, config.ScopeLocal)
	if err != nil {
		return nil, err
	}

	globalInstalled, globalPath, err := IsInstalled(item, tool, config.ScopeGlobal)
	if err != nil {
		return nil, err
	}

	return &InstallStatus{
		LocalInstalled:  localInstalled,
		LocalPath:       localPath,
		GlobalInstalled: globalInstalled,
		GlobalPath:      globalPath,
	}, nil
}
