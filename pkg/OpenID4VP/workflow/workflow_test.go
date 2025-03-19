package workflow

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.SetTestTimeout(10 * time.Minute)
	os.Setenv("TOKEN", "test_token")
	os.Setenv("SMTP_HOST", "test_host")
	os.Setenv("SMTP_PORT", "1000")
	os.Setenv("MAIL_SENDER", "test@example.org")
	// Mock activity implementation
	env.OnActivity(GenerateYAMLActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(RunStepCIJSProgramActivity, mock.Anything, mock.Anything).Return(StepCIRunnerResponse{}, nil)
	env.OnActivity(SendMailActivity, mock.Anything, mock.Anything).Return(nil)
	fakeData := SignalData{
		Success: false,
		Reason:  "Test message",
	}
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("wallet-test-signal", fakeData)
	}, time.Minute)
	env.ExecuteWorkflow(OpenIDTestWorkflow, WorkflowInput{Variant: "test", UserMail: "user@example.org"})

	var result WorkflowResponse
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Workflow terminated with a failure message: Test message", result.Message)
}
