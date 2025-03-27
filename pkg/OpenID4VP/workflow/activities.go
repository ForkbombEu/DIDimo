package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/forkbombeu/didimo/pkg/internal/stepci"
	"gopkg.in/gomail.v2"
)

// GenerateYAML generates a YAML file based on provided variant and form
func GenerateYAMLActivity(ctx context.Context, input GenerateYAMLInput) error {

	schemasPath := os.Getenv("SCHEMAS_PATH")
	if schemasPath == "" {
		return fmt.Errorf("SCHEMAS_PATH environment variable not set")
	}

	testPlanResponseSchema, err := stepci.ConvertJSONToMap(fmt.Sprintf("%s/OpenID4VP_Wallet/responses/create_test_plan.json", schemasPath))
	if err != nil {
		return err
	}
	startRunnerResponseSchema, err := stepci.ConvertJSONToMap(fmt.Sprintf("%s/OpenID4VP_Wallet/responses/start_runner.json", schemasPath))
	if err != nil {
		return err
	}
	getInfoResponseSchema, err := stepci.ConvertJSONToMap(fmt.Sprintf("%s/OpenID4VP_Wallet/responses/get_info.json", schemasPath))
	if err != nil {
		return err
	}

	steps := []stepci.Step{
		{
			Name: "Create Test Plan",
			HTTP: stepci.HTTP{
				Method: "POST",
				URL:    "https://www.certification.openid.net/api/plan",
				Params: map[string]string{
					"planName": "oid4vp-id2-wallet-test-plan",
					"variant":  input.Variant,
				},
				Auth: struct {
					Ref string `yaml:"$ref"`
				}{Ref: "#/components/token"},
				Headers: map[string]string{
					"accept":       "application/json",
					"Content-Type": "application/json",
				},
				JSON: input.Form,
				Captures: map[string]struct {
					JSONPath string `yaml:"jsonpath"`
				}{
					"plan_id": {JSONPath: "$.id"},
				},
				Check: struct {
					Status int `yaml:"status,omitempty"`
					Schema any `yaml:"schema,omitempty"`
				}{
					Status: 201,
					Schema: testPlanResponseSchema,
				},
			},
		},
		{
			Name: "Start Test Runner",
			HTTP: stepci.HTTP{
				Method: "POST",
				URL:    "https://www.certification.openid.net/api/runner",
				Auth: struct {
					Ref string `yaml:"$ref"`
				}{Ref: "#/components/token"},
				Params: map[string]string{
					"test":    "oid4vp-id2-wallet-happy-flow-no-state",
					"plan":    "${{captures.plan_id}}",
					"variant": "{}",
				},
				Headers: map[string]string{
					"accept":       "application/json",
					"Content-Type": "application/json",
				},
				Captures: map[string]struct {
					JSONPath string `yaml:"jsonpath"`
				}{
					"id": {JSONPath: "$.id"},
				},
				Check: struct {
					Status int `yaml:"status,omitempty"`
					Schema any `yaml:"schema,omitempty"`
				}{
					Status: 201,
					Schema: startRunnerResponseSchema,
				},
			},
		},
		{
			Name: "Get Runner Info",
			HTTP: stepci.HTTP{
				Method: "GET",
				URL:    "https://www.certification.openid.net/api/runner/${{captures.id}}",
				Auth: struct {
					Ref string `yaml:"$ref"`
				}{Ref: "#/components/token"},
				Captures: map[string]struct {
					JSONPath string `yaml:"jsonpath"`
				}{
					"result": {JSONPath: "$.browser.urls[0]"},
					"rid":    {JSONPath: "$.id"},
				},
				Check: struct {
					Status int `yaml:"status,omitempty"`
					Schema any `yaml:"schema,omitempty"`
				}{
					Status: 200,
					Schema: getInfoResponseSchema,
				},
			},
		},
	}
	config := stepci.Config{
		Version: "1.0",
		Components: map[string]interface{}{
			"token": map[string]interface{}{
				"bearer": map[string]string{
					"token": "${{secrets.token}}",
				},
			},
		},
		Tests: map[string]stepci.Test{
			"OPENID4VP": {Steps: steps},
		},
	}

	err = stepci.GenerateYAML(config, input.FilePath)
	if err != nil {
		return fmt.Errorf("error generating YAML: %w", err)
	}

	return nil
}

// RunStepCIJSProgram executes the JavaScript program and returns a generic JSON response.
func RunStepCIJSProgramActivity(ctx context.Context, input StepCIRunnerInput) (StepCIRunnerResponse, error) {
	runStepCIPath := os.Getenv("RUN_STEPCI_PATH")
	if runStepCIPath == "" {
		return StepCIRunnerResponse{}, fmt.Errorf("RUN_STEPCI_PATH environment variable not set")
	}

	// Set up the command
	cmd := exec.CommandContext(ctx, "bun", "run", runStepCIPath, "-p", input.FilePath, "-s", "token="+input.Token)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return StepCIRunnerResponse{}, fmt.Errorf("error executing JS program: %w\nOutput: %s", err, string(output))
	}

	// Decode JSON output into a generic map
	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		return StepCIRunnerResponse{}, fmt.Errorf("failed to parse JSON output: %w\nRaw Output: %s", err, string(output))
	}

	return StepCIRunnerResponse{Result: result}, nil
}

// SendMailActivity sends an email
func SendMailActivity(ctx context.Context, config EmailConfig) error {

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", config.SenderEmail)
	m.SetHeader("To", config.ReceiverEmail)
	m.SetHeader("Subject", config.Subject)
	m.SetBody("text/html", config.Body)
	for filename, attachedBytes := range config.Attachments {
		attached := gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachedBytes)
			return err
		},
		)
		m.Attach(filename, attached)
	}

	// Send email
	d := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.Username,
		config.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// GetLogActivity performs the GET request to fetch the log
func GetLogsActivity(ctx context.Context, input LogActivitytyInput) ([]map[string]any, error) {

	baseURL, err := url.Parse(input.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	baseURL.Path += url.PathEscape(input.RID)

	query := baseURL.Query()
	query.Set("public", "false")
	baseURL.RawQuery = query.Encode()

	finalURL := baseURL.String()

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", input.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
