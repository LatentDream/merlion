package create

import (
	"merlion/internal/model"
	"merlion/internal/store"
	"merlion/internal/styles"
	"merlion/internal/styles/components"
	taginput "merlion/internal/styles/components/tagInput"
	"merlion/internal/ui/navigation"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width           int
	height          int
	title           textinput.Model
	tagInput        taginput.Model
	isFavoriteInput components.RadioInput
	isWorkLogInput  components.RadioInput
	themeManager    *styles.ThemeManager
	storeManager    *store.Manager
}

func (m Model) SetClient(storeManager *store.Manager) tea.Cmd {
	m.storeManager = storeManager
	return nil
}

func NewModel(
	storeManager *store.Manager,
	themeManager *styles.ThemeManager,
) navigation.View {
	title := textinput.New()
	title.Placeholder = "Note title"
	title.Focus()
	title.CharLimit = 156
	title.Width = 40
	isFavoriteInput := components.NewRadioInput("Favorite", themeManager)
	isWorkLogInput := components.NewRadioInput("Work Log", themeManager)

	// Find all tags
	tags := storeManager.GetTags()

	// Initialize tag input with some sample tags
	tagInput := taginput.New(tags, themeManager, false)

	return Model{
		title:           title,
		isFavoriteInput: isFavoriteInput,
		isWorkLogInput:  isWorkLogInput,
		tagInput:        tagInput,
		themeManager:    themeManager,
		storeManager:    storeManager,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tea.WindowSize(),
	)
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	var cmd tea.Cmd

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.title.Focused() {
				m.title.Blur()
				m.tagInput.Focus()
			} else if m.tagInput.Focused() {
				m.tagInput.Blur()
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
			} else if m.tagInput.Focused() {
				m.tagInput.Blur()
				m.title.Focus()
			} else if m.isWorkLogInput.Focused {
				m.isWorkLogInput.Blur()
				m.isFavoriteInput.Focus()
			} else if m.isFavoriteInput.Focused {
				m.isFavoriteInput.Blur()
				m.tagInput.Focus()
			}
			return m, nil

		case "esc", "q":
			return m, navigation.SwitchUICmd(navigation.NoteUI)

		case "enter":
			if m.isFavoriteInput.Focused || m.isWorkLogInput.Focused {
				// Update radio inputs to handle their own enter key
				var cmd tea.Cmd
				m.isFavoriteInput, cmd = m.isFavoriteInput.Update(msg)
				m.isWorkLogInput, _ = m.isWorkLogInput.Update(msg)
				return m, cmd
			}
			if m.tagInput.Focused() {
				m.tagInput, cmd = m.tagInput.Update(msg)
				return m, cmd
			}
			if m.title.Focused() {
				// TODO: input validation - need a title
				note := model.Note{
					Title:      m.title.Value(),
					IsFavorite: m.isFavoriteInput.IsChecked(),
					IsWorkLog:  m.isWorkLogInput.IsChecked(),
					Tags:       m.tagInput.GetTags(),
				}
				// TODO: Handle potential Error returned
				m.storeManager.CreateNote(note.ToCreateRequest())
				return m, navigation.SwitchUICmd(navigation.NoteUI)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.tagInput, cmd = m.tagInput.Update(msg)
	cmds = append(cmds, cmd)

	m.title, cmd = m.title.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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
				styles.Input.Render("Title:"),
				m.title.View(),
				"",
				m.tagInput.View(),
				"",
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
