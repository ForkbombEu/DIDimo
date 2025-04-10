// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package OpenID4VP

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/forkbombeu/credimi/pkg/OpenID4VP/workflow"
	temporalclient "github.com/forkbombeu/credimi/pkg/internal/temporal_client"
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

func startWorkflow(input OpenID4VPTestInputFile, userMail, appURL string, namespace string, memo map[string]interface{}) error {
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
		Memo:      memo,
	}

	// Start the workflow execution.
	_, err = c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.OpenIDTestWorkflow, workflowInput)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}

	return nil
}

func StartWorkflow(input OpenID4VPTestInputFile, userMail, appURL string) error {
	if err := startWorkflow(input, userMail, appURL, "default", map[string]interface{}{}); err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}

	return nil
}

func StartWorkflowWithNamespace(input OpenID4VPTestInputFile, userMail, appURL, namespace string) error {
	if err := startWorkflow(input, userMail, appURL, namespace, map[string]interface{}{}); err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}
	return nil
}

func StartWorkflowWithNamespaceAndMemo(input OpenID4VPTestInputFile, userMail, appURL, namespace string, memo map[string]interface{}) error {
	if err := startWorkflow(input, userMail, appURL, namespace, memo); err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}
	return nil
}

func StartWorkflowWithMemo(input OpenID4VPTestInputFile, userMail, appURL string, memo map[string]interface{}) error {
	if err := startWorkflow(input, userMail, appURL, "default", memo); err != nil {
		return fmt.Errorf("failed to start workflow: %v", err)
	}
	return nil
}
