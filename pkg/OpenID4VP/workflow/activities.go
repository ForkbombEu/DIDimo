package workflow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/forkbombeu/didimo/pkg/OPENID4VP/testdata"
	"github.com/forkbombeu/didimo/pkg/internal/stepci"
	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"github.com/yeqown/go-qrcode/writer/terminal"
)

// GenerateYAML generates a YAML file based on provided variant and jsonPayload
func GenerateYAMLActivity(ctx context.Context, variant string, jsonPayload testdata.JSONPayload, filePath string) error {

	testPlanResponseSchema, err := stepci.ConvertJSONToMap("schemas/OpenID4VPTest/responses/create_test_plan.json")
	if err != nil {
		return err
	}
	startRunnerResponseSchema, err := stepci.ConvertJSONToMap("schemas/OpenID4VPTest/responses/start_runner.json")
	if err != nil {
		return err
	}
	getInfoResponseSchema, err := stepci.ConvertJSONToMap("schemas/OpenID4VPTest/responses/get_info.json")
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
					"variant":  variant,
				},
				Auth: struct {
					Ref string `yaml:"$ref"`
				}{Ref: "#/components/token"},
				Headers: map[string]string{
					"accept":       "application/json",
					"Content-Type": "application/json",
				},
				JSON: jsonPayload,
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

	err = stepci.GenerateYAML(config, filePath)
	if err != nil {
		return fmt.Errorf("error generating YAML: %w", err)
	}

	return nil
}

// RunStepCIJSProgram executes the JavaScript program and returns a generic JSON response.
func RunStepCIJSProgramActivity(ctx context.Context, yamlFilePath, token string) (map[string]any, error) {
	// Set up the command
	cmd := exec.CommandContext(ctx, "bun", "run", "pkg/OpenID4VP/stepci/runStepCI.js", "-p", yamlFilePath, "-s", "token="+token)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing JS program: %w\nOutput: %s", err, string(output))
	}

	// Decode JSON output into a generic map
	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w\nRaw Output: %s", err, string(output))
	}

	return result, nil
}

type writeCloserWrapper struct {
	*bytes.Buffer
}

// Close is a no-op to satisfy io.WriteCloser
func (w *writeCloserWrapper) Close() error {
	return nil
}

// GenerateQRCodeActivity takes a URL and returns a base64-encoded QR code.
func GenerateQRCodeActivity(ctx context.Context, url string) (string, error) {

	qr, err := qrcode.New(url)
	if err != nil {
		return "", fmt.Errorf("failed to create QR code: %w", err)
	}

	// Use a buffer to store the PNG output
	var buf bytes.Buffer
	w := standard.NewWithWriter(&writeCloserWrapper{&buf}) // Use custom wrapper

	// Encode and write QR code to buffer
	if err := qr.Save(w); err != nil {
		return "", fmt.Errorf("failed to generate QR code image: %w", err)
	}

	// Convert to base64 string
	qrBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return qrBase64, nil
}

func PrintQRCodeACtivity(ctx context.Context, url string) error {

	qr, err := qrcode.New(url)
	if err != nil {
		return fmt.Errorf("failed to create QR code: %w", err)
	}

	// Create a terminal writer to print the QR code to the terminal
	tw := terminal.New()

	// Write the QR code to the terminal
	err = qr.Save(tw)
	if err != nil {
		return fmt.Errorf("failed to print QR code to terminal: %w", err)
	}
	return nil
}
