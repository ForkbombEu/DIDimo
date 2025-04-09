package activities

import (
	"bytes"
	"fmt"
	"io"
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

func TestStepCIlActivity_Configure(t *testing.T) {
	activity := &StepCIWorkflowActivity{}

	tests := []struct {
		name          string
		config        map[string]string
		payload       map[string]interface{}
		templateBody  string
		expectedYAML  string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - valid template",
			config: map[string]string{
				"token": "secret-value",
			},
			payload: map[string]interface{}{
				"name": "world",
			},
			templateBody: `hello: [[ .name ]]`,
			expectedYAML: "hello: world",
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
			errorContains: "failed to open template file",
		},
		{
			name:   "Failure - invalid template syntax",
			config: map[string]string{},
			payload: map[string]interface{}{
				"name": "bad",
			},
			templateBody:  `[[ .name ]`, // malformed
			expectError:   true,
			errorContains: "failed to render YAML",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
			}
		})
	}
}

func TestStepCIActivity_Execute(t *testing.T) {
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
		config           map[string]string
		expectedError    bool
		expectedErrorMsg string
		expectedInOutput any
	}{
		{
			name: "Success - valid execution",
			payload: map[string]interface{}{
				"yaml": `
version: "1.1"
tests:
  example:
    steps:
      - name: Notfound test
        http:
          url: "${{secrets.test_secret}}"
          method: GET
          check:
            status: 404

      - name: GET request
        http:
          url: https://jsonplaceholder.typicode.com/posts/1
          method: GET
          check:
            jsonpath:
              $.id: 1
          captures:
            test:
              jsonpath: $.id
`,
			},
			config:           map[string]string{"test_secret": "https://httpbin.org/status/404"},
			expectedInOutput: map[string]any{"test": float64(1)},
		},
		{
			name: "Failure - missing runner binary",
			payload: map[string]interface{}{
				"yaml": "version: 1.0",
			},
			expectedError:    true,
			expectedErrorMsg: "stepci runner failed",
		},
		{
			name:             "Failure - incorrect secrets",
			payload:          map[string]any{"yaml": "version: 1.0"},
			config:           map[string]string{"wrongToken": "invalid-token"},
			expectedError:    true,
			expectedErrorMsg: "stepci runner failed",
		},
		{
			name: "Failure - stepCI fails",
			payload: map[string]any{
				"yaml": `
version: "1.1"
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
            status: 200
          captures:
            test:
              jsonpath: $
`,
			},
			config:           map[string]string{"test_secret": "https://httpbin.org/status/404"},
			expectedError:    true,
			expectedErrorMsg: "Workflow failed. Details:\n\n🔴 Step Failed: Notfound test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpYAMLFile, err := os.CreateTemp("", "test-*.yaml")
			require.NoError(t, err, "Failed to create temporary YAML file")
			defer os.Remove(tmpYAMLFile.Name()) // Ensure the file is removed after the test

			_, err = tmpYAMLFile.WriteString(tc.payload["yaml"].(string))
			require.NoError(t, err, "Failed to write to temporary YAML file")
			activity := &StepCIWorkflowActivity{}
			input := workflowengine.ActivityInput{
				Payload: map[string]interface{}{
					"yaml": tmpYAMLFile.Name(),
				},
				Config: tc.config,
			}
			var result workflowengine.ActivityResult
			future, err := env.ExecuteActivity(activity.Execute, input)
			if tc.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)
				future.Get(&result)
				require.Equal(t, tc.expectedInOutput, result.Output)
			}
		})
	}
}

func TestRenderYAML(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		data     map[string]any
		expected string
	}{
		{
			name:     "Simple string",
			tmpl:     "Hello, [[.Name]]!",
			data:     map[string]any{"Name": "Alice"},
			expected: "Hello, Alice!",
		},
		{
			name: "Complex object",
			tmpl: `test: [[ .Test | toJSON ]]
nested:
  [[ .Nested | toYAML | nindent 2 | trim ]]
nested2: [[ .Nested2 ]]`,
			data: map[string]any{
				"Test": map[string]string{
					"Username": "jdoe",
					"Email":    "jdoe@example.com",
				},
				"Nested": map[string]any{
					"nested": map[string]any{
						"test1": map[string]string{"Key": "value"},
						"test2": map[string]string{"Key2": "value2"},
					},
				},
				"Nested2": "nested_value2",
			},
			expected: `test: {"Email":"jdoe@example.com","Username":"jdoe"}
nested:
  nested:
      test1:
          Key: value
      test2:
          Key2: value2
nested2: nested_value2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reader := io.NopCloser(bytes.NewBufferString(tt.tmpl))
			output, err := RenderYAML(reader, tt.data)
			require.NoError(t, err, "RenderYAML should not return an error")
			require.Equal(t, strings.TrimSpace(tt.expected), strings.TrimSpace(output), "Rendered output should match expected")

		})
	}
}
