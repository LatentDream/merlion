package renderer

import (
	"fmt"
	"merlion/internal/model"
	"merlion/internal/store"
	"merlion/internal/styles"
	"merlion/internal/ui/navigation"
	"merlion/internal/utils"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/latentdream/merlion/lib/glamour"
	"github.com/latentdream/merlion/lib/glamour/ansi"
)

type Postion int

const (
	Bottom Postion = iota
	Top
)

// Only responsability is to Render a note
// All operations on a note need to be handle
// by the caller
type Model struct {
	Note         *model.Note
	width        int
	height       int
	viewport     viewport.Model
	renderer     *glamour.TermRenderer
	infoHide     bool
	infoPos      Postion
	storeManager *store.Manager
	themeManager *styles.ThemeManager
	spinner      spinner.Model
}

func New(themeManager *styles.ThemeManager, storeManager *store.Manager) Model {

	// Initialize glamour for markdown rendering
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(themeManager.GetRendererStyle()),
		glamour.WithWordWrap(int(themeManager.Theme.WordWrap)),
	)
	if err != nil {
		log.Fatalf("failed to initialize markdown renderer: %v", err)
	}

	// Viewport for the Content
	vp := viewport.New(0, 0)

	// Loading Spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(themeManager.Current().Primary)

	position := Top
	if themeManager.Config.InfoBottom {
		position = Bottom
	}

	return Model{
		Note:         nil,
		storeManager: storeManager,
		themeManager: themeManager,
		renderer:     renderer,
		viewport:     vp,
		spinner:      sp,
		infoHide:     themeManager.Config.InfoHidden,
		infoPos:      position,
	}
}

func (m *Model) ToggleHideInfo() {
	hideInfo := !m.infoHide
	m.infoHide = hideInfo
	m.themeManager.SetInfoHidden(hideInfo)
}

func (m *Model) ToggleHidePosition() {
	if m.infoPos == Bottom {
		m.infoPos = Top
		m.themeManager.SetInfoBottom(false)
	} else {
		m.infoPos = Bottom
		m.themeManager.SetInfoBottom(true)
	}
	if m.infoHide {
		m.ToggleHideInfo()
	}
}

func (m *Model) SetNote(note *model.Note) {
	m.renderer.ClearSelector()
	m.Note = note
}

func (m *Model) Render() {
	if m.Note == nil {
		welcomeMsg := m.themeManager.Styles().Text.Render("Welcome to Merlion")
		instructionMsg := m.themeManager.Styles().Help.Render("Select a Note or [c]reate one")

		welcomeLength := lipgloss.Width(welcomeMsg)
		instructionLength := lipgloss.Width(instructionMsg)

		verticalPadding := (m.viewport.Height - 2) / 2 // -2 for two message lines
		welcomePadding := strings.Repeat(" ", (m.viewport.Width-welcomeLength)/2)
		instructionPadding := strings.Repeat(" ", (m.viewport.Width-instructionLength)/2)

		content := strings.Repeat("\n", verticalPadding) +
			welcomePadding + welcomeMsg + "\n" +
			instructionPadding + instructionMsg

		m.viewport.SetContent(content)
		return
	}
	if m.Note.Content == nil {
		welcomeMsg := m.themeManager.Styles().Muted.Render("No Content")
		msgLenght := lipgloss.Width(welcomeMsg)
		verticalPadding := (m.viewport.Height - 2) / 2 // -2 for two message lines
		msgPadding := strings.Repeat(" ", (m.viewport.Width-msgLenght)/2)
		content := strings.Repeat("\n", verticalPadding) +
			msgPadding + welcomeMsg + "\n"

		m.viewport.SetContent(content)
		return
	}

	rendered, err := m.renderer.Render(*m.Note.Content)
	if err != nil {
		m.SetErrorMessage(fmt.Sprintf("Error rendering markdown: %v", err))
	} else {
		m.viewport.SetContent(rendered)
	}
}

func (m *Model) SetErrorMessage(msg string) {
	m.viewport.SetContent(msg)
}

func (m *Model) RefreshTheme() {

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(m.themeManager.GetRendererStyle()),
		glamour.WithWordWrap(int(m.themeManager.Theme.WordWrap)),
	)
	if err != nil {
		log.Errorf("Error while creating new renderer %v", err)
		return
	}
	m.renderer = renderer
	m.Render()
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.viewport.Width = width
	m.height = height
	m.viewport.Height = height - map[bool]int{false: 2, true: 0}[m.infoHide]
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) HandleSelector(selector *ansi.Selector) tea.Cmd {
	if selector == nil {
		log.Debug("selector is nil")
		return nil
	}
	if selector.Title != "" {
		log.Debug("Opening: ", "title", selector)
		note := m.storeManager.SearchByTitle(selector.Title)
		if note != nil {
			m.SetNote(note)
			m.Render()
		} else {
			// Create if doesn't exist
			createArgs := []any{selector.Title}
			cmd := navigation.SwitchUICmd(navigation.CreateUI, createArgs)
			return cmd
		}
	} else {
		log.Debug("Opening: ", "link", selector)
		cmd := exec.Command("open", selector.Link)
		err := cmd.Start()
		if err != nil {
			log.Error("Failed to open link", "error", err)
		}
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			nextIdx := m.renderer.GetIdxToShowAsSelected() + 1
			m.renderer.SetIdxToShowAsSelected(nextIdx)
			m.Render()
		case "shift+tab":
			nextIdx := m.renderer.GetIdxToShowAsSelected()
			if nextIdx > 0 {
				nextIdx--
			}
			m.renderer.SetIdxToShowAsSelected(nextIdx)
			m.Render()
		case "enter":
			selector := m.renderer.GetSelectorInfo()
			cmd = m.HandleSelector(selector)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	styles := m.themeManager.Styles()

	// Case: Viewport ~~
	if m.infoHide {
		m.viewport.Width = m.width
		m.viewport.Height = m.height
		return m.viewport.View()
	}

	// Case: Viewport + info ~~~~~~~~~~
	var styleWithTopBorder = styles.
		Container.
		BorderLeft(false).
		BorderRight(false).
		Width(m.width)

	var infoBar string
	if m.Note != nil {
		tags := ""
		if len(m.Note.Tags) > 0 {
			tags += " | Tags:"
			for _, tag := range m.Note.Tags {
				tags += " " + utils.UpperFirst(tag)
			}
		}
		worklog := ""
		if m.Note.IsWorkLog {
			worklog = " | Work Log"
		}
		favorite := ""
		if m.Note.IsFavorite {
			favorite = " | ★ "
		}
		infoBar = lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TitleMuted.Render(m.Note.Title),
			styles.Muted.Render(tags),
			styles.Muted.Render(worklog),
			styles.Muted.Render(favorite),
		)
	} else {
		infoBar = ""
	}

	m.viewport.Width = m.width
	m.viewport.Height = m.height - 2
	if m.infoPos == Top {
		styleWithTopBorder = styleWithTopBorder.BorderTop(false)
		return lipgloss.JoinVertical(
			lipgloss.Left,
			styleWithTopBorder.Render(infoBar),
			m.viewport.View(),
		)
	} else {
		styleWithTopBorder = styleWithTopBorder.BorderBottom(false)
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.viewport.View(),
			styleWithTopBorder.Render(infoBar),
		)
	}
}
