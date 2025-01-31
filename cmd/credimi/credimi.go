package main

import (
	"fmt"
	"io"
	"os"

	crIss "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"github.com/spf13/cobra"
)

//go:generate go run  github.com/atombender/go-jsonschema@latest   -p credentialissuer ../../schemas/openid-credential-issuer.schema.json -o ../../pkg/credential_issuer/openid-credential-issuer.schema.go
func main() {
	var outputFile string

	rootCmd := &cobra.Command{
		Use:   "credimi [url]",
		Short: "Fetch and parse .well-known credential issuer metadata",
		Args:  cobra.ExactArgs(1), // Ensure exactly one positional argument is provided
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			// Fetch metadata
			issuerMetadata, err := crIss.FetchCredentialIssuer(url)
			if err != nil {
				return fmt.Errorf("error fetching metadata: %v", err)
			}

			var writer io.Writer
			switch outputFile {
			case "", "stdout":
				writer = os.Stdout
			default:
				file, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("error creating file: %v", err)
				}
				defer file.Close()
				writer = file
			}

			// Output metadata
			if err := crIss.PrintCredentialIssuer(issuerMetadata, writer); err != nil {
				return fmt.Errorf("error writing metadata: %v", err)
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output destination (e.g., stdout or file path)")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
