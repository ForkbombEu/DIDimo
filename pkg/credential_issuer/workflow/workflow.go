package workflow

import (
	"errors"
	"fmt"
	"time"

	credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// WorkflowInput defines the input for the Temporal workflow.
type WorkflowInput struct {
	BaseURL  string // Base URL for the credential issuer
	IssuerID string // ID of the credentials issuer from PB
}

// CredentialWorkflow validates the schema, parses metadata, and prints it.
func CredentialWorkflow(ctx workflow.Context, input WorkflowInput) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Credential Workflow", "BaseURL", input.BaseURL)

	// Define retry policy
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 5,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    0,
	}

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 10,
		StartToCloseTimeout:    time.Minute * 5,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Fetch credential issuer
	var issuerData *credentialissuer.OpenidCredentialIssuerSchemaJson
	err := workflow.ExecuteActivity(ctx, FetchCredentialIssuerActivity, input.BaseURL).Get(ctx, &issuerData)
	if err != nil {
		logger.Error("FetchCredentialIssuerActivity failed", "error", err)
		return "", err
	}

	dbPath := getDBPath()

	// Store credentials
	err = workflow.ExecuteActivity(ctx, StoreCredentialsActivity, issuerData, input.IssuerID, dbPath).Get(ctx, nil)
	if err != nil {
		var appErr *temporal.ApplicationError
		if errors.As(err, &appErr) {
			errType := appErr.Type()
			logger.Warn("StoreCredentialsActivity failed", "errorType", errType, "error", err)

			if errType == "RestartFromFetch" {
				logger.Warn("Restarting workflow from FetchCredentialIssuerActivity")
				return "", workflow.NewContinueAsNewError(ctx, CredentialWorkflow, input)
			}
		}
		return "", err
	}

	successMessage := fmt.Sprintf("Credentials Workflow completed successfully for URL: %s", input.BaseURL)
	logger.Info(successMessage)
	return successMessage, nil
}
