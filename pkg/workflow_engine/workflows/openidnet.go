// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflows

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	temporalclient "github.com/forkbombeu/didimo/pkg/internal/temporal_client"
	"github.com/forkbombeu/didimo/pkg/utils"
	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/activities"
)

type SignalData struct {
	Success bool
	Reason  string
}

const (
	OpenIDNetTaskQueue          = "OpenIDNetTaskQueue"
	OpenIDNetStepCITemplatePath = "pkg/workflow_engine/workflows/openidnet_config/stepci_wallet_template.yaml"
)

type OpenIDNetWorkflow struct{}

func (OpenIDNetWorkflow) Name() string {
	return "Conformance check on https://www.certification.openid.net"
}

func (OpenIDNetWorkflow) GetOptions() workflow.ActivityOptions {
	return ActivityOptions
}

func (w *OpenIDNetWorkflow) Workflow(
	ctx workflow.Context,
	input workflowengine.WorkflowInput,
) (workflowengine.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	ctx = workflow.WithActivityOptions(ctx, w.GetOptions())

	stepCIWorkflowActivity := activities.StepCIWorkflowActivity{}
	stepCIInput := workflowengine.ActivityInput{
		Payload: map[string]any{
			"variant": input.Payload["variant"],
			"form":    input.Payload["form"],
		},
		Config: map[string]string{
			"template": input.Config["template"].(string),
			"token":    utils.GetEnvironmentVariable("OPENIDNET_TOKEN", nil, true),
		},
	}
	var stepCIResult workflowengine.ActivityResult
	err := stepCIWorkflowActivity.Configure(context.Background(), &stepCIInput)
	if err != nil {
		logger.Error(" StepCI configure failed", "error", err)
		return workflowengine.WorkflowResult{}, err
	}
	err = workflow.ExecuteActivity(ctx, stepCIWorkflowActivity.Name(), stepCIInput).
		Get(ctx, &stepCIResult)
	if err != nil {
		logger.Error("StepCIExecution failed", "error", err)
		return workflowengine.WorkflowResult{}, err
	}
	result, ok := stepCIResult.Output.(map[string]any)["result"].(string)
	if !ok {
		result = ""
	}
	baseURL := input.Payload["app_url"].(string) + "/tests/wallet"
	u, err := url.Parse(baseURL)
	if err != nil {
		return workflowengine.WorkflowResult{}, fmt.Errorf("unexpected error parsing URL: %v", err)
	}
	query := u.Query()
	query.Set("workflow-id", workflow.GetInfo(ctx).WorkflowExecution.ID)
	query.Set("qr", result)
	u.RawQuery = query.Encode()
	emailActivity := activities.SendMailActivity{}

	emailInput := workflowengine.ActivityInput{
		Config: map[string]string{
			"recipient": input.Payload["user_mail"].(string),
		},
		Payload: map[string]any{
			"subject": "[CREDIMI] Action required to continue your conformance checks",
			"body": fmt.Sprintf(`
		<html>
			<body>
				<p>Please click on the following link:</p>
				<p><a href="%s" target="_blank" rel="noopener">%s</a></p>
			</body>
		</html>
	`, u.String(), u.String()),
		},
	}
	err = emailActivity.Configure(context.Background(), &emailInput)
	if err != nil {
		logger.Error("Email activity configure failed", "error", err)
		return workflowengine.WorkflowResult{}, err
	}
	err = workflow.ExecuteActivity(ctx, emailActivity.Name(), emailInput).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send mail to user ", "error", err)
		return workflowengine.WorkflowResult{}, err
	}

	rid, ok := stepCIResult.Output.(map[string]any)["rid"].(string)
	if !ok {
		rid = ""
	}

	childCtx, cancelHandler := workflow.WithCancel(ctx)
	defer cancelHandler()

	child := OpenIDNetLogsWorkflow{}
	childCtx = child.Configure(childCtx)
	// Execute child workflow asynchronously
	logsWorkflow := workflow.ExecuteChildWorkflow(
		childCtx,
		child.Name(),
		workflowengine.WorkflowInput{
			Payload: map[string]any{
				"rid":     rid,
				"token":   os.Getenv("OPENIDNET_TOKEN"),
				"app_url": input.Payload["app_url"].(string),
			},
			Config: map[string]any{
				"interval": time.Second,
			},
		},
	)

	// Wait for either signal or child completion
	selector := workflow.NewSelector(ctx)
	var subWorkflowResponse workflowengine.WorkflowResult
	var data SignalData

	selector.AddFuture(logsWorkflow, func(f workflow.Future) {
		f.Get(ctx, &subWorkflowResponse)
	})
	var signalSent bool
	signalChan := workflow.GetSignalChannel(ctx, "wallet-test-signal")
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, _ bool) {
		signalSent = true
		c.Receive(ctx, &data)
		cancelHandler()
		logsWorkflow.Get(ctx, &subWorkflowResponse)
	})
	for !signalSent {
		selector.Select(ctx)
	}

	// Process the signal data
	if !data.Success {
		return workflowengine.WorkflowResult{
			Message: fmt.Sprintf("Workflow terminated with a failure message: %s", data.Reason),
			Log:     subWorkflowResponse.Log,
		}, nil
	}

	return workflowengine.WorkflowResult{
		Message: "Workflow completed successfully",
		Log:     subWorkflowResponse.Log,
	}, nil
}

