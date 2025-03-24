package grouplist

import (
	"fmt"
	"merlion/internal/controls"
	"merlion/internal/styles"
	"merlion/internal/styles/components/delegate"
	"merlion/internal/utils"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
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
	groupOffset   int
	visibleGroups int

	// Control
	keys controls.KeyMap

	// Styling
	width        int
	height       int
	themeManager *styles.ThemeManager
	delegate     delegate.StyledDelegate
	paginator    paginator.Model
}

const (
	POPULATED_VIEW_HEIGHT = 16
	GROUP_HEIGHT          = 3
	HEADER_HEIGHT         = 6
)

func New(groups []Group, delegate *delegate.StyledDelegate, tm *styles.ThemeManager) Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(tm.Current().Primary).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(tm.Current().MutedColor).Render("•")

	return Model{
		Groups:        groups,
		opennedGroup:  nil,
		themeManager:  tm,
		keys:          controls.Keys,
		delegate:      *delegate,
		page:          0,
		groupOffset:   0,
		visibleGroups: 0,
		paginator:     p,
	}
}

func (m Model) Width() int {
	return m.width
}

func (m Model) CurrentItemIdx() *int {
	return m.selectedItem
}

func (m Model) SelectedItem() list.Item {
	if len(m.Groups) <= m.selectedGroup {
		return nil
	}
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
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

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
			// items page
			if m.page > 0 {
				m.page--
				m.paginator.PrevPage()
			}

		case key.Matches(msg, m.keys.PageDown):
			// items page
			headerHeight := 8
			for i := 0; i <= m.selectedGroup; i++ {
				headerHeight++ // Group header
				if m.opennedGroup != nil && *m.opennedGroup == i {
					headerHeight += 3
				}
			}
			availableHeight := m.height - headerHeight
			itemsPerPage := availableHeight / (m.delegate.Height() + m.delegate.Spacing())
			maxPages := (len(m.Groups[m.selectedGroup].Items) + itemsPerPage - 1) / itemsPerPage
			if m.page < maxPages-1 {
				m.page++
				m.paginator.NextPage()
			}
		}
	}

	// Update paginator
	var paginatorCmd tea.Cmd
	m.paginator, paginatorCmd = m.paginator.Update(msg)
	if paginatorCmd != nil {
		cmds = append(cmds, paginatorCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) ensureItemVisible() {
	if m.selectedItem == nil || m.opennedGroup == nil {
		return
	}

	// Calculate how many items can fit in the view
	availableHeight := POPULATED_VIEW_HEIGHT
	itemsPerPage := availableHeight / (m.delegate.Height() + m.delegate.Spacing())

	// Update paginator total pages
	totalItems := len(m.Groups[*m.opennedGroup].Items)
	m.paginator.SetTotalPages((totalItems + itemsPerPage - 1) / itemsPerPage)

	// Calculate which page the selected item should be on
	currentPage := *m.selectedItem / itemsPerPage

	// If the selected item is not on the current page, update the page
	if currentPage != m.page {
		m.page = currentPage
		m.paginator.Page = currentPage
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
		emptyMessage := styles.Muted.Render("No items.")
		// Add padding to maintain fixed height (-2 for empty message and paginator)
		padding := strings.Repeat("\n", POPULATED_VIEW_HEIGHT-2)
		return emptyMessage + padding + "\n" + m.paginator.View()
	}

	// Calculate items that can fit in the fixed height
	availableHeight := POPULATED_VIEW_HEIGHT
	itemsPerPage := availableHeight / (m.delegate.Height() + m.delegate.Spacing())

	// Calculate pagination
	start := m.page * itemsPerPage
	end := start + itemsPerPage
	if end > len(items) {
		end = len(items)
	}

	// Update paginator settings
	totalPages := (len(items) + itemsPerPage - 1) / itemsPerPage
	m.paginator.SetTotalPages(totalPages)
	m.paginator.Page = m.page
	m.paginator.PerPage = itemsPerPage

	// Render items for current page
	for i, item := range items[start:end] {
		m.delegate.RenderGroupItems(&b, m, i+start, item)
		if i != end-start-1 {
			fmt.Fprint(&b, strings.Repeat("\n", m.delegate.Spacing()+1))
		}
	}

	// Add padding to maintain fixed height (excluding paginator)
	currentHeight := strings.Count(b.String(), "\n") + 1
	if currentHeight < POPULATED_VIEW_HEIGHT-1 { // -1 to leave space for paginator
		padding := strings.Repeat("\n", (POPULATED_VIEW_HEIGHT-1)-currentHeight)
		fmt.Fprint(&b, padding)
	}

	// Add paginator at the bottom
	paginatorStyle := lipgloss.NewStyle().
		PaddingLeft(2)
	fmt.Fprint(&b, "\n"+paginatorStyle.Render(m.paginator.View()))

	return b.String()
}

func (t Model) View() string {
	var s strings.Builder

	theme := t.themeManager.Current()
	styles := t.themeManager.Styles()
	indent := strings.Repeat(" ", 2)

	s.WriteString("\n")
	desc := fmt.Sprintf("   %d items", len(t.Groups))
	s.WriteString(styles.Muted.Render(desc) + "\n")
	s.WriteString("\n")

	// Calculate visible groups
	availableHeight := t.height - HEADER_HEIGHT
	openGroupSpace := 0
	if t.opennedGroup != nil {
		openGroupSpace = POPULATED_VIEW_HEIGHT + 4 // Add 4 for separators and spacing
	}

	// Calculate how many additional groups can be displayed
	remainingHeight := availableHeight - openGroupSpace
	t.visibleGroups = remainingHeight / GROUP_HEIGHT
	if t.visibleGroups < 1 {
		t.visibleGroups = 1 // Always show at least one group
	}

	// Ensure selected group is visible
	if t.selectedGroup < t.groupOffset {
		t.groupOffset = t.selectedGroup
	} else if t.selectedGroup >= t.groupOffset+t.visibleGroups {
		t.groupOffset = t.selectedGroup - t.visibleGroups + 1
	}

	// Ensure we don't scroll past the end
	maxOffset := len(t.Groups) - t.visibleGroups
	if maxOffset < 0 {
		maxOffset = 0
	}
	if t.groupOffset > maxOffset {
		t.groupOffset = maxOffset
	}

	// Show scroll indicators if needed
	if t.groupOffset > 0 {
		s.WriteString(styles.Muted.Render("  ▲"))
		s.WriteString("\n")
		s.WriteString("\n")
	}

	// List visible groups
	endIdx := t.groupOffset + t.visibleGroups
	if endIdx > len(t.Groups) {
		endIdx = len(t.Groups)
	}

	for i := t.groupOffset; i < endIdx; i++ {
		group := t.Groups[i]
		noteCount := len(group.Items)

		// indication on open/closed state
		prefix := indent
		if i == t.selectedGroup && t.selectedItem == nil {
			// When selected, reduce the indent by 2 to compensate for the border
			prefix = prefix[2:]
		}

		// Create title and description
		title := prefix + utils.UpperFirst(group.Name)
		desc := fmt.Sprintf("%d items", noteCount)
		if t.opennedGroup != nil && *t.opennedGroup == i {
			desc = "▼ " + desc
		} else {
			desc = "▶ " + desc
		}

		// Style based on selection state
		if i == t.selectedGroup && t.selectedItem == nil {
			titleStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(theme.Primary).
				Padding(0, 0, 0, 1)
			descStyle := titleStyle.Foreground(theme.Secondary)
			s.WriteString(titleStyle.Render(title) + "\n")
			s.WriteString(descStyle.Render(indent[2:]+desc) + "\n")
		} else {
			s.WriteString(styles.Text.Render(title) + "\n")
			s.WriteString(styles.Muted.Render(indent+desc) + "\n")
		}

		// If this tag is open, list its notes
		if t.opennedGroup != nil && *t.opennedGroup == i {
			// Add separator for open tag
			s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.width-1)) + "\n")
			s.WriteString("\n")

			// List notes for this tag
			s.WriteString(t.populatedView())

			// Add bottom separator after items
			s.WriteString("\n")
			s.WriteString(lipgloss.NewStyle().Foreground(theme.BorderColor).Render(strings.Repeat("─", t.width-1)) + "\n")
		}

		s.WriteString("\n")
	}

	// Show scroll indicators if needed
	if endIdx < len(t.Groups) {
		s.WriteString(styles.Muted.Render("  ▼\n"))
	}

	output := s.String()
	for strings.Count(output, "\n") < t.height-1 {
		output += "\n"
	}

	return output
}
