package files

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	keyTags       = "tags"
	keyIsFavorite = "favorite"
	keyIsWorkLog  = "worklog"
	keyCreatedAt  = "createdAt"
	keyUpdatedAt  = "updatedAt"
	keyWorkspace  = "workspaceId"
)

// splitFrontMatterContent splits an Obsidian file into front matter and content
func splitFrontMatterContent(content string) (map[string]any, string, error) {
	lines := strings.Split(content, "\n")

	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		return make(map[string]any), content, nil
	}

	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIndex = i
			break
		}
	}

	// If no closing delimiter found, treat as no front matter
	if endIndex == -1 {
		return make(map[string]any), content, nil
	}

	frontMatterText := strings.Join(lines[1:endIndex], "\n")
	var frontMatter map[string]any

	if err := yaml.Unmarshal([]byte(frontMatterText), &frontMatter); err != nil {
		return nil, "", fmt.Errorf("failed to parse YAML front matter: %w", err)
	}

	// Get the content after front matter
	var noteContent string
	if endIndex+1 < len(lines) {
		noteContent = strings.Join(lines[endIndex+1:], "\n")
		noteContent = strings.TrimPrefix(noteContent, "\n")
	}

	return frontMatter, noteContent, nil
}

// Helper functions for extracting typed values from front matter
func frontMatterGetString(m map[string]any, key string, defaultVal string) string {
	if m == nil {
		return defaultVal
	}
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultVal
}

// frontMatterGetInt returns the value of the key as an int, or the default value if not found
func frontMatterGetBool(m map[string]any, key string, defaultVal bool) bool {
	if m == nil {
		return defaultVal
	}
	if val, ok := m[key].(bool); ok {
		return val
	}
	return defaultVal
}

// frontMatterGetList returns the value of the key as a slice of strings, or the default value if not found
func frontMatterGetList(m map[string]any, key string) []string {
	if m == nil {
		return nil
	}
	if val, ok := m[key].([]any); ok {
		result := make([]string, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// frontMatterGetTime returns the value of the key as a time.Time, or the default value if not found
func frontMatterGetTime(m map[string]any, key string, defaultVal time.Time) time.Time {
	if m == nil {
		return defaultVal
	}

	// Try different time formats
	if val, ok := m[key].(string); ok {
		// Try RFC3339 format first
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t
		}
		// Try common date formats
		formats := []string{
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t
			}
		}
	}

	// Handle if it's already a time.Time
	if val, ok := m[key].(time.Time); ok {
		return val
	}

	return defaultVal
}
