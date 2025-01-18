package ui

import (
	"fmt"
	"os"

	"merlion/internal/api"
	"merlion/internal/styles"
	styledDelegate "merlion/internal/styles/components/delegate"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/editor"
)

type focusedPanel int

type noteContentMsg string
type errMsg struct{ err error }

const (
	noteList focusedPanel = iota
	markdown
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Select      key.Binding
	Back        key.Binding
	Edit        key.Binding
	Quit        key.Binding
	ToggleTheme key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "back to list"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "view note"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	ToggleTheme: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "toggle theme"),
	),
}

type item struct {
	note api.Note
}

func (i item) Title() string { return i.note.Title }
func (i item) Description() string {
	return fmt.Sprintf("Created: %s", i.note.CreatedAt.Format("2006-01-02"))
}
func (i item) FilterValue() string { return i.note.Title }

type Model struct {
	list         list.Model
	viewport     viewport.Model
	renderer     *glamour.TermRenderer
	spinner      spinner.Model
	keys         keyMap
	focusedPane  focusedPanel
	width        int
	height       int
	ready        bool
	loading      bool
	listDelegate *styledDelegate.StyledDelegate
	styles       *styles.Styles
	themeManager *styles.ThemeManager
	client       *api.Client
}

func NewModel(notes []api.Note, client *api.Client, themeManager *styles.ThemeManager) (Model, error) {
	s := themeManager.Styles()

	// Initialize glamour for markdown rendering
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return Model{}, fmt.Errorf("failed to initialize markdown renderer: %w", err)
	}

	// Initialize spinner with themed color
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(themeManager.Current().Primary)

	// Initialize list with themed styles
	delegate := styledDelegate.New(themeManager)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Notes"
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = s.TitleBar

	// Initialize main content viewport
	vp := viewport.New(0, 0)
	// vp.Style = s.Container

	// Initialize help viewport with themed styles
	return Model{
		list:         l,
		viewport:     vp,
		renderer:     renderer,
		spinner:      sp,
		keys:         keys,
		focusedPane:  noteList,
		loading:      true,
		listDelegate: delegate,
		styles:       s,
		themeManager: themeManager,
		client:       client,
	}, nil
}

// Message returned when editing is complete
type editorFinishedMsg struct {
	err error
}

func (m *Model) openEditor(content string) tea.Cmd {
	// Create a temporary file for editing
	tmpfile, err := os.CreateTemp("", "note-*.md")
	if err != nil {
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("could not create temp file: %w", err)}
		}
	}

	// Write content to temp file
	if _, err := tmpfile.WriteString(content); err != nil {
		os.Remove(tmpfile.Name())
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("could not write to temp file: %w", err)}
		}
	}
	tmpfile.Close()

	// Create editor command
	cmd, err := editor.Cmd("Note", tmpfile.Name())
	if err != nil {
		os.Remove(tmpfile.Name())
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("failed to create editor command: %w", err)}
		}
	}

	// Return command that will execute editor
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		defer os.Remove(tmpfile.Name())

		if err != nil {
			return editorFinishedMsg{fmt.Errorf("editor failed: %w", err)}
		}

		// Read the edited content
		newContent, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return editorFinishedMsg{fmt.Errorf("failed to read edited content: %w", err)}
		}

		// Update note content through your API
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			content := string(newContent)
			note.Content = &content

			// Update the item in the model's list
			currentIndex := m.list.Index()
			items := m.list.Items()
			items[currentIndex] = item{note: note}
			m.list.SetItems(items)

			// Update the note to
			req := note.ToCreateRequest()
			_, err := m.client.UpdateNote(note.NoteID, req)
			if err != nil {
				log.Error("Not able to save the note %s", note.NoteID)
				return editorFinishedMsg{fmt.Errorf("failed to save the edited content: %w", err)}
			}
		}

		return editorFinishedMsg{nil}
	})
}

// Create a command to handle theme toggle asynchronously
func toggleTheme(m *Model) {
	m.styles = m.themeManager.NextTheme()

	// Update only the necessary styles instead of recreating components
	m.list.Styles.Title = m.styles.TitleBar
	m.spinner.Style = lipgloss.NewStyle().Foreground(m.themeManager.Current().Primary)

	// Update the delegate's styles without recreating the entire list
	m.listDelegate.UpdateStyles(m.themeManager)
	m.list.SetDelegate(m.listDelegate)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		m.loadNotes,
	)
}

