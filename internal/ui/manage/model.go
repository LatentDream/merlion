package manage

import (
	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/ui/navigation"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width        int
	height       int
	noteId       *string
	themeManager *styles.ThemeManager
	client       *api.Client
}

func NewModel(
	client *api.Client,
	themeManager *styles.ThemeManager,
) navigation.View {
	return Model{
		client:       client,
		noteId:       nil,
		themeManager: themeManager,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	switch msg := msg.(type) {

	case navigation.OpenManageMsg:
		m.noteId = &msg.NoteId

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, navigation.SwitchUICmd(navigation.NoteUI)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	if m.noteId == nil {
		return "Please select a note to manage it's metadata"
	}
	return *m.noteId
}

func (m Model) SetClient(client *api.Client) tea.Cmd {
	m.client = client
	return nil
}
