package registry

// Tool represents a supported AI coding tool.
type Tool string

const (
	ToolOpenCode   Tool = "opencode"
	ToolClaudeCode Tool = "claude"
)

// ItemType represents the type of registry item.
type ItemType string

const (
	ItemTypeAgent ItemType = "agent"
	ItemTypeSkill ItemType = "skill"
)

// Item represents a single installable item (agent or skill).
type Item struct {
	// Name is the identifier for this item.
	Name string `yaml:"name"`

	// Description is a short description of what this item does.
	Description string `yaml:"description"`

	// Type indicates whether this is an agent or skill.
	Type ItemType `yaml:"type"`

	// Tool indicates which AI tool this item is for.
	Tool Tool `yaml:"tool"`

	// Content is the full content to be installed (markdown file content).
	Content string `yaml:"content"`

	// Filename is the target filename when installed.
	Filename string `yaml:"filename"`

	// Category for grouping in the UI (e.g., "code-quality", "documentation").
	Category string `yaml:"category"`

	// Tags for filtering.
	Tags []string `yaml:"tags,omitempty"`

	// Author of this item.
	Author string `yaml:"author,omitempty"`

	// License for this item.
	License string `yaml:"license,omitempty"`
}

// Registry holds all available items organized by tool.
type Registry struct {
	Items []Item `yaml:"items"`
}

// ByTool returns items filtered by the specified tool.
func (r *Registry) ByTool(tool Tool) []Item {
	var result []Item

	for _, item := range r.Items {
		if item.Tool == tool {
			result = append(result, item)
		}
	}

	return result
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

// ByToolAndType returns items filtered by both tool and type.
func (r *Registry) ByToolAndType(tool Tool, itemType ItemType) []Item {
	var result []Item

	for _, item := range r.Items {
		if item.Tool == tool && item.Type == itemType {
			result = append(result, item)
		}
	}

	return result
}

// GetTools returns all unique tools in the registry.
func (r *Registry) GetTools() []Tool {
	seen := make(map[Tool]bool)

	var tools []Tool

	for _, item := range r.Items {
		if !seen[item.Tool] {
			seen[item.Tool] = true
			tools = append(tools, item.Tool)
		}
	}

	return tools
}
