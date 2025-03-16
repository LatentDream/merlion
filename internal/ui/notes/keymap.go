package Notes

import "github.com/charmbracelet/lipgloss"

// Key Map Handler Helper ---
func toggleTheme(m *Model) {
	m.styles = m.themeManager.NextTheme()

	// Update only the necessary styles instead of recreating components
	m.noteList.Styles.Title = m.styles.Title
	m.spinner.Style = lipgloss.NewStyle().Foreground(m.themeManager.Current().Primary)

	// Update the delegate's styles without recreating the entire list
	m.listDelegate.UpdateStyles(m.themeManager)
	m.noteList.SetDelegate(m.listDelegate)

	// Update the renderer
	m.noteRenderer.RefreshTheme()

}
