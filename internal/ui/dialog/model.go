package dialog

import (
	"merlion/internal/store"
	"merlion/internal/styles"
	"merlion/internal/ui/navigation"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Model struct {
	title        string
	subtitle     string
	callback     func()
	level        navigation.Level
	themeManager *styles.ThemeManager
	width        int
	height       int
	storeManager *store.Manager
	confirm      bool
	returnUI     navigation.CurrentUI
}

func (m Model) SetClient(sotreManager *store.Manager) tea.Cmd {
	m.storeManager = sotreManager
	return nil
}

func NewModel(sotreManager *store.Manager, themeManager *styles.ThemeManager) navigation.View {
	return Model{
		themeManager: themeManager,
		storeManager: sotreManager,
		confirm:      false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	switch msg := msg.(type) {

	case navigation.OpenDialogMsg:
		m.title = msg.Title
		m.subtitle = msg.Subtitle
		m.level = msg.Level
		m.callback = msg.OnConfirm
		m.returnUI = msg.ReturnUI
		m.confirm = false

	case tea.KeyMsg:
		switch msg.String() {
		case "tab",
			"shift+tab",
			"l",
			"h",
			"left",
			"right":
			m.confirm = !m.confirm
			return m, nil
		case "esc", "q":
			return m, navigation.SwitchUICmd(navigation.NoteUI)
		case "enter":
			if m.confirm {
				if m.callback != nil {
					m.callback()
					return m, navigation.SwitchUICmd(m.returnUI)
				}
				log.Fatalf("No OnConfirm provided")
			} else {
				return m, navigation.SwitchUICmd(m.returnUI)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}

	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	styles := m.themeManager.Styles()
	formStyle := styles.ActiveContent.
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	// Choose style based on level
	var levelStyle lipgloss.Style
	if m.level == navigation.DangerLvl {
		levelStyle = styles.Error.PaddingTop(0).PaddingBottom(0)
	} else {
		levelStyle = styles.Text.PaddingTop(0).PaddingBottom(0)
	}

	// Render title
	title := levelStyle.Render(m.title)

	// Render subtitle if present
	var subtitleSection string
	if m.subtitle != "" {
		subtitleSection = styles.Muted.Render(m.subtitle)
	}

	// Render buttons
	var confirm string
	var cancel string
	if m.confirm {
		confirm = levelStyle.Render("[ Confirm ]")
		cancel = styles.Muted.Render("  Cancel  ")
	} else {
		confirm = styles.Muted.Render("  Confirm  ")
		cancel = styles.Text.Render("[ Cancel ]")
	}

	btnContainer := lipgloss.NewStyle().
		Width(46).
		Align(lipgloss.Center)

	btns := lipgloss.JoinHorizontal(
		lipgloss.Center,
		confirm,
		"  ",
		cancel,
	)

	// Build content sections
	sections := []string{title}
	if m.subtitle != "" {
		sections = append(sections, "", subtitleSection)
	}
	sections = append(sections, "", btnContainer.Render(btns), "")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				sections...,
			),
		),
	)
}
