package service

import (
	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/installer"
	"github.com/monke/skillsmith/internal/registry"
)

// toServiceItem converts a registry.Item to a service.Item.
func toServiceItem(item registry.Item) Item {
	compatibility := make([]Tool, len(item.Compatibility))
	for i, t := range item.Compatibility {
		compatibility[i] = Tool(t)
	}

	tags := make([]string, len(item.Tags))
	copy(tags, item.Tags)

	return Item{
		Name:          item.Name,
		Description:   item.Description,
		Type:          ItemType(item.Type),
		Category:      item.Category,
		Compatibility: compatibility,
		Tags:          tags,
		Author:        item.Author,
		Source:        item.Source,
		Body:          item.Body,
	}
}

// toServiceState converts an installer.ItemState to a service.ItemState.
func toServiceState(state installer.ItemState) ItemState {
	switch state {
	case installer.StateNotInstalled:
		return StateNotInstalled
	case installer.StateUpToDate:
		return StateUpToDate
	case installer.StateUpdateAvailable:
		return StateUpdateAvailable
	case installer.StateModified:
		return StateModified
	case installer.StateModifiedWithUpdate:
		return StateModifiedWithUpdate
	default:
		return StateNotInstalled
	}
}

// toRegistryTool converts a service.Tool to a registry.Tool.
func toRegistryTool(tool Tool) registry.Tool {
	return registry.Tool(tool)
}

// toConfigScope converts a service.Scope to a config.Scope.
func toConfigScope(scope Scope) config.Scope {
	return config.Scope(scope)
}

// toRegistryItemType converts a service.ItemType to a registry.ItemType.
func toRegistryItemType(itemType ItemType) registry.ItemType {
	return registry.ItemType(itemType)
}
