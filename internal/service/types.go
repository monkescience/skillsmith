package service

import "slices"

// Tool represents a supported AI coding tool.
type Tool string

const (
	ToolOpenCode Tool = "opencode"
	ToolClaude   Tool = "claude"
)

// Scope represents where items are installed.
type Scope string

const (
	ScopeLocal  Scope = "local"
	ScopeGlobal Scope = "global"
)

// ItemType represents the type of item.
type ItemType string

const (
	ItemTypeAgent ItemType = "agent"
	ItemTypeSkill ItemType = "skill"
)

// ItemState represents the installation state of an item.
type ItemState string

const (
	StateNotInstalled       ItemState = "not_installed"
	StateUpToDate           ItemState = "up_to_date"
	StateUpdateAvailable    ItemState = "update_available"
	StateModified           ItemState = "modified"
	StateModifiedWithUpdate ItemState = "modified_with_update"
)

// IsInstalled returns true if the item is installed (regardless of state).
func (s ItemState) IsInstalled() bool {
	return s != StateNotInstalled
}

// IsModified returns true if the item has local modifications.
func (s ItemState) IsModified() bool {
	return s == StateModified || s == StateModifiedWithUpdate
}

// HasUpdate returns true if an update is available.
func (s ItemState) HasUpdate() bool {
	return s == StateUpdateAvailable || s == StateModifiedWithUpdate
}

// Item represents an installable item (agent or skill).
type Item struct {
	Name          string
	Description   string
	Type          ItemType
	Category      string
	Compatibility []Tool
	Tags          []string
	Author        string
	Source        string // Registry source name
	Body          string // The actual content
}

// IsCompatibleWith checks if the item is compatible with a given tool.
func (i *Item) IsCompatibleWith(tool Tool) bool {
	return slices.Contains(i.Compatibility, tool)
}

// ItemWithState pairs an item with its installation state.
type ItemWithState struct {
	Item        Item
	State       ItemState
	InstallPath string
}

// RegistryInfo represents a configured registry source.
type RegistryInfo struct {
	Name    string
	Type    string // "builtin", "local", "git"
	Path    string
	URL     string
	Enabled bool
}

// InstallResult contains the outcome of an install operation.
type InstallResult struct {
	ItemName string
	Success  bool
	Path     string
	Skipped  bool
	Error    string
}

// UninstallResult contains the outcome of an uninstall operation.
type UninstallResult struct {
	ItemName string
	Success  bool
	Path     string
	Error    string
}

// ListItemsInput specifies filters for listing items.
type ListItemsInput struct {
	Tool  Tool     // Required: filter by tool compatibility
	Scope Scope    // Required: determines install paths and state
	Type  ItemType // Optional: filter by type (agent/skill)
}

// InstallInput specifies what to install.
type InstallInput struct {
	ItemName string
	Tool     Tool
	Scope    Scope
	Force    bool // Overwrite if exists
}

// UninstallInput specifies what to uninstall.
type UninstallInput struct {
	ItemName string
	Tool     Tool
	Scope    Scope
}

// AddRegistryInput specifies a new registry to add.
type AddRegistryInput struct {
	Name string
	Path string // For local registries
	// URL string // Future: for git registries
}
