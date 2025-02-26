package workflow

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP/testdata"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowInput struct {
	Variant  string
	Form     testdata.Form
	UserMail string
}

// OpenIDTestWorkflow starts and waits for user input
func OpenIDTestWorkflow(ctx workflow.Context, input WorkflowInput) (string, error) {
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
		return "", fmt.Errorf("TOKEN environment variable not set")
	}

	// Create a temporary file to pass to GenerateYAML
	tempFile, err := os.CreateTemp("", "generated-*.yaml")
	if err != nil {
		logger.Error("Failed to create temporary file", "error", err)
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure the temp file is cleaned up after workflow execution

	// Pass the temporary file path to the GenerateYAML activity
	err = workflow.ExecuteActivity(ctx, GenerateYAMLActivity, input.Variant, input.Form, tempFile.Name()).Get(ctx, nil)
	if err != nil {
		logger.Error("GenerateYAML failed", "error", err)
		return "", err
	}

	var response map[string]string
	err = workflow.ExecuteActivity(ctx, RunStepCIJSProgramActivity, tempFile.Name(), token).Get(ctx, &response)
	if err != nil {
		logger.Error("RunStepCIJSProgram failed", "error", err)
		return "", err
	}
	var qrBase64 string
	err = workflow.ExecuteActivity(ctx, GenerateQRCodeActivity, response["result"]).Get(ctx, &qrBase64)
	if err != nil {
		logger.Error("Failed to generate QR code", "error", err)
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	qrBytes, err := base64.StdEncoding.DecodeString(qrBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding QR code: %w", err)
	}

	SMTPHost := os.Getenv("SMTP_HOST")
	if SMTPHost == "" {
		return "", fmt.Errorf("SMTP_HOST environment variable not set")
	}
	SMTPPortString := os.Getenv("SMTP_PORT")
	if SMTPPortString == "" {
		return "", fmt.Errorf("SMTP_PORT environment variable not set")
	}
	SMTPPort, err := strconv.Atoi(SMTPPortString)
	if err != nil {
		return "", fmt.Errorf("SMTP_PORT environment variable not an integer")
	}
	username := os.Getenv("MAIL_USERNAME")

	password := os.Getenv("MAIL_USERNAME")

	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		return "", fmt.Errorf("MAIL_SENDER environment variable not set")
	}
	emailConfig := EmailConfig{
		SMTPHost:      SMTPHost,
		SMTPPort:      SMTPPort,
		Username:      username,
		Password:      password,
		SenderEmail:   sender,
		ReceiverEmail: input.UserMail,
		Subject:       "Test QR Code Email",
		Body:          "Here is your QR code from Mailpit.",
		Attachments: map[string][]byte{
			"qrcode.png": qrBytes, // Attach qrPNG as a gomail.FileSetting
		},
	}
	err = workflow.ExecuteActivity(ctx, SendMailActivity, emailConfig).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send mail to user ", "error", err)
		return "", fmt.Errorf("failed to print QR code to terminal: %w", err)
	}

	return "Worflow completed successfully", nil
}
