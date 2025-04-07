// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestGenerateYAMLActivity(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestActivityEnvironment()
	env.RegisterActivity(GenerateYAMLActivity)

	// Define mock Form
	mockForm := map[string]any{
		"alias":       "MOCK_ALIAS",
		"description": "Mock description for testing",
		"server": map[string]any{
			"authorization_endpoint": "mock-protocol://mock-auth",
		},
		"client": map[string]any{
			"client_id": "mock:client:id",
			"presentation_definition": map[string]any{
				"id": "mock_presentation_id",
				"input_descriptors": []map[string]any{
					{
						"id": "mock_descriptor_1",
						"constraints": map[string]any{
							"fields": []map[string]any{
								{
									"path": []string{"$.mock_field"},
									"filter": map[string]string{
										"type":  "string",
										"const": "mock_const_value",
									},
								},
							},
						},
						"format": map[string]any{
							"vc+sd-jwt": map[string][]string{
								"mock_jwt_alg_values": {"MOCK_ALG_1", "MOCK_ALG_2"},
								"mock_sd_jwt_values":  {"MOCK_ALG_3"},
							},
						},
					},
				},
			},
			"jwks": map[string]any{
				"keys": []map[string]any{
					{
						"kty": "MOCK_KTY",
						"alg": "MOCK_ALG",
						"crv": "MOCK_CRV",
						"d":   "MOCK_D_VALUE",
						"x":   "MOCK_X_VALUE",
						"y":   "MOCK_Y_VALUE",
					},
				},
			},
		},
	}

	testCases := []struct {
		name          string
		variant       string
		form          any
		setupEnv      func()
		expectedError bool
	}{
		{
			name:          "Success Valid input",
			variant:       "testVariant",
			form:          mockForm,
			setupEnv:      func() { os.Setenv("SCHEMAS_PATH", "../../../schemas") },
			expectedError: false,
		},
		{
			name:          "Failure  Missing SCHEMAS_PATH",
			variant:       "testVariant",
			form:          mockForm,
			setupEnv:      func() { os.Unsetenv("SCHEMAS_PATH") },
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupEnv()

			tmpFile, err := os.CreateTemp("", "test-*.yaml")
			require.NoError(t, err, "Failed to create temporary YAML file")
			defer os.Remove(tmpFile.Name())
			YAMLInput := GenerateYAMLInput{
				Variant:  tc.variant,
				Form:     tc.form,
				FilePath: tmpFile.Name(),
			}

			_, err = env.ExecuteActivity(GenerateYAMLActivity, YAMLInput)
			if tc.expectedError {
				require.Error(t, err, "Expected an error but did not receive one")
				return
			}
			require.NoError(t, err, "Activity execution failed unexpectedly")

			yamlData, err := os.ReadFile(tmpFile.Name())
			require.NoError(t, err, "Failed to read generated YAML file")

			expectedYAML := `version: "1.0"
components:
    token:
        bearer:
            token: ${{secrets.token}}
tests:
    OPENID4VP:
        steps:
            - name: Create Test Plan
              http:
                method: POST
                url: https://www.certification.openid.net/api/plan
                params:
                    planName: oid4vp-id2-wallet-test-plan
                    variant: testVariant
                auth:
                    $ref: '#/components/token'
                headers:
                    Content-Type: application/json
                    accept: application/json
                json:
                    alias: MOCK_ALIAS
                    description: Mock description for testing
                    server:
                        authorization_endpoint: mock-protocol://mock-auth
                    client:
                        client_id: mock:client:id
                        presentation_definition:
                            id: mock_presentation_id
                            input_descriptors:
                                - id: mock_descriptor_1
                                  constraints:
                                    fields:
                                        - path:
                                            - $.mock_field
                                          filter:
                                            const: mock_const_value
                                            type: string
                                  format:
                                    vc+sd-jwt:
                                        mock_jwt_alg_values:
                                            - MOCK_ALG_1
                                            - MOCK_ALG_2
                                        mock_sd_jwt_values:
                                            - MOCK_ALG_3
                        jwks:
                            keys:
                                - kty: MOCK_KTY
                                  alg: MOCK_ALG
                                  crv: MOCK_CRV
                                  d: MOCK_D_VALUE
                                  x: MOCK_X_VALUE
                                  "y": MOCK_Y_VALUE
                captures:
                    plan_id:
                        jsonpath: $.id
                check:
                    status: 201
                    schema:
                        properties:
                            id:
                                type: string
                            modules:
                                items:
                                    properties:
                                        instances:
                                            items:
                                                type: object
                                            type: array
                                        testModule:
                                            type: string
                                        variant:
                                            type: object
                                    required:
                                        - testModule
                                        - variant
                                        - instances
                                    type: object
                                type: array
                            name:
                                type: string
                        required:
                            - name
                            - id
                            - modules
                        type: object
            - name: Start Test Runner
              http:
                method: POST
                url: https://www.certification.openid.net/api/runner
                params:
                    plan: ${{captures.plan_id}}
                    test: oid4vp-id2-wallet-happy-flow-no-state
                    variant: '{}'
                auth:
                    $ref: '#/components/token'
                headers:
                    Content-Type: application/json
                    accept: application/json
                captures:
                    id:
                        jsonpath: $.id
                check:
                    status: 201
                    schema:
                        properties:
                            id:
                                type: string
                            name:
                                type: string
                            url:
                                format: uri
                                type: string
                        required:
                            - name
                            - id
                            - url
                        type: object
            - name: Get Runner Info
              http:
                method: GET
                url: https://www.certification.openid.net/api/runner/${{captures.id}}
                auth:
                    $ref: '#/components/token'
                captures:
                    result:
                        jsonpath: $.browser.urls[0]
                    rid:
                        jsonpath: $.id
                check:
                    status: 200
                    schema:
                        properties:
                            browser:
                                properties:
                                    browserApiRequests:
                                        items:
                                            type: object
                                        type: array
                                    runners:
                                        items:
                                            type: object
                                        type: array
                                    show_qr_code:
                                        type: boolean
                                    urls:
                                        items:
                                            format: uri
                                            type: string
                                        type: array
                                    urlsWithMethod:
                                        items:
                                            properties:
                                                method:
                                                    type: string
                                                url:
                                                    format: uri
                                                    type: string
                                            required:
                                                - url
                                                - method
                                            type: object
                                        type: array
                                    visited:
                                        items:
                                            type: object
                                        type: array
                                    visitedUrlsWithMethod:
                                        items:
                                            type: object
                                        type: array
                                required:
                                    - browserApiRequests
                                    - urls
                                    - show_qr_code
                                    - visited
                                    - visitedUrlsWithMethod
                                type: object
                            created:
                                format: date-time
                                type: string
                            error:
                                nullable: true
                                type: string
                            exposed:
                                properties:
                                    client_id:
                                        type: string
                                    nonce:
                                        type: string
                                    response_uri:
                                        format: uri
                                        type: string
                                    state:
                                        nullable: true
                                        type: string
                                required:
                                    - response_uri
                                type: object
                            id:
                                type: string
                            name:
                                type: string
                            owner:
                                properties:
                                    iss:
                                        format: uri
                                        type: string
                                    sub:
                                        type: string
                                required:
                                    - sub
                                    - iss
                                type: object
                            updated:
                                format: date-time
                                type: string
                        required:
                            - owner
                            - created
                            - browser
                            - name
                            - exposed
                            - id
                            - updated
                        type: object

`
			require.YAMLEq(t, expectedYAML, string(yamlData), "Generated YAML does not match expected output")
		})
	}
}

