/* We want to be able to edit:
 * - [x] Title
 * - [x] Favorite
 * - [x] Worklog
 * - [ ] Lags on the return to Notes View
 * - [ ] Tags
 * - [ ] Workspace/folder (coming soon)
 */
package manage

import (
	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/styles/components"
	"merlion/internal/ui/navigation"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Model struct {
	width           int
	height          int
	note            *api.Note
	themeManager    *styles.ThemeManager
	client          *api.Client
	spinner         spinner.Model
	isLoading       bool
	title           textinput.Model
	isFavoriteInput components.RadioInput
	isWorkLogInput  components.RadioInput
}

func NewModel(
	client *api.Client,
	themeManager *styles.ThemeManager,
) navigation.View {
	title := textinput.New()
	title.Placeholder = "Note title"
	title.Prompt = "Title: "
	title.Focus()
	title.CharLimit = 156
	title.Width = 40
	isFavoriteInput := components.NewRadioInput("Favorite", themeManager)
	isWorkLogInput := components.NewRadioInput("Work Log", themeManager)

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(themeManager.Current().Primary)

	return Model{
		isLoading:       false,
		title:           title,
		isFavoriteInput: isFavoriteInput,
		isWorkLogInput:  isWorkLogInput,
		client:          client,
		note:            nil,
		themeManager:    themeManager,
		spinner:         sp,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

type fetchNoteResultMsg struct {
	NoteId string
	Note   *api.Note
	Error  *error
}

func fetchNote(client *api.Client, noteId string) tea.Cmd {
	// Needed as we want the content to be fetch, so we don't delete it by error
	// TODO: should handle the content.isNone nil in the backend
	return func() tea.Msg {
		res, err := client.GetNote(noteId)
		if err != nil {
			return fetchNoteResultMsg{NoteId: noteId, Error: &err}
		}
		return fetchNoteResultMsg{NoteId: res.NoteID, Note: res}
	}
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		return m, spinnerCmd

	case navigation.OpenManageMsg:
		m.isLoading = true
		cmd := fetchNote(m.client, msg.NoteId)
		return m, tea.Batch(spinner.Tick, cmd)

	case fetchNoteResultMsg:
		m.isLoading = false
		if msg.Error != nil {
			log.Fatalf("Not able to fetch note: %v", &msg.Error)
		}
		m.note = msg.Note
		m.isFavoriteInput.SetChecked(msg.Note.IsFavorite)
		m.isWorkLogInput.SetChecked(msg.Note.IsWorkLog)
		m.title.SetValue(msg.Note.Title)
		cmd := m.title.Focus()
		return m, cmd

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

		case "esc", "q":
			return m, navigation.SwitchUICmd(navigation.NoteUI)

		case "enter":
			if m.note == nil {
				log.Fatal("Trying to Update a nil note - Shouldn't be possible")
			}
			if m.isFavoriteInput.Focused || m.isWorkLogInput.Focused {
				// Update radio inputs to handle their own enter key
				var cmd tea.Cmd
				m.isFavoriteInput, cmd = m.isFavoriteInput.Update(msg)
				m.isWorkLogInput, _ = m.isWorkLogInput.Update(msg)
				return m, cmd
			}
			if m.title.Focused() {
				// TODO: input validation - need a title
				m.note.Title = m.title.Value()
				m.note.IsFavorite = m.isFavoriteInput.IsChecked()
				m.note.IsWorkLog = m.isWorkLogInput.IsChecked()
				// TODO: Handle potential Error returned + loading state
				m.client.UpdateNote(m.note.NoteID, m.note.ToCreateRequest())
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

	if m.isLoading {
		loadingStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)
		return loadingStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.spinner.View(),
				" Loading...",
			),
		)
	}

	styles := m.themeManager.Styles()

	if m.note == nil {
		return "Select a note to Manage it's Metadata"
	}

	formStyle := styles.ActiveContent.
		Padding(1, 2).
		Width(50)

	title := styles.Title.Render("Manage Note")
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

func (m Model) SetClient(client *api.Client) tea.Cmd {
	m.client = client
	return nil
}
