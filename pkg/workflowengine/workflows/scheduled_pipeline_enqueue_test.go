// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflows

import (
	"context"
	"testing"
	"time"

	pipelineinternal "github.com/forkbombeu/credimi/pkg/internal/pipeline"
	"github.com/forkbombeu/credimi/pkg/workflowengine"
	"github.com/forkbombeu/credimi/pkg/workflowengine/activities"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
)

// TestScheduledPipelineEnqueueWorkflowEnqueuesDeterministicTicket verifies deterministic ticket IDs and enqueue input.
func TestScheduledPipelineEnqueueWorkflowEnqueuesDeterministicTicket(t *testing.T) {
	pipelineYAML := `
name: Scheduled Pipeline
steps:
  - id: step-1
    use: mobile-automation
    with:
      device_id: runner-b/device-b
      action_id: action-1
  - id: step-2
    use: mobile-automation
    with:
      device_id: runner-a/device-a
      action_id: action-2
`

	input := workflowengine.WorkflowInput{
		Payload: ScheduledPipelineEnqueueWorkflowInput{
			PipelineIdentifier:  "pipeline-123",
			OwnerNamespace:      "org-1",
			MaxPipelinesInQueue: 3,
		},
		Config: map[string]any{
			"app_url":          "https://example.test",
			"internal_app_url": "http://credimi:8090",
		},
	}

	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()

	httpAct := activities.NewHTTPActivity()
	env.RegisterActivityWithOptions(httpAct.Execute, activity.RegisterOptions{
		Name: httpAct.Name(),
	})
	internalHTTPAct := activities.NewInternalHTTPActivity()
	env.RegisterActivityWithOptions(internalHTTPAct.Execute, activity.RegisterOptions{
		Name: internalHTTPAct.Name(),
	})

	var capturedPayload activities.EnqueuePipelineRunTicketActivityInput
	var capturedHTTPURL string
	var capturedInternalHTTPURL string
	env.RegisterActivityWithOptions(
		func(_ context.Context, input workflowengine.ActivityInput) (workflowengine.ActivityResult, error) {
			payload, err := workflowengine.DecodePayload[activities.EnqueuePipelineRunTicketActivityInput](
				input.Payload,
			)
			require.NoError(t, err)
			capturedPayload = payload
			return workflowengine.ActivityResult{Output: map[string]any{"status": "queued"}}, nil
		},
		activity.RegisterOptions{Name: activities.EnqueuePipelineRunTicketActivityName},
	)

	env.OnActivity(httpAct.Name(), mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			input, ok := args.Get(1).(workflowengine.ActivityInput)
			require.True(t, ok)
			payload, err := workflowengine.DecodePayload[activities.HTTPActivityPayload](
				input.Payload,
			)
			require.NoError(t, err)
			capturedHTTPURL = payload.URL
		}).
		Return(workflowengine.ActivityResult{
			Output: map[string]any{
				"body": map[string]any{
					"record": map[string]any{
						"published": true,
						"yaml":      pipelineYAML,
					},
				},
			},
		}, nil)
	env.OnActivity(internalHTTPAct.Name(), mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			input, ok := args.Get(1).(workflowengine.ActivityInput)
			require.True(t, ok)
			payload, err := workflowengine.DecodePayload[activities.InternalHTTPActivityPayload](
				input.Payload,
			)
			require.NoError(t, err)
			capturedInternalHTTPURL = payload.URL
		}).
		Return(workflowengine.ActivityResult{
			Output: map[string]any{
				"body": map[string]any{"valid": true},
			},
		}, nil)

	w := NewScheduledPipelineEnqueueWorkflow()
	env.ExecuteWorkflow(w.Workflow, input)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	require.Equal(
		t,
		"sched/default-test-workflow-id/default-test-run-id",
		capturedPayload.TicketID,
	)
	require.Equal(t, "org-1", capturedPayload.OwnerNamespace)
	require.Equal(t, "pipeline-123", capturedPayload.PipelineIdentifier)
	require.Equal(t, pipelineYAML, capturedPayload.YAML)
	require.Equal(t, 3, capturedPayload.MaxPipelinesInQueue)
	require.ElementsMatch(
		t,
		[]string{"runner-a/device-a", "runner-b/device-b"},
		capturedPayload.DeviceIDs,
	)

	config := capturedPayload.PipelineConfig
	require.Equal(t, "org-1", config["namespace"])
	require.Equal(t, "https://example.test", config["app_url"])
	require.Equal(t, "http://credimi:8090", config["internal_app_url"])
	require.Equal(
		t,
		"http://credimi:8090/api/canonify/identifier/validate",
		capturedHTTPURL,
	)
	require.Equal(
		t,
		"http://credimi:8090/api/mobile-runner/validate-access",
		capturedInternalHTTPURL,
	)

	require.Equal(t, "pipeline-run", capturedPayload.Memo["test"])
	require.Equal(
		t,
		pipelineinternal.RunTypeScheduled,
		capturedPayload.Memo[pipelineinternal.RunTypeMemoKey],
	)

	require.False(t, capturedPayload.EnqueuedAt.IsZero())
	require.Equal(t, time.UTC, capturedPayload.EnqueuedAt.Location())
}

func TestCollectDeviceIDsAndNeedsGlobal(t *testing.T) {
	steps := []scheduledPipelineStep{
		{
			Use: "mobile-automation",
			With: map[string]any{
				"device_id": "runner-a",
			},
			OnError: []scheduledPipelineStep{
				{
					Use: "mobile-automation",
				},
			},
		},
	}

	deviceIDs := map[string]struct{}{}
	needsGlobal := false

	collectDeviceIDs(steps, deviceIDs, &needsGlobal)

	_, ok := deviceIDs["runner-a"]
	require.True(t, ok)
	require.True(t, needsGlobal)
}

func TestDeviceIDsWithGlobal(t *testing.T) {
	info := scheduledPipelineRunnerInfo{
		DeviceIDs:         []string{"runner-b"},
		NeedsGlobalRunner: true,
	}

	ids := deviceIDsWithGlobal(info, "runner-a")
	require.Equal(t, []string{"runner-a", "runner-b"}, ids)

	ids = deviceIDsWithGlobal(info, "runner-b")
	require.Equal(t, []string{"runner-b"}, ids)
}