func TestRunStepCIJSProgramActivity(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestActivityEnvironment()
	env.RegisterActivity(RunStepCIJSProgramActivity)
	testCases := []struct {
		name             string
		yamlContent      string // YAML content to write to the temp file
		token            string
		setupEnv         func()
		expectedError    bool
		expectedJSON     string
		expectedErrorMsg string
	}{
		{
			name: "Success Valid execution",
			yamlContent: `version: 1
tests:
  example:
    steps:
      - name: GET request
        http:
          url: https://jsonplaceholder.typicode.com/posts/1
          method: GET
          check:
            status: 200
          captures:
            result:
              jsonpath: $ `,
			setupEnv:      func() { os.Setenv("RUN_STEPCI_PATH", "../stepci/RunStepCI.js") },
			expectedError: false,
			expectedJSON: `{
  "result": {
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
  }
}`,
			expectedErrorMsg: "",
		},
		{
			name:             "Failure Missing RUN_STEPCI_PATH",
			yamlContent:      `mock: yaml content`,
			token:            "mockToken",
			setupEnv:         func() { os.Unsetenv("RUN_STEPCI_PATH") },
			expectedError:    true,
			expectedJSON:     "",
			expectedErrorMsg: "RUN_STEPCI_PATH environment variable not set",
		},
		{
			name: "Failure StepCI error",
			yamlContent: `version: 1
tests:
  example:
    steps:
      - name: GET request
        http:
          url: https://httpbin.org/status/200
          method: GET
          check:
            status: 500`,
			token:            "",
			setupEnv:         func() { os.Setenv("RUN_STEPCI_PATH", "../stepci/RunStepCI.js") },
			expectedError:    true,
			expectedJSON:     "",
			expectedErrorMsg: "Output: ❌ Workflow failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up environment variables
			tc.setupEnv()

			tmpYAMLFile, err := os.CreateTemp("", "test-*.yaml")
			require.NoError(t, err, "Failed to create temporary YAML file")
			defer os.Remove(tmpYAMLFile.Name()) // Ensure the file is removed after the test

			_, err = tmpYAMLFile.WriteString(tc.yamlContent)
			require.NoError(t, err, "Failed to write to temporary YAML file")
			stepCIInput := StepCIRunnerInput{
				FilePath: tmpYAMLFile.Name(),
				Token:    tc.token,
			}
			var result StepCIRunnerResponse
			future, err := env.ExecuteActivity(RunStepCIJSProgramActivity, stepCIInput)

			if tc.expectedError {
				require.Error(t, err, "Expected an error but did not receive one")
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err, "Activity execution failed unexpectedly")
				err := future.Get(&result)
				require.NoError(t, err, "Failed to get activity result")

				expected := make(map[string]any)
				err = json.Unmarshal([]byte(tc.expectedJSON), &expected)
				require.NoError(t, err, "Failed to unmarshal JSON")
				require.Equal(t, expected, result.Result)
			}
		})
	}
}

