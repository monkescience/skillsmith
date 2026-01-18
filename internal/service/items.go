package service

import (
	"errors"
	"fmt"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/installer"
)

// ErrItemNotFound is returned when an item is not found in the registry.
var ErrItemNotFound = errors.New("item not found")

// ListItems returns items matching the criteria with their installation state.
func (s *Service) ListItems(input ListItemsInput) ([]ItemWithState, error) {
	regTool := toRegistryTool(input.Tool)
	cfgScope := toConfigScope(input.Scope)

	// Get items filtered by tool
	items := s.registry.ByTool(regTool)

	// Optionally filter by type
	if input.Type != "" {
		regType := toRegistryItemType(input.Type)
		filtered := make([]ItemWithState, 0, len(items))

		for _, item := range items {
			if item.Type != regType {
				continue
			}

			state, path, _ := installer.GetItemState(item, regTool, cfgScope)

			filtered = append(filtered, ItemWithState{
				Item:        toServiceItem(item),
				State:       toServiceState(state),
				InstallPath: path,
			})
		}

		return filtered, nil
	}

	// Return all items for the tool
	result := make([]ItemWithState, 0, len(items))

	for _, item := range items {
		state, path, _ := installer.GetItemState(item, regTool, cfgScope)

		result = append(result, ItemWithState{
			Item:        toServiceItem(item),
			State:       toServiceState(state),
			InstallPath: path,
		})
	}

	return result, nil
}

// GetItem returns a single item by name.
func (s *Service) GetItem(name string) (*Item, error) {
	for _, item := range s.registry.Items {
		if item.Name == name {
			svcItem := toServiceItem(item)

			return &svcItem, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrItemNotFound, name)
}

// Install installs an item.
func (s *Service) Install(input InstallInput) (*InstallResult, error) {
	// Find the item
	var regItem *itemRef

	for i := range s.registry.Items {
		if s.registry.Items[i].Name == input.ItemName {
			regItem = &itemRef{index: i}

			break
		}
	}

	if regItem == nil {
		return nil, fmt.Errorf("%w: %s", ErrItemNotFound, input.ItemName)
	}

	item := s.registry.Items[regItem.index]
	regTool := toRegistryTool(input.Tool)
	cfgScope := toConfigScope(input.Scope)

	// Check compatibility
	if !item.IsCompatibleWith(regTool) {
		return &InstallResult{
			ItemName: input.ItemName,
			Success:  false,
			Skipped:  true,
			Error:    fmt.Sprintf("item not compatible with %s", input.Tool),
		}, nil
	}

	// Get path for result
	path, err := config.GetInstallPath(item, regTool, cfgScope)
	if err != nil {
		return nil, fmt.Errorf("get install path: %w", err)
	}

	// Install
	result, err := installer.Install(item, regTool, cfgScope, input.Force)
	if err != nil {
		return nil, fmt.Errorf("install: %w", err)
	}

	if !result.Success && !input.Force {
		return &InstallResult{
			ItemName: input.ItemName,
			Success:  false,
			Skipped:  true,
			Path:     path,
		}, nil
	}

	return &InstallResult{
		ItemName: input.ItemName,
		Success:  result.Success,
		Path:     path,
	}, nil
}

// itemRef holds a reference to an item in the registry.
type itemRef struct {
	index int
}

// InstallMany installs multiple items, returning results for each.
// Continues on error and returns all results.
func (s *Service) InstallMany(inputs []InstallInput) ([]InstallResult, error) {
	results := make([]InstallResult, 0, len(inputs))

	for _, input := range inputs {
		result, err := s.Install(input)
		if err != nil {
			results = append(results, InstallResult{
				ItemName: input.ItemName,
				Success:  false,
				Error:    err.Error(),
			})

			continue
		}

		results = append(results, *result)
	}

	return results, nil
}

// Uninstall removes an installed item.
func (s *Service) Uninstall(input UninstallInput) (*UninstallResult, error) {
	// Find the item
	var regItem *itemRef

	for i := range s.registry.Items {
		if s.registry.Items[i].Name == input.ItemName {
			regItem = &itemRef{index: i}

			break
		}
	}

	if regItem == nil {
		return nil, fmt.Errorf("%w: %s", ErrItemNotFound, input.ItemName)
	}

	item := s.registry.Items[regItem.index]
	regTool := toRegistryTool(input.Tool)
	cfgScope := toConfigScope(input.Scope)

	// Get path for result
	path, err := config.GetInstallPath(item, regTool, cfgScope)
	if err != nil {
		return nil, fmt.Errorf("get install path: %w", err)
	}

	// Uninstall
	result, err := installer.Uninstall(item, regTool, cfgScope)
	if err != nil {
		return nil, fmt.Errorf("uninstall: %w", err)
	}

	return &UninstallResult{
		ItemName: input.ItemName,
		Success:  result.Success,
		Path:     path,
	}, nil
}

// GetInstallPath returns the path where an item would be installed.
func (s *Service) GetInstallPath(itemName string, tool Tool, scope Scope) (string, error) {
	// Find the item
	for _, item := range s.registry.Items {
		if item.Name == itemName {
			regTool := toRegistryTool(tool)
			cfgScope := toConfigScope(scope)

			path, err := config.GetInstallPath(item, regTool, cfgScope)
			if err != nil {
				return "", fmt.Errorf("get install path: %w", err)
			}

			return path, nil
		}
	}

	return "", fmt.Errorf("%w: %s", ErrItemNotFound, itemName)
}
