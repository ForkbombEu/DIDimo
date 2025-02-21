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
		Use: "credimi [url]",
		Short: "\033[38;2;255;100;0m      dP oo       dP oo                     \033[0m\n" +
			"\033[38;2;255;71;43m      88          88                        \033[0m\n" +
			"\033[38;2;255;43;86m.d888b88 dP .d888b88 dP 88d8b.d8b. .d8888b. \033[0m\n" +
			"\033[38;2;255;14;129m88'  `88 88 88'  `88 88 88'`88'`88 88'  `88 \033[0m\n" +
			"\033[38;2;236;0;157m88.  .88 88 88.  .88 88 88  88  88 88.  .88 \033[0m\n" +
			"\033[38;2;197;0;171m`88888P8 dP `88888P8 dP dP  dP  dP `88888P' \033[0m\n" +
			"\033[38;2;159;0;186m                                             \033[0m\n" +
			"                   \033[48;2;0;0;139m\033[38;2;255;255;255m           :(){ :|:& };: \033[0m\n" +
			"                   \033[48;2;0;0;139m\033[38;2;255;255;255m by The Forkbomb Company \033[0m\n",
		Args: cobra.ExactArgs(1), // Ensure exactly one positional argument is provided
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