func TestSendMailActivity(t *testing.T) {
	// Start mock SMTP server
	mockServer := smtpmock.New(smtpmock.ConfigurationAttr{
		PortNumber:  2525, // Mock SMTP listens on 2525
		LogToStdout: false,
	})
	if err := mockServer.Start(); err != nil {
		t.Fatalf("failed to start mock SMTP server: %v", err)
	}
	defer mockServer.Stop()

	// Test email configuration
	emailConfig := EmailConfig{
		SMTPHost:      "localhost",
		SMTPPort:      2525, // Use the mock server's port
		Username:      "testuser",
		Password:      "testpassword",
		SenderEmail:   "sender@example.com",
		ReceiverEmail: "receiver@example.com",
		Subject:       "Test Email",
		Body:          "This is a test email.",
		Attachments: map[string][]byte{
			"test.txt": []byte("Test attachment"),
		},
	}

	// Run the activity with mock server
	err := SendMailActivity(context.Background(), emailConfig)
	require.NoError(t, err, "Expected email to be sent without error")

	// Check if email was sent
	if len(mockServer.Messages()) != 1 {
		t.Errorf("Expected email to be sent, but mock server received none")
	}
}

func TestGetLogsActivity(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestActivityEnvironment()
	env.RegisterActivity(GetLogsActivity)
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectedError  string
		expectedResult []map[string]interface{}
	}{
		{
			name: "Valid Logs Response",
			serverResponse: `[
				{"log_id": "1", "message": "log entry 1"},
				{"log_id": "2", "message": "log entry 2"}
			]`,
			statusCode:    http.StatusOK,
			expectedError: "",
			expectedResult: []map[string]interface{}{
				{"log_id": "1", "message": "log entry 1"},
				{"log_id": "2", "message": "log entry 2"},
			},
		},
		{
			name:           "Invalid JSON Response",
			serverResponse: `{"log_id": "1", "message"`,
			statusCode:     http.StatusOK,
			expectedError:  "failed to decode response",
			expectedResult: nil,
		},
		{
			name:           "Non-200 Status Code",
			serverResponse: "",
			statusCode:     http.StatusUnauthorized,
			expectedError:  "unexpected status code: 401",
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			}))
			defer server.Close()

			// Call the GetLogsActivity with the mock server URL
			input := GetLogsActivityInput{
				BaseURL: server.URL,
				RID:     "test-rid",
				Token:   "test-token",
			}

			future, err := env.ExecuteActivity(GetLogsActivity, input)

			// Validate the result
			if tt.expectedError == "" {
				require.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected an error")
				require.Contains(t, err.Error(), tt.expectedError)
			}

			// Check if the result matches the expected logs
			if tt.expectedResult != nil {
				var logs []map[string]any
				err := future.Get(&logs)
				require.NoError(t, err, "Failed to get activity result")
				require.ElementsMatch(t, tt.expectedResult, logs)
			}
		})
	}
}

func TestTriggerLogsUpdateActivity(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestActivityEnvironment()
	env.RegisterActivity(TriggerLogsUpdateActivity)

	tests := []struct {
		name           string
		statusCode     int
		serverResponse string
		expectedError  string
	}{
		{
			name:           "Successful Log Update",
			statusCode:     http.StatusOK,
			serverResponse: "{}",
			expectedError:  "",
		},
		{
			name:           "Non-200 Response",
			statusCode:     http.StatusInternalServerError,
			serverResponse: "Internal Server Error",
			expectedError:  "failed to send log update, received status: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.serverResponse)
			}))
			defer server.Close()

			input := TriggerLogsUpdateActivityInput{
				WorkflowID: "test-workflow",
				Logs:       []map[string]any{{"log_id": "1", "message": "test log"}},
				AppURL:     server.URL,
			}

			_, err := env.ExecuteActivity(TriggerLogsUpdateActivity, input)

			if tt.expectedError == "" {
				require.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected an error")
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
