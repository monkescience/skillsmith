package installer

import (
	"crypto/md5" //nolint:gosec // MD5 used for change detection, not security
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
)

const metadataFilename = ".skillsmith.json"

// ItemState represents the installation state of an item.
type ItemState int

const (
	StateNotInstalled       ItemState = iota
	StateUpToDate                     // installed hash == file hash == registry hash
	StateUpdateAvailable              // file matches installed hash, but registry changed
	StateModified                     // file was modified locally
	StateModifiedWithUpdate           // file modified AND registry has update
)

// String returns a human-readable name for the state.
func (s ItemState) String() string {
	switch s {
	case StateNotInstalled:
		return "not installed"
	case StateUpToDate:
		return "up to date"
	case StateUpdateAvailable:
		return "update available"
	case StateModified:
		return "modified"
	case StateModifiedWithUpdate:
		return "modified with update"
	default:
		return "unknown"
	}
}

// IsInstalled returns true if the item is installed (any state except NotInstalled).
func (s ItemState) IsInstalled() bool {
	return s != StateNotInstalled
}

// HasUpdate returns true if an update is available.
func (s ItemState) HasUpdate() bool {
	return s == StateUpdateAvailable || s == StateModifiedWithUpdate
}

// IsModified returns true if the file was modified locally.
func (s ItemState) IsModified() bool {
	return s == StateModified || s == StateModifiedWithUpdate
}

// InstalledItem tracks metadata for an installed item.
type InstalledItem struct {
	Hash        string    `json:"hash"`
	InstalledAt time.Time `json:"installed_at"`
}

// Metadata stores installation state for all items.
type Metadata struct {
	Installed map[string]InstalledItem `json:"installed"`
}

// NewMetadata creates an empty metadata struct.
func NewMetadata() *Metadata {
	return &Metadata{
		Installed: make(map[string]InstalledItem),
	}
}

// Get returns the installed item metadata if it exists.
func (m *Metadata) Get(itemName string) (InstalledItem, bool) {
	item, ok := m.Installed[itemName]

	return item, ok
}

// Set stores metadata for an installed item.
func (m *Metadata) Set(itemName string, item InstalledItem) {
	if m.Installed == nil {
		m.Installed = make(map[string]InstalledItem)
	}

	m.Installed[itemName] = item
}

// Remove deletes metadata for an item.
func (m *Metadata) Remove(itemName string) {
	delete(m.Installed, itemName)
}

// GetMetadataPath returns the path to the metadata file for a tool and scope.
func GetMetadataPath(tool registry.Tool, scope config.Scope) (string, error) {
	paths, err := config.GetPaths(tool)
	if err != nil {
		return "", fmt.Errorf("get paths: %w", err)
	}

	var baseDir string

	if scope == config.ScopeGlobal {
		baseDir = paths.GlobalDir
	} else {
		baseDir = paths.LocalDir
	}

	return filepath.Join(baseDir, metadataFilename), nil
}

// LoadMetadata loads metadata from disk, or returns empty metadata if file doesn't exist.
func LoadMetadata(tool registry.Tool, scope config.Scope) (*Metadata, error) {
	path, err := GetMetadataPath(tool, scope)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path) //nolint:gosec // path is constructed internally
	if err != nil {
		if os.IsNotExist(err) {
			return NewMetadata(), nil
		}

		return nil, fmt.Errorf("read metadata: %w", err)
	}

	var meta Metadata

	err = json.Unmarshal(data, &meta)
	if err != nil {
		// If metadata is corrupted, return empty metadata
		return NewMetadata(), nil //nolint:nilerr // intentional: corrupted metadata is treated as empty
	}

	if meta.Installed == nil {
		meta.Installed = make(map[string]InstalledItem)
	}

	return &meta, nil
}

// SaveMetadata writes metadata to disk.
func SaveMetadata(tool registry.Tool, scope config.Scope, meta *Metadata) error {
	path, err := GetMetadataPath(tool, scope)
	if err != nil {
		return err
	}

	// If metadata is empty, remove the file instead of saving empty JSON
	if len(meta.Installed) == 0 {
		err = os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove empty metadata: %w", err)
		}

		return nil
	}

	// Ensure directory exists
	err = config.EnsureDir(path)
	if err != nil {
		return fmt.Errorf("ensure dir: %w", err)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	err = os.WriteFile(path, data, filePermissions)
	if err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	return nil
}

// ComputeHash computes an MD5 hash of the content.
func ComputeHash(content string) string {
	hash := md5.Sum([]byte(content)) //nolint:gosec // MD5 used for change detection, not security

	return hex.EncodeToString(hash[:])
}

// ComputeFileHash reads a file and computes its MD5 hash.
func ComputeFileHash(path string) (string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is from trusted source
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	return ComputeHash(string(data)), nil
}
