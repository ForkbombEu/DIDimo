package activities

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

// Helper to create a basic test template
func createTempTemplate(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "template-*.yaml")
	require.NoError(t, err)
	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())
	return tmpFile.Name()
}

func TestConfigure(t *testing.T) {
	activity := &StepCIWorkflowActivity{}

	tests := []struct {
		name           string
		config         map[string]string
		payload        map[string]interface{}
		templateBody   string
		expectedYAML   string
		expectedSecret string
		expectError    bool
		errorContains  string
	}{
		{
			name: "Success - valid template",
			config: map[string]string{
				"token": "secret-value",
			},
			payload: map[string]interface{}{
				"name": "world",
			},
			templateBody:   `hello: {{ .name }}`,
			expectedYAML:   "hello: world",
			expectedSecret: "secret-value",
		},
		{
			name:          "Failure - missing template path",
			config:        map[string]string{},
			expectError:   true,
			errorContains: "missing required config",
		},
		{
			name: "Failure - invalid template path",
			config: map[string]string{
				"template": "/not/found.yaml",
			},
			expectError:   true,
			errorContains: "failed to read template file",
		},
		{
			name:   "Failure - invalid template syntax",
			config: map[string]string{},
			payload: map[string]interface{}{
				"name": "bad",
			},
			templateBody:  `{{ .name }`, // malformed
			expectError:   true,
			errorContains: "failed to parse template",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file if template body is provided but no path given
			if tc.templateBody != "" && tc.config["template"] == "" {
				tmp := createTempTemplate(t, tc.templateBody)
				defer os.Remove(tmp)
				tc.config["template"] = tmp
			}

			input := &workflowengine.ActivityInput{
				Config:  tc.config,
				Payload: tc.payload,
			}

			err := activity.Configure(input)

			if tc.expectError {
				require.Error(t, err)
				if tc.errorContains != "" {
					require.ErrorContains(t, err, tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				yaml, ok := input.Payload["yaml"].(string)
				require.True(t, ok, "expected payload to contain string field 'yaml'")
				require.Equal(t, tc.expectedYAML, strings.TrimSpace(yaml))
				require.Equal(t, tc.expectedSecret, activity.secrets["token"])
			}
		})
	}
}

func TestExecute(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestActivityEnvironment()
	activity := &StepCIWorkflowActivity{}
	env.RegisterActivity(activity.Execute)

	tmpBinDir := t.TempDir()
	binPath := fmt.Sprintf("%s/stepci-captured-runner", tmpBinDir)

	// Determine the platform and architecture
	OS := runtime.GOOS
	arch := runtime.GOARCH

	// Construct the binary download URL
	url := fmt.Sprintf("https://github.com/ForkbombEu/stepci-captured-runner/releases/latest/download/stepci-captured-runner-%s-%s", OS, arch)

	// Download the binary from GitHub
	cmd := exec.Command("wget", url, "-O", binPath)
	cmd.Dir = tmpBinDir // Set working directory to the temporary binary directory

	t.Logf("Downloading binary from: %s", url)
	err := cmd.Run()
	require.NoError(t, err, "Failed to download binary")

	// Make the binary executable
	err = os.Chmod(binPath, 0755)
	require.NoError(t, err, "Failed to make binary executable")

	// Set environment variable to point to the binary directory
	os.Setenv("BIN", tmpBinDir)

	tests := []struct {
		name             string
		payload          map[string]interface{}
		prepareBinary    bool
		secrets          map[string]string
		expectedError    bool
		expectedErrorMsg string
		expectedInOutput string
	}{
		{
			name: "Success - valid execution",
			payload: map[string]interface{}{
				"yaml": `version: "1.1"

tests:
  example:
    steps:
      - name: GET request
        http:
          url: https://httpbin.org/status/200
          method: GET
          check:
            status: 200
      - name: Notfound test
        http:
          url: "${{secrets.test_secret}}"
          method: GET
          check:
            status: 404
          captures:
            test:
              jsonpath: $`,
			},
			prepareBinary:    true,
			secrets:          map[string]string{"test_secret": "https://httpbin.org/status/404"},
			expectedInOutput: "{}",
		},
		{
			name:             "Failure - incorrect secrets",
			payload:          map[string]interface{}{"yaml": "version: 1.0"},
			secrets:          map[string]string{"wrongToken": "invalid-token"},
			expectedError:    true,
			expectedErrorMsg: "stepci runner failed", // Adjust to match your actual error
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prepareBinary {
				_, err := os.Stat(binPath)
				if err != nil {
					require.Fail(t, "Binary not found", "Ensure the real binary exists at: %s", binPath)
				}
			}
			tmpYAMLFile, err := os.CreateTemp("", "test-*.yaml")
			require.NoError(t, err, "Failed to create temporary YAML file")
			defer os.Remove(tmpYAMLFile.Name())

			_, err = tmpYAMLFile.WriteString(tc.payload["yaml"].(string))
			require.NoError(t, err, "Failed to write to temporary YAML file")
			activity := &StepCIWorkflowActivity{secrets: tc.secrets}
			input := workflowengine.ActivityInput{
				Payload: map[string]interface{}{
					"yaml": tmpYAMLFile.Name(),
				},
			}

			result, err := activity.Execute(context.Background(), input)

			if tc.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)
				require.Contains(t, result.Output, tc.expectedInOutput)
			}
		})
	}
}
