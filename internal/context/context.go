package context

import (
	"merlion/internal/styles"
)

type Context struct {
	ThemeManager   *styles.ThemeManager
	DefaultToCloud bool
}

type ContextOption func(*Context)

// WithSaveOnChange: Allows to remove the auto onsave when changing the config element
func WithSaveOnChange(saveOnChange bool) ContextOption {
	return func(c *Context) {
		c.ThemeManager.SetSaveOnChange(saveOnChange)
	}
}

// WithCompactViewStart: if the user only has access to the compact view
func WithCompactViewStart(isStartingInCompactView bool) ContextOption {
	return func(c *Context) {
		if isStartingInCompactView {
			c.ThemeManager.Config.CompactView = true
		}
	}
}

// WithLocalFirst: if the user should start with local notes as default view
func WithLocalFirst(isLocalFirst bool) ContextOption {
	return func(c *Context) {
		c.DefaultToCloud = !isLocalFirst
	}
}

func NewContext(configDir string, options ...ContextOption) (*Context, error) {
	tm, err := styles.NewThemeManager(configDir)
	if err != nil {
		return nil, err
	}

	context := &Context{
		ThemeManager:   tm,
		DefaultToCloud: tm.Config.DefaultToCloud,
	}

	// Apply all options
	for _, option := range options {
		option(context)
	}

	return context, nil
}
