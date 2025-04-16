// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package workflows

import (
	"encoding/json"
	"testing"

	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/activities"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/workflows/credentials_config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	var CheckActivity activities.CheckCredentialsIssuerActivity
	env.RegisterActivityWithOptions(CheckActivity.Execute, activity.RegisterOptions{
		Name: CheckActivity.Name(),
	})
	var JSONActivity activities.ParseJSONActivity
	env.RegisterActivityWithOptions(JSONActivity.Execute, activity.RegisterOptions{
		Name: JSONActivity.Name(),
	})
	var HTTPActivity activities.HTTPActivity
	env.RegisterActivityWithOptions(HTTPActivity.Execute, activity.RegisterOptions{
		Name: HTTPActivity.Name(),
	})
	var credentialWorkflow CredentialsIssuersWorkflow
	var issuerData credentials_config.OpenidCredentialIssuerSchemaJson
	jsonBytes, err := json.Marshal(issuerData)
	require.NoError(t, err)
	var JSONOutput map[string]any
	err = json.Unmarshal(jsonBytes, &JSONOutput)
	require.NoError(t, err)
	// Mock activity implementation
	env.OnActivity(CheckActivity.Name(), mock.Anything, mock.Anything).Return(workflowengine.ActivityResult{Output: map[string]any{"rawJSON": "{json:}", "base_url": "testURL"}}, nil)
	env.OnActivity(JSONActivity.Name(), mock.Anything, mock.Anything).Return(workflowengine.ActivityResult{Output: JSONOutput}, nil)
	env.OnActivity(HTTPActivity.Name(), mock.Anything, mock.Anything).Return(workflowengine.ActivityResult{Output: map[string]any{"body": map[string]any{"key": "test-credential"}}}, nil)
	env.ExecuteWorkflow(credentialWorkflow.Workflow, workflowengine.WorkflowInput{
		Config: map[string]any{
			"app_url": "test.app",
		},
		Payload: map[string]any{
			"issuerID": "test_issuer",
			"base_url": "test_url",
		},
	})

	var result workflowengine.WorkflowResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Successfully retrieved and stored and update credentials", result.Message)
	require.Equal(t, map[string]any{"RemovedCredentials": []any{"test-credential"}}, result.Log)
}
