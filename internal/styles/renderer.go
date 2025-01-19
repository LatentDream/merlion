package styles

import "github.com/charmbracelet/glamour/ansi"

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func uintPtr(u uint) *uint {
	return &u
}

func (tm *ThemeManager) GetRendererStyle() ansi.StyleConfig {
	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:      stringPtr(string(tm.theme.Foreground)),
			},
			Margin: uintPtr(tm.theme.Margin),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         uintPtr(1),
			IndentToken:    stringPtr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: tm.theme.ListIndent,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:      stringPtr(string(tm.theme.Primary)),
				Bold:       boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           stringPtr(string(tm.theme.Primary)),
				BackgroundColor: stringPtr(string(tm.theme.Secondary)),
				Bold:            boolPtr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  stringPtr(string(tm.theme.MutedColor)),
				Bold:   boolPtr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: boolPtr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: boolPtr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: boolPtr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  stringPtr(string(tm.theme.MutedColor)),
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{},
			Ticked:         "[✓] ",
			Unticked:       "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     stringPtr(string(tm.theme.Primary)),
			Underline: boolPtr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: stringPtr(string(tm.theme.Secondary)),
			Bold:  boolPtr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     stringPtr(string(tm.theme.Tertiary)),
			Underline: boolPtr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  stringPtr(string(tm.theme.MutedColor)),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           stringPtr(string(tm.theme.Primary)),
				BackgroundColor: stringPtr(string(tm.theme.Background)),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Foreground)),
				},
				Margin: uintPtr(tm.theme.Margin),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Foreground)),
				},
				Error: ansi.StylePrimitive{
					Color:           stringPtr(string(tm.theme.Foreground)),
					BackgroundColor: stringPtr(string(tm.theme.Error)),
				},
				Comment: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Comment)),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Secondary)),
				},
				Keyword: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Primary)),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Secondary)),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Tertiary)),
				},
				KeywordType: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Primary)),
				},
				Operator: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Secondary)),
				},
				Punctuation: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.MutedColor)),
				},
				Name: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Foreground)),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Secondary)),
				},
				NameTag: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Primary)),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Secondary)),
				},
				NameClass: ansi.StylePrimitive{
					Color:     stringPtr(string(tm.theme.Primary)),
					Underline: boolPtr(true),
					Bold:      boolPtr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Tertiary)),
				},
				NameFunction: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Success)),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Success)),
				},
				LiteralString: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Secondary)),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Primary)),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Error)),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: boolPtr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.Success)),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: boolPtr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: stringPtr(string(tm.theme.MutedColor)),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: stringPtr(string(tm.theme.Background)),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}
}
