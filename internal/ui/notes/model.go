package Notes

import (
	"fmt"

	"merlion/internal/api"
	"merlion/internal/styles"
	styledDelegate "merlion/internal/styles/components/delegate"
	"merlion/internal/ui/create"
	"merlion/internal/ui/navigation"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type focusedPanel int

const (
	noteList focusedPanel = iota
	markdown
)

const ViewRatio = 4

type item struct {
	note api.Note
}

type ViewState string

const (
	MainView   ViewState = "main"
	CreateView ViewState = "create"
)

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
	createModel  create.Model
}

func NewModel(client *api.Client, themeManager *styles.ThemeManager) Model {
	s := themeManager.Styles()

	// Initialize glamour for markdown rendering
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(themeManager.GetRendererStyle()),
		glamour.WithWordWrap(int(themeManager.Theme.WordWrap)),
	)
	if err != nil {
		log.Fatal("failed to initialize markdown renderer: %v", err)
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
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.TitleBar

	// Initialize main content viewport
	vp := viewport.New(0, 0)

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
	}
}

func (m Model) Init() tea.Cmd {
	if m.client != nil {
		return tea.Batch(
			spinner.Tick,
			m.loadNotes(),
		)
	}
	return spinner.Tick
}

func (m Model) SetClient(client *api.Client) tea.Cmd {
	m.client = client
	return m.loadNotes()
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)

	case notesLoadedMsg:
		if msg.Err != nil {
			return m, nil
		}
		items := make([]list.Item, len(msg.Notes))
		for i, note := range msg.Notes {
			items[i] = item{note: note}
		}
		m.list.SetItems(items)
		m.loading = false
		return m, nil

	case list.FilterMatchesMsg:
		m.list, cmd = m.list.Update(msg)
		return m, cmd

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
		listWidth := availableWidth / ViewRatio
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
		// If we're actively filtering, don't handle any other keypresses
		if m.list.FilterState() == list.Filtering {
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Left):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Right):
			if m.focusedPane == noteList {
				m.focusedPane = markdown
			}

		case key.Matches(msg, m.keys.ClearFilter):
			m.list.ResetFilter()
			return m, nil

		case key.Matches(msg, m.keys.NewNote):
			cmd = navigation.SwitchUICmd(navigation.CreateUI)
			return m, cmd

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
						// We don't have the content locally.. fetch
						m.loading = true
						// Fetch the note content
						return m, fetchNoteContent(m.client, note.NoteID)
					} else {
						// We have the content, render..
						rendered, err := m.renderer.Render(*note.Content)
						if err != nil {
							log.Error("Error rendering markdown: %v", err)
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

			// Update the item in the model's list
			currentIndex := m.list.Index()
			items := m.list.Items()
			items[currentIndex] = item{note: note}
			m.list.SetItems(items)

			rendered, err := m.renderer.Render(content)
			if err != nil {
				log.Errorf("Error rendering markdown: %v", err)
				m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
			} else {
				m.viewport.SetContent(rendered)
			}
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (i item) Title() string { return i.note.Title }
func (i item) Description() string {
	return fmt.Sprintf("Created: %s", i.note.CreatedAt.Format("2006-01-02"))
}
func (i item) FilterValue() string { return i.note.Title }

func (m Model) View() string {
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
		listStyle = m.styles.ActiveContent.Width(m.width / ViewRatio)
	} else {
		listStyle = m.styles.InactiveContent.Width(m.width / ViewRatio)
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
