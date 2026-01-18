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

// Screen represents the current screen/state.
type Screen int

const (
	ScreenToolSelect Screen = iota
	ScreenScopeSelect
	ScreenBrowser
	ScreenActionMenu
)

// KeyMap defines the keybindings.
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Enter       key.Binding
	Space       key.Binding
	SelectAll   key.Binding
	DeselectAll key.Binding
	Filter      key.Binding
	Back        key.Binding
	Quit        key.Binding
	Help        key.Binding
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
	Filter: key.NewBinding(
		key.WithKeys("f"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
	),
}

// BrowserItem represents an item in the browser list.
type BrowserItem struct {
	Item      registry.Item
	Selected  bool
	Installed bool
}

// MenuOption represents an option in the action menu.
type MenuOption struct {
	Label   string
	Action  string
	Count   int
	Enabled bool
}

// Model is the main application model.
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
	browserOffset int  // scroll offset for visible window
	showAll       bool // false = only compatible, true = all

	// Action menu state
	menuOptions []MenuOption
	menuCursor  int

	// Messages
	message      string
	messageStyle lipgloss.Style
}

// NewModel creates a new Model.
func NewModel(reg *registry.Registry) *Model {
	return &Model{
		registry: reg,
		screen:   ScreenToolSelect,
		tools:    registry.AllTools(),
		scopes:   []config.Scope{config.ScopeLocal, config.ScopeGlobal},
	}
}

// Init initializes the model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
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

// View renders the UI.
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

// renderLayout renders a full-screen layout with pinned header and footer.
func (m *Model) renderLayout(header, content, footer string) string {
	// Calculate dimensions
	// Box chrome: border (2) + padding (2 top/bottom)
	boxChrome := 4
	innerWidth := m.width - boxPadding - boxChrome
	innerHeight := max(m.height-boxChrome, 1)

	// Count lines in header and footer
	headerLines := strings.Count(header, "\n") + 1
	footerLines := strings.Count(footer, "\n") + 1

	// Calculate content area height
	contentHeight := max(innerHeight-headerLines-footerLines, 1)

	// Build the layout
	var sb strings.Builder

	// Header (pinned top)
	sb.WriteString(header)
	sb.WriteString("\n")

	// Content (fills middle, with vertical padding)
	contentLines := strings.Count(content, "\n") + 1
	sb.WriteString(content)

	// Add padding to push footer to bottom
	padding := contentHeight - contentLines
	for range padding {
		sb.WriteString("\n")
	}

	// Footer (pinned bottom)
	sb.WriteString("\n")
	sb.WriteString(footer)

	return boxStyle.
		Width(innerWidth + boxInnerPadding*2).
		Height(innerHeight).
		Render(sb.String())
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
	// Header
	header := titleStyle.Render("skillsmith")

	// Content
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

		line := fmt.Sprintf("%s%s", cursor, tool)
		stats := dimStyle.Render(fmt.Sprintf("  %d agents, %d skills", agents, skills))

		if i == m.toolCursor {
			content.WriteString(selectedStyle.Render(line))
			content.WriteString(stats)
		} else {
			content.WriteString(normalStyle.Render(line))
			content.WriteString(stats)
		}

		content.WriteString("\n")
	}

	// Footer
	footer := helpStyle.Render("[enter] select  [q] quit")

	return m.renderLayout(header, content.String(), footer)
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
	// Header with breadcrumb
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))

	// Content
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

	// Footer
	footer := helpStyle.Render("[enter] select  [esc] back  [q] quit")

	return m.renderLayout(header.String(), content.String(), footer)
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
	count := 0
	items := m.registry.ByTool(m.selectedTool)

	for _, item := range items {
		installed, _, _ := installer.IsInstalled(item, m.selectedTool, scope)
		if installed {
			count++
		}
	}

	return count
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
		installed, _, _ := installer.IsInstalled(item, m.selectedTool, m.selectedScope)
		m.browserItems = append(m.browserItems, BrowserItem{
			Item:      item,
			Selected:  false,
			Installed: installed,
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
	case key.Matches(msg, keys.Enter):
		m.openActionMenu()
	case key.Matches(msg, keys.Filter):
		m.showAll = !m.showAll
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
	// Header with breadcrumb
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(m.getScopeLabel()))

	// Content
	var content strings.Builder

	// Calculate visible range
	visible := m.visibleItemCount()
	totalItems := len(m.browserItems)

	// Render visible items with scroll indicator
	m.renderVisibleItems(&content, visible)

	// Scroll indicator
	if totalItems > visible {
		end := min(m.browserOffset+visible, totalItems)
		scrollInfo := fmt.Sprintf("[%d-%d of %d]", m.browserOffset+1, end, totalItems)
		content.WriteString(dimStyle.Render(scrollInfo))
		content.WriteString("\n")
	}

	// Footer with status and help
	var footer strings.Builder

	selected, installedCount, newCount := m.countSelected()

	if selected > 0 {
		status := fmt.Sprintf("%d selected (%d installed, %d new)", selected, installedCount, newCount)
		footer.WriteString(normalStyle.Render(status))
	} else {
		footer.WriteString(dimStyle.Render("No items selected"))
	}

	footer.WriteString("\n")

	// Current item path
	if m.browserCursor < len(m.browserItems) {
		item := m.browserItems[m.browserCursor].Item
		path, _ := config.GetInstallPath(item, m.selectedTool, m.selectedScope)
		footer.WriteString(pathStyle.Render(path))
	}

	footer.WriteString("\n\n")
	footer.WriteString(helpStyle.Render("[space] toggle  [a] all  [d] none  [enter] actions  [esc] back  [q] quit"))

	return m.renderLayout(header.String(), content.String(), footer.String())
}

