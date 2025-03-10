package Tabs

import (
	"merlion/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
)

const TabsHeight = 1

type Displayable interface {
	String() string
}

// Tabs represents a list with tabs
type Tabs[T Displayable] struct {
	Tabs         []T
	ActiveTab    int
	width        int
	height       int
	ShowArrows   bool
	themeManager *styles.ThemeManager
}

func (t *Tabs[T]) CurrentTab() T {
	return t.Tabs[t.ActiveTab]
}

// New creates a new TabbedList
func New[T Displayable](
	tabs []T,
	themeManager *styles.ThemeManager,
) Tabs[T] {
	return Tabs[T]{
		Tabs:         tabs,
		ActiveTab:    0,
		themeManager: themeManager,
		ShowArrows:   true,
	}
}

// SetSize sets the size of the tabbed list
func (t *Tabs[T]) SetWidth(width int) {
	t.width = width
}

// Update handles the Bubble Tea update loop
func (t *Tabs[T]) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.SetWidth(msg.Width)
		return nil
	}

	var cmd tea.Cmd
	return cmd
}

func (t *Tabs[T]) renderTabs() string {
	if t.width == 0 {
		return ""
	}
	style := t.themeManager.Styles()
	availableWidth := t.width - 4 // borders and padding
	var visibleTabs []T
	var startIdx, endIdx int
	currentWidth := 0

	// Start with active tab centered
	startIdx = max(0, t.ActiveTab-2)

	// First pass: calculate initial visible range
	for i := startIdx; i < len(t.Tabs); i++ {
		tabWidth := len(t.Tabs[i].String()) + 4 // Add padding
		if currentWidth+tabWidth > availableWidth {
			break
		}
		currentWidth += tabWidth
		endIdx = i + 1
	}

	// Adjust if active tab would be hidden off the right
	if t.ActiveTab >= endIdx {
		hiddenTabs := t.ActiveTab - endIdx + 1
		startIdx += hiddenTabs
		currentWidth = 0
		endIdx = startIdx

		// Recalculate visible tabs from new start
		for i := startIdx; i < len(t.Tabs); i++ {
			tabWidth := len(t.Tabs[i].String()) + 4
			if currentWidth+tabWidth > availableWidth {
				break
			}
			currentWidth += tabWidth
			endIdx = i + 1
		}
	}

	// Calculate total width and center if possible
	totalWidth := 0
	for i := startIdx; i < endIdx; i++ {
		totalWidth += len(t.Tabs[i].String()) + 4
	}

	// Adjust start index to center tabs if possible
	if endIdx-startIdx > 3 {
		leftSpace := (availableWidth - totalWidth) / 2
		if leftSpace > 0 {
			// Try to shift tabs left to center them
			for i := startIdx; i > 0; i-- {
				tabWidth := len(t.Tabs[i-1].String()) + 4
				if leftSpace < tabWidth {
					break
				}
				startIdx = i - 1
				totalWidth += tabWidth
				leftSpace -= tabWidth
			}
		}
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
			tabBar += style.ActiveTab.Render(tab.String())
		} else {
			tabBar += style.InactiveTab.Render(tab.String())
		}
	}

	if t.ShowArrows && endIdx < len(t.Tabs) {
		tabBar += style.InactiveTab.Render("→")
	}

	return tabBar
}

// View renders the component
func (t *Tabs[T]) View() string {
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
func (t *Tabs[T]) NextTab() T {
	if t.ActiveTab < len(t.Tabs)-1 {
		t.ActiveTab++
	} else {
		t.ActiveTab = 0
	}
	return t.Tabs[t.ActiveTab]
}

// PrevTab moves to the previous tab if available
func (t *Tabs[T]) PrevTab() T {
	if t.ActiveTab > 0 {
		t.ActiveTab--
	} else {
		t.ActiveTab = len(t.Tabs) - 1
	}
	return t.Tabs[t.ActiveTab]
}
