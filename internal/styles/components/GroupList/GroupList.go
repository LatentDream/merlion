package grouplist

import (
	"fmt"
	"merlion/internal/controls"
	"merlion/internal/styles"
	"merlion/internal/styles/components/delegate"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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
	width        int
	height       int
	themeManager *styles.ThemeManager
	delegate     delegate.StyledDelegate
}

func New(groups []Group, delegate *delegate.StyledDelegate, tm *styles.ThemeManager) Model {
	return Model{
		Groups:       groups,
		opennedGroup: nil,
		themeManager: tm,
		keys:         controls.Keys,
		delegate:     *delegate,
	}
}

func (m Model) Width() int {
	return m.width
}

func (m Model) CurrentItemIdx() *int {
	return m.selectedItem
}

func (m Model) SelectedItem() list.Item {
	if m.selectedItem != nil {
		return m.Groups[m.selectedGroup].Items[*m.selectedItem]
	}
	return nil
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetHeight(h int) {
	m.height = h
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

func (m Model) populatedView() string {
	if m.opennedGroup == nil {
		return ""
	}
	styles := m.themeManager.Styles()
	items := m.Groups[*m.opennedGroup].Items

	var b strings.Builder

	// Empty states
	if len(items) == 0 {
		return styles.Muted.Render("No items.")
	}

	if len(items) > 0 {
		// TODO: start, end := m.Paginator.GetSliceBounds(len(items))
		start := 0
		end := len(items) - 1
		docs := items[start:end]

		for i, item := range docs {
			m.delegate.RenderGroupItems(&b, m, i+start, item)
			if i != len(docs)-1 {
				fmt.Fprint(&b, strings.Repeat("\n", m.delegate.Spacing()+1))
			}
		}
	}

	// TODO:
	// If there aren't enough items to fill up this page (always the last page)
	// then we need to add some newlines to fill up the space where items would
	// have been.
	// itemsOnPage := m.Paginator.ItemsOnPage(len(items))
	// if itemsOnPage < m.Paginator.PerPage {
	// 	n := (m.Paginator.PerPage - itemsOnPage) * (m.delegate.Height() + m.delegate.Spacing())
	// 	if len(items) == 0 {
	// 		n -= m.delegate.Height() - 1
	// 	}
	// 	fmt.Fprint(&b, strings.Repeat("\n", n))
	// }

	return b.String()
}

func (t Model) View() string {
	var s strings.Builder

	theme := t.themeManager.Current()
	styles := t.themeManager.Styles()
	indent := strings.Repeat(" ", int(theme.ListIndent))

	// Header
	s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.width)) + "\n")

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
			s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.width-len(noteIndent))) + "\n")

			// List notes for this tag
			if len(tag.Items) == 0 {
				s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.MutedColor).Render("No notes for this tag") + "\n")
			} else {
				for i, _ := range tag.Items {
					// log.Debug(i)
					// s.WriteString(t.populatedView())

					// Prepare note display (truncate title if too long)
					// noteTitle := fmt.Sprintf("%d", i)
					// maxTitleLength := t.width - len(noteIndent) - 5 // Account for indent, bullet and spacing
					// // Apply word wrap if configured in theme
					// if theme.WordWrap > 0 && uint(len(noteTitle)) > theme.WordWrap {
					// 	noteTitle = noteTitle[:theme.WordWrap] + "..."
					// } else if len(noteTitle) > maxTitleLength {
					// 	noteTitle = noteTitle[:maxTitleLength] + "..."
					// }
					// // Style based on selection state
					// if t.selectedItem != nil && i == *t.selectedItem {
					// 	s.WriteString(noteIndent + styles.Highlight.Render("• "+noteTitle) + "\n")
					// } else {
					// 	s.WriteString(noteIndent + styles.Text.Render("• "+noteTitle) + "\n")
					// }
				}
			}
			s.WriteString("\n")
		}
	}

	return s.String()
}
