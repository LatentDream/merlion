package styles

import "github.com/charmbracelet/lipgloss"

var (
    AppStyle = lipgloss.NewStyle().
        Margin(1, 2)

    TitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        Padding(0, 1)

    StatusBarStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(lipgloss.Color("#666666")).
        Padding(0, 1)
)
