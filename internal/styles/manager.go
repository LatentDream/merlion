package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Config struct {
	Theme      string `json:"theme"`
	InfoHidden bool   `json:"infoHidden"`
	InfoBottom bool   `json:"infoBottom"`
	FullScreen bool   `json:"fullScreen"`
}

type ThemeManager struct {
	configDir string
	Theme     Theme
	Config    Config
}

func NewThemeManager(configDir string) (*ThemeManager, error) {
	tm := &ThemeManager{
		configDir: configDir,
		Config: Config{
			Theme:      "neotokyo",
			InfoHidden: false,
			InfoBottom: false,
			FullScreen: false,
		},
		Theme: NeoTokyo,
	}

	// Load saved config if exists
	if err := tm.loadConfig(); err != nil {
		// If no saved config, save the default
		if err := tm.SaveConfig(); err != nil {
			return nil, fmt.Errorf("saving default config: %w", err)
		}
	}

	return tm, nil
}

func (tm *ThemeManager) loadConfig() error {
	data, err := os.ReadFile(filepath.Join(tm.configDir, "config.json"))
	if err != nil {
		log.Errorf("Error loading config: %w", err)
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("unmarshaling config: %w", err)
	}

	tm.Config = config

	// Set the theme based on the config value
	switch tm.Config.Theme {
	case "gruvbox":
		tm.Theme = Gruvbox
	case "neotokyo":
		tm.Theme = NeoTokyo
	default:
		return fmt.Errorf("unknown theme: %s", tm.Config.Theme)
	}

	return nil
}

func (tm *ThemeManager) SaveConfig() error {
	data, err := json.Marshal(tm.Config)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(tm.configDir, "config.json"),
		data,
		0600,
	); err != nil {
		return fmt.Errorf("writing config file: %w", err)
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