func (w *OpenIDNetWorkflow) Start(
	input workflowengine.WorkflowInput,
) (result workflowengine.WorkflowResult, err error) {
	// Load environment variables.
	godotenv.Load()
	namespace := "default"
	if input.Config["namespace"] != nil {
		namespace = input.Config["namespace"].(string)
	}
	c, err := temporalclient.GetTemporalClientWithNamespace(
		namespace,
	)
	if err != nil {
		return workflowengine.WorkflowResult{}, fmt.Errorf("unable to create client: %v", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "OpenIDTestWorkflow" + uuid.NewString(),
		TaskQueue: OpenIDNetTaskQueue,
	}
	if input.Config["memo"] != nil {
		workflowOptions.Memo = input.Config["memo"].(map[string]any)
	}

	// Start the workflow execution.
	_, err = c.ExecuteWorkflow(context.Background(), workflowOptions, w.Name(), input)
	if err != nil {
		return workflowengine.WorkflowResult{}, fmt.Errorf("failed to start workflow: %v", err)
	}

	return workflowengine.WorkflowResult{}, nil
}

type OpenIDNetLogsWorkflow struct{}

func (OpenIDNetLogsWorkflow) Name() string {
	return "Drain logs from https://www.certification.openid.net"
}

func (OpenIDNetLogsWorkflow) GetOptions() workflow.ActivityOptions {
	return ActivityOptions
}

func (w *OpenIDNetLogsWorkflow) Workflow(
	ctx workflow.Context,
	input workflowengine.WorkflowInput,
) (workflowengine.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	subCtx := workflow.WithActivityOptions(ctx, w.GetOptions())

	var logs []map[string]interface{}
	var timerCancel workflow.CancelFunc
	var isPolling bool
	var timerFuture workflow.Future

	signalChanStart := workflow.GetSignalChannel(ctx, "wallet-test-start-log-update")
	signalChanStop := workflow.GetSignalChannel(ctx, "wallet_test-stop-log-update")

	// Timer setup function
	startTimer := func() {
		timerCtx, cancel := workflow.WithCancel(ctx)
		timerCancel = cancel
		interval := time.Duration(input.Config["interval"].(float64)) * time.Nanosecond
		timerFuture = workflow.NewTimer(timerCtx, interval)
	}

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(signalChanStart, func(c workflow.ReceiveChannel, _ bool) {
		logger.Info("Received wallet-test-start-log-update signal")
		isPolling = true
	})
	selector.AddReceive(signalChanStop, func(c workflow.ReceiveChannel, _ bool) {
		logger.Info("Received wallet-test-stop-log-update signal")
		isPolling = false
		if timerCancel != nil {
			timerCancel()
			timerCancel = nil
			timerFuture = nil
		}
	})

	for {
		// If workflow context is canceled, exit
		if errors.Is(ctx.Err(), workflow.ErrCanceled) {
			logger.Info("Workflow canceled, returning collected logs")
			return workflowengine.WorkflowResult{Log: logs}, nil
		}

		if isPolling {
			startTimer()

			getLogsInput := workflowengine.ActivityInput{
				Config: map[string]string{
					"method": "GET",
					"url": fmt.Sprintf(
						"https://www.certification.openid.net/api/log/%s",
						url.PathEscape(input.Payload["rid"].(string)),
					),
				},
				Payload: map[string]any{
					"headers": map[string]any{
						"Authorization": fmt.Sprintf("Bearer %s", input.Payload["token"].(string)),
					},
					"query_params": map[string]any{
						"public": "false",
					},
				},
			}

			var httpActivity activities.HttpActivity
			var httpResponse workflowengine.ActivityResult

			// Execute the HTTP request to fetch logs
			err := workflow.ExecuteActivity(subCtx, httpActivity.Name(), getLogsInput).Get(subCtx, &httpResponse)
			if err != nil {
				logger.Error("Failed to get logs", "error", err)
				return workflowengine.WorkflowResult{}, err
			}

			logs = AsSliceOfMaps(httpResponse.Output.(map[string]any)["body"])

			// Prepare and send logs if polling is still active
			triggerLogsInput := workflowengine.ActivityInput{
				Config: map[string]string{
					"method": "POST",
					"url":    fmt.Sprintf("%s/wallet-test/send-log-update", input.Payload["app_url"].(string)),
				},
				Payload: map[string]any{
					"headers": map[string]any{
						"Content-Type": "application/json",
					},
					"body": map[string]any{
						"workflow_id": strings.TrimSuffix(workflow.GetInfo(ctx).WorkflowExecution.ID, "-log"),
						"logs":        logs,
					},
				},
			}

			err = workflow.ExecuteActivity(subCtx, httpActivity.Name(), triggerLogsInput).Get(subCtx, nil)
			if err != nil {
				logger.Error("Failed to send logs", "error", err)
			}
			// Check if we reached a terminal condition
			if len(logs) > 0 {
				if result, ok := logs[len(logs)-1]["result"].(string); ok {
					if result == "INTERRUPTED" || result == "FINISHED" {
						logger.Info("Workflow completed with terminal log result")
						return workflowengine.WorkflowResult{Log: logs}, nil
					}
				}
			}
			selector.AddFuture(timerFuture, func(f workflow.Future) {
				timerFuture = nil
			})
		}

		// Wait for either a signal or timer event

		selector.Select(ctx)
	}
}

func (w *OpenIDNetLogsWorkflow) Configure(ctx workflow.Context) workflow.Context {
	childOptions := workflow.ChildWorkflowOptions{
		WorkflowID:        workflow.GetInfo(ctx).WorkflowExecution.ID + "-log",
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE,
	}
	return workflow.WithChildOptions(ctx, childOptions)
}

func AsSliceOfMaps(val any) []map[string]any {
	if v, ok := val.([]map[string]any); ok {
		return v
	}
	if arr, ok := val.([]any); ok {
		res := make([]map[string]any, 0, len(arr))
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				res = append(res, m)
			}
		}
		return res
	}
	return nil
}
