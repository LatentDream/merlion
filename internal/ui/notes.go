package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"merlion/internal/api"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	contentStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1A1A1A"))

	controlsStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type item struct {
	note api.Note
}

func (i item) Title() string       { return i.note.Title }
func (i item) Description() string { return fmt.Sprintf("Created: %s", i.note.CreatedAt.Format("2006-01-02")) }
func (i item) FilterValue() string { return i.note.Title }

type Model struct {
	list     list.Model
	viewport viewport.Model
	help     viewport.Model
	keys     keyMap
	width    int
	height   int
	ready    bool
}

func NewModel(notes []api.Note) Model {
	items := make([]list.Item, len(notes))
	for i, note := range notes {
		items[i] = item{note: note}
	}

	// Initialize list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Notes"
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	// Initialize main content viewport
	vp := viewport.New(0, 0)
	vp.Style = contentStyle

	// Initialize help viewport
	help := viewport.New(0, 0)
	help.Style = controlsStyle
	helpContent := strings.Join([]string{
		"↑/k: Move up",
		"↓/j: Move down",
		"enter: View note",
		"q: Quit",
	}, "\n")
	help.SetContent(helpContent)

	return Model{
		list:     l,
		viewport: vp,
		help:     help,
		keys:     keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}
		m.width = msg.Width
		m.height = msg.Height

		// Split the view
		listWidth := m.width / 3
		contentWidth := m.width - listWidth - 4 // Account for borders and padding
		
		// Left side divisions
		listHeight := m.height - 10 // Reserve space for help section
		helpHeight := 6             // Height for help section
		
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		
		m.help.Width = listWidth
		m.help.Height = helpHeight
		
		m.viewport.Width = contentWidth
		m.viewport.Height = m.height - 2

		// Update content when size changes
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			if note.Content != nil {
				m.viewport.SetContent(*note.Content)
			} else {
				m.viewport.SetContent("No content available")
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Select):
			if i := m.list.SelectedItem(); i != nil {
				note := i.(item).note
				if note.Content != nil {
					m.viewport.SetContent(*note.Content)
				} else {
					m.viewport.SetContent("No content available")
				}
			}
		}

		// Handle list navigation
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

		// Handle viewport scrolling
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Combine list and help section vertically
	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		m.list.View(),
		m.help.View(),
	)

	// Join left side with content viewport horizontally
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSide,
		m.viewport.View(),
	)
}
