package tui

import "github.com/charmbracelet/lipgloss"

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

const (
	LabelLocal  = "Local"
	LabelGlobal = "Global"
)

const (
	ActionInstall   = "install"
	ActionUpdate    = "update"
	ActionUninstall = "uninstall"
)

var (
	white     = lipgloss.Color("#FFFFFF")
	lightGray = lipgloss.Color("#AAAAAA")
	gray      = lipgloss.Color("#666666")
	darkGray  = lipgloss.Color("#444444")
	black     = lipgloss.Color("#000000")
	blue      = lipgloss.Color("#5F87FF")
	green     = lipgloss.Color("#00AA00")
	yellow    = lipgloss.Color("#AAAA00")
	cyan      = lipgloss.Color("#00AAAA")
	red       = lipgloss.Color("#AA0000")
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

	accentStyle = lipgloss.NewStyle().
			Foreground(blue)

	selectedStyle = lipgloss.NewStyle().
			Foreground(black).
			Background(white)

	normalStyle = lipgloss.NewStyle().
			Foreground(lightGray)

	dimStyle = lipgloss.NewStyle().
			Foreground(gray)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

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

	selectedCheckStyle = lipgloss.NewStyle().
				Foreground(green)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(darkGray).
			Padding(1, sidebarPadding)

	sidebarTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lightGray)

	previewHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	previewDividerStyle = lipgloss.NewStyle().
				Foreground(darkGray)

	previewBodyStyle = lipgloss.NewStyle().
				Foreground(lightGray)

	bulletStyle = lipgloss.NewStyle().
			Foreground(gray)

	pathStyle = lipgloss.NewStyle().
			Foreground(gray)

	helpStyle = lipgloss.NewStyle().
			Foreground(gray)
)

const (
	SymbolSelected   = "[x]"
	SymbolUnselected = "[ ]"
	SymbolInstalled  = "+"
	SymbolUpdate     = "^"
	SymbolModified   = "*"
	SymbolCursor     = ">"
	SymbolBullet     = "*"
)
