package main

import (
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Item represents a list item
type Item struct {
    title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type model struct {
    lists     []list.Model
    activeTab int
    tabs      []string
    delegate  list.DefaultDelegate
}

func initialModel() model {
    var lists []list.Model
    delegate := list.NewDefaultDelegate()
    
    // Create main list items
    mainItems := []list.Item{
        Item{title: "config.go", desc: "Main configuration file"},
        Item{title: "server.go", desc: "Server implementation"},
        Item{title: "handlers.go", desc: "HTTP handlers"},
        Item{title: "database.go", desc: "Database operations"},
    }
    
    // Create targets list items
    targetItems := []list.Item{
        Item{title: "build", desc: "Build the project"},
        Item{title: "test", desc: "Run all tests"},
        Item{title: "deploy", desc: "Deploy to production"},
        Item{title: "lint", desc: "Run linters"},
    }
    
    // Create main list
    mainList := list.New(mainItems, delegate, 0, 0)
    mainList.SetShowTitle(true)
    mainList.SetShowStatusBar(true)
    mainList.SetFilteringEnabled(true)
    mainList.Title = "main"
    
    // Create targets list
    targetsList := list.New(targetItems, delegate, 0, 0)
    targetsList.SetShowTitle(true)
    targetsList.SetShowStatusBar(true)
    targetsList.SetFilteringEnabled(true)
    targetsList.Title = "targets"
    
    lists = append(lists, mainList, targetsList)
    
    return model{
        lists:     lists,
        activeTab: 0,
        tabs:      []string{"main", "targets"},
        delegate:  delegate,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.activeTab = (m.activeTab + 1) % len(m.tabs)
            return m, nil
        case "shift+tab":
            m.activeTab--
            if m.activeTab < 0 {
                m.activeTab = len(m.tabs) - 1
            }
            return m, nil
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    }

    var cmd tea.Cmd
    m.lists[m.activeTab], cmd = m.lists[m.activeTab].Update(msg)
    return m, cmd
}

func (m model) View() string {
    // Style definitions
    inactiveTabStyle := lipgloss.NewStyle().
        Background(lipgloss.Color("#3C3C3C")).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 2)

    activeTabStyle := lipgloss.NewStyle().
        Background(lipgloss.Color("#8A2BE2")). // Purple for active tab
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 2)

    // Build tab bar
    var tabBar string
    for i, tab := range m.tabs {
        if i == m.activeTab {
            tabBar += activeTabStyle.Render(tab)
        } else {
            tabBar += inactiveTabStyle.Render(tab)
        }
    }

    // Add a separator line below tabs
    separator := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#3C3C3C")).
        Render("─")

    // Return the complete view
    return lipgloss.JoinVertical(
        lipgloss.Left,
        tabBar,
        separator,
        m.lists[m.activeTab].View(),
    )
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
