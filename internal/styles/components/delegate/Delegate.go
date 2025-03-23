package delegate

import (
	"fmt"
	"io"
	"merlion/internal/styles"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	bullet   = "•"
	ellipsis = "…"
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

func (d StyledDelegate) RenderGroupItems(w io.Writer, m InputModel, index int, item list.Item) {
	var (
		title, desc string
		s           = &d.Styles
	)

	if i, ok := item.(list.DefaultItem); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.Height()-1 {
				break
			}
			lines = append(lines, ansi.Truncate(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	// Conditions
	var isSelected = false
	if m.CurrentItemIdx() != nil && index == *m.CurrentItemIdx() {
		isSelected = true
	}

	if isSelected {
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Render(desc)
	} else {
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc)
		return
	}
	fmt.Fprintf(w, "%s", title)
}

type InputModel interface {
	Width() int
	CurrentItemIdx() *int
}
