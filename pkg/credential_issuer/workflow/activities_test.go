package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
	_ "modernc.org/sqlite"
	_ "modernc.org/sqlite/lib"
)

var wellKnownJSON = `{
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
}`

func TestFetchCredentialsIssuerActivity(t *testing.T) {
	testCases := []struct {
		name         string
		mockResponse string
		expectError  bool
	}{
		{
			name:         "Valid Response",
			mockResponse: wellKnownJSON,
			expectError:  false,
		},
		{
			name:         "Invalid JSON",
			mockResponse: `invalid-json`,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(FetchCredentialIssuerActivity)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, tc.mockResponse)
			}))
			defer server.Close()

			val, err := env.ExecuteActivity(FetchCredentialIssuerActivity, server.URL)
			if tc.expectError {
				assert.Error(t, err, "Expected an error")
				return
			}

			assert.NoError(t, err, "Did not expect an error")

			var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
			err = val.Get(&issuerData)
			assert.NoError(t, err, "Expected no error when retrieving activity result")
		})
	}
}

func TestStoreCredentialsActivity(t *testing.T) {
	testCases := []struct {
		name         string
		mockData     string
		dbPath       string
		expectError  bool
		expectedRows int
	}{
		{
			name:         "Valid Data",
			mockData:     wellKnownJSON,
			expectError:  false,
			expectedRows: 1,
		},
		{
			name:         "Invalid JSON",
			mockData:     `{invalid-json}`,
			expectError:  true,
			expectedRows: 0,
		},
		{
			name:         "Fail to Open DB",
			mockData:     wellKnownJSON,
			dbPath:       "/invalid/path/to/db.sqlite", // Simulate a non-writable or non-existent path
			expectError:  true,
			expectedRows: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(StoreCredentialsActivity)

			dbPath := tc.dbPath
			var tempFile *os.File
			var err error

			// Only create a temp file if no custom path is set (to simulate DB failure)
			if dbPath == "" {
				tempFile, err = os.CreateTemp("", "test_db_*.sqlite")
				assert.NoError(t, err, "Expected no error creating temp file")
				defer os.Remove(tempFile.Name())
				dbPath = tempFile.Name()
			}

			db, err := sql.Open("sqlite", dbPath)
			assert.NoError(t, err, "Expected no error opening database")
			defer db.Close()

			_, err = db.Exec(`
				CREATE TABLE credentials (
					format TEXT,
					issuer_name TEXT,
					name TEXT,
					locale TEXT,
					logo TEXT,
					credential_issuer TEXT
				);
			`)
			if !tc.expectError {
				assert.NoError(t, err, "Expected no error creating test database schema")
			}

			var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
			err = json.Unmarshal([]byte(tc.mockData), &issuerData)
			if tc.expectError && tc.dbPath == "" {
				assert.Error(t, err, "Expected JSON unmarshalling error")
				return
			}

			_, err = env.ExecuteActivity(StoreCredentialsActivity, &issuerData, "Test_Issuer", dbPath)
			if tc.expectError {
				assert.Error(t, err, "Expected an error from StoreCredentialsActivity")
				return
			}

			assert.NoError(t, err, "Did not expect an error from StoreCredentialsActivity")

			var count int
			err = db.QueryRow(`SELECT COUNT(*) FROM credentials`).Scan(&count)
			assert.NoError(t, err, "Expected no error querying database")
			assert.Equal(t, tc.expectedRows, count, "Unexpected row count in credentials table")
		})
	}
}
