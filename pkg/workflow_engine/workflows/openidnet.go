package workflows

import (
	"context"
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
	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/activities"
)

type SignalData struct {
	Success bool
	Reason  string
}

const (
	OpenIDTestTaskQueue         = "OpenIDTestTaskQueue"
	OpenIDNetStepCITemplatePath = "pkg/workflow_engine/workflows/openidnet_config/stepci_wallet_template.yaml"
)

type OpenIDNetWorkflow struct{}

func (OpenIDNetWorkflow) Name() string {
	return "OpenIDNetWorkflow"
}

func (w *OpenIDNetWorkflow) Workflow(
	ctx workflow.Context,
	input workflowengine.WorkflowInput,
) (workflowengine.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	ctx = workflow.WithActivityOptions(ctx, *input.ActivityOptions)

	stepCIWorkflowActivity := activities.StepCIWorkflowActivity{}
	stepCIInput := workflowengine.ActivityInput{
		Payload: map[string]any{
			"variant": input.Payload["variant"],
			"form":    input.Payload["form"],
		},
		Config: map[string]string{
			"template": input.Config["template"].(string),
			"token":    os.Getenv("TOKEN"),
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
			"subject": "Test QR Code Email",
			"body": fmt.Sprintf(`
		<html>
			<body>
				<p>Here is your link:</p>
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
				"token":   os.Getenv("TOKEN"),
				"app_url": input.Payload["app_url"].(string),
			},
			Config: map[string]any{
				"interval": time.Second,
			},
			ActivityOptions: input.ActivityOptions,
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
	c, err := temporalclient.GetTemporalClient()
	if err != nil {
		return workflowengine.WorkflowResult{}, fmt.Errorf("unable to create client: %v", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "OpenIDTestWorkflow" + uuid.NewString(),
		TaskQueue: OpenIDTestTaskQueue,
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
	return "OpenIDNetLogsChildWorkflow"
}

func (w *OpenIDNetLogsWorkflow) Workflow(
	ctx workflow.Context,
	input workflowengine.WorkflowInput,
) (workflowengine.WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	subCtx := workflow.WithActivityOptions(ctx, *input.ActivityOptions)

	GetLogsInput := workflowengine.ActivityInput{
		Config: map[string]string{
			"method": "GET",
			"url": fmt.Sprintf(
				"%s/%s",
				"https://www.certification.openid.net/api/log/",
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
	var logs []map[string]any

	signalChan := workflow.GetSignalChannel(ctx, "wallet-test-start-log-update")
	selector := workflow.NewSelector(ctx)
	var triggerEnabled bool // Persistent flag to enable activity executions

	// Timer setup
	var timerFuture workflow.Future
	startTimer := func() {
		timerCtx, _ := workflow.WithCancel(ctx)
		timerFuture = workflow.NewTimer(timerCtx, time.Duration(input.Config["interval"].(float64))*time.Nanosecond)
	}

	// Initialize the timer
	startTimer()

	for {
		if ctx.Err() != nil {
			logger.Info("Workflow canceled, returning collected logs")
			return workflowengine.WorkflowResult{Log: logs}, nil
		}
		var HTTPActivity activities.HTTPActivity
		var HTTPResponse workflowengine.ActivityResult
		// Fetch logs
		err := workflow.ExecuteActivity(subCtx, HTTPActivity.Name(), GetLogsInput).
			Get(subCtx, &HTTPResponse)
		if err != nil {
			logger.Error("Failed to get logs", "error", err)
			return workflowengine.WorkflowResult{}, err
		}

		logs = asSliceOfMaps(HTTPResponse.Output.(map[string]any)["body"])
		// Check termination condition
		if len(logs) > 0 {
			if result, ok := logs[len(logs)-1]["result"].(string); ok {
				if result == "INTERRUPTED" || result == "FINISHED" {
					return workflowengine.WorkflowResult{Log: logs}, nil
				}
			}
		}

		// Register events in the selector
		selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, _ bool) {
			var signalVal interface{}
			c.Receive(ctx, &signalVal)
			triggerEnabled = true // Enable activity execution
		})
		selector.AddFuture(timerFuture, func(f workflow.Future) {
			// Timer expired; reset the timer for the next interval
			startTimer()
		})

		// Wait for either a signal or the timer to expire
		selector.Select(ctx)

		// Execute TriggerLogsUpdateActivity if enabled
		if triggerEnabled {
			triggerLogsInput := workflowengine.ActivityInput{
				Config: map[string]string{
					"method": "POST",
					"url": fmt.Sprintf(
						"%s/%s",
						input.Payload["app_url"].(string),
						"wallet-test/send-log-update",
					),
				},
				Payload: map[string]any{
					"headers": map[string]any{
						"Content-Type": "application/json",
					},
					"body": map[string]any{
						"workflow_id": strings.TrimSuffix(
							workflow.GetInfo(ctx).WorkflowExecution.ID,
							"-log",
						),
						"logs": logs,
					},
				},
			}

			err = workflow.ExecuteActivity(subCtx, HTTPActivity.Name(), triggerLogsInput).
				Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to send logs", "error", err)
				return workflowengine.WorkflowResult{}, err
			}
		}
	}
}

func (w *OpenIDNetLogsWorkflow) Configure(ctx workflow.Context) workflow.Context {
	childOptions := workflow.ChildWorkflowOptions{
		WorkflowID:        workflow.GetInfo(ctx).WorkflowExecution.ID + "-log",
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE,
	}
	return workflow.WithChildOptions(ctx, childOptions)
}

func asSliceOfMaps(val any) []map[string]any {
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
