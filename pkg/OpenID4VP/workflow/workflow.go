package workflow

import (
	"fmt"
	"os"
	"time"

	"github.com/forkbombeu/didimo/pkg/OPENID4VP/testdata"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowInput struct {
	Variant     string
	JSONPayload testdata.JSONPayload
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

	token := "jgpr3AEU5FnmaLhzvk53BBZShu/sGRKOsH9XBPuQpcEELqKoU63VOdh+piMtSshF9dRHYC+OXkjjHwRoLYJXzw=="

	// Create a temporary file to pass to GenerateYAML
	tempFile, err := os.CreateTemp("", "generated-*.yaml")
	if err != nil {
		logger.Error("Failed to create temporary file", "error", err)
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure the temp file is cleaned up after workflow execution

	// Pass the temporary file path to the GenerateYAML activity
	err = workflow.ExecuteActivity(ctx, GenerateYAMLActivity, input.Variant, input.JSONPayload, tempFile.Name()).Get(ctx, nil)
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
	err = workflow.ExecuteActivity(ctx, PrintQRCodeACtivity, response["result"]).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to print QR code to terminal", "error", err)
		return "", fmt.Errorf("failed to print QR code to terminal: %w", err)
	}

	logger.Info("QR Code printed successfully!")
	return "Worflow completed successfully", nil
}
