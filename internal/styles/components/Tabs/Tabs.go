package Tabs

import (
	"merlion/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
)

const TabsHeight = 1

// Tabs represents a list with tabs
type Tabs struct {
	Tabs         []string
	ActiveTab    int
	width        int
	height       int
	ShowArrows   bool
	themeManager *styles.ThemeManager
}

// New creates a new TabbedList
func New(
	tabs []string,
	themeManager *styles.ThemeManager,
) Tabs {
	return Tabs{
		Tabs:         tabs,
		ActiveTab:    0,
		themeManager: themeManager,
		ShowArrows:   true,
	}
}

// SetSize sets the size of the tabbed list
func (t *Tabs) SetWidth(width int) {
	t.width = width
}

// Update handles the Bubble Tea update loop
func (t *Tabs) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.SetWidth(msg.Width)
		return nil
	case tea.KeyMsg:
		return t.HandleKeyMsg(msg)
	}

	var cmd tea.Cmd
	return cmd
}

func (t *Tabs) renderTabs() string {
	if t.width == 0 {
		return ""
	}

	style := t.themeManager.Styles()

	availableWidth := t.width - 4 // borders and padding
	var visibleTabs []string
	var startIdx, endIdx int
	currentWidth := 0

	// Center the active tab when possible
	startIdx = max(0, t.ActiveTab-2)
	for i := startIdx; i < len(t.Tabs); i++ {
		tabWidth := len(t.Tabs[i]) + 4 // Add padding
		if currentWidth+tabWidth > availableWidth {
			break
		}
		currentWidth += tabWidth
		endIdx = i + 1
	}

	visibleTabs = t.Tabs[startIdx:endIdx]

	// Build tab bar
	var tabBar string
	if t.ShowArrows && startIdx > 0 {
		tabBar += style.InactiveTab.Render("←")
	}

	for i, tab := range visibleTabs {
		actualIdx := startIdx + i
		if actualIdx == t.ActiveTab {
			tabBar += style.ActiveTab.Render(tab)
		} else {
			tabBar += style.InactiveTab.Render(tab)
		}
	}

	if t.ShowArrows && endIdx < len(t.Tabs) {
		tabBar += style.InactiveTab.Render("→")
	}

	return tabBar
}

// View renders the component
func (t *Tabs) View() string {
	if t.width == 0 {
		return ""
	}
	return t.renderTabs()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NextTab moves to the next tab if available
func (t *Tabs) NextTab() {
	if t.ActiveTab < len(t.Tabs)-1 {
		t.ActiveTab++
	} else {
		t.ActiveTab = 0
	}
}

// PrevTab moves to the previous tab if available
func (t *Tabs) PrevTab() {
	if t.ActiveTab > 0 {
		t.ActiveTab--
	} else {
		t.ActiveTab = len(t.Tabs) - 1
	}
}

// HandleKeyMsg handles keyboard messages
func (t *Tabs) HandleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	// TODO: Should be taken in input / Or handle by parent, but ok for now
	switch msg.String() {
	case "u":
		if t.ActiveTab > 0 {
			t.ActiveTab--
			return nil
		}
	case "i":
		if t.ActiveTab < len(t.Tabs)-1 {
			t.ActiveTab++
			return nil
		}
	}

	var cmd tea.Cmd
	return cmd
}
