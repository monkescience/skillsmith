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

// loadBrowserItems populates the browser with items for the selected tool/scope.
func (m *Model) loadBrowserItems() {
	m.browser.Items = nil
	m.browser.Offset = 0

	items := m.mgr.ListItemsWithState(m.selectedTool, m.selectedScope, "")

	for _, item := range items {
		m.browser.Items = append(m.browser.Items, BrowserItem{
			Item:     item.Item,
			Selected: false,
			Status:   item.State,
		})
	}
}

// updateBrowser handles input for the browser screen.
func (m *Model) updateBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.browser.Cursor > 0 {
			m.browser.Cursor--
			m.ensureCursorVisible()
		}
	case key.Matches(msg, keys.Down):
		if m.browser.Cursor < len(m.browser.Items)-1 {
			m.browser.Cursor++
			m.ensureCursorVisible()
		}
	case key.Matches(msg, keys.Space):
		if m.browser.Cursor < len(m.browser.Items) {
			m.browser.Items[m.browser.Cursor].Selected = !m.browser.Items[m.browser.Cursor].Selected
		}
	case key.Matches(msg, keys.SelectAll):
		for i := range m.browser.Items {
			m.browser.Items[i].Selected = true
		}
	case key.Matches(msg, keys.DeselectAll):
		for i := range m.browser.Items {
			m.browser.Items[i].Selected = false
		}
	case key.Matches(msg, keys.UpdateAll):
		m.updateAllInstalled()
	case key.Matches(msg, keys.Enter):
		m.openActionMenu()
	case key.Matches(msg, keys.Back):
		m.screen = ScreenScopeSelect
	}

	return m, nil
}

// visibleItemCount returns how many items can fit in the visible area.
func (m *Model) visibleItemCount() int {
	// Reserve lines for: header(2) + section headers(2) + status(1) + path(1) + help(2) + box borders(2) + padding(2)
	// Total overhead ~12 lines
	overhead := 12
	available := m.height - overhead

	return max(available, minVisibleItems)
}

// ensureCursorVisible adjusts browser offset to keep cursor in view.
func (m *Model) ensureCursorVisible() {
	visible := m.visibleItemCount()

	// Cursor above visible area
	if m.browser.Cursor < m.browser.Offset {
		m.browser.Offset = m.browser.Cursor
	}

	// Cursor below visible area
	if m.browser.Cursor >= m.browser.Offset+visible {
		m.browser.Offset = m.browser.Cursor - visible + 1
	}

	// Clamp offset
	maxOffset := max(len(m.browser.Items)-visible, 0)

	m.browser.Offset = min(m.browser.Offset, maxOffset)
	m.browser.Offset = max(m.browser.Offset, 0)
}

// viewBrowser renders the browser screen.
func (m *Model) viewBrowser() string {
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(accentStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))
	header.WriteString(accentStyle.Render(" > "))
	header.WriteString(normalStyle.Render(m.getScopeLabel()))

	showPreview := m.width >= minWidthForPreview

	var content string
	if showPreview {
		content = m.renderSplitView()
	} else {
		content = m.renderListOnly()
	}

	var footer strings.Builder

	selected, installedCount, newCount := m.countSelected()

	if selected > 0 {
		status := fmt.Sprintf("%d selected (%d installed, %d new)", selected, installedCount, newCount)
		footer.WriteString(normalStyle.Render(status))
	} else {
		footer.WriteString(dimStyle.Render("No items selected"))
	}

	footer.WriteString("\n\n")

	helpText := "[space] toggle  [a/d] all/none  [u] update  [enter] actions  [esc] back  [q] quit"
	footer.WriteString(helpStyle.Render(helpText))

	return m.renderLayout(header.String(), content, footer.String())
}

// renderListOnly renders the browser list without a preview sidebar.
func (m *Model) renderListOnly() string {
	var content strings.Builder

	visible := m.visibleItemCount()
	totalItems := len(m.browser.Items)

	m.renderVisibleItems(&content, visible)

	if totalItems > visible {
		end := min(m.browser.Offset+visible, totalItems)
		scrollInfo := fmt.Sprintf("[%d-%d of %d]", m.browser.Offset+1, end, totalItems)
		content.WriteString(dimStyle.Render(scrollInfo))
		content.WriteString("\n")
	}

	return lipgloss.NewStyle().
		MarginLeft(mainLeftPadding).
		Render(content.String())
}

// renderSplitView renders the browser with list on left and preview on right.
func (m *Model) renderSplitView() string {
	// Calculate widths: list 40%, sidebar 60%
	availableWidth := m.width - mainLeftPaddingTotal
	listWidth := (availableWidth * listWidthPercent) / percentDivisor
	sidebarWidth := availableWidth - listWidth - 1 // -1 for gap between panels

	sidebarInnerWidth := sidebarWidth - sidebarBorderWidth - sidebarPaddingTotal

	var listContent strings.Builder

	visible := m.visibleItemCount()
	totalItems := len(m.browser.Items)

	m.renderVisibleItemsCompact(&listContent, visible, listWidth)

	if totalItems > visible {
		end := min(m.browser.Offset+visible, totalItems)
		scrollInfo := fmt.Sprintf("[%d-%d of %d]", m.browser.Offset+1, end, totalItems)
		listContent.WriteString(dimStyle.Render(scrollInfo))
	}

	// Max height accounts for header (~2), footer (~4), sidebar chrome (~4), margins (~2)
	previewOverhead := 12
	maxPreviewHeight := max(m.height-previewOverhead, minVisibleItems)
	previewContent := m.renderPreview(sidebarInnerWidth, maxPreviewHeight)

	listPanel := lipgloss.NewStyle().
		Width(listWidth).
		MarginLeft(mainLeftPadding).
		Render(listContent.String())

	sidebarPanel := sidebarStyle.
		Width(sidebarWidth - sidebarBorderWidth).
		Height(maxPreviewHeight).
		Render(previewContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, listPanel, " ", sidebarPanel)
}

