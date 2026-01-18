package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/service"
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
		m.renderToolOption(&content, i, tool)
	}

	footer := helpStyle.Render("[enter] select  [q] quit")
	paddedContent := lipgloss.NewStyle().
		MarginLeft(mainLeftPadding).
		Render(content.String())

	return m.renderLayout(header, paddedContent, footer)
}

// renderToolOption renders a single tool option in the tool selection screen.
func (m *Model) renderToolOption(content *strings.Builder, idx int, tool service.Tool) {
	cursor := "  "
	if idx == m.toolSelect.Cursor {
		cursor = SymbolCursor + " "
	}

	agents, skills := m.countItemTypesForTool(tool)

	localInstalled, localUpdates := m.countInstalledForTool(tool, service.ScopeLocal)
	globalInstalled, globalUpdates := m.countInstalledForTool(tool, service.ScopeGlobal)
	totalUpdates := localUpdates + globalUpdates

	line := fmt.Sprintf("%s%s", cursor, tool)
	stats := dimStyle.Render(fmt.Sprintf("  %d agents, %d skills", agents, skills))
	installedInfo := m.formatInstalledInfo(localInstalled, globalInstalled, totalUpdates)

	if idx == m.toolSelect.Cursor {
		content.WriteString(selectedStyle.Render(line))
	} else {
		content.WriteString(normalStyle.Render(line))
	}

	content.WriteString(stats)
	content.WriteString(installedInfo)
	content.WriteString("\n")
}

// countItemTypesForTool returns the count of agents and skills for a tool.
func (m *Model) countItemTypesForTool(tool service.Tool) (int, int) {
	items, _ := m.svc.ListItems(service.ListItemsInput{
		Tool:  tool,
		Scope: service.ScopeLocal,
	})

	var agents, skills int

	for _, item := range items {
		if item.Item.Type == service.ItemTypeAgent {
			agents++
		} else {
			skills++
		}
	}

	return agents, skills
}

// formatInstalledInfo formats the installed/update count info string.
func (m *Model) formatInstalledInfo(localInstalled, globalInstalled, totalUpdates int) string {
	if totalUpdates > 0 {
		return fmt.Sprintf("  (%d local, %d global, ", localInstalled, globalInstalled) +
			updateStyle.Render(fmt.Sprintf("%d updates", totalUpdates)) +
			dimStyle.Render(")")
	}

	return dimStyle.Render(fmt.Sprintf("  (%d local, %d global)", localInstalled, globalInstalled))
}
