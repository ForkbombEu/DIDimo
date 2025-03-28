package template_engine

import (
	"io"
	"reflect"
	"strings"

	"testing"
)

func TestExtractPlaceholders_WithAllMetadataFields(t *testing.T) {
	content := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!
		Your age is {{.credimiPlaceHolder('Age', '5678', 'age_label', 'age_description', 'number', '30')}}.`

	expected := []PlaceholderMetadata{
		{
			FieldName:      "Name",
			CredimiID:      "1234",
			LabelKey:       "name_label",
			DescriptionKey: "name_description",
			Type:           "string",
			Example:        stringPtr("John Doe"),
		},
		{
			FieldName:      "Age",
			CredimiID:      "5678",
			LabelKey:       "age_label",
			DescriptionKey: "age_description",
			Type:           "number",
			Example:        stringPtr("30"),
		},
	}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestExtractPlaceholders_MultipleWithVaryingMetadata(t *testing.T) {
	content := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!
		Your age is {{.credimiPlaceHolder('Age', '5678', 'age_label', 'age_description', 'number')}}.
		Your email is {{.credimiPlaceHolder('Email', '9012', 'email_label', 'email_description', 'string', 'john@example.com')}}.`

	expected := []PlaceholderMetadata{
		{
			FieldName:      "Name",
			CredimiID:      "1234",
			LabelKey:       "name_label",
			DescriptionKey: "name_description",
			Type:           "string",
			Example:        stringPtr("John Doe"),
		},
		{
			FieldName:      "Age",
			CredimiID:      "5678",
			LabelKey:       "age_label",
			DescriptionKey: "age_description",
			Type:           "number",
			Example:        nil,
		},
		{
			FieldName:      "Email",
			CredimiID:      "9012",
			LabelKey:       "email_label",
			DescriptionKey: "email_description",
			Type:           "string",
			Example:        stringPtr("john@example.com"),
		},
	}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestRenderTemplate_WithCredimiPlaceholder(t *testing.T) {
	content := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!`
	data := map[string]interface{}{
		"Name": "Alice",
	}
	expected := "Hello Alice!"

	result, err := RenderTemplate(strings.NewReader(content), data)

	if err != nil {
		t.Errorf("RenderTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("RenderTemplate() = %q, want %q", result, expected)
	}
}

func TestPreprocessTemplate_MultiplePlaceholdersWithVaryingWhitespace(t *testing.T) {
	input := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string')}}!
		Your age is {{ .credimiPlaceHolder('Age','5678','age_label','age_description','number') }}.
		Welcome to {{  .credimiPlaceHolder(  'City',  '9012',  'city_label',  'city_description',  'string') }}!`

	expected := `Hello {{ .Name }}!
		Your age is {{ .Age }}.
		Welcome to {{ .City }}!`

	result := PreprocessTemplate(input)

	if result != expected {
		t.Errorf("PreprocessTemplate() = %v, want %v", result, expected)
	}
}

func TestPreprocessTemplate_WithNoPlaceholders(t *testing.T) {
	input := "Hello, world!"
	expected := "Hello, world!"

	result := PreprocessTemplate(input)

	if result != expected {
		t.Errorf("PreprocessTemplate() = %v, want %v", result, expected)
	}
}

func TestPreprocessTemplate_WithSproutFunctions(t *testing.T) {
	input := `Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string') | uppercase}}!
		Your age in 5 years will be {{.Age | add 5}}.`

	expected := `Hello {{ .Name | uppercase}}!
		Your age in 5 years will be {{.Age | add 5}}.`

	result := PreprocessTemplate(input)

	if result != expected {
		t.Errorf("PreprocessTemplate() = %v, want %v", result, expected)
	}
}

func TestGetPlaceholders_WithDuplicateCredimiIDs(t *testing.T) {
	reader1 := strings.NewReader("Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!")
	reader2 := strings.NewReader("Your age is {{.credimiPlaceHolder('Age', '1234', 'age_label', 'age_description', 'number', '30')}}.")
	reader3 := strings.NewReader("Welcome to {{.credimiPlaceHolder('City', '5678', 'city_label', 'city_description', 'string', 'New York')}}!")

	readers := []io.Reader{reader1, reader2, reader3}

	result, err := GetPlaceholders(readers)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of the expected type")
	}

	if len(normalizedFields) != 1 {
		t.Errorf("Expected 1 normalized field, got %d", len(normalizedFields))
	}

	expectedNormalizedField := map[string]interface{}{
		"CredimiID":      "1234",
		"Type":           "string",
		"DescriptionKey": "name_description",
		"LabelKey":       "name_label",
	}

	if !reflect.DeepEqual(normalizedFields[0], expectedNormalizedField) {
		t.Errorf("normalized_fields[0] = %v, want %v", normalizedFields[0], expectedNormalizedField)
	}

	specificFields, ok := result["specific_fields"].(map[string][]PlaceholderMetadata)
	if !ok {
		t.Errorf("specific_fields is not of the expected type")
	}

	if len(specificFields) != 3 {
		t.Errorf("Expected 3 entries in specific_fields, got %d", len(specificFields))
	}

	for content, placeholders := range specificFields {
		if len(placeholders) != 1 {
			t.Errorf("Expected 1 placeholder for content %q, got %d", content, len(placeholders))
		}
	}
}

