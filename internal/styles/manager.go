package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type ThemeManager struct {
	configDir string
	Theme     Theme
}

func NewThemeManager(configDir string) (*ThemeManager, error) {
	tm := &ThemeManager{
		configDir: configDir,
		Theme:     NeoTokyo, // Default theme
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
		tm.Theme = Gruvbox
	case "neotokyo":
		tm.Theme = NeoTokyo
	default:
		return fmt.Errorf("unknown theme: %s", themeName)
	}

	return nil
}

func (tm *ThemeManager) SaveTheme() error {
	data, err := json.Marshal(tm.Theme.Name)
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
	return tm.Theme
}

func (tm *ThemeManager) NextTheme() *Styles {
	var nextTheme string
	switch tm.Current().Name {
	case "gruvbox":
		nextTheme = "neotokyo"
	case "neotokyo":
		nextTheme = "gruvbox"
	}
	err := tm.SetTheme(nextTheme)
	if err != nil {
		log.Fatal("Failed to toggle theme %s", err)
	}
	return tm.Styles()
}

func (tm *ThemeManager) SetTheme(name string) error {
	switch name {
	case "gruvbox":
		tm.Theme = Gruvbox
	case "neotokyo":
		tm.Theme = NeoTokyo
	default:
		return fmt.Errorf("unknown theme: %s", name)
	}

	return tm.SaveTheme()
}

// Create the styles that use the current theme
func (tm *ThemeManager) Styles() *Styles {
	theme := tm.Current()

	// Create base container style that will be inherited
	baseContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderColor).
		PaddingLeft(1).
		PaddingRight(1).
		Width(0)

	return &Styles{
		App: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Background(theme.Selection).
			PaddingLeft(1).
			PaddingRight(1),

		Container: baseContainer,

		Highlight: lipgloss.NewStyle().
			BorderForeground(theme.BorderColor).
			Foreground(theme.Tertiary),

		Input: lipgloss.NewStyle().
			BorderForeground(theme.BorderColor).
			Foreground(theme.Secondary),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			PaddingTop(1).
			PaddingBottom(1),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success).
			PaddingTop(1).
			PaddingBottom(1),

		Muted: lipgloss.NewStyle().
			Foreground(theme.MutedColor),

		TitleBar: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Background).
			Background(theme.Primary).
			PaddingLeft(1).
			PaddingRight(1),

		SelectedItem: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		ActiveContent: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Secondary).
			PaddingLeft(1).
			PaddingRight(1).
			Width(0).
			Inherit(baseContainer),

		InactiveContent: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.BorderColor).
			PaddingLeft(1).
			PaddingRight(1).
			Width(0).
			Inherit(baseContainer),

		Controls: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.BorderColor).
			PaddingLeft(2).
			PaddingRight(2).
			PaddingTop(1).
			PaddingBottom(1),

		Help: lipgloss.NewStyle().
			Foreground(theme.MutedColor),
	}
}
