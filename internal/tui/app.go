package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/installer"
	"github.com/monke/skillsmith/internal/registry"
)

// KeyMap defines the keybindings for the app.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Local    key.Binding
	Global   key.Binding
	Toggle   key.Binding
	Search   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
}

var defaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "collapse"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "expand"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "install"),
	),
	Local: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "install local"),
	),
	Global: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "install global"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next pane"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev pane"),
	),
}

// TreeNode represents a node in the tree view.
type TreeNode struct {
	Name     string
	Item     *registry.Item // nil for category nodes
	Children []*TreeNode
	Expanded bool
	Level    int
}

// Focus represents which pane has focus.
type Focus int

const (
	FocusList Focus = iota
	FocusPreview
)

// Model is the main Bubbletea model.
type Model struct {
	registry     *registry.Registry
	tree         []*TreeNode
	flatList     []*TreeNode // flattened visible nodes
	cursor       int
	focus        Focus
	viewport     viewport.Model
	width        int
	height       int
	keys         KeyMap
	message      string
	messageStyle lipgloss.Style
	showHelp     bool
	ready        bool
}

// NewModel creates a new Model.
func NewModel(reg *registry.Registry) *Model {
	m := &Model{
		registry: reg,
		keys:     defaultKeyMap,
		focus:    FocusList,
	}
	m.buildTree()

	return m
}

// Init initializes the model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)

	case tea.KeyMsg:
		if cmd := m.handleKeyMsg(msg); cmd != nil {
			return m, cmd
		}
	}

	return m, nil
}

// View renders the UI.
func (m *Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.showHelp {
		return m.renderHelp()
	}

	// Calculate dimensions
	listWidth := m.width / listWidthRatio
	previewWidth := m.width - listWidth - borderPadding

	// Render header
	header := titleStyle.Render("skillsmith")
	header += "  " + subtitleStyle.Render("Install agents & skills for AI coding tools")

	// Render list
	listContent := m.renderList(listWidth - borderMargin)
	listStyle := listBorderStyle.Width(listWidth - borderMargin).Height(m.height - headerHeight)

	if m.focus == FocusList {
		listStyle = listStyle.BorderForeground(primaryColor)
	}

	list := listStyle.Render(listContent)

	// Render preview
	previewStyle := previewBorderStyle.Width(previewWidth - borderMargin).Height(m.height - headerHeight)

	if m.focus == FocusPreview {
		previewStyle = previewStyle.BorderForeground(primaryColor)
	}

	preview := previewStyle.Render(m.viewport.View())

	// Combine list and preview
	content := lipgloss.JoinHorizontal(lipgloss.Top, list, preview)

	// Render status bar
	statusBar := m.renderStatusBar()

	// Combine all
	return lipgloss.JoinVertical(lipgloss.Left, header, content, statusBar)
}

func (m *Model) handleWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height

	// Initialize or update viewport
	headerHeight := 3
	footerHeight := 2
	listWidth := m.width / listWidthRatio
	previewWidth := m.width - listWidth - borderPadding

	if !m.ready {
		m.viewport = viewport.New(previewWidth, m.height-headerHeight-footerHeight)
		m.viewport.Style = previewContentStyle
		m.ready = true
	} else {
		m.viewport.Width = previewWidth
		m.viewport.Height = m.height - headerHeight - footerHeight
	}

	m.updatePreview()
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	// Clear message on any key press
	m.message = ""

	switch {
	case key.Matches(msg, m.keys.Quit):
		return tea.Quit

	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp

	case key.Matches(msg, m.keys.Tab):
		m.toggleFocus()

	case key.Matches(msg, m.keys.Up):
		m.handleUp()

	case key.Matches(msg, m.keys.Down):
		m.handleDown()

	case key.Matches(msg, m.keys.Left):
		m.handleLeft()

	case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.Toggle):
		m.handleRight()

	case key.Matches(msg, m.keys.Enter):
		m.installCurrent(config.ScopeLocal)

	case key.Matches(msg, m.keys.Global):
		m.installCurrent(config.ScopeGlobal)
	}

	return nil
}

func (m *Model) toggleFocus() {
	if m.focus == FocusList {
		m.focus = FocusPreview
	} else {
		m.focus = FocusList
	}
}

func (m *Model) handleUp() {
	if m.focus == FocusList {
		if m.cursor > 0 {
			m.cursor--
			m.updatePreview()
		}
	} else {
		m.viewport.ScrollUp(1)
	}
}