// renderVisibleItemsCompact renders items without descriptions (for split view).
func (m *Model) renderVisibleItemsCompact(sb *strings.Builder, visible int, maxWidth int) {
	start := m.browser.Offset
	end := min(start+visible, len(m.browser.Items))

	lastType := registry.ItemType("")

	for i := start; i < end; i++ {
		bi := m.browser.Items[i]

		if bi.Item.Type != lastType {
			if lastType != "" {
				sb.WriteString("\n")
			}

			switch bi.Item.Type {
			case registry.ItemTypeAgent:
				sb.WriteString(headerStyle.Render("Agents:"))
			case registry.ItemTypeSkill:
				sb.WriteString(headerStyle.Render("Skills:"))
			}

			sb.WriteString("\n")

			lastType = bi.Item.Type
		}

		m.renderItem(sb, i, false, maxWidth)
	}
}

// renderVisibleItems renders items with descriptions (for list-only view).
func (m *Model) renderVisibleItems(sb *strings.Builder, visible int) {
	start := m.browser.Offset
	end := min(start+visible, len(m.browser.Items))
	lastType := registry.ItemType("")

	for i := start; i < end; i++ {
		bi := m.browser.Items[i]

		if bi.Item.Type != lastType {
			if lastType != "" {
				sb.WriteString("\n")
			}

			switch bi.Item.Type {
			case registry.ItemTypeAgent:
				sb.WriteString(headerStyle.Render("Agents:"))
			case registry.ItemTypeSkill:
				sb.WriteString(headerStyle.Render("Skills:"))
			}

			sb.WriteString("\n")

			lastType = bi.Item.Type
		}

		m.renderItem(sb, i, true, m.width)
	}
}

// renderPreview renders the preview sidebar for the currently selected item.
func (m *Model) renderPreview(width, maxHeight int) string {
	var sb strings.Builder

	sb.WriteString(sidebarTitleStyle.Render("Preview"))
	sb.WriteString("\n\n")

	if m.browser.Cursor >= len(m.browser.Items) {
		sb.WriteString(dimStyle.Render("No item selected"))

		return sb.String()
	}

	bi := m.browser.Items[m.browser.Cursor]
	item := bi.Item

	sb.WriteString(previewHeaderStyle.Render(item.Name))
	sb.WriteString("\n")

	m.renderPreviewMetadata(&sb, bi)
	m.renderPreviewContent(&sb, item, width)

	content := sb.String()
	lines := strings.Split(content, "\n")

	if len(lines) > maxHeight {
		lines = lines[:maxHeight-1]
		lines = append(lines, dimStyle.Render("..."))

		return strings.Join(lines, "\n")
	}

	return content
}

// renderPreviewMetadata renders the metadata section of the preview.
func (m *Model) renderPreviewMetadata(sb *strings.Builder, bi BrowserItem) {
	bullet := bulletStyle.Render(SymbolBullet) + " "

	// Source
	sb.WriteString(bullet)
	sb.WriteString(dimStyle.Render("source: "))

	source := bi.Item.Source
	if source == "" {
		source = BuiltinSourceName
	}

	if source == BuiltinSourceName {
		sb.WriteString(dimStyle.Render(source))
	} else {
		sb.WriteString(accentStyle.Render(source))
	}

	sb.WriteString("\n")

	// Type
	sb.WriteString(bullet)
	sb.WriteString(dimStyle.Render("type: "))

	switch bi.Item.Type {
	case registry.ItemTypeAgent:
		sb.WriteString(accentStyle.Render("agent"))
	case registry.ItemTypeSkill:
		sb.WriteString(modifiedStyle.Render("skill"))
	}

	sb.WriteString("\n")
	sb.WriteString(bullet)
	sb.WriteString(dimStyle.Render("status: "))
	sb.WriteString(getStatusLabel(bi.Status))
	sb.WriteString("\n")

	path, _ := config.GetInstallPath(bi.Item.Name, bi.Item.Type, m.selectedTool, m.selectedScope)

	sb.WriteString(bullet)
	sb.WriteString(dimStyle.Render("path: "))
	sb.WriteString(pathStyle.Render(path))
	sb.WriteString("\n\n")
}

// renderPreviewContent renders the description and content sections of the preview.
func (m *Model) renderPreviewContent(sb *strings.Builder, item registry.Item, width int) {
	divider := previewDividerStyle.Render(strings.Repeat("-", min(width, previewDividerLen)))

	if item.Description != "" {
		sb.WriteString(sectionHeaderStyle.Render("Description"))
		sb.WriteString("\n")
		sb.WriteString(divider)
		sb.WriteString("\n")

		desc := wrapText(item.Description, width)
		sb.WriteString(normalStyle.Render(desc))
		sb.WriteString("\n\n")
	}

	if item.Body != "" {
		sb.WriteString(sectionHeaderStyle.Render("Content"))
		sb.WriteString("\n")
		sb.WriteString(divider)
		sb.WriteString("\n")

		body := wrapText(item.Body, width)
		lines := strings.Split(body, "\n")

		if len(lines) > previewMaxLines {
			lines = lines[:previewMaxLines]
			lines = append(lines, dimStyle.Render("..."))
		}

		sb.WriteString(previewBodyStyle.Render(strings.Join(lines, "\n")))
	}
}
