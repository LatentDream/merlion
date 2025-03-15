package TagList

import (
	"merlion/internal/api"
	"merlion/internal/styles"

	// tea "github.com/charmbracelet/bubbletea"
)

type Tag struct {
	name  string
	notes []*api.Note
}

type TagList struct {
	// Logic
	Tags            []Tag
	opennedTag      *int
	highlightedTag  int
	highlightedNote int

	// Styling
	Width        int
	Height       int
	themeManager *styles.ThemeManager
}

func New(tags []Tag, tm *styles.ThemeManager) TagList {

	return TagList{
		Tags:         tags,
		opennedTag:   nil,
		themeManager: tm,
	}

}

func (t *TagList) View() string {
	// List all the tags

	// If the tags is open, list the element from it

	return ""
}
