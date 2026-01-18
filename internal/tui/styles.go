package tui

import "github.com/charmbracelet/lipgloss"

// Layout constants.
const (
	boxPadding         = 4
	boxInnerPadding    = 2
	boxBorderWidth     = 4 // Total border width (left + right)
	descPaddingWidth   = 50
	minVisibleItems    = 3
	minWidthForPreview = 100 // Minimum terminal width to show preview pane
	listWidthPercent   = 40  // Percentage of width for list in split view
	percentDivisor     = 100
	dividerSpacing     = 3  // Space for divider and padding
	previewPadding     = 2  // Padding inside preview
	previewDividerLen  = 40 // Max length of preview divider line
	previewMaxLines    = 15 // Max lines to show in preview body
	itemPrefixWidth    = 10 // Width for cursor, checkbox, status
	descPrefixWidth    = 12 // Width for item prefix before description
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
	darkGray  = lipgloss.Color("#444444")
	black     = lipgloss.Color("#000000")
	blue      = lipgloss.Color("#5F87FF") // Accent color
	green     = lipgloss.Color("#00AA00") // For success/up-to-date indicators
	yellow    = lipgloss.Color("#AAAA00") // For update available
	cyan      = lipgloss.Color("#00AAAA") // For modified locally
	red       = lipgloss.Color("#AA0000") // For errors
)

// Base styles.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue)

	// Accent style for highlights.
	accentStyle = lipgloss.NewStyle().
			Foreground(blue)

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
			Foreground(blue)

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

	// Selected checkbox style.
	selectedCheckStyle = lipgloss.NewStyle().
				Foreground(green)

	// Box/border styles.
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(darkGray).
			Padding(1, boxInnerPadding)

	menuBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blue).
			Padding(1, boxInnerPadding)

	// Preview pane styles.
	previewHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(blue)

	previewDividerStyle = lipgloss.NewStyle().
				Foreground(darkGray)

	previewBodyStyle = lipgloss.NewStyle().
				Foreground(lightGray)

	// Path display.
	pathStyle = lipgloss.NewStyle().
			Foreground(gray)

	// Help text.
	helpStyle = lipgloss.NewStyle().
			Foreground(gray)

	// Badge styles.
	badgeAgentStyle = lipgloss.NewStyle().
			Foreground(black).
			Background(blue).
			Padding(0, 1)

	badgeSkillStyle = lipgloss.NewStyle().
			Foreground(black).
			Background(cyan).
			Padding(0, 1)
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
