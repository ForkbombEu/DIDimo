package metadata

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	json "github.com/neilotoole/jsoncolor"
)

var credentialIssuerEndpoint string = "/.well-known/openid-credential-issuer"

// FetchURL fetches the metadata from the provided URL and validates the credential issuer.
func FetchURL(baseURL string) (*OpenidCredentialIssuerSchemaJson, error) {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	credentialIssuerURL := strings.TrimRight(baseURL, "/") + credentialIssuerEndpoint

	resp, err := http.Get(credentialIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach credential issuer URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not a credential issuer: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var metadata OpenidCredentialIssuerSchemaJson
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &metadata, nil
}

// Output metadata to the provided writer
func PrintJSON(metadata *OpenidCredentialIssuerSchemaJson, writer io.Writer) error {
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
	if err := enc.Encode(metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %v", err)
	}

	return nil
}
