package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	list       list.Model
	filterTabs []string
	activeTab  int
	width     int
	height    int
}

// Define your list item
type Item struct {
	title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

// Initialize the model
func initialModel() Model {
	// Create some sample items for the list
	items := []list.Item{
		Item{title: "Item 1", desc: "Description 1"},
		Item{title: "Item 2", desc: "Description 2"},
		Item{title: "Item 3", desc: "Description 3"},
		// Add more items as needed
	}

	// Initialize the list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.SetShowTitle(false)
	l.Title = "My List"

	// Create tabs
	tabs := []string{
		"All Items",
		"Active",
		"Completed",
		"Important",
		"Archived",
		"Flagged",
		"Recent",
		"Shared",
		"Personal",
		"Work",
	}

	return Model{
		list:       l,
		filterTabs: tabs,
		activeTab:  0,
		width:      50, // Default width
		height:     30,  // Default height
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 3) // Reserve space for tabs
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			if m.activeTab > 0 {
				m.activeTab--
			}
		case "right", "l":
			if m.activeTab < len(m.filterTabs)-1 {
				m.activeTab++
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Style definitions
	inactiveTabStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	activeTabStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#8A2BE2")). // Purple for active tab
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	// Calculate visible tabs
	availableWidth := m.width - 4 // Account for borders and padding
	var visibleTabs []string
	var startIdx, endIdx int
	currentWidth := 0

	// First, try to center the active tab
	startIdx = max(0, m.activeTab-2)
	for i := startIdx; i < len(m.filterTabs); i++ {
		tabWidth := len(m.filterTabs[i]) + 4 // Add padding
		if currentWidth+tabWidth > availableWidth {
			break
		}
		currentWidth += tabWidth
		endIdx = i + 1
	}

	visibleTabs = m.filterTabs[startIdx:endIdx]

	// Build tab bar
	var tabBar string
	if startIdx > 0 {
		tabBar += inactiveTabStyle.Render("←")
	}

	for i, tab := range visibleTabs {
		actualIdx := startIdx + i
		if actualIdx == m.activeTab {
			tabBar += activeTabStyle.Render(tab)
		} else {
			tabBar += inactiveTabStyle.Render(tab)
		}
	}

	if endIdx < len(m.filterTabs) {
		tabBar += inactiveTabStyle.Render("→")
	}

	// Combine tab bar with list
	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabBar,
		m.list.View(),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
