package workflow

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	os.Setenv("TOKEN", "test_token")
	os.Setenv("SMTP_HOST", "test_host")
	os.Setenv("SMTP_PORT", "1000")
	os.Setenv("MAIL_SENDER", "test@example.org")
	// Mock activity implementation
	env.OnActivity(GenerateYAMLActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(RunStepCIJSProgramActivity, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	env.OnActivity(GenerateQRCodeActivity, mock.Anything, mock.Anything).Return("", nil)
	env.OnActivity(SendMailActivity, mock.Anything, mock.Anything).Return(nil)
	env.ExecuteWorkflow(OpenIDTestWorkflow, WorkflowInput{Variant: "test", UserMail: "user@example.org"})

	var result string
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Worflow completed successfully", result)
}
