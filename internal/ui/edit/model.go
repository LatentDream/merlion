package edit

import (
	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/ui/navigation"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width        int
	height       int
	themeManager *styles.ThemeManager
	client       *api.Client
}

func NewModel(
	client *api.Client,
	themeManager *styles.ThemeManager,
) navigation.View {
	return Model{
		client:       client,
		themeManager: themeManager,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}

func (m Model) SetClient(client *api.Client) tea.Cmd {
	m.client = client
	return nil
}
