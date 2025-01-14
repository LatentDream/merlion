// var (
// 	titleStyle = lipgloss.NewStyle().
// 			Bold(true).
// 			Foreground(lipgloss.Color("#FAFAFA")).
// 			Background(lipgloss.Color("#7D56F4")).
// 			Padding(0, 1)
//
// 	selectedItemStyle = lipgloss.NewStyle().
// 				Foreground(lipgloss.Color("#7D56F4")).
// 				Bold(true)
//
// 	activeContentStyle = lipgloss.NewStyle().
// 				Border(lipgloss.RoundedBorder()).
// 				BorderForeground(lipgloss.Color("#7D56F4"))
//
// 	inactiveContentStyle = lipgloss.NewStyle().
// 				Border(lipgloss.RoundedBorder()).
// 				BorderForeground(lipgloss.Color("#1A1A1A"))
//
// 	controlsStyle = lipgloss.NewStyle().
// 			BorderStyle(lipgloss.RoundedBorder()).
// 			BorderForeground(lipgloss.Color("#1A1A1A")).
// 			Padding(1, 2)
//
// 	helpStyle = lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("#626262"))
// )

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"merlion/internal/api"
	"merlion/internal/styles"
)

type focusedPanel int

const (
	noteList focusedPanel = iota
	markdown
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Select   key.Binding
	Back     key.Binding
	Quit     key.Binding
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
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "back to list"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "view note"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type item struct {
	note api.Note
}

func (i item) Title() string { return i.note.Title }
func (i item) Description() string {
	return fmt.Sprintf("Created: %s", i.note.CreatedAt.Format("2006-01-02"))
}
func (i item) FilterValue() string { return i.note.Title }

type Model struct {
	list         list.Model
	viewport     viewport.Model
	help         viewport.Model
	renderer     *glamour.TermRenderer
	spinner      spinner.Model
	keys         keyMap
	focusedPane  focusedPanel
	width        int
	height       int
	ready        bool
	loading      bool
	styles       *styles.Styles
	themeManager *styles.ThemeManager
}

func NewModel(notes []api.Note, themeManager *styles.ThemeManager) (Model, error) {
	// Get styles from theme manager
	s := themeManager.Styles()

	// Initialize glamour for markdown rendering
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return Model{}, fmt.Errorf("failed to initialize markdown renderer: %w", err)
	}

	// Initialize spinner with themed color
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(themeManager.Current().Primary)

	// Initialize list with themed styles
	delegate := list.NewDefaultDelegate()
	// delegate.Styles.SelectedTitle = s.SelectedItem
	// delegate.Styles.SelectedDesc = s.SelectedItem.
	// 	Foreground(themeManager.Current().Secondary)
	// delegate.Styles.SelectedTitle = s.SelectedItem.
	// 	Border(lipgloss.Border{
	// 		Left: "┃", // This is the pipe character
	// 	}).
	// 	BorderForeground(themeManager.Current().Primary).
	// 	Padding(0, 1)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Notes"
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = s.TitleBar

	// Initialize main content viewport
	vp := viewport.New(0, 0)
	vp.Style = s.Container

	// Initialize help viewport with themed styles
	help := viewport.New(0, 0)
	help.Style = s.Controls
	helpContent := strings.Join([]string{
		"↑/k, ↓/j: Navigate",
		"←/h, →/l: Switch pane",
		"enter: Select",
		"pgup/pgdn: Scroll",
		"esc: Back to list",
		"q: Quit",
	}, "\n")
	help.SetContent(helpContent)

	return Model{
		list:         l,
		viewport:     vp,
		help:         help,
		renderer:     renderer,
		spinner:      sp,
		keys:         keys,
		focusedPane:  noteList,
		loading:      true,
		styles:       s,
		themeManager: themeManager,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		m.loadNotes,
	)
}

// NotesLoadedMsg is sent when notes are loaded
type NotesLoadedMsg struct {
	Notes []api.Note
}

// Internal alias for better readability in the package
type notesLoadedMsg = NotesLoadedMsg

func (m Model) loadNotes() tea.Msg {
	// Simulate loading delay (remove this in production)
	// time.Sleep(2 * time.Second)

	items := make([]list.Item, len(m.list.Items()))
	for i, item := range m.list.Items() {
		items[i] = item
	}
	return notesLoadedMsg{Notes: nil} // Replace nil with actual notes
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case notesLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg.Notes))
		for i, note := range msg.Notes {
			items[i] = item{note: note}
		}
		m.list.SetItems(items)
		return m, nil
	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}
		m.width = msg.Width
		m.height = msg.Height

		// Split the view
		listWidth := m.width / 3
		contentWidth := m.width - listWidth - 4

		// Left side divisions
		listHeight := m.height - 10
		helpHeight := 6

		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)

		m.help.Width = listWidth
		m.help.Height = helpHeight

		m.viewport.Width = contentWidth
		m.viewport.Height = m.height - 2

		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			if note.Content != nil {
				rendered, err := m.renderer.Render(*note.Content)
				if err != nil {
					m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
				} else {
					m.viewport.SetContent(rendered)
				}
			} else {
				m.viewport.SetContent("No content available")
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Left):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Right):
			if m.focusedPane == noteList {
				m.focusedPane = markdown
			}

		case key.Matches(msg, m.keys.Back):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Select):
			if m.focusedPane == noteList {
				if i := m.list.SelectedItem(); i != nil {
					note := i.(item).note
					if note.Content != nil {
						rendered, err := m.renderer.Render(*note.Content)
						if err != nil {
							m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
						} else {
							m.viewport.SetContent(rendered)
						}
					} else {
						m.viewport.SetContent("No content available")
					}
					m.focusedPane = markdown
				}
			}
		}

		// Handle navigation based on focused pane
		if m.focusedPane == noteList {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			// Handle markdown viewport scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var contentStyle lipgloss.Style
	if m.focusedPane == markdown {
		contentStyle = m.styles.ActiveContent
	} else {
		contentStyle = m.styles.InactiveContent
	}

	var listView string
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Width(m.list.Width()).
			Height(m.list.Height()).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)

		listView = loadingStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.spinner.View(),
				" Loading notes...",
			),
		)
	} else {
		listView = m.list.View()
	}

	// Combine list and help section vertically
	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		listView,
		m.help.View(),
	)

	// Join left side with content viewport horizontally
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSide,
		contentStyle.Render(m.viewport.View()),
	)
}
