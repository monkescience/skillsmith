// Package domain contains shared domain types used across packages.
package domain

// Tool represents a supported AI coding tool.
type Tool string

const (
	ToolOpenCode Tool = "opencode"
	ToolClaude   Tool = "claude"
)

// AllTools returns all supported tools.
func AllTools() []Tool {
	return []Tool{ToolOpenCode, ToolClaude}
}

// ItemType represents the type of registry item.
type ItemType string

const (
	ItemTypeAgent ItemType = "agent"
	ItemTypeSkill ItemType = "skill"
)

// Scope represents where to install items.
type Scope string

const (
	ScopeLocal  Scope = "local"
	ScopeGlobal Scope = "global"
)

// AllScopes returns all supported scopes.
func AllScopes() []Scope {
	return []Scope{ScopeLocal, ScopeGlobal}
}
