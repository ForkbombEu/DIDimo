package workflow

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

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
		result = ""
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

	signal := workflow.GetSignalChannel(ctx, "wallet-test-signal")
	var data SignalData
	signal.Receive(ctx, &data)

	// Process the signal data
	if !data.Success {
		return WorkflowResponse{Message: fmt.Sprintf("Workflow terminated with a failure message: %s", data.Reason)}, nil
	}

	return WorkflowResponse{Message: "Worflow completed successfully"}, nil
}
