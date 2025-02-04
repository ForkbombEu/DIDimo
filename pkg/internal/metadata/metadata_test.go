package metadata

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Example struct for unmarshalling into
type Example struct {
	Name string `json:"name"`
}

func TestFetchJSON(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectedError  string
		expectedData   *Example
	}{
		{
			name:           "success",
			serverResponse: `{"name": "Test"}`,
			statusCode:     http.StatusOK,
			expectedError:  "",
			expectedData:   &Example{Name: "Test"},
		},
		{
			name:           "non_200_status",
			serverResponse: "",
			statusCode:     http.StatusNotFound,
			expectedError:  "HTTP error:",
			expectedData:   nil,
		},
		{
			name:           "malformed_json",
			serverResponse: `{"name": "Test"`, // Malformed JSON
			statusCode:     http.StatusOK,
			expectedError:  "parse error",
			expectedData:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/test/example" {
					http.NotFound(w, r)
					return
				}
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			}))
			defer server.Close()

			// Call FetchJSON with the test server URL
			data, err := FetchJSON[Example](server.URL, "/test/example")

			// If we expect an error, ensure it matches the expected error
			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %v", tt.expectedError, err)
				}
			} else {
				// If no error is expected, ensure the data matches the expected result
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if data == nil || *data != *tt.expectedData {
					t.Errorf("expected data %v, got %v", tt.expectedData, data)
				}
			}
		})
	}
}

// TestPrintJSON tests the PrintJSON function
func TestPrintJSON(t *testing.T) {
	// Create an example object to print
	data := Example{Name: "Test"}

	// Create a buffer to capture the printed JSON
	var buf bytes.Buffer

	// Call PrintJSON to print data to the buffer
	err := PrintJSON(data, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Read the output from the buffer
	output := buf.String()

	// Check that the output contains the expected value
	if !strings.Contains(output, `"name": "Test"`) {
		t.Errorf("expected JSON to contain 'name: Test', got %v", output)
	}
}
