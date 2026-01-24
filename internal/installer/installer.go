package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
	"github.com/monke/skillsmith/internal/transformer"
)

// filePermissions is the default permission for created files.
const filePermissions = 0o600

// Result represents the outcome of an installation.
type Result struct {
	Success bool
}

// GetInstallPath returns the full path where an item should be installed.
func GetInstallPath(item registry.Item, tool registry.Tool, scope config.Scope) (string, error) {
	paths, err := config.GetPaths(string(tool))
	if err != nil {
		return "", fmt.Errorf("get paths: %w", err)
	}

	var baseDir string
	if scope == config.ScopeGlobal {
		baseDir = paths.GlobalDir
	} else {
		baseDir = paths.LocalDir
	}

	filename := item.Name + ".md"

	switch item.Type {
	case registry.ItemTypeAgent:
		if paths.AgentsSubdir == "" {
			// For tools without agent subdirs (like Claude), use skills instead
			skillDir := filepath.Join(baseDir, paths.SkillsSubdir, item.Name)

			return filepath.Join(skillDir, "SKILL.md"), nil
		}

		return filepath.Join(baseDir, paths.AgentsSubdir, filename), nil

	case registry.ItemTypeSkill:
		// Skills go in skills/<name>/SKILL.md
		skillDir := filepath.Join(baseDir, paths.SkillsSubdir, item.Name)

		return filepath.Join(skillDir, "SKILL.md"), nil

	default:
		return filepath.Join(baseDir, filename), nil
	}
}

// Install installs an item for a specific tool to the specified scope.
func Install(item registry.Item, tool registry.Tool, scope config.Scope, force bool) (*Result, error) {
	// Check compatibility
	if !item.IsCompatibleWith(tool) {
		return &Result{Success: false}, nil
	}

	path, err := GetInstallPath(item, tool, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get install path: %w", err)
	}

	// Check if file already exists
	if config.Exists(path) && !force {
		return &Result{Success: false}, nil
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

	// Save hash to metadata
	meta, err := LoadMetadata(tool, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	meta.Set(item.Name, InstalledItem{
		Hash:        ComputeHash(content),
		InstalledAt: time.Now(),
	})

	err = SaveMetadata(tool, scope, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return &Result{Success: true}, nil
}

// Uninstall removes an installed item for a specific tool.
func Uninstall(item registry.Item, tool registry.Tool, scope config.Scope) (*Result, error) {
	path, err := GetInstallPath(item, tool, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get install path: %w", err)
	}

	if !config.Exists(path) {
		return &Result{Success: false}, nil
	}

	err = os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("failed to remove file: %w", err)
	}

	// Remove from metadata (best effort, file is already removed)
	meta, _ := LoadMetadata(tool, scope)
	if meta != nil {
		meta.Remove(item.Name)
		_ = SaveMetadata(tool, scope, meta)
	}

	return &Result{Success: true}, nil
}

// GetItemState determines the installation state of an item.
func GetItemState(item registry.Item, tool registry.Tool, scope config.Scope) (ItemState, string, error) {
	path, err := GetInstallPath(item, tool, scope)
	if err != nil {
		return StateNotInstalled, "", fmt.Errorf("get install path: %w", err)
	}

	// Check if file exists
	if !config.Exists(path) {
		return StateNotInstalled, path, nil
	}

	// Load metadata
	meta, metaErr := LoadMetadata(tool, scope)
	if metaErr != nil {
		// If metadata can't be loaded, assume file exists but state unknown
		// Treat as modified since we don't know the original hash
		return StateModified, path, nil //nolint:nilerr // intentional: treat as modified
	}

	// Get the stored hash for this item
	installedInfo, hasMetadata := meta.Get(item.Name)

	// Compute current file hash
	fileHash, hashErr := ComputeFileHash(path)
	if hashErr != nil {
		return StateModified, path, nil //nolint:nilerr // intentional: treat as modified
	}

	// Compute what the registry version would look like
	registryContent, transformErr := transformer.Transform(item, tool)
	if transformErr != nil {
		return StateModified, path, nil //nolint:nilerr // intentional: treat as modified
	}

	registryHash := ComputeHash(registryContent)

	// If no metadata, we don't know the original installed version
	// Compare file to registry to make best guess
	if !hasMetadata {
		if fileHash == registryHash {
			return StateUpToDate, path, nil
		}
		// File exists but doesn't match registry and no metadata
		// Could be modified or could be old version - treat as modified to be safe
		return StateModified, path, nil
	}

	// We have metadata - compare all three hashes
	installedHash := installedInfo.Hash
	fileMatchesInstalled := fileHash == installedHash
	registryMatchesInstalled := registryHash == installedHash

	switch {
	case fileMatchesInstalled && registryMatchesInstalled:
		// All three match - up to date
		return StateUpToDate, path, nil
	case fileMatchesInstalled && !registryMatchesInstalled:
		// File unchanged, but registry has new version
		return StateUpdateAvailable, path, nil
	case !fileMatchesInstalled && registryMatchesInstalled:
		// File was modified, registry unchanged
		return StateModified, path, nil
	default:
		// File modified AND registry has new version
		return StateModifiedWithUpdate, path, nil
	}
}
