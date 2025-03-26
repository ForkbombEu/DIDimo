package template_engine

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/group/all"
)

type PlaceholderMetadata struct {
	Field         string
	Descriptions  string 
	CredimiID     string 
	Type          string
}

func parseMetadata(metaStr string) (string, string, string) {
	descriptions := ""
	credimiID := ""
	placeholderType := ""

	parts := strings.Split(metaStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, "description:") {
			descriptions = strings.TrimSpace(strings.TrimPrefix(part, "description:"))
		}

		if strings.HasPrefix(part, "credimi_id:") {
			credimiID = strings.TrimSpace(strings.TrimPrefix(part, "credimi_id:"))
			continue
		}

		if strings.HasPrefix(part, "type:") {
			placeholderType = strings.TrimSpace(strings.TrimPrefix(part, "type:"))
			continue
		}	
	}

	return descriptions, credimiID, placeholderType
}

func RemoveNewlinesAndBackslashes(input string) string {
	output := strings.ReplaceAll(input, "\n", "")
	output = strings.ReplaceAll(output, "\\", "")
	output = strings.ReplaceAll(output, "\"", "'")
	return output
}


func ExtractPlaceholders(content string, normalized ...bool) []PlaceholderMetadata {
	norm := true
	if len(normalized) > 0 {
		norm = normalized[0]
	}

	regex := regexp.MustCompile(`{{\s*\.([a-zA-Z0-9_.]+)(?:\s*§\s*([^§]+?)\s*§)?\s*}}`)
	matches := regex.FindAllStringSubmatch(content, -1)

	var placeholders []PlaceholderMetadata
	unique := make(map[string]bool)

	for _, match := range matches {
		field := match[1]
		descriptions := ""
		credimiID := ""
		placeholderType := ""

		if len(match) >= 3 && match[2] != "" {
			descriptions, credimiID, placeholderType = parseMetadata(match[2])
		}

		if norm {
			if !unique[field] {
				unique[field] = true
				placeholders = append(placeholders, PlaceholderMetadata{
					Field:         field,
					Descriptions:  descriptions,
					CredimiID:     credimiID,
					Type:          placeholderType,
				})
			}
		} else {
			placeholders = append(placeholders, PlaceholderMetadata{
				Field:         field,
				Descriptions:  descriptions,
				CredimiID:     credimiID,
				Type:          placeholderType,
			})
		}
	}

	return placeholders
}

func stripMetadataFromTemplate(content string) string {
	regex := regexp.MustCompile(`({{\s*\.[a-zA-Z0-9_.]+)\s*§\s*[^§]+\s*§\s*}}`)
	return regex.ReplaceAllString(content, "$1}}")
}

func RenderTemplate(reader io.Reader, data map[string]interface{}) (string, error) {
	handler := sprout.New(
		sprout.WithGroups(all.RegistryGroup()),
	)
	funcs := handler.Build()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return "", err
	}

	templateContent := buf.String()
	cleanContent := stripMetadataFromTemplate(templateContent)

	tmpl, err := template.New("tmpl").Funcs(funcs).Parse(cleanContent)
	if err != nil {
		return "", err
	}

	buf.Reset()
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetPlaceholders(readers []io.Reader, normalized ...bool) ([]PlaceholderMetadata, error) {
	norm := true
	if len(normalized) > 0 {
		norm = normalized[0]
	}
	var allPlaceholders []PlaceholderMetadata
	unique := make(map[string]bool)

	for _, r := range readers {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return nil, err
		}
		content := buf.String()
		placeholders := ExtractPlaceholders(content, norm)
		for _, ph := range placeholders {
			if norm {
				if !unique[ph.Field] {
					unique[ph.Field] = true
					allPlaceholders = append(allPlaceholders, ph)
				}
			} else {
				allPlaceholders = append(allPlaceholders, ph)
			}
		}
	}

	return allPlaceholders, nil
}

