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
	}

	// Quiet theme (light)
	Quiet = Theme{
		Name:           "quiet",
		Background:     lipgloss.Color("#ffffff"),
		Foreground:     lipgloss.Color("#333333"),
		Selection:      lipgloss.Color("#f0f0f0"),
		Comment:        lipgloss.Color("#787c99"),
		Primary:        lipgloss.Color("#4a76cd"),
		Secondary:      lipgloss.Color("#8a5cf5"),
		Tertiary:       lipgloss.Color("#d17e3d"),
		Error:          lipgloss.Color("#cc3333"),
		Warning:        lipgloss.Color("#cc9933"),
		Success:        lipgloss.Color("#339933"),
		BorderColor:    lipgloss.Color("#e0e0e0"),
		HighlightColor: lipgloss.Color("#4a76cd"),
		MutedColor:     lipgloss.Color("#999999"),
	}

	// Purple Dream theme
	PurpleDream = Theme{
		Name:           "purpledream",
		Background:     lipgloss.Color("#FFFFFF"),
		Foreground:     lipgloss.Color("#333333"),
		Selection:      lipgloss.Color("#F0F0FA"),
		Comment:        lipgloss.Color("#9B8BC0"),
		Primary:        lipgloss.Color("#7D56F4"),
		Secondary:      lipgloss.Color("#9B78FF"),
		Tertiary:       lipgloss.Color("#c39ac9"),
		Error:          lipgloss.Color("#FF0000"),
		Warning:        lipgloss.Color("#FFB86C"),
		Success:        lipgloss.Color("#50FA7B"),
		BorderColor:    lipgloss.Color("#7D56F4"),
		HighlightColor: lipgloss.Color("#FAFAFA"),
		MutedColor:     lipgloss.Color("#B8B8B8"),
	}
)
