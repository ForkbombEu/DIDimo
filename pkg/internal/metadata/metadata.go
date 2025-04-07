// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package metadata

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	json "github.com/neilotoole/jsoncolor"
)

type NetworkError struct {
	Message string
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %s", e.Message)
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d %s", e.StatusCode, e.Message)
}

type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error: %s", e.Message)
}

// FetchJSON fetches the metadata from the provided URL and unmarshals it into the target structure.
func FetchJSON[T any](baseURL string, endpoint string) (*T, error) {
	fetchURL := strings.TrimRight(baseURL, "/") + endpoint

	resp, err := http.Get(fetchURL)
	if err != nil {
		return nil, &NetworkError{Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Message: resp.Status}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	var target T
	if err := json.Unmarshal(body, &target); err != nil {
		return nil, &ParseError{Message: err.Error()}
	}

	return &target, nil
}

// PrintJSON outputs the provided data structure as JSON to the writer.
// It uses colorized output if the writer supports it.
func PrintJSON[T any](data T, writer io.Writer) error {
	var enc *json.Encoder

	// Note: this check will fail if running inside Goland (and
	// other IDEs?) as IsColorTerminal will return false.
	if json.IsColorTerminal(writer) {
		// Safe to use color
		enc = json.NewEncoder(writer)

		// DefaultColors are similar to jq
		clrs := json.DefaultColors()
		enc.SetColors(clrs)
	} else {
		// Can't use color; but the encoder will still work
		enc = json.NewEncoder(writer)
	}
	enc.SetIndent("", "  ")
	// Encode the metadata into the writer
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to write metadata: %v", err)
	}

	return nil
}
