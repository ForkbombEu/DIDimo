// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflow

import (
	"errors"
	"testing"

	credentialissuer "github.com/forkbombeu/credimi/pkg/credential_issuer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	var issuerData credentialissuer.OpenidCredentialIssuerSchemaJson
	// Mock activity implementation
	env.OnActivity(FetchCredentialIssuerActivity, mock.Anything, mock.Anything).Return(&issuerData, nil)
	credential := Credential(issuerData.CredentialConfigurationsSupported["discount_from_voucher"])
	inputStoreOrUpdate := StoreCredentialsActivityInput{
		IssuerData: &issuerData,
		IssuerID:   mock.Anything,
		DBPath:     mock.Anything,
		CredKey:    mock.Anything,
		IssuerName: mock.Anything,
		Credential: credential,
	}
	env.OnActivity(StoreOrUpdateCredentialsActivity, inputStoreOrUpdate).Return(nil)
	env.OnActivity(CleanupCredentialsActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	env.ExecuteWorkflow(CredentialWorkflow, CredentialWorkflowInput{BaseURL: "example@test.com"})

	var result CredentialWorkflowResponse

	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Credentials Workflow completed successfully for URL: example@test.com", result.Message)
}

func Test_SuccessfulFetchIssuersWorkflows(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	issuers := []string{"issuer1", "issuer2", "issuer3"}

	env.OnActivity(FetchIssuersActivity, mock.Anything).Return(FetchIssuersActivityResponse{
		Issuers: issuers,
	}, nil)
	env.OnActivity(CreateCredentialIssuersActivity, mock.Anything, CreateCredentialIssuersInput{
		Issuers: issuers,
		DBPath:  mock.Anything,
	}).Return(nil)
	env.ExecuteWorkflow(FetchIssuersWorkflow)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_UnsuccessfulFetchIssuersWorkflows(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.OnActivity(FetchIssuersActivity, mock.Anything).Return(FetchIssuersActivityResponse{}, errors.New("error"))
	env.ExecuteWorkflow(FetchIssuersWorkflow)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
