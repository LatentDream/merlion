package TagList

import (
	"fmt"
	"merlion/internal/api"
	"merlion/internal/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Tag struct {
	name  string
	notes []*api.Note
}

type TagList struct {
	// Logic
	Tags            []Tag
	opennedTag      *int
	highlightedTag  int
	highlightedNote int

	// Styling
	Width        int
	Height       int
	themeManager *styles.ThemeManager
}

func New(tags []Tag, tm *styles.ThemeManager) TagList {

	return TagList{
		Tags:         tags,
		opennedTag:   nil,
		themeManager: tm,
	}

}

func (t *TagList) View() string {
	var s strings.Builder

	theme := t.themeManager.Current()
	titleStyle := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	highlightedStyle := lipgloss.NewStyle().Foreground(theme.HighlightColor).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(theme.Foreground)

	// Use theme's list indent
	indent := strings.Repeat(" ", int(theme.ListIndent))

	// Header
	s.WriteString(titleStyle.Render("TAGS\n"))
	s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.Width)) + "\n")

	// List all tags
	for i, tag := range t.Tags {
		// Prepare tag display with note count
		noteCount := len(tag.notes)
		tagDisplay := fmt.Sprintf("%s (%d)", tag.name, noteCount)

		// Prefix for indicating open/closed state
		prefix := indent
		if t.opennedTag != nil && *t.opennedTag == i {
			prefix = indent[:len(indent)-1] + "▼ "
		} else {
			prefix = indent[:len(indent)-1] + "▶ "
		}

		// Style based on selection state
		if i == t.highlightedTag {
			s.WriteString(highlightedStyle.Render(prefix+tagDisplay) + "\n")
		} else {
			s.WriteString(normalStyle.Render(prefix+tagDisplay) + "\n")
		}

		// If this tag is open, list its notes
		if t.opennedTag != nil && *t.opennedTag == i {
			s.WriteString("\n")

			// Sub-header for notes
			noteIndent := indent + "  "
			s.WriteString(noteIndent + titleStyle.Render("NOTES") + "\n")
			s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.Width-len(noteIndent))) + "\n")

			// List notes for this tag
			if len(tag.notes) == 0 {
				s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.MutedColor).Render("No notes for this tag") + "\n")
			} else {
				for j, note := range tag.notes {
					// Prepare note display (truncate title if too long)
					noteTitle := note.Title
					maxTitleLength := t.Width - len(noteIndent) - 5 // Account for indent, bullet and spacing

					// Apply word wrap if configured in theme
					if theme.WordWrap > 0 && uint(len(noteTitle)) > theme.WordWrap {
						noteTitle = noteTitle[:theme.WordWrap] + "..."
					} else if len(noteTitle) > maxTitleLength {
						noteTitle = noteTitle[:maxTitleLength] + "..."
					}

					// Style based on selection state
					if i == t.highlightedTag && j == t.highlightedNote {
						s.WriteString(noteIndent + highlightedStyle.Render("• "+noteTitle) + "\n")
					} else {
						s.WriteString(noteIndent + normalStyle.Render("• "+noteTitle) + "\n")
					}
				}
			}
			s.WriteString("\n")
		}
	}

	return s.String()
}
