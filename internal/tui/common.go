package tui

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/monke/skillsmith/internal/service"
)

// Screen represents the current screen state.
type Screen int

const (
	ScreenToolSelect Screen = iota
	ScreenScopeSelect
	ScreenBrowser
	ScreenActionMenu
)

// KeyMap defines all keyboard shortcuts.
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

// BrowserItem represents an item in the browser list.
type BrowserItem struct {
	Item     service.Item
	Selected bool
	Status   service.ItemState
}

// MenuOption represents an option in the action menu.
type MenuOption struct {
	Label   string
	Action  string
	Enabled bool
}

// ToolSelectState holds state for the tool selection screen.
type ToolSelectState struct {
	Tools  []service.Tool
	Cursor int
}

// ScopeSelectState holds state for the scope selection screen.
type ScopeSelectState struct {
	Scopes []service.Scope
	Cursor int
}

// BrowserState holds state for the browser screen.
type BrowserState struct {
	Items  []BrowserItem
	Cursor int
	Offset int // scroll offset for visible window
}

// ActionMenuState holds state for the action menu screen.
type ActionMenuState struct {
	Options []MenuOption
	Cursor  int
}
