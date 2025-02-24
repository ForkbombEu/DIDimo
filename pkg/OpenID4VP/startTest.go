package OpenID4VP

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP/testdata"
	"github.com/forkbombeu/didimo/pkg/OpenID4VP/workflow"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
)

// OpenID4VPTestInputFile represents the structure of the JSON file
// containing both the variant and the form payload.
type OpenID4VPTestInputFile struct {
	Variant json.RawMessage `json:"variant"`
	Form    testdata.Form   `json:"form"`
}

// startWorkflow starts the Temporal workflow
func StartWorkflow(input OpenID4VPTestInputFile, userMail string) error {
	// Load environment variables.
	godotenv.Load()
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = "localhost:7233"
	}

	// Create a Temporal client.
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		return fmt.Errorf("unable to create client: %v", err)
	}
	defer c.Close()

	// Convert the variant (JSON) to a string.
	variantStr := string(input.Variant)

	// Prepare workflow input.
	workflowInput := workflow.WorkflowInput{
		Variant:  variantStr,
		Form:     input.Form,
		UserMail: userMail,
	}

	// Define workflow options.
	workflowOptions := client.StartWorkflowOptions{
		ID:        "OpenIDTestWorkflow" + uuid.NewString(),
		TaskQueue: "openid-test-task-queue",
	}

	// Start the workflow execution.
	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.OpenIDTestWorkflow, workflowInput)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}

	// Wait for the workflow result.
	var result string
	if err := workflowRun.Get(context.Background(), &result); err != nil {
		return fmt.Errorf("failed to get workflow result: %v", err)
	}

	// Print the result.
	fmt.Printf("Workflow result: %s\n", result)
	return nil
}
