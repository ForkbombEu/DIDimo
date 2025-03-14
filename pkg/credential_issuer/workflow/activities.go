package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	ctx context.Context, input StoreCredentialsActivityInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Storing or updating credential", "issuerID", input.IssuerID, "credKey", input.CredKey)

	db, err := sql.Open("sqlite", input.DBPath)
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
	if len(input.Credential.Display) > 0 {
		if input.Credential.Display[0].Name != "" {
			credName = input.Credential.Display[0].Name
		}
		if input.Credential.Display[0].Locale != nil {
			credLocale = *input.Credential.Display[0].Locale
		}
		if input.Credential.Display[0].Logo != nil && input.Credential.Display[0].Logo.Uri != "" {
			credLogo = input.Credential.Display[0].Logo.Uri
		}
	}
	credJSON, err := json.Marshal(input.Credential)
	if err != nil {
		logger.Warn("Failed to marshal credential JSON", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Failed to marshal credential JSON: %v", err),
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
	_, err = db.ExecContext(ctx, upsertSQL, input.Credential.Format, input.IssuerName, credName, credLocale, credLogo, credJSON, input.CredKey, input.IssuerID)
	if err != nil {
		logger.Warn("SQL UPSERT failed", "error", err)
		return temporal.NewApplicationError(
			fmt.Sprintf("Database UPSERT failed: %v", err),
			"RestartFromFetch",
			err,
		)
	}

	logger.Info("Successfully stored or updated credential", "credKey", input.CredKey)
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

func FetchIssuersActivity(ctx context.Context) (FetchIssuersActivityResponse, error) {
	// Start with offset 0.
	hrefs, err := fetchIssuersRecursive(ctx, 0)
	if err != nil {
		return FetchIssuersActivityResponse{}, err
	}
	return FetchIssuersActivityResponse{Issuers: hrefs}, nil
}

func fetchIssuersRecursive(ctx context.Context, after int) ([]string, error) {
	var url string
	if after > 0 {
		url = fmt.Sprintf("%s&page=%d", FidesIssuersUrl, after)
	} else {
		url = FidesIssuersUrl
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var root FidesResponse
	if err = json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	hrefs, err := extractHrefsFromApiResponse(root)
	if err != nil {
		return nil, fmt.Errorf("failed to extract hrefs: %w", err)
	}

	if root.Page.Number >= root.Page.TotalPages || len(hrefs) < 200 {
		return hrefs, nil
	}

	nextHrefs, err := fetchIssuersRecursive(ctx, after+1)
	if err != nil {
		return nil, err
	}

	return append(hrefs, nextHrefs...), nil
}

func extractHrefsFromApiResponse(root FidesResponse) ([]string, error) {
	var hrefs []string
	for _, item := range root.Content {
		trimmedHref := RemoveWellKnownSuffix(item.IssuanceURL)
		hrefs = append(hrefs, trimmedHref)
	}
	return hrefs, nil
}

func CreateCredentialIssuersActivity(ctx context.Context, input CreateCredentialIssuersInput) error {
	db, err := sql.Open("sqlite", input.DBPath)
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
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM credential_issuers WHERE url = ?", url).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to query database: %w", err)
	}
	return count > 0, nil
}

func RemoveWellKnownSuffix(urlStr string) string {
	const suffix = "/.well-known/openid-credential-issuer"
	if strings.HasSuffix(urlStr, suffix) {
		return strings.TrimSuffix(urlStr, suffix)
	}
	return urlStr
}
