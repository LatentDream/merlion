package create

import (
	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/styles/components"
	"merlion/internal/ui/navigation"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title           textinput.Model
	width           int
	height          int
	isFavoriteInput components.RadioInput
	isWorkLogInput  components.RadioInput
	themeManager    *styles.ThemeManager
	client          *api.Client
}

func (m Model) SetClient(client *api.Client) tea.Cmd {
	m.client = client
	return nil
}

func NewModel(client *api.Client, themeManager *styles.ThemeManager) navigation.View {
	title := textinput.New()
	title.Placeholder = "Note title"
	title.Prompt = "Title: "
	title.Focus()
	title.CharLimit = 156
	title.Width = 40
	isFavoriteInput := components.NewRadioInput("Favorite", themeManager)
	isWorkLogInput := components.NewRadioInput("Work Log", themeManager)

	return Model{
		title:           title,
		isFavoriteInput: isFavoriteInput,
		isWorkLogInput:  isWorkLogInput,
		themeManager:    themeManager,
		client:          client,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tea.WindowSize(),
	)
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.title.Focused() {
				m.title.Blur()
				m.isFavoriteInput.Focus()
			} else if m.isFavoriteInput.Focused {
				m.isFavoriteInput.Blur()
				m.isWorkLogInput.Focus()
			} else if m.isWorkLogInput.Focused {
				m.isWorkLogInput.Blur()
				m.title.Focus()
			}
			return m, nil
		case "shift+tab":
			if m.title.Focused() {
				m.title.Blur()
				m.isWorkLogInput.Focus()
			} else if m.isWorkLogInput.Focused {
				m.isWorkLogInput.Blur()
				m.isFavoriteInput.Focus()
			} else if m.isFavoriteInput.Focused {
				m.isFavoriteInput.Blur()
				m.title.Focus()
			}
			return m, nil
		case "esc":
			return m, navigation.SwitchUICmd(navigation.NoteUI)
		case "enter":

			if m.isFavoriteInput.Focused || m.isWorkLogInput.Focused {
				// Update radio inputs to handle their own enter key
				var cmd tea.Cmd
				m.isFavoriteInput, cmd = m.isFavoriteInput.Update(msg)
				m.isWorkLogInput, _ = m.isWorkLogInput.Update(msg)
				return m, cmd
			}
			if m.title.Focused() {
				// TODO: input validation - need a title
				note := api.Note{
					Title:      m.title.Value(),
					IsFavorite: m.isFavoriteInput.IsChecked(),
					IsWorkLog:  m.isWorkLogInput.IsChecked(),
				}
				// TODO: Handle potential Error returned
				m.client.CreateNote(note.ToCreateRequest())
				return m, navigation.SwitchUICmd(navigation.NoteUI)
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
	styles := m.themeManager.Styles()

	formStyle := styles.ActiveContent.
		Padding(1, 2).
		Width(50)

	title := styles.Title.Render("Create Note")
	help := styles.Help.Render("enter: save/toggle • tab: next • esc: cancel")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				m.title.View(),
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.isFavoriteInput.View(),
					m.isWorkLogInput.View(),
				),
				help,
			),
		),
	)
}
