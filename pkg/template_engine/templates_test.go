package template_engine

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestExtractPlaceholders_SinglePlaceholder(t *testing.T) {
	content := "Hello, {{.Name}}!"
	expected := []string{"Name"}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestExtractPlaceholders_MultipleUniquePlaceholders(t *testing.T) {
	content := "Hello, {{.Name}}! Your age is {{.Age}} and your email is {{.Email}}. {{.Name}} appears twice."
	expected := []string{"Name", "Age", "Email"}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestExtractPlaceholders_NoPlaceholders(t *testing.T) {
	t.Skip()
	content := "Hello, World! This is a plain text without any placeholders."
	expected := []string{}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

func TestExtractPlaceholders_PlaceholdersWithNumbers(t *testing.T) {
	content := "Hello, {{.Name1}}! Your ID is {{.User2ID}} and your score is {{.Score3}}."
	expected := []string{"Name1", "User2ID", "Score3"}

	result := ExtractPlaceholders(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractPlaceholders() = %v, want %v", result, expected)
	}
}

//

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

	if !reflect.DeepEqual(sortStrings(result), sortStrings(expected)) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}
}

func sortStrings(strs []string) []string {
	sorted := make([]string, len(strs))
	copy(sorted, strs)
	sort.Strings(sorted)
	return sorted
}

func TestGetPlaceholders_ErrorOnReaderFailure(t *testing.T) {
	reader1 := strings.NewReader("Hello, {{.Name}}!")
	reader2 := &errorReader{err: fmt.Errorf("read error")}
	reader3 := strings.NewReader("Welcome, {{.User}}!")

	readers := []io.Reader{reader1, reader2, reader3}

	result, err := GetPlaceholders(readers)

	if err == nil {
		t.Error("GetPlaceholders() should have returned an error")
	}

	if result != nil {
		t.Errorf("GetPlaceholders() = %v, want nil", result)
	}
}

type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func TestGetPlaceholders_LargeInput(t *testing.T) {
	// Create a large input string with many placeholders
	largeInput := strings.Repeat("Hello, {{.Name}}! Your age is {{.Age}}. ", 10000)

	// Create a custom reader that simulates reading from a large file
	largeReader := &customReader{content: largeInput}

	readers := []io.Reader{largeReader}

	result, err := GetPlaceholders(readers)

	if err != nil {
		t.Errorf("GetPlaceholders() returned an error: %v", err)
	}

	expected := []string{"Name", "Age"}
	if !reflect.DeepEqual(sortStrings(result), sortStrings(expected)) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}

	// Check if memory usage is within acceptable limits
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.Alloc > 10*1024*1024 { // 100 MB
		t.Errorf("Memory usage too high: %d bytes", m.Alloc)
	}
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

	if !reflect.DeepEqual(sortStrings(result), sortStrings(expected)) {
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

	if !reflect.DeepEqual(sortStrings(result), sortStrings(expected)) {
		t.Errorf("GetPlaceholders() = %v, want %v", result, expected)
	}
}
