package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
	_ "modernc.org/sqlite/lib"

	credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// FetchCredentialIssuerActivity fetches credential issuer metadata from a URL.
func FetchCredentialIssuerActivity(ctx context.Context, baseURL string) (*credentialissuer.OpenidCredentialIssuerSchemaJson, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting FetchCredentialIssuerActivity", "baseURL", baseURL)

	issuerMetadata, err := credentialissuer.FetchCredentialIssuer(baseURL)
	if err != nil {
		logger.Warn("Error fetching credential issuer, retrying", "error", err)
		return nil, temporal.NewApplicationError(
			fmt.Sprintf("Error fetching credential issuer: %v", err),
			"RetryableError",
			err,
		)
	}

	logger.Info("Successfully fetched credential issuer metadata")
	return issuerMetadata, nil
}

// StoreOrUpdateCredentialsActivity inserts or updates credential issuer data in the database.
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
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to open database: %v", err),
			"RetryStoreCredentials",
			err,
		)
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
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to marshal credential JSON: %v", err),
			"RetryStoreCredentials",
			err,
		)
	}

	// Check if the credential exists
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM credentials WHERE key = ? AND credential_issuer = ?", credKey, issuerID).Scan(&count)
	if err != nil {
		logger.Warn("Failed to check existing credential", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Database query failed while checking existing credentials: %v", err),
			"RetryStoreCredentials",
			err,
		)
	}

	// UPSERT SQL query
	upsertSQL := `
	INSERT INTO credentials (
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
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(key, credential_issuer) DO UPDATE SET
		format = excluded.format,
		issuer_name = excluded.issuer_name,
		name = excluded.name,
		locale = excluded.locale,
		logo = excluded.logo,
		json = excluded.json,
		updated = CURRENT_TIMESTAMP;`

	// Execute the UPSERT query
	_, err = db.ExecContext(ctx, upsertSQL, credential.Format, issuerName, credName, credLocale, credLogo, credJSON, credKey, issuerID)
	if err != nil {
		logger.Warn("SQL UPSERT failed", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Database UPSERT failed: %v", err),
			"RestartFromFetch",
			err,
		)
	}

	logger.Info("Successfully stored or updated credential", "credKey", credKey)
	return nil
}

// CleanupCredentialsActivity removes outdated credentials from the database.
func CleanupCredentialsActivity(ctx context.Context, issuerID, dbPath string, validKeys []string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Cleaning up outdated credentials for issuer", "issuerID", issuerID)

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Warn("Failed to open database", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to open database: %v", err),
			"RetryCleanupCredentials",
			err,
		)
	}
	defer db.Close()

	// Fetch existing credential keys for this issuer
	existingKeys := make(map[string]bool)
	rows, err := db.QueryContext(ctx, "SELECT key FROM credentials WHERE credential_issuer = ?", issuerID)
	if err != nil {
		logger.Warn("Failed to fetch existing credential keys", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Database query failed while fetching existing keys: %v", err),
			"RetryCleanupCredentials",
			err,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			logger.Warn("Failed to scan credential key", "error", err)
			return temporal.NewApplicationError(
				fmt.Sprintf("Failed to scan credential key: %v", err),
				"RetryCleanupCredentials",
				err,
			)
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
				return temporal.NewApplicationError(
					fmt.Sprintf("Database delete failed: %v", err),
					"RetryCleanupCredentials",
					err,
				)
			}
			logger.Info("Deleted outdated credential", "credKey", key)
		}
	}

	logger.Info("Cleanup completed successfully")
	return nil
}
