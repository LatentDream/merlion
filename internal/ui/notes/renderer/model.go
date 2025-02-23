package renderer

import (
	"fmt"
	"merlion/internal/api"
	"merlion/internal/styles"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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
	note         *api.Note
	width        int
	height       int
	viewport     viewport.Model
	renderer     *glamour.TermRenderer
	infoHide     bool
	infoPos      Postion
	themeManager *styles.ThemeManager
	spinner      spinner.Model
}

func New(themeManager *styles.ThemeManager) Model {

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

	return Model{
		note:         nil,
		themeManager: themeManager,
		renderer:     renderer,
		viewport:     vp,
		spinner:      sp,
	}
}

func (m *Model) ToggleHideInfo() {
	m.infoHide = !m.infoHide
}

func (m *Model) SetNote(note *api.Note) {
	m.note = note
}

func (m *Model) Render() {
	if m.note == nil {
		m.viewport.SetContent("Welcome to Merlion")
		return
	}
	if m.note.Content == nil {
		m.viewport.SetContent("No Content Availalble")
		return
	}
	rendered, err := m.renderer.Render(*m.note.Content)
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Info("Msg: ", msg)
	log.Info("Viewport heights before:", "content", m.viewport.ScrollPercent(), "visible", m.viewport.Height)
	// Handle regular viewport updates
	m.viewport, cmd = m.viewport.Update(msg)
	log.Info("Viewport heights after:", "content", m.viewport.ScrollPercent(), "visible", m.viewport.Height)
	return m, cmd
}

func upperFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
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
	if m.note != nil {
		tags := ""
		if len(m.note.Tags) > 0 {
			tags += " | Tags:"
			for _, tag := range m.note.Tags {
				tags += " " + upperFirst(tag)
			}
		}
		worklog := ""
		if m.note.IsWorkLog {
			worklog = " | Work Log"
		}
		favorite := ""
		if m.note.IsFavorite {
			favorite = " | ★ "
		}
		infoBar = lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TitleMuted.Render(m.note.Title),
			styles.Muted.Render(tags),
			styles.Muted.Render(worklog),
			styles.Muted.Render(favorite),
		)
	} else {
		infoBar = "Please select a note"
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
