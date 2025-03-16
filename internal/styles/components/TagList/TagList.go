package TagList

import (
	"fmt"
	"merlion/internal/api"
	"merlion/internal/controls"
	"merlion/internal/styles"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	highlightedNote *int

	// Control
	keys controls.KeyMap

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
		keys:         controls.Keys,
	}
}

func (m TagList) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m TagList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.highlightedTag > 0 {
				if m.highlightedTag == (*m.opennedTag+1) && len(m.Tags[*m.opennedTag].notes) == 0 {
					m.highlightedTag = *m.opennedTag
				} else if m.highlightedTag == (*m.opennedTag + 1) {
					m.highlightedTag = *m.opennedTag
					*m.highlightedNote = len(m.Tags[*m.opennedTag].notes) - 1
				} else if m.highlightedTag == *m.opennedTag {
					if len(m.Tags[*m.opennedTag].notes) == 0 {
						m.highlightedTag -= 0
					} else if m.highlightedNote == nil {
						*m.highlightedNote = len(m.Tags[*m.opennedTag].notes) - 1
					} else if *m.highlightedNote == 0 {
						m.highlightedNote = nil
						m.highlightedTag -= 1
					} else {
						*m.highlightedNote -= 1
					}
				} else {
					m.highlightedTag -= 1
				}
			}
		case key.Matches(msg, m.keys.Down):
			if m.highlightedTag <= len(m.Tags) {
				if m.highlightedTag == *m.opennedTag {
					if m.highlightedNote == nil && len(m.Tags[m.highlightedTag].notes) == 0 {
						m.highlightedTag += 1
					} else if m.highlightedNote == nil {
						*m.highlightedNote = 0
					} else if *m.highlightedNote == (len(m.Tags[m.highlightedTag].notes) - 1) {
						m.highlightedNote = nil
						m.highlightedTag += 1
					} else {
						*m.highlightedNote += 1
					}
				} else {
					m.highlightedTag += 1
				}
			}
		case key.Matches(msg, m.keys.Select):
			if m.highlightedNote == nil {
				m.opennedTag = &m.highlightedTag
				m.highlightedNote = nil
			}
		}
	}

	return m, cmd
}

func (t TagList) View() string {
	var s strings.Builder

	theme := t.themeManager.Current()
	styles := t.themeManager.Styles()
	indent := strings.Repeat(" ", int(theme.ListIndent))

	// Header
	s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.Width)) + "\n")

	// List all tags
	for i, tag := range t.Tags {
		noteCount := len(tag.notes)
		tagDisplay := fmt.Sprintf("%s (%d)", tag.name, noteCount)

		// indication on open/closed state
		prefix := indent
		if t.opennedTag != nil && *t.opennedTag == i {
			prefix = indent[:len(indent)-1] + "▼ "
		} else {
			prefix = indent[:len(indent)-1] + "▶ "
		}

		// Style based on selection state
		if i == t.highlightedTag {
			s.WriteString(styles.Highlight.Render(prefix+tagDisplay) + "\n")
		} else {
			s.WriteString(styles.Text.Render(prefix+tagDisplay) + "\n")
		}

		// If this tag is open, list its notes
		if t.opennedTag != nil && *t.opennedTag == i {
			s.WriteString("\n")

			// Sub-header for notes
			noteIndent := indent + "  "
			s.WriteString(noteIndent + styles.Subtitle.Render("NOTES") + "\n")
			s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.Width-len(noteIndent))) + "\n")

			// List notes for this tag
			if len(tag.notes) == 0 {
				s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.MutedColor).Render("No notes for this tag") + "\n")
			} else {
				for _, note := range tag.notes {
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
					if i == t.highlightedTag && t.highlightedNote != nil {
						s.WriteString(noteIndent + styles.Highlight.Render("• "+noteTitle) + "\n")
					} else {
						s.WriteString(noteIndent + styles.Text.Render("• "+noteTitle) + "\n")
					}
				}
			}
			s.WriteString("\n")
		}
	}

	return s.String()
}
