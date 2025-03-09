package Notes

import (
	"fmt"
	"sort"

	"merlion/internal/api"
	"merlion/internal/styles"
	"merlion/internal/styles/components/Tabs"
	styledDelegate "merlion/internal/styles/components/delegate"
	"merlion/internal/ui/create"
	"merlion/internal/ui/navigation"
	"merlion/internal/ui/notes/renderer"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	noteRenderer renderer.Model
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
	l.Styles.Title = s.Title

	filterTabs := []TabKind{AllNotes, Favorites, WorkLogs}
	tabs := Tabs.New(filterTabs, themeManager)

	noteRenderer := renderer.New(themeManager)

	// Initialize help viewport with themed styles
	return Model{
		noteList:     l,
		fileterTabs:  tabs,
		noteRenderer: noteRenderer,
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
		sort.Slice(filteredNotes, func(i, j int) bool {
			return filteredNotes[j].CreatedAt.Before(filteredNotes[i].CreatedAt)
		})
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
			m.noteRenderer.SetErrorMessage(fmt.Sprintf("Error editing note: %v", msg.err))
			return m, nil
		}

		// Refresh the viewport content after successful edit
		if i := m.noteList.SelectedItem(); i != nil {
			// NOTE: Content should be edited on the master list only
			// -> and refresh the list after
			note := i.(item).note
			m.noteRenderer.SetNote(&note)
			m.noteRenderer.Render()
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
			m.noteRenderer.SetSize(contentWidth, m.height)

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
			m.noteRenderer.SetSize(contentWidth, m.height)
		}
		// Update content if selected
		m.noteRenderer.Render()

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
			log.Info("Next called")
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

		case key.Matches(msg, m.keys.ToggleInfo):
			m.noteRenderer.ToggleHideInfo()

		case key.Matches(msg, m.keys.ToggleInfoPosition):
			log.Info("Toggle called")
			m.noteRenderer.ToggleHidePosition()

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
			var noteToDelete api.Note
			if m.focusedPane == noteList {
				if i := m.noteList.SelectedItem(); i != nil {
					noteToDelete = i.(item).note
				}
			} else {
				if m.noteRenderer.Note != nil {
					noteToDelete = *m.noteRenderer.Note
				}
			}
			if noteToDelete.NoteID != "" {
				cmd = navigation.AskConfirmationCmd(
					"Are you sure you want to delete this note ?",
					noteToDelete.Title,
					navigation.DangerLvl,
					func() {
						m.client.DeleteNote(noteToDelete.NoteID)
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
						return m, fetchNoteContent(m.client, note.NoteID)
					}
					m.noteRenderer.SetNote(&note)
					m.noteRenderer.Render()
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
			m.noteRenderer, cmd = m.noteRenderer.Update(msg)
			cmds = append(cmds, cmd)
		}

	case noteContentMsg:
		if i := m.noteList.SelectedItem(); i != nil {
			note := i.(item).note
			if note.NoteID != msg.NoteId {
				log.Fatalf("Receive a Content of a un-selected note")
			}
			content := msg.Content
			note.Content = &content
			m.noteRenderer.SetNote(&note)

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
			m.focusedPane = markdown
			m.noteRenderer.Render()
		}
		m.loading = false
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

	combinedView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.fileterTabs.View(),
		m.noteList.View(),
	)

	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		listStyle.Render(combinedView),
	)
	rightSide := rendererStyle.Render(
		m.noteRenderer.View(),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSide,
		rightSide,
	)

}

func (m Model) mobileView() string {
	var style lipgloss.Style
	style = m.styles.MobileContent.
		Width(m.width)

	if m.focusedPane == markdown {
		return style.Render(m.noteRenderer.View())
	} else {
		combinedView := lipgloss.JoinVertical(
			lipgloss.Left,
			m.fileterTabs.View(),
			m.noteList.View(),
		)
		return style.Render(combinedView)
	}
}

func (m Model) View() string {

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
				"Loading...",
			),
		)
	}

	if m.viewType == large {
		return m.desktopView()
	} else {
		return m.mobileView()
	}
}