func (m *Model) renderVisibleItems(sb *strings.Builder, visible int) {
	// Determine which items are visible
	start := m.browserOffset
	end := min(start+visible, len(m.browserItems))

	// Track section transitions
	lastType := registry.ItemType("")

	for i := start; i < end; i++ {
		bi := m.browserItems[i]

		// Render section header on type change
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

		m.renderItem(sb, i)
	}
}

func (m *Model) renderItem(sb *strings.Builder, idx int) {
	bi := m.browserItems[idx]

	// Selection checkbox
	checkbox := SymbolUnselected
	if bi.Selected {
		checkbox = SymbolSelected
	}

	// Installed indicator
	installed := " "
	if bi.Installed {
		installed = SymbolInstalled
	}

	// Cursor
	cursor := "  "
	if idx == m.browserCursor {
		cursor = SymbolCursor + " "
	}

	// Build line
	name := bi.Item.Name
	desc := bi.Item.Description

	// Truncate description
	maxDescLen := m.width - descPaddingWidth
	if maxDescLen > 0 && len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}

	line := fmt.Sprintf("%s%s %s %-20s", cursor, checkbox, installed, name)
	descPart := dimStyle.Render(desc)

	switch {
	case idx == m.browserCursor:
		sb.WriteString(selectedStyle.Render(line))
	case bi.Installed:
		sb.WriteString(installedStyle.Render(line))
	default:
		sb.WriteString(normalStyle.Render(line))
	}

	sb.WriteString(" ")
	sb.WriteString(descPart)
	sb.WriteString("\n")
}

func (m *Model) countSelected() (int, int, int) {
	var total, installed, newItems int

	for _, bi := range m.browserItems {
		if !bi.Selected {
			continue
		}

		total++

		if bi.Installed {
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
	// If nothing selected, select current item
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

	selected, installedCount, newCount := m.countSelected()

	if newCount > 0 {
		m.menuOptions = append(m.menuOptions, MenuOption{
			Label:   fmt.Sprintf("Install new (%d)", newCount),
			Action:  "install_new",
			Count:   newCount,
			Enabled: true,
		})
	}

	if selected > 0 {
		m.menuOptions = append(m.menuOptions, MenuOption{
			Label:   fmt.Sprintf("Reinstall all (%d)", selected),
			Action:  "reinstall",
			Count:   selected,
			Enabled: true,
		})
	}

	if installedCount > 0 {
		m.menuOptions = append(m.menuOptions, MenuOption{
			Label:   fmt.Sprintf("Uninstall (%d)", installedCount),
			Action:  "uninstall",
			Count:   installedCount,
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
	case "install_new":
		m.installSelected(false)
	case "reinstall":
		m.installSelected(true)
	case "uninstall":
		m.uninstallSelected()
	}
}

func (m *Model) installSelected(force bool) {
	installed := 0
	skipped := 0

	for i, bi := range m.browserItems {
		if !bi.Selected {
			continue
		}

		// Skip already installed unless forcing
		if bi.Installed && !force {
			skipped++

			continue
		}

		result, err := installer.Install(bi.Item, m.selectedTool, m.selectedScope, force)
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageStyle = errorMsgStyle
			m.screen = ScreenBrowser

			return
		}

		if result.Success {
			installed++
			m.browserItems[i].Installed = true
		}

		m.browserItems[i].Selected = false
	}

	if installed > 0 {
		m.message = fmt.Sprintf("Installed %d items", installed)
		m.messageStyle = successMsgStyle
	} else if skipped > 0 {
		m.message = fmt.Sprintf("Skipped %d already installed", skipped)
		m.messageStyle = dimStyle
	}

	m.screen = ScreenBrowser
}

func (m *Model) uninstallSelected() {
	uninstalled := 0

	for i, bi := range m.browserItems {
		if !bi.Selected || !bi.Installed {
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
			m.browserItems[i].Installed = false
		}

		m.browserItems[i].Selected = false
	}

	if uninstalled > 0 {
		m.message = fmt.Sprintf("Uninstalled %d items", uninstalled)
		m.messageStyle = successMsgStyle
	}

	m.screen = ScreenBrowser
}

func (m *Model) viewActionMenu() string {
	// Header with breadcrumb
	var header strings.Builder

	header.WriteString(titleStyle.Render("skillsmith"))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(string(m.selectedTool)))
	header.WriteString(dimStyle.Render(" > "))
	header.WriteString(normalStyle.Render(m.getScopeLabel()))

	// Content - menu box
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

	content := menuBoxStyle.Render(menuContent.String())

	// Footer
	footer := helpStyle.Render("[enter] confirm  [esc] cancel")

	return m.renderLayout(header.String(), content, footer)
}
