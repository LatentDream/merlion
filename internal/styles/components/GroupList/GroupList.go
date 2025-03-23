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
	page          int
	scrollOffset  int

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
		page:         0,
		scrollOffset: 0,
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
	m.page = 0
	m.scrollOffset = 0
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

		case key.Matches(msg, m.keys.PageUp):
			if m.page > 0 {
				m.page--
			}

		case key.Matches(msg, m.keys.PageDown):
			// Calculate available height for items
			headerHeight := 0
			for i := 0; i <= m.selectedGroup; i++ {
				headerHeight++ // Group header
				if m.opennedGroup != nil && *m.opennedGroup == i {
					headerHeight += 3 // Notes header + border + spacing
				}
			}
			availableHeight := m.height - headerHeight
			itemsPerPage := availableHeight / (m.delegate.Height() + m.delegate.Spacing())
			maxPages := (len(m.Groups[m.selectedGroup].Items) + itemsPerPage - 1) / itemsPerPage
			if m.page < maxPages-1 {
				m.page++
			}
		}
	}

	return m, cmd
}

func (m *Model) ensureItemVisible() {
	if m.selectedItem == nil || m.opennedGroup == nil {
		return
	}

	// Calculate available height for items
	headerHeight := 0
	for i := 0; i <= m.selectedGroup; i++ {
		headerHeight++ // Group header
		if m.opennedGroup != nil && *m.opennedGroup == i {
			headerHeight += 3 // Notes header + border + spacing
		}
	}
	availableHeight := m.height - headerHeight
	itemsPerPage := availableHeight / (m.delegate.Height() + m.delegate.Spacing())
	
	// Calculate the position of the selected item relative to the current page
	relativePos := *m.selectedItem - (m.page * itemsPerPage)
	
	// If the item is below the visible area
	if relativePos >= itemsPerPage {
		m.page = *m.selectedItem / itemsPerPage
		m.scrollOffset = 0
	}
	
	// If the item is above the visible area
	if relativePos < 0 {
		m.page = *m.selectedItem / itemsPerPage
		m.scrollOffset = 0
	}
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

	// Calculate available height for items
	// Account for:
	// - Group headers (1 line each)
	// - Notes header (2 lines)
	// - Borders (2 lines)
	// - Spacing between groups (1 line)
	headerHeight := 0
	for i := 0; i <= m.selectedGroup; i++ {
		headerHeight++ // Group header
		if m.opennedGroup != nil && *m.opennedGroup == i {
			headerHeight += 3 // Notes header + border + spacing
		}
	}
	availableHeight := m.height - headerHeight

	// Calculate pagination based on available height
	itemsPerPage := availableHeight / (m.delegate.Height() + m.delegate.Spacing())
	start := m.page * itemsPerPage
	end := start + itemsPerPage
	if end > len(items) {
		end = len(items)
	}

	// Render items for current page
	for i, item := range items[start:end] {
		m.delegate.RenderGroupItems(&b, m, i+start, item)
		if i != end-start-1 {
			fmt.Fprint(&b, strings.Repeat("\n", m.delegate.Spacing()+1))
		}
	}

	// Add padding if needed
	itemsOnPage := end - start
	if itemsOnPage < itemsPerPage {
		n := (itemsPerPage - itemsOnPage) * (m.delegate.Height() + m.delegate.Spacing())
		fmt.Fprint(&b, strings.Repeat("\n", n))
	}

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

		// indication on open/closed state
		prefix := indent
		if i == t.selectedGroup {
			// When selected, reduce the indent by 2 to compensate for the border
			prefix = prefix[2:]
		}
		if t.opennedGroup != nil && *t.opennedGroup == i {
			prefix = prefix[:len(prefix)-1] + "▼ "
		} else {
			prefix = prefix[:len(prefix)-1] + "▶ "
		}

		// Create title and description
		title := prefix + tag.Name
		desc := fmt.Sprintf("%d items", noteCount)

		// Style based on selection state
		if i == t.selectedGroup {
			titleStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(theme.Primary).
				Padding(0, 0, 0, 1)
			descStyle := titleStyle.Foreground(theme.Secondary)
			s.WriteString(titleStyle.Render(title) + "\n")
			s.WriteString(descStyle.Render(indent[2:] + desc) + "\n")
		} else {
			s.WriteString(styles.Text.Render(title) + "\n")
			s.WriteString(styles.Muted.Render(indent + desc) + "\n")
		}

		// If this tag is open, list its notes
		if t.opennedGroup != nil && *t.opennedGroup == i {
			// Add separator for open tag
			s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.width)) + "\n")
			s.WriteString("\n")

			// Sub-header for notes
			noteIndent := indent + "  "

			// List notes for this tag
			if len(tag.Items) == 0 {
				s.WriteString(noteIndent + lipgloss.NewStyle().Foreground(theme.MutedColor).Render("No notes for this tag") + "\n")
			} else {
				s.WriteString(t.populatedView())
			}
			s.WriteString("\n")
			// Add bottom separator after items
			s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.width)) + "\n")
		}
	}

	return s.String()
}
