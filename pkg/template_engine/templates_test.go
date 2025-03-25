package template_engine

import (
	// "fmt"
	"io"
	"reflect"

	// "runtime"
	"sort"
	"strings"
	"testing"
)

func TestExtractPlaceholders_SinglePlaceholder(t *testing.T) {
	content := "Hello, {{.Name}}!"
	expected := []PlaceholderMetadata{
		{
			Field:        "Name",
			Translations: map[string]string{},
			Descriptions: map[string]string{},
			CredimiID:    "",
		},
	}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestExtractPlaceholders_MultipleUniquePlaceholders(t *testing.T) {
	content := "Hello, {{.Name}}! Your age is {{.Age}} and your email is {{.Email}}. {{.Name}} appears twice."
	expected := []PlaceholderMetadata{
		{Field: "Name", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
		{Field: "Age", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
		{Field: "Email", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
	}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestExtractPlaceholders_WithCompleteMetadata(t *testing.T) {
	content := "Hello, {{.Name|en:name, it:nome, description:en:Enter your name, description:it:Inserisci il tuo nome, credimi_id:1234}}!"
	expected := []PlaceholderMetadata{
		{
			Field:        "Name",
			Translations: map[string]string{"en": "name", "it": "nome"},
			Descriptions: map[string]string{"en": "Enter your name", "it": "Inserisci il tuo nome"},
			CredimiID:    "1234",
		},
	}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestExtractPlaceholders_PlaceholdersWithNumbers(t *testing.T) {
	content := "Hello, {{.Name1}}! Your ID is {{.User2ID}} and your score is {{.Score3}}."
	expected := []PlaceholderMetadata{
		{Field: "Name1", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
		{Field: "User2ID", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
		{Field: "Score3", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
	}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestRenderTemplate_ValidContentAndData(t *testing.T) {
	content := "Hello, {{.Name}}! Your age is {{.Age}}."
	data := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}
	expected := "Hello, Alice! Your age is 30."

	result, err := RenderTemplate(strings.NewReader(content), data)

	if err != nil {
		t.Errorf("RenderTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("RenderTemplate() = %v, want %v", result, expected)
	}
}

func TestRenderTemplate_EmptyContent(t *testing.T) {
	content := ""
	data := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}
	expected := ""

	result, err := RenderTemplate(strings.NewReader(content), data)

	if err != nil {
		t.Errorf("renderTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("renderTemplate() = %q, want %q", result, expected)
	}
}

func TestRenderTemplate_EmptyDataMap(t *testing.T) {
	content := "Hello, {{.Name}}! Welcome to {{.Place}}."
	data := map[string]interface{}{}
	expected := "Hello, <no value>! Welcome to <no value>."

	result, err := RenderTemplate(strings.NewReader(content), data)

	if err != nil {
		t.Errorf("RenderTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("RenderTemplate() = %q, want %q", result, expected)
	}
}

func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	content := "Hello, {{.Name}! Your age is {{.Age}}."
	data := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}

	result, err := RenderTemplate(strings.NewReader(content), data)

	if err == nil {
		t.Error("RenderTemplate() should have returned an error for invalid template")
	}

	if result != "" {
		t.Errorf("RenderTemplate() = %q, want empty string", result)
	}
}

type customReader struct {
	content string
	pos     int
}

func (r *customReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.content) {
		return 0, io.EOF
	}

	n = copy(p, r.content[r.pos:])
	r.pos += n
	return n, nil
}

func TestRenderTemplate_WithCustomReader(t *testing.T) {
	content := "Hello, {{.Name}}! Your age is {{.Age}}."
	data := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}
	expected := "Hello, Alice! Your age is 30."

	custom := &customReader{content: content}

	result, err := RenderTemplate(custom, data)

	if err != nil {
		t.Errorf("RenderTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("RenderTemplate() = %v, want %v", result, expected)
	}
}

func TestRenderTemplate_WithSproutFunctions(t *testing.T) {
	content := "Hello, {{.Name}}! Your age in 5 years will be {{.Age | add 5}}."
	data := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}
	expected := "Hello, Alice! Your age in 5 years will be 35."

	result, err := RenderTemplate(strings.NewReader(content), data)

	if err != nil {
		t.Errorf("RenderTemplate() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("RenderTemplate() = %q, want %q", result, expected)
	}
}

func TestGetPlaceholders_MultipleReadersWithDifferentContent(t *testing.T) {
	reader1 := strings.NewReader("Hello, {{.Name}}! Your age is {{.Age}}.")
	reader2 := strings.NewReader("Welcome to {{.City}}, {{.Name}}! The weather is {{.Temperature}}.")
	reader3 := strings.NewReader("Your favorite color is {{.Color}}.")

	readers := []io.Reader{reader1, reader2, reader3}

	expected := []string{"Name", "Age", "City", "Temperature", "Color"}

	result, err := GetPlaceholders(readers)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	resultFields := make([]string, len(result))
	for i, placeholder := range result {
		resultFields[i] = placeholder.Field
	}
	if !reflect.DeepEqual(sortStrings(resultFields), sortStrings(expected)) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}
}

func sortStrings(strs []string) []string {
	sorted := make([]string, len(strs))
	copy(sorted, strs)
	sort.Strings(sorted)
	return sorted
}

func TestGetPlaceholders_UniqueWithSproutFunctions(t *testing.T) {
	reader1 := strings.NewReader("Hello, {{.Name}}! Your age is {{.Age}}.")
	reader2 := strings.NewReader("Welcome, {{.Name | capitalize}}! In 5 years, you'll be {{.Age | add 5}}.")
	reader3 := strings.NewReader("Goodbye, {{.Name | lowercase}}! Your birth year is {{.Age | subtract 2023}}.")

	readers := []io.Reader{reader1, reader2, reader3}

	expected := []string{"Name", "Age"}

	result, err := GetPlaceholders(readers)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	resultFields := make([]string, len(result))
	for i, placeholder := range result {
		resultFields[i] = placeholder.Field
	}
	if !reflect.DeepEqual(sortStrings(resultFields), sortStrings(expected)) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}
}

func TestGetPlaceholders_WithNoNormalization(t *testing.T) {
	reader1 := strings.NewReader("Hello, {{.Name}}! Your age is {{.Age}}.")
	reader2 := strings.NewReader("Welcome, {{.Name}}! Your age is {{.Age}}.")
	reader3 := strings.NewReader("Goodbye, {{.Name}}! Your age is {{.Age}}.")

	readers := []io.Reader{reader1, reader2, reader3}

	expected := []string{"Name", "Age", "Name", "Age", "Name", "Age"}

	result, err := GetPlaceholders(readers, false)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	resultFields := make([]string, len(result))
	for i, placeholder := range result {
		resultFields[i] = placeholder.Field
	}
	if !reflect.DeepEqual(sortStrings(resultFields), sortStrings(expected)) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}
}

func TestGetPlaceholders_WithAllMetadataSetAndMultipleReadersAndSprouts(t *testing.T) {
	reader1 := strings.NewReader("Hello, {{.Name|en:name, it:nome, description:en:Enter your name, description:it:Inserisci il tuo nome, credimi_id:1234, type: string}}! Your age is {{.Age }}.")
	reader2 := strings.NewReader("Welcome, {{.Name|en:name, it:nome, description:en:Enter your name, description:it:Inserisci il tuo nome, credimi_id:1234, type: string}}! Your age is {{.Age}}.")
	reader3 := strings.NewReader("Goodbye, {{.Name|en:name, it:nome, description:en:Enter your name, description:it:Inserisci il tuo nome, credimi_id:1234, type: string}}! Your age is {{.Age}}.")

	readers := []io.Reader{reader1, reader2, reader3}

	expected := []PlaceholderMetadata{
		{
			Field:        "Name",
			Translations: map[string]string{"en": "name", "it": "nome"},
			Descriptions: map[string]string{"en": "Enter your name", "it": "Inserisci il tuo nome"},
			CredimiID:    "1234",
			Type:         "string",
		},
		{Field: "Age", Translations: map[string]string{}, Descriptions: map[string]string{}, CredimiID: ""},
	}

	result, err := GetPlaceholders(readers, true)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}
}

func TestParseMetadata_EmptyInput(t *testing.T) {
	metaStr := ""
	expectedTranslations := make(map[string]string)
	expectedDescriptions := make(map[string]string)
	expectedCredimiID := ""
	expectedPlaceholderType := ""

	translations, descriptions, credimiID, placeHolderType := parseMetadata(metaStr)

	if !reflect.DeepEqual(translations, expectedTranslations) {
		t.Errorf("parseMetadata() translations = %v, want %v", translations, expectedTranslations)
	}

	if !reflect.DeepEqual(descriptions, expectedDescriptions) {
		t.Errorf("parseMetadata() descriptions = %v, want %v", descriptions, expectedDescriptions)
	}

	if credimiID != expectedCredimiID {
		t.Errorf("parseMetadata() credimiID = %v, want %v", credimiID, expectedCredimiID)
	}

	if placeHolderType != expectedPlaceholderType {
		t.Errorf("parseMetadata() placeHolderType = %v, want %v", placeHolderType, expectedPlaceholderType)
	}
}

func TestParseMetadata_MultipleTranslationsAndDescriptions(t *testing.T) {
	metaStr := "en:name, it:nome, description:en:Enter your name, description:it:Inserisci il tuo nome, credimi_id:1234"
	expectedTranslations := map[string]string{
		"en": "name",
		"it": "nome",
	}
	expectedDescriptions := map[string]string{
		"en": "Enter your name",
		"it": "Inserisci il tuo nome",
	}
	expectedCredimiID := "1234"

	translations, descriptions, credimiID, _ := parseMetadata(metaStr)

	if !reflect.DeepEqual(translations, expectedTranslations) {
		t.Errorf("parseMetadata() translations = %v, want %v", translations, expectedTranslations)
	}

	if !reflect.DeepEqual(descriptions, expectedDescriptions) {
		t.Errorf("parseMetadata() descriptions = %v, want %v", descriptions, expectedDescriptions)
	}

	if credimiID != expectedCredimiID {
		t.Errorf("parseMetadata() credimiID = %v, want %v", credimiID, expectedCredimiID)
	}
}

func TestParseMetadata_MultipleCredimiIDs(t *testing.T) {
	metaStr := "en:name, it:nome, description:en:Enter your name, credimi_id:1234, description:it:Inserisci il tuo nome, credimi_id:5678"
	expectedTranslations := map[string]string{
		"en": "name",
		"it": "nome",
	}
	expectedDescriptions := map[string]string{
		"en": "Enter your name",
		"it": "Inserisci il tuo nome",
	}
	expectedCredimiID := "5678"

	translations, descriptions, credimiID, _ := parseMetadata(metaStr)

	if !reflect.DeepEqual(translations, expectedTranslations) {
		t.Errorf("parseMetadata() translations = %v, want %v", translations, expectedTranslations)
	}

	if !reflect.DeepEqual(descriptions, expectedDescriptions) {
		t.Errorf("parseMetadata() descriptions = %v, want %v", descriptions, expectedDescriptions)
	}

	if credimiID != expectedCredimiID {
		t.Errorf("parseMetadata() credimiID = %v, want %v", credimiID, expectedCredimiID)
	}
}