func (m *Model) handleDown() {
	if m.focus == FocusList {
		if m.cursor < len(m.flatList)-1 {
			m.cursor++
			m.updatePreview()
		}
	} else {
		m.viewport.ScrollDown(1)
	}
}

func (m *Model) handleLeft() {
	if m.focus != FocusList {
		return
	}

	node := m.flatList[m.cursor]
	if node.Expanded && len(node.Children) > 0 {
		node.Expanded = false

		m.updateFlatList()
	}
}

func (m *Model) handleRight() {
	if m.focus != FocusList {
		return
	}

	node := m.flatList[m.cursor]
	if len(node.Children) > 0 {
		node.Expanded = !node.Expanded

		m.updateFlatList()
	}
}

// buildTree constructs the tree structure from the registry.
func (m *Model) buildTree() {
	toolNodes := make(map[registry.Tool]*TreeNode)

	// Create tool nodes
	for _, tool := range m.registry.GetTools() {
		toolNode := &TreeNode{
			Name:     string(tool),
			Expanded: true,
			Level:    0,
		}
		toolNodes[tool] = toolNode
		m.tree = append(m.tree, toolNode)
	}

	// Group items by tool and type
	for _, tool := range m.registry.GetTools() {
		toolNode := toolNodes[tool]

		// Create type groups
		agents := m.registry.ByToolAndType(tool, registry.ItemTypeAgent)
		skills := m.registry.ByToolAndType(tool, registry.ItemTypeSkill)

		if len(agents) > 0 {
			agentNode := &TreeNode{
				Name:     "Agents",
				Expanded: true,
				Level:    1,
			}

			for i := range agents {
				item := agents[i]
				agentNode.Children = append(agentNode.Children, &TreeNode{
					Name:  item.Name,
					Item:  &item,
					Level: itemLevel,
				})
			}

			toolNode.Children = append(toolNode.Children, agentNode)
		}

		if len(skills) > 0 {
			skillNode := &TreeNode{
				Name:     "Skills",
				Expanded: true,
				Level:    1,
			}

			for i := range skills {
				item := skills[i]
				skillNode.Children = append(skillNode.Children, &TreeNode{
					Name:  item.Name,
					Item:  &item,
					Level: itemLevel,
				})
			}

			toolNode.Children = append(toolNode.Children, skillNode)
		}
	}

	m.updateFlatList()
}

// updateFlatList flattens the tree based on expanded states.
func (m *Model) updateFlatList() {
	m.flatList = nil

	var flatten func(nodes []*TreeNode)

	flatten = func(nodes []*TreeNode) {
		for _, node := range nodes {
			m.flatList = append(m.flatList, node)

			if node.Expanded && len(node.Children) > 0 {
				flatten(node.Children)
			}
		}
	}

	flatten(m.tree)
}

// installCurrent installs the currently selected item.
func (m *Model) installCurrent(scope config.Scope) {
	if m.cursor >= len(m.flatList) {
		return
	}

	node := m.flatList[m.cursor]
	if node.Item == nil {
		m.message = "Select an item to install"
		m.messageStyle = warningStyle

		return
	}

	result, err := installer.Install(*node.Item, scope, false)
	if err != nil {
		m.message = fmt.Sprintf("Error: %v", err)
		m.messageStyle = errorStyle

		return
	}

	if result.Success {
		scopeStr := "locally"
		if scope == config.ScopeGlobal {
			scopeStr = "globally"
		}

		m.message = fmt.Sprintf("Installed %s %s", node.Item.Name, scopeStr)
		m.messageStyle = successStyle
	} else {
		m.message = fmt.Sprintf("%s: %s", node.Item.Name, result.Message)
		m.messageStyle = warningStyle
	}
}

// updatePreview updates the preview viewport content.
func (m *Model) updatePreview() {
	if m.cursor >= len(m.flatList) {
		return
	}

	node := m.flatList[m.cursor]
	if node.Item == nil {
		// Show category info
		m.viewport.SetContent(m.renderCategoryPreview(node))
	} else {
		m.viewport.SetContent(m.renderItemPreview(node.Item))
	}
}

// renderCategoryPreview renders preview for a category node.
func (m *Model) renderCategoryPreview(node *TreeNode) string {
	var sb strings.Builder

	title := previewTitleStyle.Render(node.Name)
	sb.WriteString(title)
	sb.WriteString("\n\n")

	childCount := len(node.Children)
	if childCount > 0 {
		sb.WriteString(fmt.Sprintf("%d items\n", childCount))
	}

	return sb.String()
}

