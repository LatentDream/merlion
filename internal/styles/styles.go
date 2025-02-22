package styles

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	App       lipgloss.Style
	Title     lipgloss.Style
	Container lipgloss.Style
	Highlight lipgloss.Style
	Input     lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Text      lipgloss.Style
	Muted     lipgloss.Style

	TitleBar        lipgloss.Style
	ActiveTab       lipgloss.Style
	InactiveTab     lipgloss.Style
	SelectedItem    lipgloss.Style
	ActiveContent   lipgloss.Style
	MobileContent   lipgloss.Style
	InactiveContent lipgloss.Style
	Controls        lipgloss.Style
	Help            lipgloss.Style
}
