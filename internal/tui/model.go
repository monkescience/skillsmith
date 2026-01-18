package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/registry"
)

// Model is the main application model for the TUI.
type Model struct {
	registry *registry.Registry
	screen   Screen
	width    int
	height   int
	ready    bool

	// Current selections (persisted across screens)
	selectedTool  registry.Tool
	selectedScope config.Scope

	// Screen-specific state
	toolSelect  ToolSelectState
	scopeSelect ScopeSelectState
	browser     BrowserState
	actionMenu  ActionMenuState

	// Messages
	message      string
	messageStyle lipgloss.Style
}

// NewModel creates a new TUI model with the given registry.
func NewModel(reg *registry.Registry) *Model {
	return &Model{
		registry: reg,
		screen:   ScreenToolSelect,
		toolSelect: ToolSelectState{
			Tools:  registry.AllTools(),
			Cursor: 0,
		},
		scopeSelect: ScopeSelectState{
			Scopes: []config.Scope{config.ScopeLocal, config.ScopeGlobal},
			Cursor: 0,
		},
	}
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
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

// View implements tea.Model.
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