// renderItemPreview renders preview for an item.
func (m *Model) renderItemPreview(item *registry.Item) string {
	var sb strings.Builder

	// Title
	title := previewTitleStyle.Render(fmt.Sprintf("%s %s", TypeIcon(string(item.Type)), item.Name))
	sb.WriteString(title)
	sb.WriteString("\n")

	// Description
	sb.WriteString(subtitleStyle.Render(item.Description))
	sb.WriteString("\n\n")

	// Metadata
	sb.WriteString(fmt.Sprintf("Type: %s\n", item.Type))
	sb.WriteString(fmt.Sprintf("Tool: %s\n", item.Tool))

	if item.Category != "" {
		sb.WriteString(fmt.Sprintf("Category: %s\n", item.Category))
	}

	if item.Author != "" {
		sb.WriteString(fmt.Sprintf("Author: %s\n", item.Author))
	}

	// Tags
	if len(item.Tags) > 0 {
		sb.WriteString("\nTags: ")

		for _, tag := range item.Tags {
			sb.WriteString(tagStyle.Render(tag))
		}

		sb.WriteString("\n")
	}

	// Install status
	status, err := installer.GetInstallStatus(*item)
	if err == nil {
		sb.WriteString("\n")

		if status.LocalInstalled {
			sb.WriteString(successStyle.Render("Installed locally"))
			sb.WriteString("\n")
		}

		if status.GlobalInstalled {
			sb.WriteString(successStyle.Render("Installed globally"))
			sb.WriteString("\n")
		}
	}

	// Content preview
	sb.WriteString("\n")
	sb.WriteString(categoryStyle.Render("Content:"))
	sb.WriteString("\n")
	sb.WriteString(item.Content)

	return sb.String()
}

// renderList renders the tree list.
func (m *Model) renderList(width int) string {
	var sb strings.Builder

	visibleHeight := m.height - listPadding
	start := 0
	end := len(m.flatList)

	// Scroll if needed
	if m.cursor >= visibleHeight {
		start = m.cursor - visibleHeight + 1
	}

	if end > start+visibleHeight {
		end = start + visibleHeight
	}

	for i := start; i < end; i++ {
		node := m.flatList[i]

		// Indentation
		indent := strings.Repeat("  ", node.Level)

		// Icon
		var icon string

		if len(node.Children) > 0 {
			if node.Expanded {
				icon = "▼ "
			} else {
				icon = "▶ "
			}
		} else {
			icon = "  "
		}

		// Name with type icon
		name := node.Name
		if node.Item != nil {
			name = TypeIcon(string(node.Item.Type)) + " " + name
		}

		line := indent + icon + name

		// Truncate if too long
		if len(line) > width-borderMargin {
			line = line[:width-5] + "..."
		}

		// Style based on selection
		switch {
		case i == m.cursor:
			line = selectedItemStyle.Width(width).Render(line)
		case node.Item == nil:
			// Category styling
			line = categoryStyle.Render(line)
		default:
			line = normalItemStyle.Render(line)
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderStatusBar renders the status bar.
func (m *Model) renderStatusBar() string {
	var left string

	if m.message != "" {
		left = m.messageStyle.Render(m.message)
	} else {
		left = helpStyle.Render("enter: install local | g: install global | ?: help | q: quit")
	}

	right := helpStyle.Render(fmt.Sprintf("%d/%d", m.cursor+1, len(m.flatList)))

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - borderMargin
	gap = max(gap, statusBarGapMin)

	return statusBarStyle.Width(m.width).Render(left + strings.Repeat(" ", gap) + right)
}

// renderHelp renders the help screen.
func (m *Model) renderHelp() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Keybindings"))
	sb.WriteString("\n\n")

	bindings := []struct {
		key  string
		desc string
	}{
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"←/h", "Collapse folder"},
		{"→/l", "Expand folder"},
		{"space", "Toggle expand/collapse"},
		{"enter", "Install to local project"},
		{"g", "Install globally"},
		{"tab", "Switch pane focus"},
		{"/", "Search (coming soon)"},
		{"?", "Toggle help"},
		{"q", "Quit"},
	}

	for _, b := range bindings {
		sb.WriteString(helpKeyStyle.Render(fmt.Sprintf("%8s", b.key)))
		sb.WriteString("  ")
		sb.WriteString(helpStyle.Render(b.desc))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("Press ? to close help"))

	return lipgloss.NewStyle().Padding(helpPadding).Render(sb.String())
}
