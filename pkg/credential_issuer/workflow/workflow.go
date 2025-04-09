// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflow

import (
	"errors"
	"fmt"
	"os"
	"time"

	credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CredentialWorkflow validates the schema, parses metadata, and prints it.
func CredentialWorkflow(ctx workflow.Context, input CredentialWorkflowInput) (CredentialWorkflowResponse, error) {
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
		return CredentialWorkflowResponse{Message: ""}, err
	}

	dbPath := os.Getenv("DATA_DB_PATH")
	var validKeys []string
	// Store credentials
	for credKey, credential := range issuerData.CredentialConfigurationsSupported {
		castedCredential := Credential(credential)
		activityInput := StoreCredentialsActivityInput{
			IssuerData: issuerData,
			IssuerID:   input.IssuerID,
			DBPath:     dbPath,
			CredKey:    credKey,
			IssuerName: *issuerData.Display[0].Name,
			Credential: castedCredential,
		}
		err := workflow.ExecuteActivity(ctx, StoreOrUpdateCredentialsActivity, activityInput).Get(ctx, nil)
		if err != nil {
			var appErr *temporal.ApplicationError
			if errors.As(err, &appErr) {
				errType := appErr.Type()
				logger.Warn("StoreCredentialsActivity failed", "errorType", errType, "error", err)

				if errType == "RestartFromFetch" {
					logger.Warn("Restarting workflow from FetchCredentialIssuerActivity")
					// return "", workflow.NewContinueAsNewError(ctx, CredentialWorkflow, input)
					return CredentialWorkflowResponse{Message: ""}, workflow.NewContinueAsNewError(ctx, CredentialWorkflow, input)
				}
			}
			return CredentialWorkflowResponse{Message: ""}, err
		}
		validKeys = append(validKeys, credKey)
	}
	err = workflow.ExecuteActivity(ctx, CleanupCredentialsActivity, input.IssuerID, dbPath, validKeys).Get(ctx, nil)
	if err != nil {
		logger.Error("FCleanupCredentialsActivity failed", "error", err)
		return CredentialWorkflowResponse{Message: ""}, err
	}

	successMessage := fmt.Sprintf("Credentials Workflow completed successfully for URL: %s", input.BaseURL)
	logger.Info(successMessage)
	return CredentialWorkflowResponse{Message: successMessage}, nil
}

func FetchIssuersWorkflow(ctx workflow.Context) error {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    500,
	}

	options := workflow.ActivityOptions{
		TaskQueue:           FetchIssuersTaskQueue,
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var response FetchIssuersActivityResponse

	err := workflow.ExecuteActivity(ctx, FetchIssuersActivity).Get(ctx, &response)

	if err != nil {
		return err
	}

	if len(response.Issuers) == 0 {
		return errors.New("no issuers found")
	}
	input := CreateCredentialIssuersInput{
		Issuers: response.Issuers,
		DBPath:  os.Getenv("DATA_DB_PATH"),
	}

	errCreate := workflow.ExecuteActivity(ctx, CreateCredentialIssuersActivity, input).Get(ctx, nil)

	if errCreate != nil {
		return errCreate
	}
	return nil
}
