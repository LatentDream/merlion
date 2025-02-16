package components

import (
	"merlion/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RadioInput struct {
	label        string
	checked      bool
	Focused      bool
	themeManager *styles.ThemeManager
	onChange     func(bool)
}

func NewRadioInput(label string, tm *styles.ThemeManager) RadioInput {
	return RadioInput{
		label:        label,
		checked:      false,
		themeManager: tm,
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
	r.Focused = true
}

// Blur removes focus from the radio input
func (r *RadioInput) Blur() {
	r.Focused = false
}

func (r RadioInput) Update(msg tea.Msg) (RadioInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !r.Focused {
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
	style := r.themeManager.Styles()
	checkbox := "[ ]"
	if r.checked {
		checkbox = "[✓]"
	}
	checkboxStyle := style.Input.BorderStyle(lipgloss.HiddenBorder())
	labelStyle := style.Text

	if r.Focused {
		checkboxStyle = checkboxStyle.Foreground(style.SelectedItem.GetForeground())
		labelStyle = style.SelectedItem
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		checkboxStyle.Render(checkbox),
		labelStyle.Render(r.label),
	)
}
