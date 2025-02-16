package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RadioInput struct {
	label    string
	checked  bool
	focused  bool
	style    lipgloss.Style
	onChange func(bool)
}

func NewRadioInput(label string) RadioInput {
	return RadioInput{
		label:   label,
		checked: false,
		style:   lipgloss.NewStyle(),
	}
}

// SetChecked sets the checked state
func (r *RadioInput) SetChecked(checked bool) {
	r.checked = checked
	if r.onChange != nil {
		r.onChange(checked)
	}
}

// IsChecked returns the current checked state
func (r RadioInput) IsChecked() bool {
	return r.checked
}

// OnChange sets the callback function for when the radio state changes
func (r *RadioInput) OnChange(fn func(bool)) {
	r.onChange = fn
}

// Focus focuses the radio input
func (r *RadioInput) Focus() {
	r.focused = true
}

// Blur removes focus from the radio input
func (r *RadioInput) Blur() {
	r.focused = false
}

// IsFocused returns whether the radio input is focused
func (r RadioInput) IsFocused() bool {
	return r.focused
}

// WithStyle sets the style for the radio input
func (r RadioInput) WithStyle(style lipgloss.Style) RadioInput {
	r.style = style
	return r
}

func (r RadioInput) Update(msg tea.Msg) (RadioInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !r.focused {
			return r, nil
		}

		switch msg.String() {
		case " ", "enter":
			r.checked = !r.checked
			if r.onChange != nil {
				r.onChange(r.checked)
			}
		}
	}
	return r, nil
}

func (r RadioInput) View() string {
	checkbox := "[ ]"
	if r.checked {
		checkbox = "[✓]"
	}

	// If focused, we can add a different style
	if r.focused {
		return r.style.Copy().Bold(true).Render(checkbox + " " + r.label)
	}

	return r.style.Render(checkbox + " " + r.label)
}
