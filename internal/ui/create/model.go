package create

import (
	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/ui/navigation"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title        textinput.Model
	width        int
	height       int
	themeManager *styles.ThemeManager
	client       *api.Client
}

func (m Model) SetClient(client *api.Client) {
	m.client = client
}

type DoneMsg struct {
	Note api.Note
	Err  error
}

func NewModel(client *api.Client, themeManager *styles.ThemeManager) navigation.View {
	title := textinput.New()
	title.Placeholder = "Note title"
	title.Focus()
	title.CharLimit = 156
	title.Width = 40

	return Model{
		title:        title,
		themeManager: themeManager,
		client:       client,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return DoneMsg{Err: nil}
			}
		case "enter":
			// TODO: input validation - need a title
			note := api.Note{
				Title: m.title.Value(),
			}
			newNote, err := m.client.CreateNote(note.ToCreateRequest())
			if newNote != nil {
				note = *newNote
			}
			return m, func() tea.Msg {
				return DoneMsg{Note: note, Err: err}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.title, cmd = m.title.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.themeManager.Current().Primary).
		Padding(1, 2).
		Width(50)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Create New Note",
				"",
				"Title:",
				m.title.View(),
				"",
				"enter: save • esc: cancel",
			),
		),
	)
}
