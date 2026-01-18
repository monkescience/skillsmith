package tui

import "github.com/charmbracelet/lipgloss"

// Layout constants.
const (
	boxPadding       = 4
	boxInnerPadding  = 2
	descPaddingWidth = 50
	minVisibleItems  = 3
)

// Scope labels.
const (
	LabelLocal  = "Local"
	LabelGlobal = "Global"
)

// Color palette.
var (
	white     = lipgloss.Color("#FFFFFF")
	lightGray = lipgloss.Color("#AAAAAA")
	gray      = lipgloss.Color("#666666")
	black     = lipgloss.Color("#000000")
	green     = lipgloss.Color("#00AA00") // For success/up-to-date indicators
	yellow    = lipgloss.Color("#AAAA00") // For update available
	cyan      = lipgloss.Color("#00AAAA") // For modified locally
	red       = lipgloss.Color("#AA0000") // For errors
)

// Base styles.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

	// List styles.
	selectedStyle = lipgloss.NewStyle().
			Foreground(black).
			Background(white)

	normalStyle = lipgloss.NewStyle().
			Foreground(lightGray)

	dimStyle = lipgloss.NewStyle().
			Foreground(gray)

	// Category/header styles.
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

	// Status indicators.
	installedStyle = lipgloss.NewStyle().
			Foreground(green)

	updateStyle = lipgloss.NewStyle().
			Foreground(yellow)

	modifiedStyle = lipgloss.NewStyle().
			Foreground(cyan)

	errorMsgStyle = lipgloss.NewStyle().
			Foreground(red)

	successMsgStyle = lipgloss.NewStyle().
			Foreground(green)

	// Box/border styles.
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(1, boxInnerPadding)

	menuBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lightGray).
			Padding(1, boxInnerPadding)

	// Path display.
	pathStyle = lipgloss.NewStyle().
			Foreground(lightGray)

	// Help text.
	helpStyle = lipgloss.NewStyle().
			Foreground(gray)
)

// Symbols.
const (
	SymbolSelected   = "[x]"
	SymbolUnselected = "[ ]"
	SymbolInstalled  = "✓"
	SymbolUpdate     = "↑"
	SymbolModified   = "*"
	SymbolCursor     = ">"
)
