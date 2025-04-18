// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package credentialissuer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchURL(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectedError  string
	}{
		{
			name: "Valid Metadata Response",
			serverResponse: `{
  "credential_configurations_supported": {
    "pid_jwt_vc_json": {
      "credential_signing_alg_values_supported": [
        "EdDSA",
        "ES256",
        "ES256K",
        "RSA"
      ],
      "cryptographic_binding_methods_supported": [
        "did:ebsi",
        "did:web",
        "did:jwk"
      ],
      "display": [
        {
          "name": "test ID"
        },
        {
          "locale": "nl-NL",
          "name": "Test ID"
        },
        {
          "locale": "en-GB",
          "name": "Personal ID"
        }
      ],
      "format": "jwt_vc_json"
    }
  },
  "credential_endpoint": "https://example/credential",
  "credential_issuer": "https://example.org",
  "deferred_credential_endpoint": "https://example.org/credential_deferred"
}`,
			statusCode:    http.StatusOK,
			expectedError: "",
		},
		{
			name:           "Invalid JSON Response",
			serverResponse: `{"issuer": "https://example.com", "name": `,
			statusCode:     http.StatusOK,
			expectedError:  "parse error",
		},
		{
			name:           "Non-200 Status Code",
			serverResponse: "",
			statusCode:     http.StatusNotFound,
			expectedError:  "not a credential issuer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock TLS server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			}))
			defer server.Close()

			// Call FetchCredentialIssuer
			metadata, err := FetchCredentialIssuer(server.URL)

			// Validate result
			if tt.expectedError == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedError != "" && (err == nil || !strings.Contains(err.Error(), tt.expectedError)) {
				t.Fatalf("expected error containing '%s', got: %v", tt.expectedError, err)
			}

			if tt.expectedError == "" && metadata == nil {
				t.Fatalf("expected valid metadata, got nil")
			}
		})
	}
}
