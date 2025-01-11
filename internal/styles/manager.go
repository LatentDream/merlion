package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

type ThemeManager struct {
	configDir string
	current   Theme
}

func NewThemeManager(configDir string) (*ThemeManager, error) {
	tm := &ThemeManager{
		configDir: configDir,
		current:   NeoTokyo, // Default theme
	}

	// Load saved theme if exists
	if err := tm.loadTheme(); err != nil {
		// If no saved theme, save the default
		if err := tm.SaveTheme(); err != nil {
			return nil, fmt.Errorf("saving default theme: %w", err)
		}
	}

	return tm, nil
}

func (tm *ThemeManager) loadTheme() error {
	data, err := os.ReadFile(filepath.Join(tm.configDir, "theme.json"))
	if err != nil {
		return err
	}

	var themeName string
	if err := json.Unmarshal(data, &themeName); err != nil {
		return fmt.Errorf("unmarshaling theme: %w", err)
	}

	switch themeName {
	case "gruvbox":
		tm.current = Gruvbox
	case "neotokyo":
		tm.current = NeoTokyo
	case "quiet":
		tm.current = Quiet
	case "purpledream":
		tm.current = PurpleDream
	default:
		return fmt.Errorf("unknown theme: %s", themeName)
	}

	return nil
}

func (tm *ThemeManager) SaveTheme() error {
	data, err := json.Marshal(tm.current.Name)
	if err != nil {
		return fmt.Errorf("marshaling theme: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(tm.configDir, "theme.json"),
		data,
		0600,
	); err != nil {
		return fmt.Errorf("writing theme file: %w", err)
	}

	return nil
}

func (tm *ThemeManager) Current() Theme {
	return tm.current
}

func (tm *ThemeManager) SetTheme(name string) error {
	switch name {
	case "gruvbox":
		tm.current = Gruvbox
	case "neotokyo":
		tm.current = NeoTokyo
	case "quiet":
		tm.current = Quiet
	case "purpledream":
		tm.current = PurpleDream
	default:
		return fmt.Errorf("unknown theme: %s", name)
	}

	return tm.SaveTheme()
}

// Now let's create the styles that use the current theme
func (tm *ThemeManager) Styles() *Styles {
	theme := tm.Current()

	return &Styles{
		App: lipgloss.NewStyle().
			Background(theme.Background).
			Foreground(theme.Foreground),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Background(theme.Selection).
			Padding(0, 1),

		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.BorderColor).
			Padding(1, 2),

		Input: lipgloss.NewStyle().
			BorderForeground(theme.BorderColor).
			Foreground(theme.Foreground),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Padding(1, 0),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success).
			Padding(1, 0),

		Muted: lipgloss.NewStyle().
			Foreground(theme.MutedColor),
	}
}
