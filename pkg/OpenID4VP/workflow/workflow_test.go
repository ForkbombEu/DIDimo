package workflow

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflows(t *testing.T) {
	testCases := []struct {
		name           string
		mockActivities func(env *testsuite.TestWorkflowEnvironment)
		signalDelay    time.Duration
		signalData     SignalData
		expectedMsg    string
	}{
		{
			name: "Signal before child completes",
			mockActivities: func(env *testsuite.TestWorkflowEnvironment) {
				env.OnActivity(GenerateYAMLActivity, mock.Anything, mock.Anything).Return(nil)
				env.OnActivity(RunStepCIJSProgramActivity, mock.Anything, mock.Anything).
					Return(StepCIRunnerResponse{Result: map[string]interface{}{"rid": "12345"}}, nil)
				env.OnActivity(SendMailActivity, mock.Anything, mock.Anything).Return(nil)
				env.OnActivity(GetLogsActivity, mock.Anything, mock.Anything).
					Return([]map[string]interface{}{{"result": "RUNNING"}}, nil)
			},
			signalDelay: time.Minute,
			signalData:  SignalData{Success: false, Reason: "Test failure"},
			expectedMsg: "Workflow terminated with a failure message: Test failure",
		},
		{
			name: "Child terminates before signal",
			mockActivities: func(env *testsuite.TestWorkflowEnvironment) {
				env.OnActivity(GenerateYAMLActivity, mock.Anything, mock.Anything).Return(nil)
				env.OnActivity(RunStepCIJSProgramActivity, mock.Anything, mock.Anything).
					Return(StepCIRunnerResponse{Result: map[string]interface{}{"rid": "12345"}}, nil)
				env.OnActivity(SendMailActivity, mock.Anything, mock.Anything).Return(nil)
				env.OnActivity(GetLogsActivity, mock.Anything, mock.Anything).
					Return([]map[string]interface{}{{"result": "FINISHED"}}, nil)
			},
			signalDelay: 2 * time.Minute,
			signalData:  SignalData{Success: true},
			expectedMsg: "Workflow completed successfully",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestWorkflowEnvironment()
			env.SetTestTimeout(10 * time.Minute)

			// Set environment variables
			os.Setenv("TOKEN", "test_token")
			os.Setenv("SMTP_HOST", "test_host")
			os.Setenv("SMTP_PORT", "1000")
			os.Setenv("MAIL_SENDER", "test@example.org")

			tc.mockActivities(env)

			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("wallet-test-signal", tc.signalData)
			}, tc.signalDelay)
			env.RegisterWorkflow(LogSubWorkflow)
			// Execute workflow
			env.ExecuteWorkflow(OpenIDTestWorkflow, WorkflowInput{Variant: "test", UserMail: "user@example.org"})

			var result WorkflowResponse
			assert.NoError(t, env.GetWorkflowResult(&result))
			assert.Equal(t, tc.expectedMsg, result.Message)
		})
	}
}
func Test_LogSubWorkflow(t *testing.T) {
	testCases := []struct {
		name          string
		mockResponse  []map[string]interface{}
		expectRunning bool
	}{
		{
			name:          "Workflow completes when result is FINISHED",
			mockResponse:  []map[string]interface{}{{"result": "FINISHED"}},
			expectRunning: false,
		},
		{
			name:          "Workflow runs indefinitely when result is RUNNING",
			mockResponse:  []map[string]interface{}{{"result": "RUNNING"}},
			expectRunning: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestWorkflowEnvironment()

			callCount := 0
			env.OnActivity(GetLogsActivity, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					callCount++
				}).
				Return(tc.mockResponse, nil)

			done := make(chan struct{})
			go func() {
				env.ExecuteWorkflow(LogSubWorkflow, LogWorkflowInput{
					RID:      "test-rid",
					Token:    "test-token",
					Interval: time.Second * 10,
				})
				close(done)
			}()

			if tc.expectRunning {
				env.RegisterDelayedCallback(env.CancelWorkflow, time.Second*45)

				<-done
				assert.NoError(t, env.GetWorkflowError())
				assert.Greater(t, callCount, 1) // Expecting multiple activity calls
			} else {
				<-done
				var result LogWorkflowResponse
				assert.NoError(t, env.GetWorkflowResult(&result))
				assert.NotEmpty(t, result.Logs)
				assert.Equal(t, callCount, 1) // Only one activity call (no looping)
			}
		})
	}
}
