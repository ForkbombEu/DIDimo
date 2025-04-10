package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/group/all"

	"github.com/forkbombeu/didimo/pkg/utils"
	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
)

type StepCIWorkflowActivity struct{}

// Configure injects the parsed template and token into the payload
func (a *StepCIWorkflowActivity) Configure(
	ctx context.Context,
	input *workflowengine.ActivityInput,
) error {
	templatePath := input.Config["template"]
	if templatePath == "" {
		return errors.New("missing required config: 'template'")
	}

	file, err := os.Open(templatePath)
	if err != nil {
		return fmt.Errorf("failed to open template file: %w", err)
	}
	defer file.Close()

	rendered, err := RenderYAML(file, input.Payload)
	if err != nil {
		return fmt.Errorf("failed to render YAML: %w", err)
	}

	input.Payload["yaml"] = rendered

	return nil
}

func (a *StepCIWorkflowActivity) Execute(
	ctx context.Context,
	input workflowengine.ActivityInput,
) (workflowengine.ActivityResult, error) {
	var result workflowengine.ActivityResult

	yamlContent, ok := input.Payload["yaml"].(string)
	if !ok || yamlContent == "" {
		return workflowengine.Fail(&result, "missing rendered YAML in payload")
	}

	var secretString bytes.Buffer
	for k, v := range input.Config {
		if k != "template" {
			if secretString.Len() > 0 {
				secretString.WriteString(" ")
			}
			secretString.WriteString(fmt.Sprintf("%s=%s", k, v))
		}
	}

	binDir := utils.GetEnvironmentVariable("BIN", ".bin", false)
	binName := "stepci-captured-runner"
	binPath := fmt.Sprintf("%s/%s", binDir, binName)
	// Build the arguments for the command
	args := []string{yamlContent, "-s", secretString.String()}

	cmd := exec.CommandContext(ctx, binPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return workflowengine.Fail(
			&result,
			fmt.Sprintf("stepci runner failed: %v\nOutput: %s", err, output),
		)
	}
	var outputJSON map[string]any

	if err := json.Unmarshal(output, &outputJSON); err != nil {
		return workflowengine.Fail(&result, fmt.Sprintf("failed to unmarshal JSON output: %v", err))
	}
	result.Output = outputJSON
	return result, nil
}

func RenderYAML(reader io.Reader, data map[string]interface{}) (string, error) {
	handler := sprout.New(
		sprout.WithGroups(all.RegistryGroup()),
	)
	funcs := handler.Build()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return "", err
	}

	templateContent := buf.String()

	tmpl, err := template.New("yaml").Delims("[[", "]]").
		Funcs(funcs).
		Parse(templateContent)
	if err != nil {
		return "", err
	}

	buf.Reset()
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// Decode HTML entities from the rendered string
	result := html.UnescapeString(buf.String())

	// Remove any leading/trailing whitespace or extra newlines from the result
	return strings.TrimSpace(result), nil
}
