package extension

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type WikiLink struct {
	ast.BaseInline
	Title string
}

func (n *WikiLink) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindWikiLink = ast.NewNodeKind("WikiLink")

func (n *WikiLink) Kind() ast.NodeKind { return KindWikiLink }

// Custom parser for wiki links [[Title Link]]
type wikiLinkParser struct{}

var WikiLinkParser = &wikiLinkParser{}

func (p *wikiLinkParser) Trigger() []byte {
	return []byte{'['}
}

func (p *wikiLinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	nbOpenBracket := 2
	nbCloseBracket := 2
	nbOfBracket := nbOpenBracket + nbCloseBracket
	if len(line) < nbOfBracket {
		return nil
	}

	if line[0] != '[' || line[1] != '[' {
		return nil
	}

	// Find closing ]]
	var titleEnd int
	found := false
	for i := 2; i < len(line)-1; i++ {
		if line[i] == ']' && line[i+1] == ']' {
			titleEnd = i
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	title := string(line[nbOpenBracket:titleEnd])
	block.Advance(titleEnd + nbCloseBracket)
	return &WikiLink{
		Title: title,
	}
}

// Create a renderer for wiki links
type wikiLinkRenderer struct{}

func WikiLinkHtmlRenderer() renderer.NodeRenderer {
	return &wikiLinkRenderer{}
}

func (r *wikiLinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikiLink, r.renderWikiLink)
}

func (r *wikiLinkRenderer) renderWikiLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*WikiLink)
	// Render as a link to wiki page
	_, _ = w.WriteString(fmt.Sprintf(`<a href="/wiki/%s" class="wiki-link">%s</a>`, strings.ReplaceAll(n.Title, " ", "_"), n.Title))
	return ast.WalkContinue, nil
}
