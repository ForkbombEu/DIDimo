// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package template_engine

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

// func TestExtractPlaceholders_WithAllMetadataFields(t *testing.T) {
// 	content := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!
// 		Your age is {{.credimiPlaceHolder('Age', '5678', 'age_label', 'age_description', 'number', '30')}}.`

// 	expected := []PlaceholderMetadata{
// 		{
// 			FieldName:      "Name",
// 			CredimiID:      "1234",
// 			LabelKey:       "name_label",
// 			DescriptionKey: "name_description",
// 			Type:           "string",
// 			Example:        stringPtr("John Doe"),
// 		},
// 		{
// 			FieldName:      "Age",
// 			CredimiID:      "5678",
// 			LabelKey:       "age_label",
// 			DescriptionKey: "age_description",
// 			Type:           "number",
// 			Example:        stringPtr("30"),
// 		},
// 	}

// 	result := ExtractPlaceholders(content)

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
// 	}
// }

// func stringPtr(s string) *string {
// 	return &s
// }

// func TestExtractPlaceholders_MultipleWithVaryingMetadata(t *testing.T) {
// 	content := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!
// 		Your age is {{.credimiPlaceHolder('Age', '5678', 'age_label', 'age_description', 'number')}}.
// 		Your email is {{.credimiPlaceHolder('Email', '9012', 'email_label', 'email_description', 'string', 'john@example.com')}}.`

// 	expected := []PlaceholderMetadata{
// 		{
// 			FieldName:      "Name",
// 			CredimiID:      "1234",
// 			LabelKey:       "name_label",
// 			DescriptionKey: "name_description",
// 			Type:           "string",
// 			Example:        "John Doe",
// 		},
// 		{
// 			FieldName:      "Age",
// 			CredimiID:      "5678",
// 			LabelKey:       "age_label",
// 			DescriptionKey: "age_description",
// 			Type:           "number",
// 			Example:        "",
// 		},
// 		{
// 			FieldName:      "Email",
// 			CredimiID:      "9012",
// 			LabelKey:       "email_label",
// 			DescriptionKey: "email_description",
// 			Type:           "string",
// 			Example:        "john@example.com",
// 		},
// 	}

// 	result := ExtractPlaceholders(content)

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
// 	}
// }

// func TestRenderTemplate_WithCredimiPlaceholder(t *testing.T) {
// 	content := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!`
// 	data := map[string]interface{}{
// 		"Name": "Alice",
// 	}
// 	expected := "Hello Alice!"

// 	result, err := RenderTemplate(strings.NewReader(content), data)

// 	if err != nil {
// 		t.Errorf("RenderTemplate() returned an error: %v", err)
// 	}

// 	if result != expected {
// 		t.Errorf("RenderTemplate() = %q, want %q", result, expected)
// 	}
// }

