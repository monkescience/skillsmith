package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
)

// updateToolSelect handles input for the tool selection screen.
func (m *Model) updateToolSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.toolSelect.Cursor > 0 {
			m.toolSelect.Cursor--
		}
	case key.Matches(msg, keys.Down):
		if m.toolSelect.Cursor < len(m.toolSelect.Tools)-1 {
			m.toolSelect.Cursor++
		}
	case key.Matches(msg, keys.Enter):
		m.selectedTool = m.toolSelect.Tools[m.toolSelect.Cursor]
		m.screen = ScreenScopeSelect
		m.scopeSelect.Cursor = 0
	}

	return m, nil
}

// viewToolSelect renders the tool selection screen.
func (m *Model) viewToolSelect() string {
	header := titleStyle.Render("skillsmith")

	var content strings.Builder

	content.WriteString("Select target tool:\n\n")

	for i, tool := range m.toolSelect.Tools {
		cursor := "  "
		if i == m.toolSelect.Cursor {
			cursor = SymbolCursor + " "
		}

		// Count compatible items
		items := m.registry.ByTool(tool)
		agents := 0
		skills := 0

		for _, item := range items {
			if item.Type == registry.ItemTypeAgent {
				agents++
			} else {
				skills++
			}
		}

		// Count installed items per scope
		localInstalled, localUpdates := m.countInstalledForTool(tool, config.ScopeLocal)
		globalInstalled, globalUpdates := m.countInstalledForTool(tool, config.ScopeGlobal)
		totalUpdates := localUpdates + globalUpdates

		line := fmt.Sprintf("%s%s", cursor, tool)
		stats := dimStyle.Render(fmt.Sprintf("  %d agents, %d skills", agents, skills))

		var installedInfo string
		if totalUpdates > 0 {
			installedInfo = fmt.Sprintf("  (%d local, %d global, ", localInstalled, globalInstalled)
			installedInfo += updateStyle.Render(fmt.Sprintf("%d updates", totalUpdates))
			installedInfo += dimStyle.Render(")")
		} else {
			installedInfo = dimStyle.Render(fmt.Sprintf("  (%d local, %d global)", localInstalled, globalInstalled))
		}

		if i == m.toolSelect.Cursor {
			content.WriteString(selectedStyle.Render(line))
			content.WriteString(stats)
			content.WriteString(installedInfo)
		} else {
			content.WriteString(normalStyle.Render(line))
			content.WriteString(stats)
			content.WriteString(installedInfo)
		}

		content.WriteString("\n")
	}

	footer := helpStyle.Render("[enter] select  [q] quit")
	paddedContent := lipgloss.NewStyle().
		MarginLeft(mainLeftPadding).
		Render(content.String())

	return m.renderLayout(header, paddedContent, footer)
}
