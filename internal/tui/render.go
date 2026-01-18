package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/monke/skillsmith/internal/service"
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
func getStatusIndicator(status service.ItemState) (string, lipgloss.Style) {
	switch status {
	case service.StateNotInstalled:
		return " ", normalStyle
	case service.StateUpToDate:
		return SymbolInstalled, installedStyle
	case service.StateUpdateAvailable:
		return SymbolUpdate, updateStyle
	case service.StateModified:
		return SymbolModified, modifiedStyle
	case service.StateModifiedWithUpdate:
		return SymbolUpdate + SymbolModified, updateStyle
	default:
		return " ", normalStyle
	}
}

// getStatusLabel returns a full status label for the preview sidebar.
func getStatusLabel(status service.ItemState) string {
	switch status {
	case service.StateNotInstalled:
		return dimStyle.Render("not installed")
	case service.StateUpToDate:
		return installedStyle.Render("installed")
	case service.StateUpdateAvailable:
		return updateStyle.Render("update available")
	case service.StateModified:
		return modifiedStyle.Render("modified locally")
	case service.StateModifiedWithUpdate:
		return updateStyle.Render("modified + update")
	default:
		return dimStyle.Render("unknown")
	}
}

// getStatusShortLabel returns a short status label for display in the action preview.
func getStatusShortLabel(status service.ItemState) string {
	switch status {
	case service.StateNotInstalled:
		return "new"
	case service.StateUpToDate:
		return "installed"
	case service.StateUpdateAvailable:
		return "update available"
	case service.StateModified:
		return "modified"
	case service.StateModifiedWithUpdate:
		return "modified + update"
	default:
		return "unknown"
	}
}

// getSourceTag returns a formatted source tag for non-builtin items.
func getSourceTag(source string) string {
	if source != "" && source != BuiltinSourceName {
		return fmt.Sprintf(" [%s]", source)
	}

	return ""
}

// renderItem renders a single browser list item.
func (m *Model) renderItem(sb *strings.Builder, idx int, showDesc bool, maxWidth int) {
	bi := m.browser.Items[idx]

	checkbox := dimStyle.Render(SymbolUnselected)
	if bi.Selected {
		checkbox = selectedCheckStyle.Render(SymbolSelected)
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

	namePart := fmt.Sprintf(" %-*s", nameWidth, bi.Item.Name)
	sourceTag := getSourceTag(bi.Item.Source)

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

	if sourceTag != "" {
		sb.WriteString(dimStyle.Render(sourceTag))
	}

	if showDesc {
		m.renderItemDescription(sb, bi.Item.Description, maxWidth, nameWidth, len(sourceTag))
	}

	sb.WriteString("\n")
}

// renderItemDescription renders the description portion of a list item.
func (m *Model) renderItemDescription(sb *strings.Builder, desc string, maxWidth, nameWidth, sourceTagLen int) {
	maxDescLen := maxWidth - nameWidth - itemPrefixWidth - descPaddingExtra - sourceTagLen

	if maxDescLen > 10 && len(desc) > 0 {
		if len(desc) > maxDescLen {
			desc = desc[:maxDescLen-3] + "..."
		}

		sb.WriteString(" ")
		sb.WriteString(dimStyle.Render(desc))
	}
}
