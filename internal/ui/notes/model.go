package Notes

import (
	"fmt"
	"sort"
	"strings"

	"merlion/internal/controls"
	"merlion/internal/model"
	"merlion/internal/vault"
	"merlion/internal/vault/cloud"
	"merlion/internal/styles"
	styledDelegate "merlion/internal/styles/components/delegate"
	grouplist "merlion/internal/styles/components/groupList"
	tabs "merlion/internal/styles/components/tabs"
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
	note model.Note
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
	Tags
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
	case Tags:
		return "Tags"
	default:
		return "Unknown"
	}
}

type Model struct {
	noteList     list.Model
	fileterTabs  tabs.Tabs[TabKind]
	noteRenderer renderer.Model
	spinner      spinner.Model
	keys         controls.KeyMap
	focusedPane  focusedPanel
	width        int
	height       int
	ready        bool
	loading      bool
	listDelegate *styledDelegate.StyledDelegate
	styles       *styles.Styles
	themeManager *styles.ThemeManager
	storeManager *vault.Manager
	viewType     ViewType
	compactView  bool
	tagsList     grouplist.Model
}

func NewModel(storeManager *vault.Manager, themeManager *styles.ThemeManager, firstTab string) Model {
	s := themeManager.Styles()

	// Initialize spinner with themed color
	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
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
	l.SetShowHelp(false)
	l.Styles.Title = s.Title

	filterTabs := []TabKind{AllNotes, Favorites, WorkLogs, Tags}

	tabs := tabs.New(filterTabs, themeManager, firstTab)

	noteRenderer := renderer.New(themeManager, storeManager)

	gl := grouplist.New([]grouplist.Group{}, delegate, themeManager)

	// Initialize help viewport with themed styles
	return Model{
		noteList:     l,
		fileterTabs:  tabs,
		noteRenderer: noteRenderer,
		spinner:      sp,
		keys:         controls.Keys,
		focusedPane:  noteList,
		loading:      true,
		listDelegate: delegate,
		styles:       s,
		themeManager: themeManager,
		storeManager: storeManager,
		viewType:     large,
		compactView:  themeManager.Config.CompactView,
		tagsList:     gl,
	}
}

func (m Model) Init(args ...any) tea.Cmd {
	if m.storeManager != nil {
		return tea.Batch(
			m.spinner.Tick,
			m.loadNotes(),
		)
	}
	return m.spinner.Tick
}

func (m Model) SetCloudClient(client *cloud.Client) navigation.View {
	m.storeManager.UpdateCloudClient(client)
	return m
}

