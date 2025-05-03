package styles

import "github.com/latentdream/merlion/lib/glamour/ansi"

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
				Color:       stringPtr(string(tm.Theme.Foreground)),
			},
			Margin: uintPtr(tm.Theme.Margin),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         uintPtr(1),
			IndentToken:    stringPtr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: tm.Theme.ListIndent,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       stringPtr(string(tm.Theme.Primary)),
				Bold:        boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           stringPtr(string(tm.Theme.Background)),
				BackgroundColor: stringPtr(string(tm.Theme.Secondary)),
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
				Color:  stringPtr(string(tm.Theme.MutedColor)),
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
			Color:  stringPtr(string(tm.Theme.MutedColor)),
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
			Color:     stringPtr(string(tm.Theme.Primary)),
			Underline: boolPtr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: stringPtr(string(tm.Theme.Secondary)),
			Bold:  boolPtr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     stringPtr(string(tm.Theme.Tertiary)),
			Underline: boolPtr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  stringPtr(string(tm.Theme.MutedColor)),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           stringPtr(string(tm.Theme.Primary)),
				BackgroundColor: stringPtr(string(tm.Theme.Background)),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Foreground)),
				},
				Margin: uintPtr(tm.Theme.Margin),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Foreground)),
				},
				Error: ansi.StylePrimitive{
					Color:           stringPtr(string(tm.Theme.Foreground)),
					BackgroundColor: stringPtr(string(tm.Theme.Error)),
				},
				Comment: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Comment)),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Secondary)),
				},
				Keyword: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Primary)),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Secondary)),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Tertiary)),
				},
				KeywordType: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Primary)),
				},
				Operator: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Secondary)),
				},
				Punctuation: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.MutedColor)),
				},
				Name: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Foreground)),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Secondary)),
				},
				NameTag: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Primary)),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Secondary)),
				},
				NameClass: ansi.StylePrimitive{
					Color:     stringPtr(string(tm.Theme.Primary)),
					Underline: boolPtr(true),
					Bold:      boolPtr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Tertiary)),
				},
				NameFunction: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Success)),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Success)),
				},
				LiteralString: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Secondary)),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Primary)),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Error)),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: boolPtr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.Success)),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: boolPtr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: stringPtr(string(tm.Theme.MutedColor)),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: stringPtr(string(tm.Theme.Background)),
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
