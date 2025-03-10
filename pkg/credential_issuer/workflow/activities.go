package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
	_ "modernc.org/sqlite/lib"

	credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// FetchCredentialIssuerActivity wraps the FetchCredentialIssuer function and manages error categorization.
// FetchCredentialIssuerActivity fetches credential issuer metadata from a URL.
func FetchCredentialIssuerActivity(ctx context.Context, baseURL string) (*credentialissuer.OpenidCredentialIssuerSchemaJson, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting FetchCredentialIssuerActivity", "baseURL", baseURL)

	issuerMetadata, err := credentialissuer.FetchCredentialIssuer(baseURL)
	if err != nil {
		logger.Warn("Error fetching credential issuer, retrying", "error", err)
		// Always retry the activity
		return nil, temporal.NewApplicationError(fmt.Sprintf("Error fetching credential issuer: %v. Retry", err), "RetryableError", err)
	}

	logger.Info("Successfully fetched credential issuer metadata")
	return issuerMetadata, nil
}

// StoreCredentialsActivity inserts credential issuer data into the Pocketbase database
func StoreOrUpdateCredentialsActivity(
	ctx context.Context,
	issuerID, issuerName, credKey, dbPath string,
	credential struct {
		CredentialDefinition                 *credentialissuer.CredentialDefinition                      `json:"credential_definition,omitempty"`
		CredentialSigningAlgValuesSupported  []credentialissuer.CredentialSigningAlgValuesSupportedElem  `json:"credential_signing_alg_values_supported,omitempty"`
		CryptographicBindingMethodsSupported []credentialissuer.CryptographicBindingMethodsSupportedElem `json:"cryptographic_binding_methods_supported,omitempty"`
		Display                              []credentialissuer.DisplayElem_1                            `json:"display,omitempty"`
		Format                               string                                                      `json:"format"`
		ProofTypesSupported                  credentialissuer.ProofTypesSupported                        `json:"proof_types_supported,omitempty"`
		Scope                                *string                                                     `json:"scope,omitempty"`
	},
) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Storing or updating credential", "issuerID", issuerID, "credKey", credKey)

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Warn("Failed to open database", "error", err)
		return temporal.NewApplicationError("Failed to open database", "RetryStoreCredentials", err)
	}
	defer db.Close()

	// Extract credential details
	credName, credLocale, credLogo := "", "", ""
	if len(credential.Display) > 0 {
		if credential.Display[0].Name != "" {
			credName = credential.Display[0].Name
		}
		if credential.Display[0].Locale != nil {
			credLocale = *credential.Display[0].Locale
		}
		if credential.Display[0].Logo != nil && credential.Display[0].Logo.Uri != "" {
			credLogo = credential.Display[0].Logo.Uri
		}
	}
	credJSON, err := json.Marshal(credential)
	if err != nil {
		logger.Warn("Failed to marshal credential JSON", "error", err)
		return temporal.NewApplicationError("Failed to marshal credential JSON", "RetryStoreCredentials", err)
	}

	// Check if the credential exists
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM credentials WHERE key = ? AND credential_issuer = ?", credKey, issuerID).Scan(&count)
	if err != nil {
		logger.Warn("Failed to check existing credential", "error", err)
		return temporal.NewApplicationError("Database query failed while checking existing credentials", "RetryStoreCredentials", err)
	}

	if count > 0 {
		// Update existing credential
		updateSQL := `
        UPDATE credentials 
        SET format = ?, name = ?, locale = ?, logo = ?, json = ?, updated = CURRENT_TIMESTAMP 
        WHERE key = ? AND credential_issuer = ?;`
		_, err = db.ExecContext(ctx, updateSQL,
			credential.Format,
			credName,
			credLocale,
			credLogo,
			credJSON,
			credKey,
			issuerID,
		)
		if err != nil {
			logger.Warn("SQL update failed", "error", err)
			return temporal.NewApplicationError("Database update failed", "RestartFromFetch", err)
		}
		logger.Info("Updated existing credential", "credKey", credKey)
	} else {
		// Insert new credential
		insertSQL := `
        INSERT INTO credentials(
            format,
			issuer_name,
            name, 
            locale,
            logo,
            json,
            key,
            credential_issuer,
            created,
            updated
        ) VALUES (?,?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`
		_, err = db.ExecContext(ctx, insertSQL,
			credential.Format,
			issuerName,
			credName,
			credLocale,
			credLogo,
			credJSON,
			credKey,
			issuerID,
		)
		if err != nil {
			logger.Warn("SQL insert failed", "error", err)
			return temporal.NewApplicationError("Database insert failed", "RestartFromFetch", err)
		}
		logger.Info("Inserted new credential", "credKey", credKey)
	}

	logger.Info("Successfully stored or updated credential", "credKey", credKey)
	return nil
}

func CleanupCredentialsActivity(ctx context.Context, issuerID, dbPath string, validKeys []string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Cleaning up outdated credentials for issuer", "issuerID", issuerID)

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to open database: %v", err), "error", err)
		return temporal.NewApplicationError("Failed to open database", "RetryCleanupCredentials", err)
	}
	defer db.Close()

	// Fetch existing credential keys for this issuer
	existingKeys := make(map[string]bool)
	rows, err := db.QueryContext(ctx, "SELECT key FROM credentials WHERE credential_issuer = ?", issuerID)
	if err != nil {
		logger.Warn("Failed to fetch existing credential keys", "error", err)
		return temporal.NewApplicationError("Database query failed while fetching existing keys", "RetryCleanupCredentials", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			logger.Warn("Failed to scan credential key", "error", err)
			return temporal.NewApplicationError("Failed to scan credential key", "RetryCleanupCredentials", err)
		}
		existingKeys[key] = true
	}

	// Convert valid keys to a set for quick lookup
	validKeysSet := make(map[string]bool)
	for _, key := range validKeys {
		validKeysSet[key] = true
	}

	// Delete credentials that are no longer valid
	for key := range existingKeys {
		if !validKeysSet[key] {
			_, err := db.ExecContext(ctx, "DELETE FROM credentials WHERE key = ? AND credential_issuer = ?", key, issuerID)
			if err != nil {
				logger.Warn("SQL delete failed", "error", err)
				return temporal.NewApplicationError(fmt.Sprintf("Database delete failed: %v", err), "RetryCleanupCredentials", err)
			}
			logger.Info("Deleted old credential", "credKey", key)
		}
	}

	logger.Info("Cleanup completed successfully")
	return nil
}

func getDBPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get caller info")
	}
	sourceDir := filepath.Dir(filename)
	dbPath := filepath.Join(sourceDir, "../../../pb_data/data.db")
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatal("Failed to resolve absolute path:", err)
	}
	return absDBPath
}