// CMDs
func fetchNoteContent(client *api.Client, noteId string) tea.Cmd {
	return func() tea.Msg {
		res, err := client.GetNote(noteId)
		if err != nil {
			return errMsg{err}
		}
		return noteContentMsg(*res.Content)
	}
}

// NotesLoadedMsg is sent when notes are loaded
type NotesLoadedMsg struct {
	Notes []api.Note
}

// Internal alias for better readability in the package
type notesLoadedMsg = NotesLoadedMsg

func (m Model) loadNotes() tea.Msg {
	items := make([]list.Item, len(m.list.Items()))
	for i, item := range m.list.Items() {
		items[i] = item
	}
	return notesLoadedMsg{Notes: nil} // Replace nil with actual notes
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case notesLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg.Notes))
		for i, note := range msg.Notes {
			items[i] = item{note: note}
		}
		m.list.SetItems(items)
		return m, nil

	case editorFinishedMsg:
		if msg.err != nil {
			m.viewport.SetContent(fmt.Sprintf("Error editing note: %v", msg.err))
			return m, nil
		}

		// Refresh the viewport content after successful edit
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			rendered, err := m.renderer.Render(*note.Content)
			if err != nil {
				m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
			} else {
				m.viewport.SetContent(rendered)
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}
		m.width = msg.Width - 4
		m.height = msg.Height - 2

		// Account for padding and borders in the style
		horizontalPadding := m.styles.ActiveContent.GetHorizontalPadding() +
			m.styles.ActiveContent.GetHorizontalBorderSize()

		// Calculate available width after accounting for style spacing
		availableWidth := m.width - horizontalPadding

		// Split the view with adjusted measurements
		listWidth := availableWidth / 3
		contentWidth := availableWidth - listWidth

		// Left side height accounting for any vertical spacing
		listHeight := m.height

		// Update component dimensions
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)

		m.viewport.Width = contentWidth
		m.viewport.Height = m.height

		// Update content if selected
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			if note.Content != nil {
				rendered, err := m.renderer.Render(*note.Content)
				if err != nil {
					m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
				} else {
					m.viewport.SetContent(rendered)
				}
			} else {
				m.viewport.SetContent("No content available")
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Left):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Right):
			if m.focusedPane == noteList {
				m.focusedPane = markdown
			}

		case key.Matches(msg, m.keys.Back):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Edit):
			if i := m.list.SelectedItem(); i != nil {
				note := i.(item).note
				if note.Content != nil {
					return m, m.openEditor(*note.Content)
				}
			}

		case key.Matches(msg, m.keys.Select):
			if m.focusedPane == noteList {
				if i := m.list.SelectedItem(); i != nil {
					note := i.(item).note
					if note.Content == nil {
						// Set loading state
						m.viewport.SetContent("Loading note content...")
						m.loading = true
						// Fetch the note content
						return m, fetchNoteContent(m.client, note.NoteID)
					} else {
						rendered, err := m.renderer.Render(*note.Content)
						if err != nil {
							m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
						} else {
							m.viewport.SetContent(rendered)
						}
					}
					m.focusedPane = markdown
				}
			}

		case key.Matches(msg, m.keys.ToggleTheme):
			toggleTheme(&m)
		}

		// Handle navigation based on focused pane
		if m.focusedPane == noteList {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			// Handle markdown viewport scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case noteContentMsg:
		m.loading = false
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			content := string(msg)
			note.Content = &content
			rendered, err := m.renderer.Render(content)
			if err != nil {
				m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
			} else {
				m.viewport.SetContent(rendered)
			}
		}
		return m, nil

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Set up renderer style
	var rendererStyle lipgloss.Style
	if m.focusedPane == markdown {
		rendererStyle = m.styles.ActiveContent
	} else {
		rendererStyle = m.styles.InactiveContent
	}

	// Set up list style with explicit width
	var listStyle lipgloss.Style
	if m.focusedPane == noteList {
		listStyle = m.styles.ActiveContent.Width(m.width / 3) // Force width to be 1/3
	} else {
		listStyle = m.styles.InactiveContent.Width(m.width / 3) // Force width to be 1/3
	}

	var listView string
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Width(m.list.Width()).
			Height(m.list.Height()).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)

		listView = loadingStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.spinner.View(),
				" Loading notes...",
			),
		)
	} else {
		// Apply the style to both the container and the list
		listView = listStyle.Render(m.list.View())
	}

	// Combine list and help section vertically
	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		listView,
	)

	// Join left side with content viewport horizontally
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSide,
		rendererStyle.Render(m.viewport.View()),
	)
}
