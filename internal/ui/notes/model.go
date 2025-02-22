package Notes

import (
	"fmt"

	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/styles/components/Tabs"
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

type ViewType int

const (
	large ViewType = iota
	small
)

const ViewRatio = 4
const LargeScreenBreakpoint = 140

type item struct {
	note api.Note
}

type ViewState string

const (
	MainView   ViewState = "main"
	CreateView ViewState = "create"
)

type TabKind int

const (
	AllNotes TabKind = iota
	Favorites
	WorkLogs
)

// String implements the Displayable interface
func (t TabKind) String() string {
	switch t {
	case AllNotes:
		return "All Notes"
	case Favorites:
		return "Favorites"
	case WorkLogs:
		return "Work Logs"
	default:
		return "Unknown"
	}
}

type Model struct {
	noteList     list.Model
	allNotes     []api.Note
	fileterTabs  Tabs.Tabs[TabKind]
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
	viewType     ViewType
}

func NewModel(client *api.Client, themeManager *styles.ThemeManager) Model {
	s := themeManager.Styles()

	// Initialize glamour for markdown rendering
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(themeManager.GetRendererStyle()),
		glamour.WithWordWrap(int(themeManager.Theme.WordWrap)),
	)
	if err != nil {
		log.Fatalf("failed to initialize markdown renderer: %v", err)
	}

	// Initialize spinner with themed color
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(themeManager.Current().Primary)

	// Initialize list with themed styles
	delegate := styledDelegate.New(themeManager)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Notes"
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.TitleBar

	filterTabs := []TabKind{AllNotes, Favorites, WorkLogs}
	tabs := Tabs.New(filterTabs, themeManager)

	// Initialize main content viewport
	vp := viewport.New(0, 0)

	// Initialize help viewport with themed styles
	return Model{
		noteList:     l,
		fileterTabs:  tabs,
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
		viewType:     large,
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

func createNoteItems(notes []api.Note, filter TabKind) []list.Item {
	filteredNotes := make([]api.Note, 0)
	if filter == Favorites {
		for _, note := range notes {
			if note.IsFavorite {
				filteredNotes = append(filteredNotes, note)
			}
		}
	} else if filter == WorkLogs {
		for _, note := range notes {
			if note.IsWorkLog {
				filteredNotes = append(filteredNotes, note)
			}
		}
	} else {
		filteredNotes = notes
	}
	items := make([]list.Item, len(filteredNotes))
	for i, note := range filteredNotes {
		items[i] = item{note: note}
	}
	return items
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
		m.allNotes = msg.Notes
		items := createNoteItems(msg.Notes, AllNotes)
		m.noteList.SetItems(items)

		m.loading = false
		return m, nil

	case list.FilterMatchesMsg:
		m.noteList, cmd = m.noteList.Update(msg)
		return m, cmd

	case editorFinishedMsg:
		if msg.err != nil {
			m.viewport.SetContent(fmt.Sprintf("Error editing note: %v", msg.err))
			return m, nil
		}

		// Refresh the viewport content after successful edit
		if i := m.noteList.SelectedItem(); i != nil {
			// NOTE: Content should be edited on the master list only
			// -> and refresh the list after
			note := i.(item).note
			rendered, err := m.renderer.Render(*note.Content)
			if err != nil {
				m.viewport.SetContent(fmt.Sprintf("Error rendering markdown: %v", err))
			} else {
				m.viewport.SetContent(rendered)
			}
			updated := false
			for i, n := range m.allNotes {
				if n.NoteID == note.NoteID {
					n.Content = note.Content
					m.allNotes[i] = n
					updated = true
					break
				}
			}
			if !updated {
				log.Fatalf("Master list didn't get updated after Editor Finish")
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}
		m.width = msg.Width - 4
		m.height = msg.Height - 2

		if msg.Width >= LargeScreenBreakpoint {
			m.viewType = large

			// Account for padding and borders in the style
			horizontalPadding := m.styles.ActiveContent.GetHorizontalPadding() +
				m.styles.ActiveContent.GetHorizontalBorderSize()
			availableWidth := m.width - horizontalPadding

			// Split the view with adjusted measurements
			listWidth := availableWidth / ViewRatio
			listHeight := m.height - Tabs.TabsHeight

			m.noteList.SetWidth(listWidth)
			m.noteList.SetHeight(listHeight)
			m.fileterTabs.SetWidth(listWidth)

			contentWidth := availableWidth - listWidth
			m.viewport.Width = contentWidth
			m.viewport.Height = m.height

		} else {
			m.viewType = small

			// Account for padding and borders in the style
			horizontalPadding := m.styles.MobileContent.GetHorizontalPadding() +
				m.styles.MobileContent.GetHorizontalBorderSize()
			availableWidth := m.width - horizontalPadding

			listWidth := availableWidth
			listHeight := m.height - Tabs.TabsHeight
			m.noteList.SetWidth(listWidth)
			m.noteList.SetHeight(listHeight)
			m.fileterTabs.SetWidth(listWidth)

			contentWidth := availableWidth
			m.viewport.Width = contentWidth
			m.viewport.Height = m.height
		}
		// Update content if selected
		if i := m.noteList.SelectedItem(); i != nil {
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
		if m.noteList.FilterState() == list.Filtering {
			m.noteList, cmd = m.noteList.Update(msg)
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

		case key.Matches(msg, m.keys.NextTab):
			if m.focusedPane == noteList {
				activeTabName := m.fileterTabs.NextTab()
				items := createNoteItems(m.allNotes, activeTabName)
				m.noteList.SetItems(items)
			}

		case key.Matches(msg, m.keys.PrevTab):
			if m.focusedPane == noteList {
				activeTabName := m.fileterTabs.PrevTab()
				items := createNoteItems(m.allNotes, activeTabName)
				m.noteList.SetItems(items)
			}

		case key.Matches(msg, m.keys.ClearFilter):
			m.noteList.ResetFilter()
			if m.focusedPane == markdown {
				m.focusedPane = noteList
			}
			return m, nil

		case key.Matches(msg, m.keys.Create):
			cmd = navigation.SwitchUICmd(navigation.CreateUI)
			return m, cmd

		case key.Matches(msg, m.keys.Delete):
			if i := m.noteList.SelectedItem(); i != nil {
				note := i.(item).note
				cmd = navigation.AskConfirmationCmd(
					"Are you sure you want to delete this note ?",
					note.Title,
					navigation.DangerLvl,
					func() {
						m.client.DeleteNote(note.NoteID)
					},
					navigation.NoteUI,
				)
				return m, cmd
			}

		case key.Matches(msg, m.keys.Back):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Edit):
			if i := m.noteList.SelectedItem(); i != nil {
				note := i.(item).note
				if note.Content != nil {
					return m, m.openEditor(*note.Content)
				}
			}

		case key.Matches(msg, m.keys.Select):
			if m.focusedPane == noteList {
				if i := m.noteList.SelectedItem(); i != nil {
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
							log.Errorf("Error rendering markdown: %v", err)
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
			m.noteList, cmd = m.noteList.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			// Handle markdown viewport scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case noteContentMsg:
		m.loading = false
		if i := m.noteList.SelectedItem(); i != nil {
			note := i.(item).note
			content := string(msg)
			note.Content = &content

			// Update the item in the model's list
			currentIndex := m.noteList.Index()
			items := m.noteList.Items()
			items[currentIndex] = item{note: note}
			m.noteList.SetItems(items)
			// NOTE: Content should be edited on the master list only
			// -> and refresh the list after
			updated := false
			for i, n := range m.allNotes {
				if n.NoteID == note.NoteID {
					n.Content = &content
					m.allNotes[i] = n
					updated = true
					break
				}
			}
			if !updated {
				log.Fatalf("Master list didn't get updated after downloading content")
			}

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
	if i.note.IsFavorite {
		return fmt.Sprintf("★ Created: %s", i.note.CreatedAt.Format("2006-01-02"))
	} else {
		return fmt.Sprintf("Created: %s", i.note.CreatedAt.Format("2006-01-02"))
	}
}
func (i item) FilterValue() string { return i.note.Title }

func (m Model) desktopView() string {
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
			Width(m.noteList.Width()).
			Height(m.noteList.Height()).
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
		combinedView := lipgloss.JoinVertical(
			lipgloss.Left,
			m.fileterTabs.View(),
			m.noteList.View(),
		)
		listView = listStyle.Render(combinedView)
	}

	// Combine left container(s)
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

func (m Model) mobileView() string {
	var style lipgloss.Style
	style = m.styles.MobileContent.
		Width(m.width)

	if m.focusedPane == markdown {
		return style.Render(m.viewport.View())
	} else {
		if m.loading {
			loadingStyle := lipgloss.NewStyle().
				Width(m.width).
				Height(m.noteList.Height()).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			return loadingStyle.Render(
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					m.spinner.View(),
					" Loading notes...",
				),
			)
		} else {
			combinedView := lipgloss.JoinVertical(
				lipgloss.Left,
				m.fileterTabs.View(),
				m.noteList.View(),
			)
			return style.Render(combinedView)
		}
	}
}

func (m Model) View() string {
	if m.viewType == large {
		return m.desktopView()
	} else {
		return m.mobileView()
	}
}
