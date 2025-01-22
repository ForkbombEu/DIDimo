package metadata

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestFetchURL tests the FetchURL function with various scenarios.
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
    },
    "pid_vc+sd-jwt": {
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
          "name": "Test ID"
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
      "format": "vc+sd-jwt"
    }
  },
  "credential_endpoint": "https://example/credential",
  "credential_issuer": "https://example.org",
  "deferred_credential_endpoint": "https://example.org/credential_deferred",
  "display": [
    {
      "logo": {
        "alt_text": "example logo",
        "uri": "https://example/card_logo.png"
      },
      "name": "example B.V."
    },
    {
      "locale": "nl-NL",
      "logo": {
        "alt_text": "example logo",
        "uri": "https://example/card_logo.png"
      },
      "name": "example B.V."
    },
    {
      "locale": "en-US",
      "logo": {
        "alt_text": "example logo",
        "uri": "https://example/card_logo.png"
      },
      "name": "example B.V."
    }
  ]
}`,
			statusCode:    http.StatusOK,
			expectedError: "",
		},
		{
			name:           "Invalid JSON Response",
			serverResponse: `{"issuer": "https://example.com", "name": `,
			statusCode:     http.StatusOK,
			expectedError:  "failed to parse JSON",
		},
		{
			name: "Not conformant to JSON schema",
			serverResponse: `{
  "credential_issuer": "https://dev.issuer1.forkbomb.eu/credential_issuer",
  "credential_endpoint": "https://dev.issuer1.forkbomb.eu/credential_issuer/credential",
  "authorization_servers": [
    "https://dev.authz-server1.forkbomb.eu/authz_server"
  ],
  "display": [
    {
      "name": "Forkbomb Test Issuer",
      "locale": "en-US"
    }
  ],
  "jwks": {
    "keys": [
      {
        "kid": "did:dyne:sandbox.genericissuer:3suepGGjNHJmGDBebsCmapkdfBfXwFZzEQcEAMu7EdwA#es256_public_key",
        "crv": "P-256",
        "alg": "ES256",
        "kty": "EC"
      }
    ]
  },
  "credential_configurations_supported": {
    "discount_from_voucher": {
      "format": "vc+sd-jwt",
      "cryptographic_binding_methods_supported": [
        "jwk",
        "did:dyne:sandbox.signroom"
      ],
      "credential_signing_alg_values_supported": [
        "ES256"
      ],
      "proof_types_supported": {
        "jwt": {
          "proof_signing_alg_values_supported": [
            "ES256"
          ]
        }
      },
      "display": [
        {
          "name": "Get discount from Voucher dev",
          "locale": "en-US",
          "logo": {
            "url": "https://avatars.githubusercontent.com/u/96812851?s=200&v=4",
            "alt_text": "Get discount from Voucher dev logo",
            "uri": "https://avatars.githubusercontent.com/u/96812851?s=200&v=4"
          },
          "background_color": "#12107c",
          "text_color": "#FFFFFF",
          "description": "Get a special discount for all plans of DIDroom! Enter your voucher and get a discount credential."
        }
      ],
      "vct": "discount_from_voucher",
      "claims": {
        "has_discount_from_voucher": {
          "mandatory": true,
          "display": [
            {
              "locale": "en-US",
              "name": "Has a discount from Voucher"
            }
          ]
        }
      }
    }
  }
}`,
			statusCode:    http.StatusOK,
			expectedError: `failed to parse JSON: invalid value (expected one of []interface {}{"JWK", "jwk", "did", "did:web", "did:ebsi", "did:jwk", "did:dyne", "did:key", "cose_key"}): "did:dyne:sandbox.signroom"`,
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
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != credentialIssuerEndpoint {
					http.NotFound(w, r)
					return
				}
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			}))
			defer server.Close()

			// Call FetchURL
			metadata, err := FetchURL(server.URL)

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
