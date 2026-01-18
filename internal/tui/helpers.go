package tui

import (
	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/installer"
	"github.com/monke/skillsmith/internal/registry"
)

// getLocalPath returns the local installation path for the selected tool.
func (m *Model) getLocalPath() string {
	paths, err := config.GetPaths(m.selectedTool)
	if err != nil {
		return "./"
	}

	return paths.LocalDir + "/"
}

// getGlobalPath returns the global installation path for the selected tool.
func (m *Model) getGlobalPath() string {
	paths, err := config.GetPaths(m.selectedTool)
	if err != nil {
		return "~/"
	}

	return paths.GlobalDir + "/"
}

// countInstalled returns the number of installed items for the selected tool and scope.
func (m *Model) countInstalled(scope config.Scope) int {
	installed, _ := m.countInstalledForTool(m.selectedTool, scope)

	return installed
}

// countInstalledForTool returns the installed count and update count for a tool/scope.
func (m *Model) countInstalledForTool(tool registry.Tool, scope config.Scope) (int, int) {
	var installed, updates int

	items := m.registry.ByTool(tool)

	for _, item := range items {
		state, _, _ := installer.GetItemState(item, tool, scope)
		if state.IsInstalled() {
			installed++
		}

		if state.HasUpdate() {
			updates++
		}
	}

	return installed, updates
}

// getScopeLabel returns the display label for the current scope.
func (m *Model) getScopeLabel() string {
	if m.selectedScope == config.ScopeGlobal {
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
