// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package OpenID4VP

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP/workflow"
	temporalclient "github.com/forkbombeu/didimo/pkg/internal/temporal_client"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
)

// OpenID4VPTestInputFile represents the structure of the JSON file
// containing both the variant and the form payload.
type OpenID4VPTestInputFile struct {
	Variant json.RawMessage `json:"variant"`
	Form    any             `json:"form"`
}

// startWorkflow starts the Temporal workflow
func StartWorkflow(input OpenID4VPTestInputFile, userMail, appURL string, namespace string) error {
	// Load environment variables.
	godotenv.Load()
	c, err := temporalclient.GetTemporalClientWithNamespace(namespace)

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
		AppURL:   appURL,
	}

	// Define workflow options.
	workflowOptions := client.StartWorkflowOptions{
		ID:        "OpenIDTestWorkflow" + uuid.NewString(),
		TaskQueue: workflow.OpenIDTestTaskQueue,
	}

	// Start the workflow execution.
	_, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.OpenIDTestWorkflow, workflowInput)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}

	return nil
}
