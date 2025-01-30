package workflow

import (
	"crypto/tls"
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

func TestFetchCredentialIssuerActivity(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()

	// Register the mock activity function
	env.RegisterActivity(FetchCredentialIssuerActivity)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, wellKnownJSON)
	}))
	defer server.Close()

	originalTransport := http.DefaultTransport

	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	defer func() { http.DefaultTransport = originalTransport }()

	val, err := env.ExecuteActivity(FetchCredentialIssuerActivity, server.URL)
	assert.NoError(t, err, "Expected no error from FetchCredentialIssuerActivity")

	// Retrieve result
	var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
	err = val.Get(&issuerData)
	assert.NoError(t, err, "Expected no error when retrieving activity result")

}

// TestStoreCredentialsActivity tests StoreCredentialsActivity using an in-memory SQLite database
func TestStoreCredentialsActivity(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(StoreCredentialsActivity)
	tempFile, err := os.CreateTemp("", "test_db_*.sqlite")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tempFile.Name()) // Ensure the file is deleted after the test

	// Open the SQLite database file
	db, err := sql.Open("sqlite", tempFile.Name())
	assert.NoError(t, err, "Expected no error when opening in-memory database")
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE credentials (
            format TEXT,
            issuer_name TEXT,
            name TEXT,
            locale TEXT,
            logo TEXT
        );
    `)
	assert.NoError(t, err, "Expected no error when creating test database schema")
	var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
	err = json.Unmarshal([]byte(wellKnownJSON), &issuerData)
	assert.NoError(t, err, "Expected no error when parse json")

	_, err = env.ExecuteActivity(StoreCredentialsActivity, &issuerData, tempFile.Name())
	assert.NoError(t, err, "Expected no error from StoreCredentialsActivity")

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM credentials`).Scan(&count)
	assert.NoError(t, err, "Expected no error when querying test database")
	assert.Equal(t, 1, count, "Expected exactly one credential stored in the database")
}
