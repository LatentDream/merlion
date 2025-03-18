package grouplist

import (
	"fmt"
	"merlion/internal/controls"
	"merlion/internal/styles"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Group struct {
	Name  string
	Items []list.Item
}

type Model struct {
	// Logic
	Groups        []Group
	opennedGroup  *int
	selectedGroup int
	selectedItem  *int

	// Control
	keys controls.KeyMap

	// Styling
	Width        int
	Height       int
	themeManager *styles.ThemeManager
	delegate     list.ItemDelegate
}

func New(groups []Group, delegate list.ItemDelegate, tm *styles.ThemeManager) Model {
	return Model{
		Groups:       groups,
		opennedGroup: nil,
		themeManager: tm,
		keys:         controls.Keys,
		delegate:     delegate,
	}
}

func (m Model) SelectedItem() list.Item {
	if m.selectedItem != nil {
		return m.Groups[m.selectedGroup].Items[*m.selectedItem]
	}
	return nil
}

func (m *Model) SetWidth(w int) {
	m.Width = w
}

func (m *Model) SetHeight(h int) {
	m.Height = h
}

func (m *Model) SetGroups(groups []Group) {
	m.Groups = groups
	m.opennedGroup = nil
	m.selectedGroup = 0
	m.selectedItem = nil
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Up):
			m.handleUpNavigation()

		case key.Matches(msg, m.keys.Down):
			m.handleDownNavigation()

		case key.Matches(msg, m.keys.Select):
			m.handleSelectItem()
		}
	}

	return m, cmd
}

func (t Model) View() string {
	var s strings.Builder

	theme := t.themeManager.Current()
	styles := t.themeManager.Styles()
	indent := strings.Repeat(" ", int(theme.ListIndent))

	// Header
	s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.Width)) + "\n")

	// List all tags
	for i, tag := range t.Groups {
		noteCount := len(tag.Items)
		tagDisplay := fmt.Sprintf("%s (%d)", tag.Name, noteCount)

		// indication on open/closed state
		prefix := indent
		if t.opennedGroup != nil && *t.opennedGroup == i {
			prefix = indent[:len(indent)-1] + "▼ "
		} else {
			prefix = indent[:len(indent)-1] + "▶ "
		}

		// Style based on selection state
		if i == t.selectedGroup {
			s.WriteString(styles.Highlight.Render(prefix+tagDisplay) + "\n")
		} else {
			s.WriteString(styles.Text.Render(prefix+tagDisplay) + "\n")
		}

		// If this tag is open, list its notes
		if t.opennedGroup != nil && *t.opennedGroup == i {
			s.WriteString("\n")

			// Sub-header for notes
			noteIndent := indent + "  "
			s.WriteString(noteIndent + styles.Subtitle.Render("NOTES") + "\n")
			s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.Width-len(noteIndent))) + "\n")

			// List notes for this tag
			if len(tag.Items) == 0 {
				s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.MutedColor).Render("No notes for this tag") + "\n")
			} else {
				for i, _ := range tag.Items {
					// Prepare note display (truncate title if too long)
					noteTitle := fmt.Sprintf("%d", i)
					maxTitleLength := t.Width - len(noteIndent) - 5 // Account for indent, bullet and spacing

					// t.delegate.Render()

					// Apply word wrap if configured in theme
					if theme.WordWrap > 0 && uint(len(noteTitle)) > theme.WordWrap {
						noteTitle = noteTitle[:theme.WordWrap] + "..."
					} else if len(noteTitle) > maxTitleLength {
						noteTitle = noteTitle[:maxTitleLength] + "..."
					}

					// Style based on selection state
					if t.selectedItem != nil && i == *t.selectedItem {
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
