package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Mention Custom AST nodes
type Mention struct {
	ast.BaseInline
	Username string
}

type WikiLink struct {
	ast.BaseInline
	Title string
}

func (n *Mention) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

func (n *WikiLink) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindMention = ast.NewNodeKind("Mention")
var KindWikiLink = ast.NewNodeKind("WikiLink")

func (n *Mention) Kind() ast.NodeKind  { return KindMention }
func (n *WikiLink) Kind() ast.NodeKind { return KindWikiLink }

// Custom parser for mentions (@username)
type mentionParser struct{}

var defaultMentionParser = &mentionParser{}

func (p *mentionParser) Trigger() []byte {
	return []byte{'@'}
}

func (p *mentionParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) == 0 || line[0] != '@' {
		return nil
	}

	// Regex for @username
	re := regexp.MustCompile(`^@([a-zA-Z0-9_]+)`)
	match := re.FindSubmatch(line)
	if match == nil {
		return nil
	}

	username := string(match[1])
	block.Advance(len(match[0]))

	mention := &Mention{
		Username: username,
	}
	return mention
}

// Custom parser for wiki links ([[Title Link]])
type wikiLinkParser struct{}

var defaultWikiLinkParser = &wikiLinkParser{}

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

type customExtension struct{}

func (e *customExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(defaultMentionParser, 500),
			util.Prioritized(defaultWikiLinkParser, 101),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newMentionRenderer(), 500),
			util.Prioritized(newWikiLinkRenderer(), 101),
		),
	)
}

// Create a renderer for mentions
type mentionRenderer struct{}

func newMentionRenderer() renderer.NodeRenderer {
	return &mentionRenderer{}
}

func (r *mentionRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMention, r.renderMention)
}

func (r *mentionRenderer) renderMention(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*Mention)
	// Render as a link to user profile
	_, _ = w.WriteString(fmt.Sprintf(`<a href="/users/%s" class="mention">@%s</a>`, n.Username, n.Username))
	return ast.WalkContinue, nil
}

// Create a renderer for wiki links
type wikiLinkRenderer struct{}

func newWikiLinkRenderer() renderer.NodeRenderer {
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

func main() {
	fmt.Println("Let's play!")

	source, err := os.ReadFile("./wikilink.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a goldmark markdown parser with custom extension
	md := goldmark.New(
		goldmark.WithExtensions(
			&customExtension{},
		),
	)

	// Parse the markdown to get the AST
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	// Print the AST
	fmt.Println("AST Structure:")
	fmt.Println(strings.Repeat("-", 50))
	printAST(doc, source, 0)
	fmt.Println(strings.Repeat("-", 50))

	// Convert to HTML
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	fmt.Println("Final HTML output:")
	fmt.Println(buf.String())
}

// Updated printAST function to handle custom nodes
func printAST(node ast.Node, source []byte, depth int) {
	indent := strings.Repeat("  ", depth)

	// Print node type
	fmt.Printf("%s%T", indent, node)

	// Print additional info based on node type
	switch n := node.(type) {
	case *ast.Text:
		fmt.Printf(" [Text: %q]", n.Text(source))
	case *ast.String:
		fmt.Printf(" [String: %q]", n.Value)
	case *ast.CodeSpan:
		fmt.Printf(" [CodeSpan: %q]", n.Text(source))
	case *ast.Link:
		fmt.Printf(" [Link: %q -> %q]", n.Title, n.Destination)
	case *ast.Image:
		fmt.Printf(" [Image: %q -> %q]", n.Title, n.Destination)
	case *ast.Heading:
		fmt.Printf(" [Level: %d]", n.Level)
	case *Mention:
		fmt.Printf(" [Mention: @%s]", n.Username)
	case *WikiLink:
		fmt.Printf(" [WikiLink: [[%s]]]", n.Title)
	}

	fmt.Println()

	// Recursively print children
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		printAST(child, source, depth+1)
	}
}
