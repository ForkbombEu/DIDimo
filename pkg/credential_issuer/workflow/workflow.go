package workflow

import (
	"errors"
	"time"

	credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// WorkflowInput defines the input for the Temporal workflow.
type WorkflowInput struct {
	BaseURL string // Base URL for the credential issuer
}

// CredimiWorkflow validates the schema, parses metadata, and prints it.
func CredentialWorkflow(ctx workflow.Context, input WorkflowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Credential Workflow")

	// Define retry policy
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 5,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    0, // Infinite retries
	}

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 10,
		StartToCloseTimeout:    time.Minute * 5,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for {
		var issuerData *credentialissuer.OpenidCredentialIssuerSchemaJson
		err := workflow.ExecuteActivity(ctx, FetchCredentialIssuerActivity, input.BaseURL).Get(ctx, &issuerData)
		if err != nil {
			logger.Error("FetchCredentialIssuerActivity failed, retrying", "error", err)
			continue // Restart from FetchCredentialIssuerActivity
		}

		for {
			err = workflow.ExecuteActivity(ctx, StoreCredentialsActivity, issuerData, getDBPath()).Get(ctx, nil)
			if err != nil {
				var appErr *temporal.ApplicationError
				if errors.As(err, &appErr) {
					errType := appErr.Type()
					logger.Warn("StoreCredentialsActivity failed", "errorType", errType, "error", err)

					if errType == "RetryStoreCredentials" {
						// Retry only StoreCredentialsActivity
						continue
					} else if errType == "RestartFromFetch" {
						// Restart from FetchCredentialIssuerActivity
						break
					}
				}
			}

			// If StoreCredentialsActivity succeeds, break the loop
			break
		}

		// If everything succeeds, exit workflow loop
		break
	}

	logger.Info("Credential Workflow completed successfully")
	return nil
}
