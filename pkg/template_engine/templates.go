package template_engine

import (
	"bytes"
	"io"
	"regexp"
	"text/template"

	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/group/all"
)


func ExtractPlaceholders(content string) []string {
	placeholderRegex := regexp.MustCompile(`{{\s*\.([a-zA-Z0-9_.]+)\s*}}`)
	matches := placeholderRegex.FindAllStringSubmatch(content, -1)

	var placeholders []string
	unique := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			name := match[1]
			if !unique[name] {
				unique[name] = true
				placeholders = append(placeholders, name)
			}
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

func GetPlaceholders(reader []io.Reader) ([]string, error) {
	var placeholders []string
	unique := make(map[string]bool)

	for _, r := range reader {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return nil, err
		}

		content := buf.String()
		placeholders = append(placeholders, ExtractPlaceholders(content)...)
	}

	for _, placeholder := range placeholders {
		if !unique[placeholder] {
			unique[placeholder] = true
		}
	}

	var result []string
	for placeholder := range unique {
		result = append(result, placeholder)
	}

	return result, nil
}
	
	

