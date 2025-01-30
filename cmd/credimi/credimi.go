package main

import (
	"context"
	"fmt"
	"io"
	"os"

	crIss "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
)

//go:generate go run  github.com/atombender/go-jsonschema@latest   -p credentialissuer ../../schemas/openid-credential-issuer.schema.json -o ../../pkg/credential_issuer/openid-credential-issuer.schema.go
func main() {
	var outputFile string
	var startWorkflow bool

	rootCmd := &cobra.Command{
		Use:   "credimi [url]",
		Short: "Fetch and parse .well-known credential issuer metadata",
		Args:  cobra.ExactArgs(1), // Ensure exactly one positional argument is provided
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			if startWorkflow {
				// Start the Temporal workflow
				return runWorkflow(url)
			}

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
	rootCmd.Flags().BoolVarP(&startWorkflow, "workflow", "w", false, "Start the workflow to process metadata")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runWorkflow(url string) error {

	// Create Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		return fmt.Errorf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Set up the workflow input
	workflowInput := workflow.WorkflowInput{
		BaseURL: url,
	}

	// Execute the workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "credimi-workflow" + uuid.New().String(),
		TaskQueue: "CredimiTaskQueue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.CredentialWorkflow, workflowInput)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}

	// Wait for the workflow to complete and get the result
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		return fmt.Errorf("error running workflow: %v", err)
	}

	fmt.Printf("Workflow completed successfully: %s\n", result)
	return nil
}
