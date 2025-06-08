package controls

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up                 key.Binding
	Down               key.Binding
	Left               key.Binding
	Right              key.Binding
	ClearFilter        key.Binding
	NextTab            key.Binding
	PrevTab            key.Binding
	PageUp             key.Binding
	PageDown           key.Binding
	Select             key.Binding
	Back               key.Binding
	Edit               key.Binding
	Manage             key.Binding
	Quit               key.Binding
	ToggleTheme        key.Binding
	ToggleInfo         key.Binding
	ToggleInfoPosition key.Binding
	ToggleCompactView  key.Binding
	Create             key.Binding
	Delete             key.Binding
	ToggleStore        key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "Back to list"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "View note"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("del", "Delete"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Next Tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "Previous Tab"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup/ctrl+u", "Page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn/ctrl+d", "Page down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit"),
	),
	Manage: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Manage note info"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Clear filter"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
	ToggleTheme: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "Toggle theme"),
	),
	ToggleInfo: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "Toggle Note Info"),
	),
	ToggleInfoPosition: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "Toggle Note Info Position"),
	),
	ToggleCompactView: key.NewBinding(
		key.WithKeys("ctrl+f"),
		key.WithHelp("ctrl+f", "Toggle Compact view only (large screen only)"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "New"),
	),
	ToggleStore: key.NewBinding(
		key.WithKeys("(", ")"),
		key.WithHelp(")", "Toggle Store"),
	),
}

func (k KeyMap) ToSlice() []key.Binding {
	var bindings []key.Binding
	v := reflect.ValueOf(k)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// Check if the field is of type key.Binding
		if binding, ok := field.Interface().(key.Binding); ok {
			bindings = append(bindings, binding)
		}
	}

	return bindings
}
