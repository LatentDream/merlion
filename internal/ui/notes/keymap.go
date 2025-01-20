package Notes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	ClearFilter key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Select      key.Binding
	Back        key.Binding
	Edit        key.Binding
	Quit        key.Binding
	ToggleTheme key.Binding
	NewNote		key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "Back to list"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Clear filter"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "View note"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "Page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "Page down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
	ToggleTheme: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "Toggle theme"),
	),
	NewNote: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New"),
	),

}

// Key Map Handler Helper ---
func toggleTheme(m *Model) {
	m.styles = m.themeManager.NextTheme()

	// Update only the necessary styles instead of recreating components
	m.list.Styles.Title = m.styles.TitleBar
	m.spinner.Style = lipgloss.NewStyle().Foreground(m.themeManager.Current().Primary)

	// Update the delegate's styles without recreating the entire list
	m.listDelegate.UpdateStyles(m.themeManager)
	m.list.SetDelegate(m.listDelegate)

	// Update the renderer
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(m.themeManager.GetRendererStyle()),
		glamour.WithWordWrap(int(m.themeManager.Theme.WordWrap)),
	)
	if err != nil {
		log.Error("Error while creating new renderer %v", err)
	} else {
		// Swap renderer
		m.renderer = renderer

		// Re-render
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			content := "Please reload the notes"
			if note.Content != nil {
				content = *note.Content
			}
			rendered, err := m.renderer.Render(content)
			if err != nil {
				log.Error("Failed to render new note after theme swap %v", err)
				m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
			} else {
				m.viewport.SetContent(rendered)
			}
		}
	}

}
