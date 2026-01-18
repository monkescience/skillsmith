package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/service"
)

// openActionMenu transitions to the action menu screen.
func (m *Model) openActionMenu() {
	selected, _, _ := m.countSelected()
	if selected == 0 && m.browser.Cursor < len(m.browser.Items) {
		m.browser.Items[m.browser.Cursor].Selected = true
	}

	m.buildMenuOptions()
	m.actionMenu.Cursor = 0
	m.screen = ScreenActionMenu
}

// buildMenuOptions constructs the menu options based on selection state.
func (m *Model) buildMenuOptions() {
	m.actionMenu.Options = nil

	_, installedCount, newCount := m.countSelected()

	if newCount > 0 {
		m.actionMenu.Options = append(m.actionMenu.Options, MenuOption{
			Label:   fmt.Sprintf("Install (%d new)", newCount),
			Action:  ActionInstall,
			Enabled: true,
		})
	}

	if installedCount > 0 {
		m.actionMenu.Options = append(m.actionMenu.Options, MenuOption{
			Label:   fmt.Sprintf("Update (%d installed)", installedCount),
			Action:  ActionUpdate,
			Enabled: true,
		})

		m.actionMenu.Options = append(m.actionMenu.Options, MenuOption{
			Label:   fmt.Sprintf("Uninstall (%d)", installedCount),
			Action:  ActionUninstall,
			Enabled: true,
		})
	}
}

// updateActionMenu handles input for the action menu screen.
func (m *Model) updateActionMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		m.actionMenu.Cursor--

		for m.actionMenu.Cursor >= 0 && !m.actionMenu.Options[m.actionMenu.Cursor].Enabled {
			m.actionMenu.Cursor--
		}

		if m.actionMenu.Cursor < 0 {
			m.actionMenu.Cursor = 0

			for !m.actionMenu.Options[m.actionMenu.Cursor].Enabled && m.actionMenu.Cursor < len(m.actionMenu.Options)-1 {
				m.actionMenu.Cursor++
			}
		}
	case key.Matches(msg, keys.Down):
		m.actionMenu.Cursor++

		for m.actionMenu.Cursor < len(m.actionMenu.Options) && !m.actionMenu.Options[m.actionMenu.Cursor].Enabled {
			m.actionMenu.Cursor++
		}

		if m.actionMenu.Cursor >= len(m.actionMenu.Options) {
			m.actionMenu.Cursor = len(m.actionMenu.Options) - 1

			for !m.actionMenu.Options[m.actionMenu.Cursor].Enabled && m.actionMenu.Cursor > 0 {
				m.actionMenu.Cursor--
			}
		}
	case key.Matches(msg, keys.Enter):
		m.executeMenuAction()
	case key.Matches(msg, keys.Back):
		m.screen = ScreenBrowser
	}

	return m, nil
}

// executeMenuAction performs the selected menu action.
func (m *Model) executeMenuAction() {
	if m.actionMenu.Cursor >= len(m.actionMenu.Options) {
		return
	}

	opt := m.actionMenu.Options[m.actionMenu.Cursor]

	switch opt.Action {
	case ActionInstall:
		m.installNew()
	case ActionUpdate:
		m.updateInstalled()
	case ActionUninstall:
		m.uninstallSelected()
	}
}

// installNew installs all selected new items.
func (m *Model) installNew() {
	installed := 0

	for i, bi := range m.browser.Items {
		if !bi.Selected || bi.Status.IsInstalled() {
			continue
		}

		result, err := m.svc.Install(service.InstallInput{
			ItemName: bi.Item.Name,
			Tool:     m.selectedTool,
			Scope:    m.selectedScope,
			Force:    false,
		})
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			installed++
			m.browser.Items[i].Status = service.StateUpToDate
		}

		m.browser.Items[i].Selected = false
	}

	if installed > 0 {
		m.message = fmt.Sprintf("Installed %d items", installed)
		m.messageStyle = successMsgStyle
	}

	m.screen = ScreenBrowser
}

// updateInstalled updates all selected installed items.
func (m *Model) updateInstalled() {
	updated := 0
	skippedModified := 0

	for i, bi := range m.browser.Items {
		if !bi.Selected || !bi.Status.IsInstalled() {
			continue
		}

		if bi.Status.IsModified() {
			skippedModified++
			m.browser.Items[i].Selected = false

			continue
		}

		result, err := m.svc.Install(service.InstallInput{
			ItemName: bi.Item.Name,
			Tool:     m.selectedTool,
			Scope:    m.selectedScope,
			Force:    true,
		})
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			updated++
			m.browser.Items[i].Status = service.StateUpToDate
		}

		m.browser.Items[i].Selected = false
	}

	switch {
	case updated > 0 && skippedModified > 0:
		m.message = fmt.Sprintf("Updated %d items, skipped %d modified", updated, skippedModified)
		m.messageStyle = successMsgStyle
	case updated > 0:
		m.message = fmt.Sprintf("Updated %d items", updated)
		m.messageStyle = successMsgStyle
	case skippedModified > 0:
		m.message = fmt.Sprintf("Skipped %d modified items", skippedModified)
		m.messageStyle = modifiedStyle
	}

	m.screen = ScreenBrowser
}

