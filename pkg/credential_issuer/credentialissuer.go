package credentialissuer

import (
	"fmt"
	"io"
	"strings"

	metadata "github.com/forkbombeu/didimo/pkg/internal/metadata"
)

var credentialIssuerEndpoint string = "/.well-known/openid-credential-issuer"

// FetchCredentialIssuer fetches the metadata for a credential issuer and handles specific errors, including 404 status.
func FetchCredentialIssuer(baseURL string) (*OpenidCredentialIssuerSchemaJson, error) {
	if strings.HasPrefix(baseURL, "http://") {
		// Return an error if the URL uses HTTP
		return nil, fmt.Errorf("HTTP is not supported, only HTTPS allowed for the credential issuer URL")
	}

	// Ensure the base URL is in a valid format (https)
	if !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	// Call FetchJSON from the metadata package
	issuerMetadata, err := metadata.FetchJSON[OpenidCredentialIssuerSchemaJson](baseURL, credentialIssuerEndpoint)
	if err != nil {
		switch e := err.(type) {
		case *metadata.NetworkError:
			// Handle network errors
			return nil, fmt.Errorf("network or URL issue when accessing credential issuer: %v", e)
		case *metadata.HTTPError:
			if e.StatusCode == 404 {
				return nil, fmt.Errorf("not a credential issuer: %v", e)
			}
			return nil, fmt.Errorf("unexpected HTTP error: %v", e)
		case *metadata.ParseError:
			return nil, fmt.Errorf("not conformant JSON returned by credential issuer: %v", e)
		default:
			// Generic error fallback
			return nil, fmt.Errorf("error fetching credential issuer metadata: %v", e)
		}
	}

	return issuerMetadata, nil
}

// PrintCredentialIssuer prints the credential issuer metadata to the provided writer.
func PrintCredentialIssuer(issuerData *OpenidCredentialIssuerSchemaJson, writer io.Writer) error {
	if err := metadata.PrintJSON[OpenidCredentialIssuerSchemaJson](*issuerData, writer); err != nil {
		return fmt.Errorf("failed to print credential issuer metadata: %v", err)
	}
	return nil
}
