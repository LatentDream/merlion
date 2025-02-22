package info

import (
	"merlion/internal/api"
	"merlion/internal/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	note         *api.Note
	Width        int
	themeManager *styles.ThemeManager
}

func New(themeManager *styles.ThemeManager) Model {
	return Model{note: nil, themeManager: themeManager}
}

func (m *Model) SetNote(note *api.Note) {
	m.note = note
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func upperFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

func (m Model) View() string {
	styles := m.themeManager.Styles()
	var styleWithTopBorder = styles.
		Container.
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(false).
		Width(m.Width)

	var content string
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
		content = lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TitleMuted.Render(m.note.Title),
			styles.Muted.Render(tags),
			styles.Muted.Render(worklog),
			styles.Muted.Render(favorite),
		)
	} else {
		content = "Please select a note"
	}

	return styleWithTopBorder.Render(content)
}
