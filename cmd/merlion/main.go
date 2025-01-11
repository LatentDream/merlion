package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running merlion: %v", err)
		os.Exit(1)
	}
}

type model struct {
	content string
	width   int
	height  int
}

func initialModel() model {
	return model{
		content: "Welcome to Merlion - Markdown TUI Editor",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	// Handle window size changes
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	return fmt.Sprintf(
		"Window size: %d x %d\n\n%s\n\n(press q to quit)",
		m.width, m.height, m.content,
	)
}
