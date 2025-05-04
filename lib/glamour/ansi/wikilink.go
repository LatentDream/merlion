package ansi

import (
	"io"
	"strings"
)

// A WikiLinkElement is used to render hyperlinks.
type WikiLinkElement struct {
	Token string
}

// Render renders a WikiLinkElement.
func (e *WikiLinkElement) Render(w io.Writer, ctx RenderContext) error {
	if err := e.renderTextPart(w, ctx); err != nil {
		return err
	}
	return nil
}

func (e *WikiLinkElement) renderTextPart(w io.Writer, ctx RenderContext) error {
	style := ctx.options.Styles.ImageText
	style.Format = strings.TrimSuffix(style.Format, " →")

	if len(e.Token) > 0 {
		el := &BaseElement{
			Token:  e.Token,
			// Suffix: " →",
			Style:  ctx.options.Styles.Link,
		}
		err := el.Render(w, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
