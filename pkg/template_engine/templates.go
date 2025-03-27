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
	FieldName      string
	CredimiID      string
	LabelKey       string
	DescriptionKey string
	Type           string
	Example        *string
}

func PreprocessTemplate(content string) string {
	re := regexp.MustCompile(`\s*\.credimiPlaceHolder\(\s*'([^']+)'(?:,\s*'[^']*')*\s*\)\s*`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) > 1 {
			fieldName := submatches[1]
			return " ." + fieldName + " "
		}
		return match
	})
}

func ExtractPlaceholders(content string) []PlaceholderMetadata {
	placeholderRegex := regexp.MustCompile(`{{\s*\.credimiPlaceHolder\('([^']+)',\s*'([^']+)',\s*'([^']+)',\s*'([^']+)',\s*'([^']+)'(?:,\s*'([^']+)')?\)\s*}}`)

	matches := placeholderRegex.FindAllStringSubmatch(content, -1)

	var placeholders []PlaceholderMetadata

	for _, match := range matches {
		if len(match) >= 6 {
			metadata := PlaceholderMetadata{
				FieldName:      match[1],
				CredimiID:      match[2],
				LabelKey:       match[3],
				DescriptionKey: match[4],
				Type:           match[5],
			}
			if len(match) > 6 && match[6] != "" {
				metadata.Example = &match[6]
			}
			placeholders = append(placeholders, metadata)
		}
	}
	return placeholders
}

func RemoveNewlinesAndBackslashes(input string) string {
	output := strings.ReplaceAll(input, "\n", "")
	output = strings.ReplaceAll(output, "\\", "")
	output = strings.ReplaceAll(output, "\"", "'")
	return output
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

	processedContent := PreprocessTemplate(templateContent)

	tmpl, err := template.New("tmpl").Funcs(funcs).Parse(processedContent)
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

func GetPlaceholders(readers []io.Reader, names ...[]string) (map[string]interface{}, error) {
	hasNames := false
	if len(names) > 0 {
		hasNames = true
	}

	var allPlaceholders []PlaceholderMetadata
	specificFields := make(map[string][]PlaceholderMetadata)
	credimiIDCount := make(map[string]int)

	for i, r := range readers {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return nil, err
		}
		content := buf.String()
		placeholders := ExtractPlaceholders(content)

		for _, ph := range placeholders {
			credimiIDCount[ph.CredimiID]++
			allPlaceholders = append(allPlaceholders, ph)
		}

		if hasNames {
			specificFields[names[0][i]] = placeholders
		} else {
			specificFields[content] = placeholders
		}
	}

	normalizedFields := make([]map[string]interface{}, 0)
	seenCredimiIDs := make(map[string]bool)
	for _, ph := range allPlaceholders {
		if credimiIDCount[ph.CredimiID] > 1 && !seenCredimiIDs[ph.CredimiID] {
			seenCredimiIDs[ph.CredimiID] = true
			field := map[string]interface{}{
				"CredimiID":           ph.CredimiID,
				"Type":                 ph.Type,
				"DescriptionKey": ph.DescriptionKey,
				"LabelKey":       ph.LabelKey,
			}
			normalizedFields = append(normalizedFields, field)
		}
	}

	result := map[string]interface{}{
		"normalized_fields": normalizedFields,
		"specific_fields":   specificFields,
	}

	return result, nil
}
