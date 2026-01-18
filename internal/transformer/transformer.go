package transformer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/monke/skillsmith/internal/registry"
)

var errUnsupportedTool = errors.New("unsupported tool")

// Transform converts a generic registry item to tool-specific content.
func Transform(item registry.Item, tool registry.Tool) (string, error) {
	switch tool {
	case registry.ToolOpenCode:
		return transformOpenCode(item), nil
	case registry.ToolClaude:
		return transformClaude(item), nil
	default:
		return "", fmt.Errorf("%w: %s", errUnsupportedTool, tool)
	}
}

// transformOpenCode converts an item to OpenCode format.
func transformOpenCode(item registry.Item) string {
	var sb strings.Builder

	// Write frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("description: %s\n", item.Description))

	// Add mode based on type
	if item.Type == registry.ItemTypeAgent {
		sb.WriteString("mode: subagent\n")
	}

	// Add tools configuration
	if item.Tools.Write != nil || item.Tools.Edit != nil || item.Tools.Bash != nil {
		sb.WriteString("tools:\n")

		if item.Tools.Write != nil {
			sb.WriteString(fmt.Sprintf("  write: %t\n", *item.Tools.Write))
		}

		if item.Tools.Edit != nil {
			sb.WriteString(fmt.Sprintf("  edit: %t\n", *item.Tools.Edit))
		}

		if item.Tools.Bash != nil {
			sb.WriteString(fmt.Sprintf("  bash: %t\n", *item.Tools.Bash))
		}
	}

	sb.WriteString("---\n\n")

	// Write body
	sb.WriteString(item.Body)

	return sb.String()
}

// transformClaude converts an item to Claude Code format.
// For skills, this creates a SKILL.md with the appropriate frontmatter.
func transformClaude(item registry.Item) string {
	var sb strings.Builder

	// Write frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", item.Name))
	sb.WriteString(fmt.Sprintf("description: %s\n", item.Description))

	if item.License != "" {
		sb.WriteString(fmt.Sprintf("license: %s\n", item.License))
	}

	sb.WriteString("---\n\n")

	// Write body
	sb.WriteString(item.Body)

	return sb.String()
}
