package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
func StoreCredentialsActivity(ctx context.Context, input StoreCredentialsActivityInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Storing credential issuer data into database")

	if len(input.IssuerData.Display) == 0 {
		logger.Warn("IssuerData.Display is empty, restarting from FetchCredentialIssuerActivity")
		return temporal.NewApplicationError("Invalid issuer data: no display information", "RestartFromFetch")
	}

	logger.Info("Using database path", "dbPath", input.DBPath)

	// Attempt to open the database
	db, err := sql.Open("sqlite", input.DBPath)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to open database : %v, retrying StoreCredentialsActivity", err), "error", err)
		return temporal.NewApplicationError("Failed to open database", "RetryStoreCredentials", err)
	}
	defer db.Close()

	issuerName := input.IssuerData.Display[0].Name

	for _, cred := range input.IssuerData.CredentialConfigurationsSupported {
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
			input.IssuerID,
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
	dbPath := filepath.Join(sourceDir, "../../../pb_data/data.db")
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatal("Failed to resolve absolute path:", err)
	}
	return absDBPath
}

func FetchIssuersActivity(ctx context.Context) (FetchIssuersActivityResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", EbsiIssuersUrl, nil)
	if err != nil {
		return FetchIssuersActivityResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FetchIssuersActivityResponse{}, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchIssuersActivityResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchIssuersActivityResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var root ApiResponse
	if err = json.Unmarshal(body, &root); err != nil {
		return FetchIssuersActivityResponse{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	hrefs, err := extractHrefsFromApiResponse(root)
	if err != nil {
		return FetchIssuersActivityResponse{}, fmt.Errorf("failed to extract hrefs: %w", err)
	}

	return FetchIssuersActivityResponse{Issuers: hrefs}, nil
}

func extractHrefsFromApiResponse(root ApiResponse) ([]string, error) {
	var hrefs []string
	for _, item := range root.Items {
		hrefs = append(hrefs, item.Href)
	}
	return hrefs, nil
}

func CreateCredentialIssuersActivity(ctx context.Context, input CreateCredentialIssuersInput) error {
	db, err := sql.Open("sqlite", getDBPath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	for _, issuer := range input.Issuers {
		exists, err := checkIfCredentialIssuerExist(ctx, db, issuer)
		if err != nil {
			return fmt.Errorf("failed to check if issuer exists: %w", err)
		}
		if exists {
			continue
		}
		_, err = db.ExecContext(ctx, "INSERT INTO credential_issuers(url) VALUES (?)", issuer)
		if err != nil {
			return fmt.Errorf("failed to insert issuer into database: %w", err)
		}
	}

	return nil
}

func checkIfCredentialIssuerExist(ctx context.Context, db *sql.DB, url string) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM credential_issuer WHERE url = ?", url).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to query database: %w", err)
	}
	return count > 0, nil
}
