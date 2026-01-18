package tui

import "github.com/charmbracelet/lipgloss"

// Layout constants.
const (
	listWidthRatio  = 3
	borderPadding   = 4
	borderMargin    = 2
	headerHeight    = 6
	listPadding     = 8
	itemLevel       = 2
	helpPadding     = 2
	statusBarGapMin = 0
)

var (
	// Colors.
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	secondaryColor = lipgloss.Color("#10B981") // Green
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	errorColor     = lipgloss.Color("#EF4444") // Red
	successColor   = lipgloss.Color("#10B981") // Green
	warningColor   = lipgloss.Color("#F59E0B") // Amber

	// Base styles.
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// List styles.
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(primaryColor).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	// Category/folder styles.
	categoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor)

	// Preview pane styles.
	previewTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(mutedColor).
				MarginBottom(1).
				PaddingBottom(1)

	previewContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E5E7EB"))

	// Status bar styles.
	statusBarStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 1)

	// Border styles.
	listBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1)

	previewBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				Padding(1)

	// Tag styles.
	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(mutedColor).
			Padding(0, 1).
			MarginRight(1)

	// Message styles.
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Help styles.
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor)
)

// TypeIcon returns an icon for item types.
func TypeIcon(itemType string) string {
	switch itemType {
	case "agent":
		return ""
	case "skill":
		return ""
	default:
		return ""
	}
}
