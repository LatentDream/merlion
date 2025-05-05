package ansi

import (
	"io"
)

// A WikiLinkElement is used to render hyperlinks.
type WikiLinkElement struct {
	Token string
}

// Render renders a WikiLinkElement.
func (e *WikiLinkElement) Render(w io.Writer, ctx RenderContext) error {
	isSelected := ctx.Selector.isSelected(&Selector{Title: e.Token})
	if err := e.renderWikiLink(w, ctx, isSelected); err != nil {
		return err
	}
	return nil
}

func (e *WikiLinkElement) renderWikiLink(w io.Writer, ctx RenderContext, isSelected bool) error {
	var style StylePrimitive
	if isSelected {
		style = ctx.options.Styles.WikiLinkHighlighted
	} else {
		style = ctx.options.Styles.WikiLink
	}
	if len(e.Token) > 0 {
		el := &BaseElement{
			Token: e.Token,
			Style: style,
		}
		err := el.Render(w, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
