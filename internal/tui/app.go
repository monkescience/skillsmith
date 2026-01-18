package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/installer"
	"github.com/monke/skillsmith/internal/registry"
)

type Screen int

const (
	ScreenToolSelect Screen = iota
	ScreenScopeSelect
	ScreenBrowser
	ScreenActionMenu
)

type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Enter       key.Binding
	Space       key.Binding
	SelectAll   key.Binding
	DeselectAll key.Binding
	UpdateAll   key.Binding
	Back        key.Binding
	Quit        key.Binding
}

var keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
	),
	SelectAll: key.NewBinding(
		key.WithKeys("a"),
	),
	DeselectAll: key.NewBinding(
		key.WithKeys("d"),
	),
	UpdateAll: key.NewBinding(
		key.WithKeys("u"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
	),
}

type BrowserItem struct {
	Item     registry.Item
	Selected bool
	Status   installer.ItemState
}

type MenuOption struct {
	Label   string
	Action  string
	Enabled bool
}

type Model struct {
	registry *registry.Registry
	screen   Screen
	width    int
	height   int
	ready    bool

	// Tool selection
	tools      []registry.Tool
	toolCursor int

	// Scope selection
	scopes      []config.Scope
	scopeCursor int

	// Current selections
	selectedTool  registry.Tool
	selectedScope config.Scope

	// Browser state
	browserItems  []BrowserItem
	browserCursor int
	browserOffset int // scroll offset for visible window

	// Action menu state
	menuOptions []MenuOption
	menuCursor  int

	// Messages
	message      string
	messageStyle lipgloss.Style
}

func NewModel(reg *registry.Registry) *Model {
	return &Model{
		registry: reg,
		screen:   ScreenToolSelect,
		tools:    registry.AllTools(),
		scopes:   []config.Scope{config.ScopeLocal, config.ScopeGlobal},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tea.KeyMsg:
		// Clear message on key press
		m.message = ""

		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}

		switch m.screen {
		case ScreenToolSelect:
			return m.updateToolSelect(msg)
		case ScreenScopeSelect:
			return m.updateScopeSelect(msg)
		case ScreenBrowser:
			return m.updateBrowser(msg)
		case ScreenActionMenu:
			return m.updateActionMenu(msg)
		}
	}

	return m, nil
}

func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	switch m.screen {
	case ScreenToolSelect:
		return m.viewToolSelect()
	case ScreenScopeSelect:
		return m.viewScopeSelect()
	case ScreenBrowser:
		return m.viewBrowser()
	case ScreenActionMenu:
		return m.viewActionMenu()
	default:
		return "Unknown screen"
	}
}

// renderLayout renders a borderless layout with pinned header and footer.
func (m *Model) renderLayout(header, content, footer string) string {
	// Count lines in header and footer
	headerLines := strings.Count(header, "\n") + 1
	footerLines := strings.Count(footer, "\n") + 1

	// Calculate content area height (full terminal height minus header/footer/margins)
	margins := 2 // Top and bottom margin
	contentHeight := max(m.height-headerLines-footerLines-margins, 1)

	// Build the layout with left padding
	var sb strings.Builder

	leftPad := strings.Repeat(" ", mainLeftPadding)

	// Header (pinned top)
	sb.WriteString(leftPad)
	sb.WriteString(header)
	sb.WriteString("\n\n")

	// Content (fills middle) - content is pre-padded
	contentLines := strings.Count(content, "\n") + 1
	sb.WriteString(content)

	if !strings.HasSuffix(content, "\n") {
		sb.WriteString("\n")
	}

	// Add padding to push footer to bottom
	padding := contentHeight - contentLines
	for range padding {
		sb.WriteString("\n")
	}

	// Footer (pinned bottom) - add left padding to each line
	for i, line := range strings.Split(footer, "\n") {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(leftPad)
		sb.WriteString(line)
	}

	return sb.String()
}

// =============================================================================
// Tool Selection Screen
// =============================================================================

func (m *Model) updateToolSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.toolCursor > 0 {
			m.toolCursor--
		}
	case key.Matches(msg, keys.Down):
		if m.toolCursor < len(m.tools)-1 {
			m.toolCursor++
		}
	case key.Matches(msg, keys.Enter):
		m.selectedTool = m.tools[m.toolCursor]
		m.screen = ScreenScopeSelect
		m.scopeCursor = 0
	}

	return m, nil
}

