// Package styles contains the logic to manage the user theme
package styles

import (
	"fmt"

	"merlion/internal/config"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type ThemeManager struct {
	Theme        Theme
	Config       *config.UserConfig
	saveOnChange bool
}

func NewThemeManager() (*ThemeManager, error) {
	config := config.Load()
	tm := &ThemeManager{
		Config:       config,
		Theme:        FindThemeByName(config.Theme),
		saveOnChange: true,
	}

	return tm, nil
}

// SetSaveOnChange: Allows to remove the auto onsave when changing the config element
func (tm *ThemeManager) SetSaveOnChange(saveOnChange bool) {
	tm.saveOnChange = saveOnChange
}

func (tm *ThemeManager) SaveConfig() error {
	if !tm.saveOnChange {
		log.Info("Config not saved because saveOnChange is false")
		return nil
	}
	err := tm.Config.Save()
	return err
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
		nextTheme = "terminal"
	case "terminal":
		nextTheme = "gruvbox"
	}
	err := tm.SetTheme(nextTheme)
	if err != nil {
		log.Fatalf("Failed to toggle theme %s", err)
	}
	return tm.Styles()
}

func (tm *ThemeManager) SetTheme(name string) error {
	switch name {
	case "gruvbox":
		tm.Theme = Gruvbox
	case "neotokyo":
		tm.Theme = NeoTokyo
	case "terminal":
		tm.Theme = Terminal
	default:
		return fmt.Errorf("unknown theme: %s", name)
	}
	tm.Config.Theme = name
	return tm.SaveConfig()
}

func (tm *ThemeManager) SetInfoHidden(hide bool) error {
	tm.Config.InfoHidden = hide
	return tm.SaveConfig()
}

func (tm *ThemeManager) SetInfoBottom(hide bool) error {
	tm.Config.InfoBottom = hide
	return tm.SaveConfig()
}

func (tm *ThemeManager) SetCompactViewOnly(compact bool) error {
	tm.Config.CompactView = compact
	return tm.SaveConfig()
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
			Foreground(theme.Background).
			Background(theme.Primary).
			PaddingLeft(1).
			PaddingRight(1),

		Subtitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Background).
			Background(theme.Secondary).
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

		Text: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Muted: lipgloss.NewStyle().
			Foreground(theme.MutedColor),

		TitleMuted: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Background).
			Background(theme.MutedColor).
			PaddingLeft(1).
			PaddingRight(1),

		ActiveTab: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Background).
			Background(theme.Primary).
			Padding(0, 2),

		InactiveTab: lipgloss.NewStyle().
			// Background(lipgloss.Color(theme.MutedColor)).
			Foreground(lipgloss.Color(theme.MutedColor)).
			Padding(0, 2),

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

		MobileContent: lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder()).
			PaddingLeft(0).
			PaddingRight(0).
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
			Foreground(theme.HelpColor),
	}
}
