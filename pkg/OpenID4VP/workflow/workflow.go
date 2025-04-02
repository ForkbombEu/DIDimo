package workflow

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// OpenIDTestWorkflow starts and waits for user input
func OpenIDTestWorkflow(ctx workflow.Context, input WorkflowInput) (WorkflowResponse, error) {
	logger := workflow.GetLogger(ctx)
	// Define retry policy
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 5,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    5,
	}

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 10,
		StartToCloseTimeout:    time.Minute * 5,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	token := os.Getenv("TOKEN")
	if token == "" {
		return WorkflowResponse{}, fmt.Errorf("TOKEN environment variable not set")
	}

	// Create a temporary file to pass to GenerateYAML
	tempFile, err := os.CreateTemp("", "generated-*.yaml")
	if err != nil {
		logger.Error("Failed to create temporary file", "error", err)
		return WorkflowResponse{}, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure the temp file is cleaned up after workflow execution
	YAMLInput := GenerateYAMLInput{
		Variant:  input.Variant,
		Form:     input.Form,
		FilePath: tempFile.Name(),
	}
	// Pass the temporary file path to the GenerateYAML activity
	err = workflow.ExecuteActivity(ctx, GenerateYAMLActivity, YAMLInput).Get(ctx, nil)
	if err != nil {
		logger.Error("GenerateYAML failed", "error", err)
		return WorkflowResponse{}, err
	}
	stepCIInput := StepCIRunnerInput{
		FilePath: tempFile.Name(),
		Token:    token,
	}

	var response StepCIRunnerResponse
	err = workflow.ExecuteActivity(ctx, RunStepCIJSProgramActivity, stepCIInput).Get(ctx, &response)
	if err != nil {
		logger.Error("RunStepCIJSProgram failed", "error", err)
		return WorkflowResponse{}, err
	}

	SMTPHost := os.Getenv("SMTP_HOST")
	if SMTPHost == "" {
		return WorkflowResponse{}, fmt.Errorf("SMTP_HOST environment variable not set")
	}
	SMTPPortString := os.Getenv("SMTP_PORT")
	if SMTPPortString == "" {
		return WorkflowResponse{}, fmt.Errorf("SMTP_PORT environment variable not set")
	}
	SMTPPort, err := strconv.Atoi(SMTPPortString)
	if err != nil {
		return WorkflowResponse{}, fmt.Errorf("SMTP_PORT environment variable not an integer")
	}
	username := os.Getenv("MAIL_USERNAME")

	password := os.Getenv("MAIL_USERNAME")

	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		return WorkflowResponse{}, fmt.Errorf("MAIL_SENDER environment variable not set")
	}
	baseURL := input.AppURL + "/tests/wallet"
	u, err := url.Parse(baseURL)
	if err != nil {
		return WorkflowResponse{}, fmt.Errorf("unexpected error parsing URL: %v", err)
	}

	result, ok := response.Result["result"].(string)
	if !ok {
		result = " "
	}

	query := u.Query()
	query.Set("workflow-id", workflow.GetInfo(ctx).WorkflowExecution.ID)
	query.Set("qr", result)
	u.RawQuery = query.Encode()

	finalURL := u.String()
	emailConfig := EmailConfig{
		SMTPHost:      SMTPHost,
		SMTPPort:      SMTPPort,
		Username:      username,
		Password:      password,
		SenderEmail:   sender,
		ReceiverEmail: input.UserMail,
		Subject:       "Test QR Code Email",
		Body: fmt.Sprintf(`
		<html>
			<body>
				<p>Here is your link:</p>
				<p><a href="%s" target="_blank" rel="noopener">%s</a></p>
			</body>
		</html>
	`, finalURL, finalURL),
	}

	err = workflow.ExecuteActivity(ctx, SendMailActivity, emailConfig).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send mail to user ", "error", err)
		return WorkflowResponse{}, fmt.Errorf("failed to print QR code to terminal: %w", err)
	}

	rid, ok := response.Result["rid"].(string)
	if !ok {
		return WorkflowResponse{}, fmt.Errorf("failed to get id from response: %v", err)
	}

	childCtx, cancelHandler := workflow.WithCancel(ctx)
	defer cancelHandler()

	childOptions := workflow.ChildWorkflowOptions{
		WorkflowID:        workflow.GetInfo(ctx).WorkflowExecution.ID + "-log",
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE,
	}
	childCtx = workflow.WithChildOptions(childCtx, childOptions)

	// Execute child workflow asynchronously
	logsWorkflow := workflow.ExecuteChildWorkflow(childCtx, LogSubWorkflow, LogWorkflowInput{
		AppURL:   input.AppURL,
		RID:      rid,
		Token:    token,
		Interval: time.Second * 30,
	})

	// Wait for either signal or child completion
	selector := workflow.NewSelector(ctx)
	var subWorkflowResponse LogWorkflowResponse
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
		return WorkflowResponse{Message: fmt.Sprintf("Workflow terminated with a failure message: %s", data.Reason), Logs: subWorkflowResponse.Logs}, nil
	}

	return WorkflowResponse{Message: "Workflow completed successfully", Logs: subWorkflowResponse.Logs}, nil
}

func LogSubWorkflow(ctx workflow.Context, input LogWorkflowInput) (LogWorkflowResponse, error) {
	logger := workflow.GetLogger(ctx)

	// Configure activity options for the sub-workflow
	subActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 10,
		StartToCloseTimeout:    time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	subCtx := workflow.WithActivityOptions(ctx, subActivityOptions)

	logInput := GetLogsActivityInput{
		BaseURL: LogsBaseURL,
		RID:     input.RID,
		Token:   input.Token,
	}

	var logs []map[string]any
	var triggerLogs bool // Flag to track if we should start TriggerLogsUpdateActivity

	signalChan := workflow.GetSignalChannel(ctx, "wallet-test-start-log-update")
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, _ bool) {
		triggerLogs = true
	})

	for {

		if ctx.Err() != nil {
			logger.Info("Workflow canceled, returning collected logs")
			return LogWorkflowResponse{Logs: logs}, nil
		}

		// Get the logs
		err := workflow.ExecuteActivity(subCtx, GetLogsActivity, logInput).Get(subCtx, &logs)
		if err != nil {
			logger.Error("Failed to get log", "error", err)
			return LogWorkflowResponse{}, err
		}

		// Check if we received a signal before running TriggerLogsUpdateActivity
		selector.Select(ctx)

		if triggerLogs {
			triggerLogsInput := TriggerLogsUpdateActivityInput{
				AppURL:     input.AppURL,
				Logs:       logs,
				WorkflowID: strings.TrimSuffix(workflow.GetInfo(ctx).WorkflowExecution.ID, "-log"),
			}
			err = workflow.ExecuteActivity(subCtx, TriggerLogsUpdateActivity, triggerLogsInput).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to send logs", "error", err)
				return LogWorkflowResponse{}, err
			}
			if len(logs) > 0 {
				if result, ok := logs[len(logs)-1]["result"].(string); ok {
					if result == "INTERRUPTED" || result == "FINISHED" {
						return LogWorkflowResponse{Logs: logs}, nil
					}
				}
			}
		}

		workflow.Sleep(subCtx, input.Interval)
	}
}
