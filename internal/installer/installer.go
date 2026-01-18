package installer

import (
	"fmt"
	"os"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
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

// Install installs an item to the specified scope.
func Install(item registry.Item, scope config.Scope, force bool) (*Result, error) {
	path, err := config.GetInstallPath(item, scope)
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

	// Ensure parent directory exists
	err = config.EnsureDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the content
	err = os.WriteFile(path, []byte(item.Content), filePermissions)
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

// Uninstall removes an installed item.
func Uninstall(item registry.Item, scope config.Scope) (*Result, error) {
	path, err := config.GetInstallPath(item, scope)
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

// IsInstalled checks if an item is already installed at the given scope.
func IsInstalled(item registry.Item, scope config.Scope) (bool, string, error) {
	path, err := config.GetInstallPath(item, scope)
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

// GetInstallStatus returns installation status for both scopes.
func GetInstallStatus(item registry.Item) (*InstallStatus, error) {
	localInstalled, localPath, err := IsInstalled(item, config.ScopeLocal)
	if err != nil {
		return nil, err
	}

	globalInstalled, globalPath, err := IsInstalled(item, config.ScopeGlobal)
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
