package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/installer"
)

// renderLayout renders a borderless layout with pinned header and footer.
func (m *Model) renderLayout(header, content, footer string) string {
	headerLines := strings.Count(header, "\n") + 1
	footerLines := strings.Count(footer, "\n") + 1

	// Calculate content area height (full terminal height minus header/footer/margins)
	margins := 2 // Top and bottom margin
	contentHeight := max(m.height-headerLines-footerLines-margins, 1)

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

// getStatusIndicator returns the symbol and style for a given item state.
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

// getStatusLabel returns a full status label for the preview sidebar.
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

// renderItem renders a single browser list item.
func (m *Model) renderItem(sb *strings.Builder, idx int, showDesc bool, maxWidth int) {
	bi := m.browser.Items[idx]

	var checkbox string
	if bi.Selected {
		checkbox = selectedCheckStyle.Render(SymbolSelected)
	} else {
		checkbox = dimStyle.Render(SymbolUnselected)
	}

	statusSymbol, statusStyle := getStatusIndicator(bi.Status)

	cursor := "  "
	if idx == m.browser.Cursor {
		cursor = accentStyle.Render(SymbolCursor) + " "
	}

	nameWidth := 20
	if maxWidth > 0 && maxWidth < 40 {
		nameWidth = maxWidth - itemPrefixWidth
	}

	// Add external indicator for non-builtin sources
	name := bi.Item.Name
	if bi.Item.Source != "" && bi.Item.Source != BuiltinSourceName {
		name = bi.Item.Name + SymbolExternal
	}

	namePart := fmt.Sprintf(" %-*s", nameWidth, name)

	sb.WriteString(cursor)
	sb.WriteString(checkbox)
	sb.WriteString(" ")
	sb.WriteString(statusStyle.Render(statusSymbol))

	switch {
	case idx == m.browser.Cursor:
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
