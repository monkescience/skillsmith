package tui

import "github.com/charmbracelet/lipgloss"

// Layout constants.
const (
	mainLeftPadding      = 2  // Left margin for main content area
	mainLeftPaddingTotal = 4  // mainLeftPadding * 2 (both sides)
	sidebarPadding       = 2  // Padding inside sidebar
	sidebarPaddingTotal  = 4  // sidebarPadding * 2 (both sides)
	minVisibleItems      = 3  // Minimum items to show in list
	minWidthForPreview   = 80 // Minimum terminal width to show preview pane
	listWidthPercent     = 40 // Percentage of width for list (sidebar gets rest)
	percentDivisor       = 100
	sidebarBorderWidth   = 2  // Border takes 2 chars (left + right)
	previewMaxLines      = 25 // Max lines to show in preview body (more room now)
	previewDividerLen    = 20 // Length of section dividers in preview
	itemPrefixWidth      = 10 // Width for cursor, checkbox, status
	descPaddingExtra     = 2  // Extra padding for description calculation
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
			Foreground(white)

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

	// Selected checkbox style.
	selectedCheckStyle = lipgloss.NewStyle().
				Foreground(green)

	// Menu box style (for action menu overlay).
	menuBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blue).
			Padding(1, sidebarPadding)

	// Sidebar styles (right panel with border).
	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(darkGray).
			Padding(1, sidebarPadding)

	sidebarTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	// Section headers inside sidebar.
	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lightGray)

	// Preview pane styles.
	previewHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	previewDividerStyle = lipgloss.NewStyle().
				Foreground(darkGray)

	previewBodyStyle = lipgloss.NewStyle().
				Foreground(lightGray)

	// Bullet style.
	bulletStyle = lipgloss.NewStyle().
			Foreground(gray)

	// Path display.
	pathStyle = lipgloss.NewStyle().
			Foreground(gray)

	// Help text.
	helpStyle = lipgloss.NewStyle().
			Foreground(gray)
)

// Symbols.
const (
	SymbolSelected   = "[x]"
	SymbolUnselected = "[ ]"
	SymbolInstalled  = "+"
	SymbolUpdate     = "^"
	SymbolModified   = "*"
	SymbolCursor     = ">"
	SymbolBullet     = "*"
)
