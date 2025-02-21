package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP"
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

	rootCmd.AddCommand(OpenID4VPTestCommand())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func OpenID4VPTestCommand() *cobra.Command {
	var inputFile, userMail string

	cmd := &cobra.Command{
		Use:   "openid4vp-test",
		Short: "Start an OpenID Test workflow using a JSON input file and a user email flag",
		RunE: func(cmd *cobra.Command, args []string) error {

			file, err := os.Open(inputFile)
			if err != nil {
				return fmt.Errorf("error opening JSON file: %v", err)
			}
			defer file.Close()

			var input OpenID4VP.OpenID4VPTestInputFile
			if err := json.NewDecoder(file).Decode(&input); err != nil {
				return fmt.Errorf("error parsing JSON file: %v", err)
			}
			// Debug print: pretty-print the parsed JSON payload as JSON
			if jsonBytes, err := json.MarshalIndent(input.Form, "", "  "); err == nil {
				fmt.Printf("DEBUG: JSON Payload:\n%s\n", string(jsonBytes))
			} else {
				fmt.Printf("DEBUG: Error marshalling JSON payload: %v\n", err)
			}
			// Start the workflow using the parsed JSON and the user email.
			return OpenID4VP.StartWorkflow(input, userMail)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to JSON input file containing both variant and form")
	cmd.Flags().StringVarP(&userMail, "user-mail", "u", "", "User email for the workflow")
	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("user-mail")

	return cmd
}