// updateAllInstalled updates all installed items (triggered by 'u' key in browser).
func (m *Model) updateAllInstalled() {
	updated := 0
	skippedModified := 0

	for i, bi := range m.browser.Items {
		if !bi.Status.IsInstalled() {
			continue
		}

		if bi.Status.IsModified() {
			skippedModified++

			continue
		}

		result, err := m.svc.Install(service.InstallInput{
			ItemName: bi.Item.Name,
			Tool:     m.selectedTool,
			Scope:    m.selectedScope,
			Force:    true,
		})
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle

			return
		}

		if result.Success {
			updated++
			m.browser.Items[i].Status = service.StateUpToDate
		}

		m.browser.Items[i].Selected = false
	}

	switch {
	case updated > 0 && skippedModified > 0:
		m.message = fmt.Sprintf("Updated %d items, skipped %d modified", updated, skippedModified)
		m.messageStyle = successMsgStyle
	case updated > 0:
		m.message = fmt.Sprintf("Updated %d items", updated)
		m.messageStyle = successMsgStyle
	case skippedModified > 0:
		m.message = fmt.Sprintf("Skipped %d modified items", skippedModified)
		m.messageStyle = modifiedStyle
	default:
		m.message = "No installed items to update"
		m.messageStyle = dimStyle
	}
}

// uninstallSelected removes all selected installed items.
func (m *Model) uninstallSelected() {
	uninstalled := 0

	for i, bi := range m.browser.Items {
		if !bi.Selected || !bi.Status.IsInstalled() {
			continue
		}

		result, err := m.svc.Uninstall(service.UninstallInput{
			ItemName: bi.Item.Name,
			Tool:     m.selectedTool,
			Scope:    m.selectedScope,
		})
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			uninstalled++
			m.browser.Items[i].Status = service.StateNotInstalled
		}

		m.browser.Items[i].Selected = false
	}

	if uninstalled > 0 {
		m.message = fmt.Sprintf("Uninstalled %d items", uninstalled)
		m.messageStyle = successMsgStyle
	}

	m.screen = ScreenBrowser
}

// viewActionMenu renders the action menu screen.
func (m *Model) viewActionMenu() string {
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(accentStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))
	header.WriteString(accentStyle.Render(" > "))
	header.WriteString(normalStyle.Render(m.getScopeLabel()))

	// Calculate widths similar to browser split view
	availableWidth := m.width - mainLeftPaddingTotal
	listWidth := (availableWidth * listWidthPercent) / percentDivisor
	sidebarWidth := availableWidth - listWidth - 1

	sidebarInnerWidth := sidebarWidth - sidebarBorderWidth - sidebarPaddingTotal

	// Build menu content (left panel)
	var menuContent strings.Builder

	selected, _, _ := m.countSelected()
	menuContent.WriteString(headerStyle.Render(fmt.Sprintf("%d items selected", selected)))
	menuContent.WriteString("\n\n")

	for i, opt := range m.actionMenu.Options {
		cursor := "  "
		if i == m.actionMenu.Cursor {
			cursor = SymbolCursor + " "
		}

		switch {
		case !opt.Enabled:
			menuContent.WriteString(dimStyle.Render("  " + opt.Label))
		case i == m.actionMenu.Cursor:
			menuContent.WriteString(selectedStyle.Render(cursor + opt.Label))
		default:
			menuContent.WriteString(normalStyle.Render(cursor + opt.Label))
		}

		menuContent.WriteString("\n")
	}

	// Max height for sidebar
	previewOverhead := 12
	maxPreviewHeight := max(m.height-previewOverhead, minVisibleItems)

	// Build preview content (right panel)
	previewContent := m.renderActionPreview(sidebarInnerWidth, maxPreviewHeight)

	// Assemble split view
	listPanel := lipgloss.NewStyle().
		Width(listWidth).
		MarginLeft(mainLeftPadding).
		Render(menuContent.String())

	sidebarPanel := sidebarStyle.
		Width(sidebarWidth - sidebarBorderWidth).
		Height(maxPreviewHeight).
		Render(previewContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, " ", sidebarPanel)
	footer := helpStyle.Render("[enter] confirm  [esc] cancel")

	return m.renderLayout(header.String(), content, footer)
}

// renderActionPreview renders the right sidebar showing items affected by the current action.
func (m *Model) renderActionPreview(_, maxHeight int) string {
	var sb strings.Builder

	// Determine title based on current menu option
	title := "Affected Items"
	action := ""

	if m.actionMenu.Cursor < len(m.actionMenu.Options) {
		action = m.actionMenu.Options[m.actionMenu.Cursor].Action

		switch action {
		case ActionInstall:
			title = "Will Install"
		case ActionUpdate:
			title = "Will Update"
		case ActionUninstall:
			title = "Will Uninstall"
		}
	}

	sb.WriteString(sidebarTitleStyle.Render(title))
	sb.WriteString("\n\n")

	items := m.getItemsForAction(action)

	if len(items) == 0 {
		sb.WriteString(dimStyle.Render("No items"))

		return sb.String()
	}

	bullet := bulletStyle.Render(SymbolBullet) + " "

	for _, bi := range items {
		sb.WriteString(bullet)

		// Name
		sb.WriteString(normalStyle.Render(bi.Item.Name))

		// Type and status in parentheses
		typeStr := string(bi.Item.Type)
		statusStr := getStatusShortLabel(bi.Status)

		sb.WriteString(dimStyle.Render(fmt.Sprintf(" (%s, %s)", typeStr, statusStr)))
		sb.WriteString("\n")
	}

	content := sb.String()
	lines := strings.Split(content, "\n")

	if len(lines) > maxHeight {
		lines = lines[:maxHeight-1]
		lines = append(lines, dimStyle.Render("..."))

		return strings.Join(lines, "\n")
	}

	return content
}
