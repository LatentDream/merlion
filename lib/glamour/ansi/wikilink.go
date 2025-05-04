package ansi

import (
	"io"
)

// A WikiLinkElement is used to render hyperlinks.
type WikiLinkElement struct {
	Token      string
	isSelected bool
}

// Render renders a WikiLinkElement.
func (e *WikiLinkElement) Render(w io.Writer, ctx RenderContext) error {
	if err := e.renderTextPart(w, ctx); err != nil {
		return err
	}
	return nil
}

func (e *WikiLinkElement) renderTextPart(w io.Writer, ctx RenderContext) error {
	style := ctx.options.Styles.WikiLink
	if e.isSelected {
		style = ctx.options.Styles.WikiLinkHighlighted
	}
	if len(e.Token) > 0 {
		el := &BaseElement{
			Token: e.Token,
			// Suffix: " →",
			Style: style,
		}
		err := el.Render(w, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
