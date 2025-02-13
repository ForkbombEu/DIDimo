package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activity implementation
	env.OnActivity(FetchCredentialIssuerActivity, mock.Anything, mock.Anything).Return(nil, nil)
	env.OnActivity(StoreCredentialsActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	env.ExecuteWorkflow(CredentialWorkflow, WorkflowInput{BaseURL: "example@test.com"})

	var result string
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Credentials Workflow completed successfully for URL: example@test.com", result)
}