func TestPreprocessTemplate(t *testing.T) {
	input := "{{ credimiPlaceholder \"testalias\" \"autogenerated\" \"i18n_testalias\" \"i18n_testalias_description\" \"string\" \"\" }}"

	expected := "{{ .testalias }}"
	result, err := PreprocessTemplate(input)
	if err != nil {
		t.Errorf("PreprocessTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("PreprocessTemplate() = %v, want %v", result, expected)
	}
}

func TestPreprocessTemplate_WithExample(t *testing.T) {
	input := "{{ credimiPlaceholder \"testalias\" \"autogenerated\" \"i18n_testalias\" \"i18n_testalias_description\" \"string\" \"example\" }}"

	expected := "{{ .testalias }}"
	result, err := PreprocessTemplate(input)
	if err != nil {
		t.Errorf("PreprocessTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("PreprocessTemplate() = %v, want %v", result, expected)
	}
}

func TestPreprocessTemplate_WithSproutFunctions(t *testing.T) {
	input := "{{ credimiPlaceholder \"testalias\" \"autogenerated\" \"i18n_testalias\" \"i18n_testalias_description\" \"string\" \"\" | toUpper }}"

	expected := "{{ .TESTALIAS }}"
	result, err := PreprocessTemplate(input)
	if err != nil {
		t.Errorf("PreprocessTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("PreprocessTemplate() = %v, want %v", result, expected)
	}
}

func TestPreprocessTemplate_ExtractMetadata(t *testing.T) {
	input := "{{ credimiPlaceholder \"testalias\" \"autogenerated\" \"i18n_testalias\" \"i18n_testalias_description\" \"string\" \"\" }}"

	expected := PlaceholderMetadata{
		FieldName:      "testalias",
		CredimiID:      "autogenerated",
		LabelKey:       "i18n_testalias",
		DescriptionKey: "i18n_testalias_description",
		Type:           "string",
		Example:        "",
	}
	_, err := PreprocessTemplate(input)
	if err != nil {
		t.Errorf("PreprocessTemplate() returned an error: %v", err)
	}

	result := ExtractMetadata()

	if len(result) != 1 {
		t.Errorf("Expected 1 metadata entry, got %d", len(result))
	}

	if !reflect.DeepEqual(result[0], expected) {
		t.Errorf("ExtractMetadata() = %v, want %v", result[0], expected)
	}
}

func TestGetPlaceholders_MultipleReaders(t *testing.T) {
	reader1 := strings.NewReader("Hello {{ credimiPlaceholder \"name\" \"123\" \"name_label\" \"name_desc\" \"string\" \"\" }}")
	reader2 := strings.NewReader("Age: {{ credimiPlaceholder \"age\" \"456\" \"age_label\" \"age_desc\" \"number\" \"\" }}")
	reader3 := strings.NewReader("Hello {{ credimiPlaceholder \"name\" \"123\" \"name_label\" \"name_desc\" \"string\" \"\" }} again")

	readers := []io.Reader{reader1, reader2, reader3}
	names := []string{"template1", "template2", "template3"}

	result, err := GetPlaceholders(readers, names)
	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	// Check normalized fields
	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of expected type")
	}
	if len(normalizedFields) != 1 {
		t.Errorf("Expected 1 normalized field, got %d", len(normalizedFields))
	}
	if normalizedFields[0]["CredimiID"] != "123" {
		t.Errorf("Expected CredimiID 123, got %v", normalizedFields[0]["CredimiID"])
	}

	// Check specific fields
	specificFields, ok := result["specific_fields"].(map[string]interface{})
	if !ok {
		t.Errorf("specific_fields is not of expected type")
	}
	if len(specificFields) != 3 {
		t.Errorf("Expected 3 specific fields, got %d", len(specificFields))
	}

	// Check content of template1
	template1, ok := specificFields["template1"].(map[string]interface{})
	if !ok {
		t.Errorf("template1 is not of expected type")
	}
	if template1["content"] != "Hello {{ .name }}" {
		t.Errorf("Unexpected content for template1: %v", template1["content"])
	}

	// Check fields of template2
	template2, ok := specificFields["template2"].(map[string]interface{})

	if !ok {
		t.Errorf("template2 is not of expected type")
	}
	fields2, ok := template2["fields"].([]PlaceholderMetadata)
	if !ok {
		t.Errorf("fields of template2 is not of expected type")
	}
	if len(fields2) != 1 || fields2[0].FieldName != "age" {
		t.Errorf("Unexpected fields for template2: %v", fields2)
	}
}

func TestGetPlaceholders_WithMultipleCredimiIDs(t *testing.T) {
	reader1 := strings.NewReader("Hello {{ credimiPlaceholder \"name\" \"123\" \"name_label\" \"name_desc\" \"string\" \"\" }}")
	reader2 := strings.NewReader("Hello {{ credimiPlaceholder \"name\" \"123\" \"name_label\" \"name_desc\" \"string\" \"\" }} again")

	readers := []io.Reader{reader1, reader2}
	names := []string{"template1", "template2"}

	result, err := GetPlaceholders(readers, names)
	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of expected type")
	}
	if len(normalizedFields) != 1 {
		t.Errorf("Expected 1 normalized field, got %d", len(normalizedFields))
	}
	if normalizedFields[0]["CredimiID"] != "123" {
		t.Errorf("Expected CredimiID 123, got %v", normalizedFields[0]["CredimiID"])
	}
}

func TestGetPlaceholders_ComplexTemplate(t *testing.T) {
	templateContent := `
	{
		"variant": {
			"credential_format": "iso_mdl",
			"client_id_scheme": "did",
			"request_method": "request_uri_signed",
			"response_mode": "direct_post"
		},
		"form": {
			"alias": "{{ credimiPlaceholder "testalias" "autogenerated" "i18n_testalias" "i18n_testalias_description" "string" "" }}",
			"client": {
				"client_id": "{{ credimiPlaceholder "client_id" "id_credimi_id" "i18n_client_id" "i18n_client_id_description" "string" "" }}",
				"jwks": "{{ credimiPlaceholder "jwks" "id_credimi_id" "i18n_jwks" "i18n_jwks_description" "object" "" }}",
				"presentation_definition": "{{ credimiPlaceholder "presentation_definition" "id_credimi_id" "i18n_presentation_definition" "i18n_presentation_definition_description" "object" "" }}"
			},
			"description": "{{ credimiPlaceholder "description" "autogenerated" "i18n_description" "i18n_description_description" "string" "" }}",
			"server": {
				"authorization_endpoint": "openid-vc://"
			}
		}
	}`

	reader := strings.NewReader(templateContent)
	readers := []io.Reader{reader}
	names := []string{"complex_template"}

	result, err := GetPlaceholders(readers, names)
	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	// Check normalized fields
	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of expected type")
	}
	if len(normalizedFields) != 2 {
		t.Errorf("Expected 1 normalized field, got %d", len(normalizedFields))
	}
	if normalizedFields[0]["CredimiID"] != "autogenerated" {
		t.Errorf("Expected CredimiID autogenerated, got %v", normalizedFields[0]["autogenerated"])
	}

	// Check specific fields
	specificFields, ok := result["specific_fields"].(map[string]interface{})
	if !ok {
		t.Errorf("specific_fields is not of expected type")
	}
	if len(specificFields) != 1 {
		t.Errorf("Expected 1 specific field, got %d", len(specificFields))
	}

	// Check content of the complex template
	templateData, ok := specificFields["complex_template"].(map[string]interface{})
	if !ok {
		t.Errorf("complex_template is not of expected type")
	}
	content, ok := templateData["content"].(string)
	if !ok {
		t.Errorf("content of complex_template is not of expected type")
	}
	if !strings.Contains(content, "{{ .testalias }}") || !strings.Contains(content, "{{ .client_id }}") {
		t.Errorf("Unexpected content for complex_template: %v", content)
	}

	// Check fields of the complex template
	fields, ok := templateData["fields"].([]PlaceholderMetadata)
	if !ok {
		t.Errorf("fields of complex_template is not of expected type")
	}
	expectedFields := map[string]string{
		"testalias":               "autogenerated",
		"client_id":               "id_credimi_id",
		"jwks":                    "id_credimi_id",
		"presentation_definition": "id_credimi_id",
		"description":             "autogenerated",
	}
	if len(fields) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(fields))
	}
	for _, field := range fields {
		if expectedCredimiID, exists := expectedFields[field.FieldName]; !exists || field.CredimiID != expectedCredimiID {
			t.Errorf("Unexpected field: %v", field)
		}
	}
}
