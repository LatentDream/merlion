package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	fmt.Println("Let's play!")
	fileName := "./demo.md"
	source, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create a goldmark markdown parser
	md := goldmark.New()

	// Parse the markdown to get the AST
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	// Print the AST
	fmt.Println("AST Structure:")
	fmt.Println(strings.Repeat("-", 50))
	printAST(doc, source, 0)
	fmt.Println(strings.Repeat("-", 50))

	// Convert to HTML (like your original code)
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	fmt.Println("Final HTML output:")
	fmt.Println(buf.String())

	// If you want to render the markdown back to markdown format
	// you would need a custom renderer, but goldmark doesn't have a built-in markdown renderer
}

// Function to recursively print the AST
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
	}

	fmt.Println()

	// Recursively print children
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		printAST(child, source, depth+1)
	}
}
