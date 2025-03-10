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
		MaximumAttempts:    5,
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
	var validKeys []string
	// Store credentials
	for credKey, credential := range issuerData.CredentialConfigurationsSupported {
		err := workflow.ExecuteActivity(ctx, StoreOrUpdateCredentialsActivity, input.IssuerID, issuerData.Display[0].Name, credKey, dbPath, credential).Get(ctx, nil)
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
		validKeys = append(validKeys, credKey)
	}
	err = workflow.ExecuteActivity(ctx, CleanupCredentialsActivity, input.IssuerID, dbPath, validKeys).Get(ctx, nil)
	if err != nil {
		logger.Error("FCleanupCredentialsActivity failed", "error", err)
		return "", err
	}

	successMessage := fmt.Sprintf("Credentials Workflow completed successfully for URL: %s", input.BaseURL)
	logger.Info(successMessage)
	return successMessage, nil
}
