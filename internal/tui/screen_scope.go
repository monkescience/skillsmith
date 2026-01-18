package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/service"
)

// updateScopeSelect handles input for the scope selection screen.
func (m *Model) updateScopeSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.scopeSelect.Cursor > 0 {
			m.scopeSelect.Cursor--
		}
	case key.Matches(msg, keys.Down):
		if m.scopeSelect.Cursor < len(m.scopeSelect.Scopes)-1 {
			m.scopeSelect.Cursor++
		}
	case key.Matches(msg, keys.Enter):
		m.selectedScope = m.scopeSelect.Scopes[m.scopeSelect.Cursor]
		m.loadBrowserItems()
		m.screen = ScreenBrowser
		m.browser.Cursor = 0
	case key.Matches(msg, keys.Back):
		m.screen = ScreenToolSelect
	}

	return m, nil
}

// viewScopeSelect renders the scope selection screen.
func (m *Model) viewScopeSelect() string {
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))

	var content strings.Builder

	content.WriteString("Install location:\n\n")

	for i, scope := range m.scopeSelect.Scopes {
		cursor := "  "
		if i == m.scopeSelect.Cursor {
			cursor = SymbolCursor + " "
		}

		var label, path string

		switch scope {
		case service.ScopeLocal:
			label = LabelLocal
			path = m.getLocalPath()
		case service.ScopeGlobal:
			label = LabelGlobal
			path = m.getGlobalPath()
		}

		// Count installed for this scope
		installed := m.countInstalled(scope)

		line := fmt.Sprintf("%s%-8s", cursor, label)
		pathInfo := dimStyle.Render(fmt.Sprintf("%s  (%d installed)", path, installed))

		if i == m.scopeSelect.Cursor {
			content.WriteString(selectedStyle.Render(line))
			content.WriteString(pathInfo)
		} else {
			content.WriteString(normalStyle.Render(line))
			content.WriteString(pathInfo)
		}

		content.WriteString("\n")
	}

	footer := helpStyle.Render("[enter] select  [esc] back  [q] quit")
	paddedContent := lipgloss.NewStyle().
		MarginLeft(mainLeftPadding).
		Render(content.String())

	return m.renderLayout(header.String(), paddedContent, footer)
}