func (m *Model) viewToolSelect() string {
	header := titleStyle.Render("skillsmith")

	var content strings.Builder

	content.WriteString("Select target tool:\n\n")

	for i, tool := range m.tools {
		cursor := "  "
		if i == m.toolCursor {
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

		if i == m.toolCursor {
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

// =============================================================================
// Scope Selection Screen
// =============================================================================

func (m *Model) updateScopeSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.scopeCursor > 0 {
			m.scopeCursor--
		}
	case key.Matches(msg, keys.Down):
		if m.scopeCursor < len(m.scopes)-1 {
			m.scopeCursor++
		}
	case key.Matches(msg, keys.Enter):
		m.selectedScope = m.scopes[m.scopeCursor]
		m.loadBrowserItems()
		m.screen = ScreenBrowser
		m.browserCursor = 0
	case key.Matches(msg, keys.Back):
		m.screen = ScreenToolSelect
	}

	return m, nil
}

func (m *Model) viewScopeSelect() string {
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))

	var content strings.Builder

	content.WriteString("Install location:\n\n")

	for i, scope := range m.scopes {
		cursor := "  "
		if i == m.scopeCursor {
			cursor = SymbolCursor + " "
		}

		var label, path string

		switch scope {
		case config.ScopeLocal:
			label = LabelLocal
			path = m.getLocalPath()
		case config.ScopeGlobal:
			label = LabelGlobal
			path = m.getGlobalPath()
		}

		// Count installed for this scope
		installed := m.countInstalled(scope)

		line := fmt.Sprintf("%s%-8s", cursor, label)
		pathInfo := dimStyle.Render(fmt.Sprintf("%s  (%d installed)", path, installed))

		if i == m.scopeCursor {
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

func (m *Model) getLocalPath() string {
	paths, err := config.GetPaths(m.selectedTool)
	if err != nil {
		return "./"
	}

	return paths.LocalDir + "/"
}

func (m *Model) getGlobalPath() string {
	paths, err := config.GetPaths(m.selectedTool)
	if err != nil {
		return "~/"
	}

	return paths.GlobalDir + "/"
}

func (m *Model) countInstalled(scope config.Scope) int {
	installed, _ := m.countInstalledForTool(m.selectedTool, scope)

	return installed
}

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

func (m *Model) getScopeLabel() string {
	if m.selectedScope == config.ScopeGlobal {
		return LabelGlobal
	}

	return LabelLocal
}

// =============================================================================
// Browser Screen
// =============================================================================

func (m *Model) loadBrowserItems() {
	m.browserItems = nil
	m.browserOffset = 0

	items := m.registry.ByTool(m.selectedTool)

	for _, item := range items {
		state, _, _ := installer.GetItemState(item, m.selectedTool, m.selectedScope)
		m.browserItems = append(m.browserItems, BrowserItem{
			Item:     item,
			Selected: false,
			Status:   state,
		})
	}
}

func (m *Model) updateBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.browserCursor > 0 {
			m.browserCursor--
			m.ensureCursorVisible()
		}
	case key.Matches(msg, keys.Down):
		if m.browserCursor < len(m.browserItems)-1 {
			m.browserCursor++
			m.ensureCursorVisible()
		}
	case key.Matches(msg, keys.Space):
		if m.browserCursor < len(m.browserItems) {
			m.browserItems[m.browserCursor].Selected = !m.browserItems[m.browserCursor].Selected
		}
	case key.Matches(msg, keys.SelectAll):
		for i := range m.browserItems {
			m.browserItems[i].Selected = true
		}
	case key.Matches(msg, keys.DeselectAll):
		for i := range m.browserItems {
			m.browserItems[i].Selected = false
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

// ensureCursorVisible adjusts browserOffset to keep cursor in view.
func (m *Model) ensureCursorVisible() {
	visible := m.visibleItemCount()

	// Cursor above visible area
	if m.browserCursor < m.browserOffset {
		m.browserOffset = m.browserCursor
	}

	// Cursor below visible area
	if m.browserCursor >= m.browserOffset+visible {
		m.browserOffset = m.browserCursor - visible + 1
	}

	// Clamp offset
	maxOffset := max(len(m.browserItems)-visible, 0)

	m.browserOffset = min(m.browserOffset, maxOffset)
	m.browserOffset = max(m.browserOffset, 0)
}

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

func (m *Model) renderListOnly() string {
	var content strings.Builder

	visible := m.visibleItemCount()
	totalItems := len(m.browserItems)

	m.renderVisibleItems(&content, visible)

	if totalItems > visible {
		end := min(m.browserOffset+visible, totalItems)
		scrollInfo := fmt.Sprintf("[%d-%d of %d]", m.browserOffset+1, end, totalItems)
		content.WriteString(dimStyle.Render(scrollInfo))
		content.WriteString("\n")
	}

	// Add left padding using lipgloss
	return lipgloss.NewStyle().
		MarginLeft(mainLeftPadding).
		Render(content.String())
}

func (m *Model) renderSplitView() string {
	// Calculate widths: list 40%, sidebar 60%
	availableWidth := m.width - mainLeftPaddingTotal // Account for margins
	listWidth := (availableWidth * listWidthPercent) / percentDivisor
	sidebarWidth := availableWidth - listWidth - 1 // -1 for gap between panels

	sidebarInnerWidth := sidebarWidth - sidebarBorderWidth - sidebarPaddingTotal

	var listContent strings.Builder

	visible := m.visibleItemCount()
	totalItems := len(m.browserItems)

	m.renderVisibleItemsCompact(&listContent, visible, listWidth)

	if totalItems > visible {
		end := min(m.browserOffset+visible, totalItems)
		scrollInfo := fmt.Sprintf("[%d-%d of %d]", m.browserOffset+1, end, totalItems)
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

func (m *Model) renderVisibleItemsCompact(sb *strings.Builder, visible int, maxWidth int) {
	start := m.browserOffset
	end := min(start+visible, len(m.browserItems))

	lastType := registry.ItemType("")

	for i := start; i < end; i++ {
		bi := m.browserItems[i]

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

		m.renderItem(sb, i, false, maxWidth) // No description in split view
	}
}

func (m *Model) renderPreview(width, maxHeight int) string {
	var sb strings.Builder

	sb.WriteString(sidebarTitleStyle.Render("Preview"))
	sb.WriteString("\n\n")

	if m.browserCursor >= len(m.browserItems) {
		sb.WriteString(dimStyle.Render("No item selected"))

		return sb.String()
	}

	bi := m.browserItems[m.browserCursor]
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

func (m *Model) renderPreviewMetadata(sb *strings.Builder, bi BrowserItem) {
	bullet := bulletStyle.Render(SymbolBullet) + " "

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

	path, _ := config.GetInstallPath(bi.Item, m.selectedTool, m.selectedScope)

	sb.WriteString(bullet)
	sb.WriteString(dimStyle.Render("path: "))
	sb.WriteString(pathStyle.Render(path))
	sb.WriteString("\n\n")
}

func getStatusLabel(status installer.ItemState) string {
	switch status {
	case installer.StateNotInstalled:
		return dimStyle.Render("not installed")
	case installer.StateUpToDate:
		return installedStyle.Render("installed")
	case installer.StateUpdateAvailable:
		return updateStyle.Render("update available")
	case installer.StateModified:
		return modifiedStyle.Render("modified locally")
	case installer.StateModifiedWithUpdate:
		return updateStyle.Render("modified + update")
	default:
		return dimStyle.Render("unknown")
	}
}

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

// wrapText wraps text to the specified width while preserving existing newlines.
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder

	lines := strings.Split(text, "\n")

	for lineIdx, line := range lines {
		if lineIdx > 0 {
			result.WriteString("\n")
		}

		words := strings.Fields(line)
		lineLen := 0

		for i, word := range words {
			wordLen := len(word)

			if lineLen+wordLen+1 > width && lineLen > 0 {
				result.WriteString("\n")

				lineLen = 0
			} else if i > 0 && lineLen > 0 {
				result.WriteString(" ")

				lineLen++
			}

			result.WriteString(word)

			lineLen += wordLen
		}
	}

	return result.String()
}

func (m *Model) renderVisibleItems(sb *strings.Builder, visible int) {
	start := m.browserOffset
	end := min(start+visible, len(m.browserItems))
	lastType := registry.ItemType("")

	for i := start; i < end; i++ {
		bi := m.browserItems[i]

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

func getStatusIndicator(status installer.ItemState) (string, lipgloss.Style) {
	switch status {
	case installer.StateNotInstalled:
		return " ", normalStyle
	case installer.StateUpToDate:
		return SymbolInstalled, installedStyle
	case installer.StateUpdateAvailable:
		return SymbolUpdate, updateStyle
	case installer.StateModified:
		return SymbolModified, modifiedStyle
	case installer.StateModifiedWithUpdate:
		return SymbolUpdate + SymbolModified, updateStyle
	default:
		return " ", normalStyle
	}
}

func (m *Model) renderItem(sb *strings.Builder, idx int, showDesc bool, maxWidth int) {
	bi := m.browserItems[idx]

	var checkbox string
	if bi.Selected {
		checkbox = selectedCheckStyle.Render(SymbolSelected)
	} else {
		checkbox = dimStyle.Render(SymbolUnselected)
	}

	statusSymbol, statusStyle := getStatusIndicator(bi.Status)

	cursor := "  "
	if idx == m.browserCursor {
		cursor = accentStyle.Render(SymbolCursor) + " "
	}

	nameWidth := 20
	if maxWidth > 0 && maxWidth < 40 {
		nameWidth = maxWidth - itemPrefixWidth
	}

	namePart := fmt.Sprintf(" %-*s", nameWidth, bi.Item.Name)

	sb.WriteString(cursor)
	sb.WriteString(checkbox)
	sb.WriteString(" ")
	sb.WriteString(statusStyle.Render(statusSymbol))

	switch {
	case idx == m.browserCursor:
		sb.WriteString(selectedStyle.Render(namePart))
	case bi.Status.IsInstalled():
		sb.WriteString(statusStyle.Render(namePart))
	default:
		sb.WriteString(normalStyle.Render(namePart))
	}

	if showDesc {
		desc := bi.Item.Description
		maxDescLen := maxWidth - nameWidth - itemPrefixWidth - descPaddingExtra

		if maxDescLen > 10 && len(desc) > 0 {
			if len(desc) > maxDescLen {
				desc = desc[:maxDescLen-3] + "..."
			}

			sb.WriteString(" ")
			sb.WriteString(dimStyle.Render(desc))
		}
	}

	sb.WriteString("\n")
}

func (m *Model) countSelected() (int, int, int) {
	var total, installed, newItems int

	for _, bi := range m.browserItems {
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

// =============================================================================
// Action Menu
// =============================================================================

func (m *Model) openActionMenu() {
	selected, _, _ := m.countSelected()
	if selected == 0 && m.browserCursor < len(m.browserItems) {
		m.browserItems[m.browserCursor].Selected = true
	}

	m.buildMenuOptions()
	m.menuCursor = 0
	m.screen = ScreenActionMenu
}

func (m *Model) buildMenuOptions() {
	m.menuOptions = nil

	_, installedCount, newCount := m.countSelected()

	if newCount > 0 {
		m.menuOptions = append(m.menuOptions, MenuOption{
			Label:   fmt.Sprintf("Install (%d new)", newCount),
			Action:  ActionInstall,
			Enabled: true,
		})
	}

	if installedCount > 0 {
		m.menuOptions = append(m.menuOptions, MenuOption{
			Label:   fmt.Sprintf("Update (%d installed)", installedCount),
			Action:  ActionUpdate,
			Enabled: true,
		})

		m.menuOptions = append(m.menuOptions, MenuOption{
			Label:   fmt.Sprintf("Uninstall (%d)", installedCount),
			Action:  ActionUninstall,
			Enabled: true,
		})
	}
}

func (m *Model) updateActionMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		m.menuCursor--

		for m.menuCursor >= 0 && !m.menuOptions[m.menuCursor].Enabled {
			m.menuCursor--
		}

		if m.menuCursor < 0 {
			m.menuCursor = 0

			for !m.menuOptions[m.menuCursor].Enabled && m.menuCursor < len(m.menuOptions)-1 {
				m.menuCursor++
			}
		}
	case key.Matches(msg, keys.Down):
		m.menuCursor++

		for m.menuCursor < len(m.menuOptions) && !m.menuOptions[m.menuCursor].Enabled {
			m.menuCursor++
		}

		if m.menuCursor >= len(m.menuOptions) {
			m.menuCursor = len(m.menuOptions) - 1

			for !m.menuOptions[m.menuCursor].Enabled && m.menuCursor > 0 {
				m.menuCursor--
			}
		}
	case key.Matches(msg, keys.Enter):
		m.executeMenuAction()
	case key.Matches(msg, keys.Back):
		m.screen = ScreenBrowser
	}

	return m, nil
}

func (m *Model) executeMenuAction() {
	if m.menuCursor >= len(m.menuOptions) {
		return
	}

	opt := m.menuOptions[m.menuCursor]

	switch opt.Action {
	case ActionInstall:
		m.installNew()
	case ActionUpdate:
		m.updateInstalled()
	case ActionUninstall:
		m.uninstallSelected()
	}
}

func (m *Model) installNew() {
	installed := 0

	for i, bi := range m.browserItems {
		if !bi.Selected || bi.Status.IsInstalled() {
			continue
		}

		result, err := installer.Install(bi.Item, m.selectedTool, m.selectedScope, false)
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			installed++
			m.browserItems[i].Status = installer.StateUpToDate
		}

		m.browserItems[i].Selected = false
	}

	if installed > 0 {
		m.message = fmt.Sprintf("Installed %d items", installed)
		m.messageStyle = successMsgStyle
	}

	m.screen = ScreenBrowser
}

func (m *Model) updateInstalled() {
	updated := 0
	skippedModified := 0

	for i, bi := range m.browserItems {
		if !bi.Selected || !bi.Status.IsInstalled() {
			continue
		}

		if bi.Status.IsModified() {
			skippedModified++
			m.browserItems[i].Selected = false

			continue
		}

		result, err := installer.Install(bi.Item, m.selectedTool, m.selectedScope, true)
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			updated++
			m.browserItems[i].Status = installer.StateUpToDate
		}

		m.browserItems[i].Selected = false
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

func (m *Model) updateAllInstalled() {
	updated := 0
	skippedModified := 0

	for i, bi := range m.browserItems {
		if !bi.Status.IsInstalled() {
			continue
		}

		if bi.Status.IsModified() {
			skippedModified++

			continue
		}

		result, err := installer.Install(bi.Item, m.selectedTool, m.selectedScope, true)
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle

			return
		}

		if result.Success {
			updated++
			m.browserItems[i].Status = installer.StateUpToDate
		}

		m.browserItems[i].Selected = false
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

func (m *Model) uninstallSelected() {
	uninstalled := 0

	for i, bi := range m.browserItems {
		if !bi.Selected || !bi.Status.IsInstalled() {
			continue
		}

		result, err := installer.Uninstall(bi.Item, m.selectedTool, m.selectedScope)
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			uninstalled++
			m.browserItems[i].Status = installer.StateNotInstalled
		}

		m.browserItems[i].Selected = false
	}

	if uninstalled > 0 {
		m.message = fmt.Sprintf("Uninstalled %d items", uninstalled)
		m.messageStyle = successMsgStyle
	}

	m.screen = ScreenBrowser
}

// getItemsForAction returns the selected items that will be affected by the given action.
func (m *Model) getItemsForAction(action string) []BrowserItem {
	var items []BrowserItem

	for _, bi := range m.browserItems {
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

// renderActionPreview renders the right sidebar showing items affected by the current action.
func (m *Model) renderActionPreview(_, maxHeight int) string {
	var sb strings.Builder

	// Determine title based on current menu option
	title := "Affected Items"
	action := ""

	if m.menuCursor < len(m.menuOptions) {
		action = m.menuOptions[m.menuCursor].Action

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

// getStatusShortLabel returns a short status label for display in the action preview.
func getStatusShortLabel(status installer.ItemState) string {
	switch status {
	case installer.StateNotInstalled:
		return "new"
	case installer.StateUpToDate:
		return "installed"
	case installer.StateUpdateAvailable:
		return "update available"
	case installer.StateModified:
		return "modified"
	case installer.StateModifiedWithUpdate:
		return "modified + update"
	default:
		return "unknown"
	}
}

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

	for i, opt := range m.menuOptions {
		cursor := "  "
		if i == m.menuCursor {
			cursor = SymbolCursor + " "
		}

		switch {
		case !opt.Enabled:
			menuContent.WriteString(dimStyle.Render("  " + opt.Label))
		case i == m.menuCursor:
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
