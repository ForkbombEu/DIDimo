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
	Translations  map[string]string 
	Descriptions  map[string]string 
	CredimiID     string 
	Type          string
}

func parseMetadata(metaStr string) (map[string]string, map[string]string, string, string) {
	translations := make(map[string]string)
	descriptions := make(map[string]string)
	credimiID := ""
	inDescription := false
	placeholderType := ""

	parts := strings.Split(metaStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, "description:") {
			inDescription = true
			part = strings.TrimPrefix(part, "description:")
			part = strings.TrimSpace(part)
		}

		if strings.HasPrefix(part, "credimi_id:") {
			credimiID = strings.TrimSpace(strings.TrimPrefix(part, "credimi_id:"))
			continue
		}

		if strings.HasPrefix(part, "type:") {
			placeholderType = strings.TrimSpace(strings.TrimPrefix(part, "type:"))
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if inDescription {
				descriptions[key] = value
			} else {
				translations[key] = value
			}
		}
	}

	return translations, descriptions, credimiID, placeholderType
}

func ExtractPlaceholders(content string, normalized ...bool) []PlaceholderMetadata {
	norm := true
	if len(normalized) > 0 {
		norm = normalized[0]
	}

	regex := regexp.MustCompile(`{{\s*\.([a-zA-Z0-9_.]+)(?:\s*\|\s*([^}]+))?\s*}}`)
	matches := regex.FindAllStringSubmatch(content, -1)

	var placeholders []PlaceholderMetadata
	unique := make(map[string]bool)

	for _, match := range matches {
		field := match[1]
		translations := make(map[string]string)
		descriptions := make(map[string]string)
		credimiID := ""
		placeholderType := ""

		if len(match) >= 3 && match[2] != "" {
			translations, descriptions, credimiID, placeholderType = parseMetadata(match[2])
		}

		if norm {
			if !unique[field] {
				unique[field] = true
				placeholders = append(placeholders, PlaceholderMetadata{
					Field:         field,
					Translations:  translations,
					Descriptions:  descriptions,
					CredimiID:     credimiID,
					Type:          placeholderType,
				})
			}
		} else {
			placeholders = append(placeholders, PlaceholderMetadata{
				Field:         field,
				Translations:  translations,
				Descriptions:  descriptions,
				CredimiID:     credimiID,
				Type:          placeholderType,
			})
		}
	}

	return placeholders
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

	tmpl, err := template.New("tmpl").Funcs(funcs).Parse(buf.String())
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