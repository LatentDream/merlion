package delegate

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"merlion/internal/styles"
)

type StyledDelegate struct {
	*list.DefaultDelegate
}

func New(themeManager *styles.ThemeManager) *StyledDelegate {
	delegate := list.NewDefaultDelegate()
	styledDelegate := &StyledDelegate{
		DefaultDelegate: &delegate,
	}
	styledDelegate.UpdateStyles(themeManager)

	return styledDelegate
}

func (d *StyledDelegate) UpdateStyles(themeManager *styles.ThemeManager) {
	selectedItem := lipgloss.NewStyle()
	d.Styles.SelectedTitle = selectedItem.
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(themeManager.Current().Primary).
		Padding(0, 0, 0, 1)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.
		Foreground(themeManager.Current().Secondary)
}