func TestGetPlaceholders_WithSameFieldNameAndCredimiID(t *testing.T) {
	reader1 := strings.NewReader("Hello {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'John Doe')}}!")
	reader2 := strings.NewReader("Your name is {{.credimiPlaceHolder('Name', '1234', 'name_label', 'name_description', 'string', 'Jane Doe')}}.")

	readers := []io.Reader{reader1, reader2}

	result, err := GetPlaceholders(readers)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of the expected type")
	}

	if len(normalizedFields) != 1 {
		t.Errorf("Expected 1 normalized field, got %d", len(normalizedFields))
	}

	expectedNormalizedField := map[string]interface{}{
		"CredimiID":      "1234",
		"Type":           "string",
		"DescriptionKey": "name_description",
		"LabelKey":       "name_label",
	}

	if !reflect.DeepEqual(normalizedFields[0], expectedNormalizedField) {
		t.Errorf("normalized_fields[0] = %v, want %v", normalizedFields[0], expectedNormalizedField)
	}

	specificFields, ok := result["specific_fields"].(map[string][]PlaceholderMetadata)
	if !ok {
		t.Errorf("specific_fields is not of the expected type")
	}

	if len(specificFields) != 2 {
		t.Errorf("Expected 2 entries in specific_fields, got %d", len(specificFields))
	}

	for content, placeholders := range specificFields {
		if len(placeholders) != 1 {
			t.Errorf("Expected 1 placeholder for content %q, got %d", content, len(placeholders))
		}
		if placeholders[0].FieldName != "Name" || placeholders[0].CredimiID != "1234" {
			t.Errorf("Unexpected placeholder: %+v", placeholders[0])
		}
	}
}

func TestGetPlaceholders_WithSameFieldNameDifferentCredimiIDs(t *testing.T) {
	reader1 := strings.NewReader("Hello {{.credimiPlaceHolder('Name', '1234', 'name_label_1', 'name_description_1', 'string', 'John Doe')}}!")
	reader2 := strings.NewReader("Your name is {{.credimiPlaceHolder('Name', '5678', 'name_label_2', 'name_description_2', 'string', 'Jane Doe')}}.")

	readers := []io.Reader{reader1, reader2}

	result, err := GetPlaceholders(readers)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of the expected type")
	}

	if len(normalizedFields) != 0 {
		t.Errorf("Expected 0 normalized fields, got %d", len(normalizedFields))
	}

	specificFields, ok := result["specific_fields"].(map[string][]PlaceholderMetadata)
	if !ok {
		t.Errorf("specific_fields is not of the expected type")
	}

	if len(specificFields) != 2 {
		t.Errorf("Expected 2 entries in specific_fields, got %d", len(specificFields))
	}

	for content, placeholders := range specificFields {
		if len(placeholders) != 1 {
			t.Errorf("Expected 1 placeholder for content %q, got %d", content, len(placeholders))
		}
		if placeholders[0].FieldName != "Name" {
			t.Errorf("Expected FieldName to be 'Name', got %s", placeholders[0].FieldName)
		}
	}

	allPlaceholders := []PlaceholderMetadata{}
	for _, placeholders := range specificFields {
		allPlaceholders = append(allPlaceholders, placeholders...)
	}

	if len(allPlaceholders) != 2 {
		t.Errorf("Expected 2 total placeholders, got %d", len(allPlaceholders))
	}

	credimiIDs := map[string]bool{
		"1234": false,
		"5678": false,
	}

	for _, ph := range allPlaceholders {
		if _, exists := credimiIDs[ph.CredimiID]; !exists {
			t.Errorf("Unexpected CredimiID: %s", ph.CredimiID)
		}
		credimiIDs[ph.CredimiID] = true
	}

	for id, seen := range credimiIDs {
		if !seen {
			t.Errorf("CredimiID %s was not found in the placeholders", id)
		}
	}
}

func TestGetPlacholders_WithNames(t *testing.T) {
	reader1 := strings.NewReader("Hello {{.credimiPlaceHolder('Name', '1234', 'name_label_1', 'name_description_1', 'string', 'John Doe')}}!")
	reader2 := strings.NewReader("Your name is {{.credimiPlaceHolder('Name', '5678', 'name_label_2', 'name_description_2', 'string', 'Jane Doe')}}.")

	readers := []io.Reader{reader1, reader2}

	names := []string{
		"content1",
		"content2",
	}

	result, err := GetPlaceholders(readers, names)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	normalizedFields, ok := result["normalized_fields"].([]map[string]interface{})
	if !ok {
		t.Errorf("normalized_fields is not of the expected type")
	}

	if len(normalizedFields) != 0 {
		t.Errorf("Expected 0 normalized fields, got %d", len(normalizedFields))
	}

	specificFields, ok := result["specific_fields"].(map[string][]PlaceholderMetadata)
	if !ok {
		t.Errorf("specific_fields is not of the expected type")
	}

	if len(specificFields) != 2 {
		t.Errorf("Expected 2 entries in specific_fields, got %d", len(specificFields))
	}

	for content, placeholders := range specificFields {
		if len(placeholders) != 1 {
			t.Errorf("Expected 1 placeholder for content %q, got %d", content, len(placeholders))
		}
		if placeholders[0].FieldName != "Name" {
			t.Errorf("Expected FieldName to be 'Name', got %s", placeholders[0].FieldName)
		}
	}

	allPlaceholders := []PlaceholderMetadata{}
	for _, placeholders := range specificFields {
		allPlaceholders = append(allPlaceholders, placeholders...)
	}

	if len(allPlaceholders) != 2 {
		t.Errorf("Expected 2 total placeholders, got %d", len(allPlaceholders))
	}

	credimiIDs := map[string]bool{
		"1234": false,
		"5678": false,
	}

	for _, ph := range allPlaceholders {
		if _, exists := credimiIDs[ph.CredimiID]; !exists {
			t.Errorf("Unexpected CredimiID: %s", ph.CredimiID)
		}
		credimiIDs[ph.CredimiID] = true
	}

	for id, seen := range credimiIDs {
		if !seen {
			t.Errorf("CredimiID %s was not found in the placeholders", id)
		}
	}
}
