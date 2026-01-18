package tui

import (
	"github.com/monke/skillsmith/internal/service"
)

// getLocalPath returns the local installation path for the selected tool.
func (m *Model) getLocalPath() string {
	// Use first item to get the path, or return default
	items, err := m.svc.ListItems(service.ListItemsInput{
		Tool:  m.selectedTool,
		Scope: service.ScopeLocal,
	})
	if err != nil || len(items) == 0 {
		return "./"
	}

	// Use the install path directory from first item
	path := items[0].InstallPath
	if path == "" {
		return "./"
	}

	// Return the directory portion
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i+1]
		}
	}

	return "./"
}

// getGlobalPath returns the global installation path for the selected tool.
func (m *Model) getGlobalPath() string {
	// Use first item to get the path, or return default
	items, err := m.svc.ListItems(service.ListItemsInput{
		Tool:  m.selectedTool,
		Scope: service.ScopeGlobal,
	})
	if err != nil || len(items) == 0 {
		return "~/"
	}

	// Use the install path directory from first item
	path := items[0].InstallPath
	if path == "" {
		return "~/"
	}

	// Return the directory portion
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i+1]
		}
	}

	return "~/"
}

// countInstalled returns the number of installed items for the selected tool and scope.
func (m *Model) countInstalled(scope service.Scope) int {
	installed, _ := m.countInstalledForTool(m.selectedTool, scope)

	return installed
}

// countInstalledForTool returns the installed count and update count for a tool/scope.
func (m *Model) countInstalledForTool(tool service.Tool, scope service.Scope) (int, int) {
	var installed, updates int

	items, err := m.svc.ListItems(service.ListItemsInput{
		Tool:  tool,
		Scope: scope,
	})
	if err != nil {
		return 0, 0
	}

	for _, item := range items {
		if item.State.IsInstalled() {
			installed++
		}

		if item.State.HasUpdate() {
			updates++
		}
	}

	return installed, updates
}

// getScopeLabel returns the display label for the current scope.
func (m *Model) getScopeLabel() string {
	if m.selectedScope == service.ScopeGlobal {
		return LabelGlobal
	}

	return LabelLocal
}

// countSelected returns total selected, installed count, and new count.
func (m *Model) countSelected() (int, int, int) {
	var total, installed, newItems int

	for _, bi := range m.browser.Items {
		if !bi.Selected {
			continue
		}

		total++

		if bi.Status.IsInstalled() {
			installed++
		} else {
			newItems++
		}
	}

	return total, installed, newItems
}

// getItemsForAction returns the selected items that will be affected by the given action.
func (m *Model) getItemsForAction(action string) []BrowserItem {
	var items []BrowserItem

	for _, bi := range m.browser.Items {
		if !bi.Selected {
			continue
		}

		switch action {
		case ActionInstall:
			if !bi.Status.IsInstalled() {
				items = append(items, bi)
			}
		case ActionUpdate, ActionUninstall:
			if bi.Status.IsInstalled() {
				items = append(items, bi)
			}
		}
	}

	return items
}
