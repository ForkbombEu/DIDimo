package workflow

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activity implementation
	env.OnActivity(FetchCredentialIssuerActivity, mock.Anything, mock.Anything).Return(nil, nil)
	env.OnActivity(StoreCredentialsActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	env.ExecuteWorkflow(CredentialWorkflow, CredentialWorkflowInput{BaseURL: "example@test.com"})

	var result string
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Credentials Workflow completed successfully for URL: example@test.com", result)
}

func Test_SuccessfulFetchIssuersWorkflows(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.OnActivity(FetchIssuersActivity, mock.Anything).Return(FetchIssuersActivityResponse{}, nil)
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
