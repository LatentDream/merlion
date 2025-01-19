package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name string

	// Base colors
	Background lipgloss.Color
	Foreground lipgloss.Color
	Selection  lipgloss.Color
	Comment    lipgloss.Color

	// UI colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Tertiary  lipgloss.Color
	Error     lipgloss.Color
	Warning   lipgloss.Color
	Success   lipgloss.Color

	// Specific UI elements
	BorderColor    lipgloss.Color
	HighlightColor lipgloss.Color
	MutedColor     lipgloss.Color

	// Specific Renderer elements
	Margin     uint
	ListIndent uint
}

var (
	Gruvbox = Theme{
		Name:           "gruvbox",
		Background:     lipgloss.Color("#282828"),
		Foreground:     lipgloss.Color("#ebdbb2"),
		Selection:      lipgloss.Color("#504945"),
		Comment:        lipgloss.Color("#928374"),
		Primary:        lipgloss.Color("#fe8019"),
		Secondary:      lipgloss.Color("#83a598"),
		Tertiary:       lipgloss.Color("#b8bb26"),
		Error:          lipgloss.Color("#fb4934"),
		Warning:        lipgloss.Color("#fabd2f"),
		Success:        lipgloss.Color("#b8bb26"),
		BorderColor:    lipgloss.Color("#504945"),
		HighlightColor: lipgloss.Color("#d79921"),
		MutedColor:     lipgloss.Color("#a89984"),
		Margin:         2,
		ListIndent:     4,
	}

	// Neo Tokyo theme
	NeoTokyo = Theme{
		Name:           "neotokyo",
		Background:     lipgloss.Color("#1a1b26"),
		Foreground:     lipgloss.Color("#a9b1d6"),
		Selection:      lipgloss.Color("#2f3549"),
		Comment:        lipgloss.Color("#565f89"),
		Primary:        lipgloss.Color("#7aa2f7"),
		Secondary:      lipgloss.Color("#bb9af7"),
		Tertiary:       lipgloss.Color("#73daca"),
		Error:          lipgloss.Color("#f7768e"),
		Warning:        lipgloss.Color("#e0af68"),
		Success:        lipgloss.Color("#9ece6a"),
		BorderColor:    lipgloss.Color("#24283b"),
		HighlightColor: lipgloss.Color("#ff9e64"),
		MutedColor:     lipgloss.Color("#787c99"),
		Margin:         2,
		ListIndent:     4,
	}
)
