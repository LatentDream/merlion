package taginput

import (
	"fmt"
	"merlion/internal/styles"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Model struct {
	existingTags []string
	input        textinput.Model
	currentTags  []string
	suggestions  []string
	themeManager *styles.ThemeManager
}

func New(tags []string, tm *styles.ThemeManager, focused bool) Model {
	input := textinput.New()
	input.SetSuggestions(tags)
	input.ShowSuggestions = true
	input.Placeholder = "Note Tags"
	if focused {
		input.Focus()
	}

	return Model{
		existingTags: tags,
		input:        input,
		currentTags:  make([]string, 0),
		suggestions:  make([]string, 0),
		themeManager: tm,
	}
}

// Prefill the input with some tags
func (m *Model) SetCurrentTags(tags []string) {
	m.currentTags = tags
}

// Update the available tags proposition list
func (m *Model) SetAvailableTags(tags []string) {
	m.existingTags = tags
}

// Focus focuses the tag input
func (m *Model) Focus() {
	m.input.Focus()
}

// Blur removes focus from the tag input
func (m *Model) Blur() {
	m.input.Blur()
}

// Focused returns whether the input is focused
func (m Model) Focused() bool {
	return m.input.Focused()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			// If there are suggestions, complete with the first one
			if len(m.suggestions) > 0 {
				value := m.suggestions[0]
				m.currentTags = append(m.currentTags, value)
				m.input.SetValue("")
				m.suggestions = []string{}
				return m, nil
			}
		case "2":
			// If there are suggestions, complete with the first one
			if len(m.suggestions) > 0 {
				value := m.suggestions[1]
				m.currentTags = append(m.currentTags, value)
				m.input.SetValue("")
				m.suggestions = []string{}
				return m, nil
			}
		case "enter":
			// Add the current tag if not empty
			if value := strings.TrimSpace(m.input.Value()); value != "" {
				m.currentTags = append(m.currentTags, value)
				m.input.SetValue("")
				m.suggestions = []string{}
			}
			return m, nil
		case "backspace":
			// If input is empty and we have tags, remove the last tag
			if m.input.Value() == "" && len(m.currentTags) > 0 {
				m.currentTags = m.currentTags[:len(m.currentTags)-1]
				return m, nil
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	m.suggestions = m.getSuggestions(m.input.Value())
	return m, cmd
}

func (m Model) getSuggestions(current string) []string {
	if current == "" {
		log.Debug("empty input, no suggestions")
		return nil
	}

	// Look for matching existing tags
	var matches []string
	currentLower := strings.ToLower(current)
	log.Debug("finding suggestions", "input", current, "existing_tags", m.existingTags)

	for _, tag := range m.existingTags {
		// Skip if this tag is already in currentTags
		alreadyAdded := false
		log.Debug("Comparating: ", tag)
		for _, ct := range m.currentTags {
			if strings.EqualFold(ct, tag) {
				alreadyAdded = true
				log.Debug("skipping tag - already added", "tag", tag)
				break
			}
		}
		if alreadyAdded {
			continue
		}

		// Add tag if it matches the current input
		if strings.HasPrefix(strings.ToLower(tag), currentLower) {
			log.Debug("found matching tag", "tag", tag, "input", current)
			matches = append(matches, tag)
		}
	}
	log.Debug("suggestion results", "matches", matches)
	return matches
}

func (m Model) View() string {
	styles := m.themeManager.Styles()
	var sb strings.Builder
	sb.WriteString(styles.Input.Render("Tags:"))
	sb.WriteString("\n")

	// Show current tags
	if len(m.currentTags) > 0 {
		sb.WriteString(strings.Join(m.currentTags, ", "))
		sb.WriteString(", ")
	}

	// Show current input
	sb.WriteString(m.input.View())

	// Show suggestions if any
	if len(m.suggestions) > 0 {
		sb.WriteString(styles.Muted.Render("\nSuggestion:"))
		for i, suggestion := range m.suggestions {
			if i > 1 {
				break
			}
			sb.WriteString(styles.InactiveTab.Render(suggestion + fmt.Sprintf("[%d]", i+1)))
		}
	} else {
		sb.WriteString(styles.Muted.Render("\nSuggestion:"))
	}

	return sb.String()
}

func (m Model) GetTags() []string {
	return m.currentTags
}
