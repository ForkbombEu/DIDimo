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

type mockCredentials struct {
	CredentialDefinition                 *credentialissuer.CredentialDefinition                      `json:"credential_definition,omitempty"`
	CredentialSigningAlgValuesSupported  []credentialissuer.CredentialSigningAlgValuesSupportedElem  `json:"credential_signing_alg_values_supported,omitempty"`
	CryptographicBindingMethodsSupported []credentialissuer.CryptographicBindingMethodsSupportedElem `json:"cryptographic_binding_methods_supported,omitempty"`
	Display                              []credentialissuer.DisplayElem_1                            `json:"display,omitempty"`
	Format                               string                                                      `json:"format"`
	ProofTypesSupported                  credentialissuer.ProofTypesSupported                        `json:"proof_types_supported,omitempty"`
	Scope                                *string                                                     `json:"scope,omitempty"`
}

func TestStoreOrUpdateCredentialsActivity(t *testing.T) {

	var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
	err := json.Unmarshal([]byte(wellKnownJSON), &issuerData)
	assert.NoError(t, err, "Did not expect an error")

	testCases := []struct {
		name         string
		dbPath       string
		credential   mockCredentials
		expectError  bool
		expectedRows int
	}{
		{
			name:         "Valid Data - Insert",
			credential:   mockCredentials(issuerData.CredentialConfigurationsSupported["discount_from_voucher"]),
			expectError:  false,
			expectedRows: 1,
		},
		{
			name:         "Valid Data - Update",
			credential:   mockCredentials(issuerData.CredentialConfigurationsSupported["discount_from_voucher"]),
			expectError:  false,
			expectedRows: 1,
		},
		{
			name:         "Fail to Open DB",
			credential:   mockCredentials{},
			dbPath:       "/invalid/path/to/db.sqlite", // Simulate a non-writable or non-existent path
			expectError:  true,
			expectedRows: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(StoreOrUpdateCredentialsActivity)

			// Set up the database
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
				CREATE TABLE IF NOT EXISTS credentials (
					format TEXT,
					issuer_name TEXT,
					name TEXT,
					locale TEXT,
					logo TEXT,
					credential_issuer TEXT,
					key TEXT PRIMARY KEY,
					json TEXT,
					created TIMESTAMP,
					updated TIMESTAMP
				);
			`)
			if !tc.expectError {
				assert.NoError(t, err, "Expected no error creating test database schema")
			}

			// Insert initial credential if updating
			if tc.name == "Valid Data - Update" {
				initialCredential := struct {
					Format           string `json:"format"`
					Name             string `json:"name"`
					Locale           string `json:"locale"`
					Logo             string `json:"logo"`
					CredentialIssuer string `json:"credential_issuer"`
					Key              string `json:"key"`
					Json             string `json:"json"`
				}{
					Format:           "JSON-LD",
					Name:             "CredentialName",
					Locale:           "en-US",
					Logo:             "https://example.com/logo.png",
					CredentialIssuer: "issuerID",
					Key:              "credKey",
					Json:             `{"format": "JSON-LD"}`,
				}
				_, err = db.Exec(`
					INSERT INTO credentials(
						format, name, locale, logo, credential_issuer, key, json, created, updated
					) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				`, initialCredential.Format, initialCredential.Name, initialCredential.Locale, initialCredential.Logo,
					initialCredential.CredentialIssuer, initialCredential.Key, initialCredential.Json)
				assert.NoError(t, err, "Expected no error inserting initial credential")
				tc.credential.Display[0].Name = "UpdatedCredentialName"
			}

			// Execute the activity
			issuerID := "issuerID"
			issuerName := "issuerName"
			credKey := "credKey"

			_, err = env.ExecuteActivity(StoreOrUpdateCredentialsActivity, issuerID, issuerName, credKey, dbPath, tc.credential)
			if tc.expectError {
				assert.Error(t, err, "Expected an error from StoreOrUpdateCredentialsActivity")
				return
			}

			// Check database row count
			var count int
			err = db.QueryRow(`SELECT COUNT(*) FROM credentials WHERE key = ?`, credKey).Scan(&count)
			assert.NoError(t, err, "Expected no error querying database")
			assert.Equal(t, tc.expectedRows, count, "Unexpected row count in credentials table")

			// Additional check for update test
			if tc.name == "Valid Data - Update" {
				// Verify that the credential is updated
				var updatedName string
				err = db.QueryRow(`SELECT name FROM credentials WHERE key = ?`, credKey).Scan(&updatedName)
				assert.NoError(t, err, "Expected no error querying updated credential name")
				assert.Equal(t, "UpdatedCredentialName", updatedName, "Credential name should be updated")
			}
		})
	}
}

func TestCleanupCredentialsActivity(t *testing.T) {
	var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
	err := json.Unmarshal([]byte(wellKnownJSON), &issuerData)
	assert.NoError(t, err, "Did not expect an error")

	testCases := []struct {
		name         string
		dbPath       string
		issuerID     string
		validKeys    []string
		expectError  bool
		expectedRows int
	}{
		{
			name:        "Valid Data - Cleanup",
			issuerID:    "issuerID",
			validKeys:   []string{"validKey1", "validKey2"},
			expectError: false,
			// Assuming only valid keys remain
			expectedRows: 2,
		},
		{
			name:         "Fail to Open DB",
			issuerID:     "issuerID",
			validKeys:    []string{"validKey1", "validKey2"},
			dbPath:       "/invalid/path/to/db.sqlite", // Simulate a non-writable or non-existent path
			expectError:  true,
			expectedRows: 0,
		},
		{
			name:        "No Keys to Cleanup",
			issuerID:    "issuerID",
			validKeys:   []string{}, // No valid keys
			expectError: false,
			// Assuming all records will be deleted (if they exist)
			expectedRows: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(CleanupCredentialsActivity)

			// Set up the database
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
				CREATE TABLE IF NOT EXISTS credentials (
					format TEXT,
					issuer_name TEXT,
					name TEXT,
					locale TEXT,
					logo TEXT,
					credential_issuer TEXT,
					key TEXT PRIMARY KEY,
					json TEXT,
					created TIMESTAMP,
					updated TIMESTAMP
				);
			`)
			if !tc.expectError {
				assert.NoError(t, err, "Expected no error creating test database schema")
			}

			// Insert some initial data
			if len(tc.validKeys) > 0 && !tc.expectError {
				for _, key := range tc.validKeys {
					_, err := db.Exec(`
						INSERT INTO credentials(format, issuer_name, name, locale, logo, credential_issuer, key, json, created, updated)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					`, "JSON-LD", "issuerName", "CredentialName", "en-US", "https://example.com/logo.png", tc.issuerID, key, `{"format": "JSON-LD"}`)
					assert.NoError(t, err, "Expected no error inserting initial credential")
				}
			}

			// Execute the activity
			_, err = env.ExecuteActivity(CleanupCredentialsActivity, tc.issuerID, dbPath, tc.validKeys)
			if tc.expectError {
				assert.Error(t, err, "Expected an error from CleanupCredentialsActivity")
				return
			}

			// Check database row count
			var count int
			err = db.QueryRow(`SELECT COUNT(*) FROM credentials WHERE credential_issuer = ?`, tc.issuerID).Scan(&count)
			assert.NoError(t, err, "Expected no error querying database")
			assert.Equal(t, tc.expectedRows, count, "Unexpected row count in credentials table")
		})
	}
}
