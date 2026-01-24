package registry

import "slices"

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

// ToolConfig contains tool-specific settings for an item.
type ToolConfig struct {
	// Enabled tools (map of tool name to enabled state).
	Write *bool `yaml:"write,omitempty"`
	Edit  *bool `yaml:"edit,omitempty"`
	Bash  *bool `yaml:"bash,omitempty"`
}

// Item represents a single installable item (agent or skill).
type Item struct {
	// Name is the identifier for this item.
	Name string `yaml:"name"`

	// Description is a short description of what this item does.
	Description string `yaml:"description"`

	// Type indicates whether this is an agent or skill.
	Type ItemType `yaml:"-"` // Derived from directory, not from frontmatter

	// Category for grouping in the UI (e.g., "code-quality", "documentation").
	Category string `yaml:"category"`

	// Compatibility lists which tools this item works with.
	Compatibility []Tool `yaml:"compatibility"`

	// Tools configuration (which tools are enabled/disabled).
	Tools ToolConfig `yaml:"tools,omitempty"`

	// Tags for filtering.
	Tags []string `yaml:"tags,omitempty"`

	// Author of this item.
	Author string `yaml:"author,omitempty"`

	// License for this item.
	License string `yaml:"license,omitempty"`

	// Metadata is an arbitrary key-value map for additional properties.
	Metadata map[string]string `yaml:"metadata,omitempty"`

	// Body is the content after frontmatter (the actual prompt/instructions).
	Body string `yaml:"-"`

	// SourcePath is the path to the source file (for debugging).
	SourcePath string `yaml:"-"`

	// Source is the name of the registry source this item came from.
	Source string `yaml:"-"`
}

// IsCompatibleWith checks if the item is compatible with a given tool.
func (i *Item) IsCompatibleWith(tool Tool) bool {
	return slices.Contains(i.Compatibility, tool)
}

// Registry holds all available items.
type Registry struct {
	Items []Item
}

// ByType returns items filtered by the specified type.
func (r *Registry) ByType(itemType ItemType) []Item {
	var result []Item

	for _, item := range r.Items {
		if item.Type == itemType {
			result = append(result, item)
		}
	}

	return result
}

// ByTool returns items compatible with the specified tool.
func (r *Registry) ByTool(tool Tool) []Item {
	var result []Item

	for _, item := range r.Items {
		if item.IsCompatibleWith(tool) {
			result = append(result, item)
		}
	}

	return result
}

// ByToolAndType returns items compatible with a tool and matching a type.
func (r *Registry) ByToolAndType(tool Tool, itemType ItemType) []Item {
	var result []Item

	for _, item := range r.Items {
		if item.Type == itemType && item.IsCompatibleWith(tool) {
			result = append(result, item)
		}
	}

	return result
}

// GetTools returns all tools that have at least one compatible item.
func (r *Registry) GetTools() []Tool {
	seen := make(map[Tool]bool)

	for _, item := range r.Items {
		for _, tool := range item.Compatibility {
			seen[tool] = true
		}
	}

	var tools []Tool

	for _, tool := range AllTools() {
		if seen[tool] {
			tools = append(tools, tool)
		}
	}

	return tools
}

// GetCategories returns all unique categories.
func (r *Registry) GetCategories() []string {
	seen := make(map[string]bool)

	var categories []string

	for _, item := range r.Items {
		if item.Category != "" && !seen[item.Category] {
			seen[item.Category] = true
			categories = append(categories, item.Category)
		}
	}

	return categories
}
