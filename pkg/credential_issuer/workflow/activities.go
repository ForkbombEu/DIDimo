package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
func StoreCredentialsActivity(ctx context.Context, issuerData *credentialissuer.OpenidCredentialIssuerSchemaJson, issuerID, dbPath string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Storing credential issuer data into database")

	if len(issuerData.Display) == 0 {
		logger.Warn("IssuerData.Display is empty, restarting from FetchCredentialIssuerActivity")
		return temporal.NewApplicationError("Invalid issuer data: no display information", "RestartFromFetch")
	}

	logger.Info("Using database path", "dbPath", dbPath)

	// Attempt to open the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to open database : %v, retrying StoreCredentialsActivity", err), "error", err)
		return temporal.NewApplicationError("Failed to open database", "RetryStoreCredentials", err)
	}
	defer db.Close()

	issuerName := issuerData.Display[0].Name

	for _, cred := range issuerData.CredentialConfigurationsSupported {
		credName, credLocale, credLogo := "", "", ""
		if len(cred.Display) > 0 {
			if cred.Display[0].Name != "" {
				credName = cred.Display[0].Name
			}
			if cred.Display[0].Locale != nil {
				credLocale = *cred.Display[0].Locale
			}
			if cred.Display[0].Logo != nil && cred.Display[0].Logo.Uri != "" {
				credLogo = cred.Display[0].Logo.Uri
			}
		}
		credJSON, err := json.Marshal(cred)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to marshal credential %s : %v, retrying StoreCredentialsActivity", credName, err), "error", err)
			return temporal.NewApplicationError(fmt.Sprintf("Failed to marshal credential %s", credName), "RetryStoreCredentials", err)
		}
		insertSQL := `
        INSERT INTO credentials(
            format,
            issuer_name, 
            name, 
            locale,
            logo,
			json,
			credential_issuer,
			created,
			updated
        ) VALUES (?, ?, ?, ?, ?, ?,?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

		_, err = db.ExecContext(ctx, insertSQL,
			cred.Format,
			issuerName,
			credName,
			credLocale,
			credLogo,
			credJSON,
			issuerID,
		)
		if err != nil {
			logger.Warn("SQL execution failed, restarting from FetchCredentialIssuerActivity", "error", err)
			return temporal.NewApplicationError(fmt.Sprintf("Database query failed: %v", err), "RestartFromFetch", err)
		}

		logger.Info("Successfully stored credential issuer data")
	}

	return nil
}

func getDBPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get caller info")
	}
	sourceDir := filepath.Dir(filename)
	dbPath := getEnv("DB_PATH", filepath.Join(sourceDir, "../../../pb_data/data.db"))
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatal("Failed to resolve absolute path:", err)
	}
	return absDBPath
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