func createNoteItems(notes []model.Note, filter TabKind) []list.Item {
	filteredNotes := make([]model.Note, 0)
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

func (m *Model) refreshNotesView() {
	currTab := m.fileterTabs.CurrentTab()
	items := createNoteItems(m.storeManager.Notes, currTab)
	m.noteList.SetItems(items)
	groups := createTagGroups(m.storeManager.Notes)
	m.tagsList.SetGroups(groups)
}

func (m Model) getCurrentNote(considerRenderer bool) *model.Note {
	var currentNote *model.Note
	if m.focusedPane == markdown && considerRenderer {
		if m.noteRenderer.Note != nil {
			log.Debug("Selected Rendered note")
			currentNote = m.noteRenderer.Note
		}
	} else {
		if m.fileterTabs.CurrentTab() == Tags {
			selectedItem := m.tagsList.SelectedItem()
			if selectedItem != nil {
				if noteItem, ok := selectedItem.(item); ok {
					note := noteItem.note
					log.Debug("Selected tag note")
					currentNote = &note
				}
			}
		} else {
			if selectedItem := m.noteList.SelectedItem(); selectedItem != nil {
				if noteItem, ok := selectedItem.(item); ok {
					note := noteItem.note
					log.Debug("Selected list note")
					currentNote = &note
				}
			}
		}
	}
	return currentNote
}

func createTagGroups(notes []model.Note) []grouplist.Group {
	groups := make(map[string][]list.Item)

	for _, note := range notes {
		for _, tag := range note.Tags {
			tag = strings.ToLower(tag)
			newItem := item{note: note}

			if items, exists := groups[tag]; exists {
				groups[tag] = append(items, newItem)
			} else {
				groups[tag] = []list.Item{newItem}
			}
		}
	}

	result := make([]grouplist.Group, 0, len(groups))
	for name, items := range groups {
		result = append(result, grouplist.Group{
			Name:  name,
			Items: items,
		})
	}

	return result
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
		m.refreshNotesView()
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
		m.refreshNotesView()
		if note := m.getCurrentNote(true); note != nil {
			m.noteRenderer.SetNote(note)
			m.noteRenderer.Render()
		}
		return m, nil

	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}
		m.width = msg.Width - 4
		m.height = msg.Height - 2

		if msg.Width >= LargeScreenBreakpoint && !m.compactView {
			m.viewType = large

			// Account for padding and borders in the style
			horizontalPadding := m.styles.ActiveContent.GetHorizontalPadding() +
				m.styles.ActiveContent.GetHorizontalBorderSize()
			availableWidth := m.width - horizontalPadding

			// Split the view with adjusted measurements
			listWidth := availableWidth / ViewRatio
			listHeight := m.height - tabs.TabsHeight

			m.noteList.SetWidth(listWidth)
			m.noteList.SetHeight(listHeight)
			m.tagsList.SetWidth(listWidth)
			m.tagsList.SetHeight(listHeight)
			m.fileterTabs.Width = listWidth

			contentWidth := availableWidth - listWidth
			m.noteRenderer.SetSize(contentWidth, m.height)

		} else {
			m.viewType = small

			// Account for padding and borders in the style
			horizontalPadding := m.styles.MobileContent.GetHorizontalPadding() +
				m.styles.MobileContent.GetHorizontalBorderSize()
			availableWidth := m.width - horizontalPadding

			listWidth := availableWidth
			listHeight := m.height - tabs.TabsHeight
			m.noteList.SetWidth(listWidth)
			m.noteList.SetHeight(listHeight)
			m.tagsList.SetWidth(listWidth)
			m.tagsList.SetHeight(listHeight)
			m.fileterTabs.Width = listWidth

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
			if m.focusedPane == noteList {
				m.fileterTabs.NextTab()
				m.refreshNotesView()
			}

		case key.Matches(msg, m.keys.PrevTab):
			if m.focusedPane == noteList {
				m.fileterTabs.PrevTab()
				m.refreshNotesView()
			}

		case key.Matches(msg, m.keys.ToggleInfo):
			m.noteRenderer.ToggleHideInfo()

		case key.Matches(msg, m.keys.ToggleInfoPosition):
			m.noteRenderer.ToggleHidePosition()

		case key.Matches(msg, m.keys.ToggleCompactView):
			m.ToggleFullscreen()
			return m, tea.WindowSize()

		case key.Matches(msg, m.keys.ClearFilter):
			m.noteList.ResetFilter()
			if m.focusedPane == markdown {
				m.focusedPane = noteList
			}
			return m, nil

		case key.Matches(msg, m.keys.Create):
			cmd = navigation.SwitchUICmd(navigation.CreateUI, []any{})
			return m, cmd

		case key.Matches(msg, m.keys.Delete):
			noteToDelete := m.getCurrentNote(true)
			if noteToDelete != nil && noteToDelete.NoteID != "" {
				cmd = navigation.AskConfirmationCmd(
					"Are you sure you want to delete this note ?",
					noteToDelete.Title,
					navigation.DangerLvl,
					func() {
						m.storeManager.DeleteNote(noteToDelete.NoteID)
					},
					navigation.NoteUI,
				)
				return m, cmd
			}

		case key.Matches(msg, m.keys.Back):
			m.focusedPane = noteList

		case key.Matches(msg, m.keys.Edit):
			noteToEdit := m.getCurrentNote(true)
			if noteToEdit != nil {
				if noteToEdit.Content != nil {
					return m, m.openEditor()
				} else {
					// We don't have the content locally.. fetch
					m.loading = true
					return m, fetchNoteContent(m.storeManager, noteToEdit.NoteID)
				}
			}

		case key.Matches(msg, m.keys.Manage):
			noteToManage := m.getCurrentNote(true)
			if noteToManage != nil {
				m.loading = true
				return m, navigation.OpenManageViewCmd(noteToManage.NoteID)
			}

		case key.Matches(msg, m.keys.Select):
			if m.focusedPane == noteList {
				note := m.getCurrentNote(false)
				if note != nil {
					if note.Content == nil {
						// We don't have the content locally.. fetch
						m.loading = true
						return m, fetchNoteContent(m.storeManager, note.NoteID)
					}
					m.noteRenderer.SetNote(note)
					m.noteRenderer.Render()
					m.focusedPane = markdown
				}
			}

		case key.Matches(msg, m.keys.ToggleTheme):
			toggleTheme(&m)

		case key.Matches(msg, m.keys.ToggleStore):
			m.storeManager.NextStore()
			m.loading = true
			loadCmd := m.loadNotes()
			cmds = append(cmds, loadCmd)
			m.noteRenderer.SetNote(nil)
			m.noteRenderer.Render()
		}

		// Handle navigation based on focused pane
		if m.focusedPane == noteList && m.fileterTabs.CurrentTab() == Tags {
			m.tagsList, cmd = m.tagsList.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.focusedPane == noteList {
			m.noteList, cmd = m.noteList.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			// Handle markdown viewport scrolling
			m.noteRenderer, cmd = m.noteRenderer.Update(msg)
			cmds = append(cmds, cmd)
		}

	case noteContentMsg:
		currentNote := m.getCurrentNote(false)
		updatedNote := msg.UpdatedNote
		if currentNote != nil {
			if currentNote.NoteID != updatedNote.NoteID {
				log.Fatalf("Receive a Content of a un-selected note")
			}

			// Swap the item in the presentation list for the one with the content
			log.Debug("Swaping view note for updated note")
			currentIndex := m.noteList.Index()
			items := m.noteList.Items()
			items[currentIndex] = item{note: *updatedNote}
			m.noteList.SetItems(items)

			// Display the content in the renderer
			log.Debug("Rendering updated note")
			m.noteRenderer.SetNote(updatedNote)
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

func (m *Model) ToggleFullscreen() {
	newConf := !m.compactView
	m.compactView = newConf
	m.themeManager.SetCompactViewOnly(newConf)
}

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
	if m.fileterTabs.CurrentTab() == Tags {
		listView = m.tagsList.View()
	} else {
		listView = m.noteList.View()
	}

	combinedView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.fileterTabs.View(),
		listView,
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
		var listView string
		if m.fileterTabs.CurrentTab() == Tags {
			listView = m.tagsList.View()
		} else {
			listView = m.noteList.View()
		}

		combinedView := lipgloss.JoinVertical(
			lipgloss.Left,
			m.fileterTabs.View(),
			listView,
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
				" Loading...",
			),
		)
	}

	if m.viewType == large {
		return m.desktopView()
	} else {
		return m.mobileView()
	}
}
