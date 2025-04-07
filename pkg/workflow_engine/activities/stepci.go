package activities

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/forkbombeu/didimo/pkg/utils"
	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
)

type StepCIWorkflowActivity struct {
	secrets map[string]string
}

// Configure injects the parsed template and token into the payload
func (a *StepCIWorkflowActivity) Configure(input *workflowengine.ActivityInput) error {
	templatePath := input.Config["template"]

	if templatePath == "" {
		return errors.New("missing required config: 'template'")
	}
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}
	tmpl, err := template.New("yaml").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	rendered, err := renderTemplate(tmpl, input.Payload)
	if err != nil {
		return err
	}
	input.Payload["yaml"] = rendered
	a.secrets = make(map[string]string)
	for k, v := range input.Config {
		if k != "template" {
			a.secrets[k] = v
		}
	}

	return nil
}

func (a *StepCIWorkflowActivity) Execute(ctx context.Context, input workflowengine.ActivityInput) (workflowengine.ActivityResult, error) {
	var result workflowengine.ActivityResult

	yamlContent, ok := input.Payload["yaml"].(string)
	if !ok || yamlContent == "" {
		return fail(&result, "missing rendered YAML in payload")
	}

	tmpFile, err := os.CreateTemp("", "stepci-*.yaml")
	if err != nil {
		return fail(&result, fmt.Sprintf("failed to create temp file: %v", err))
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		return fail(&result, fmt.Sprintf("failed to write to temp file: %v", err))
	}
	tmpFile.Close()

	var secretString bytes.Buffer
	for k, v := range a.secrets {
		if secretString.Len() > 0 {
			secretString.WriteString(" ")
		}
		secretString.WriteString(fmt.Sprintf("%s=%s", k, v))
	}

	binDir := utils.GetEnvironmentVariable("BIN", "/.bin", false)
	binName := "stepci-captured-runner"
	binPath := fmt.Sprintf("%s/%s", binDir, binName)

	// Build the arguments for the command
	args := []string{"-p", tmpFile.Name(), "-s", secretString.String()}

	cmd := exec.CommandContext(ctx, binPath, args...)
	output, err := cmd.CombinedOutput()

	result.Output = string(output)
	if err != nil {
		return fail(&result, fmt.Sprintf("stepci runner failed: %v\nOutput: %s", err, result.Output))
	}

	return result, nil
}

func renderTemplate(tmpl *template.Template, data map[string]any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return buf.String(), nil
}

func fail(result *workflowengine.ActivityResult, msg string) (workflowengine.ActivityResult, error) {
	err := errors.New(msg)
	result.Errors = append(result.Errors, err)
	return *result, err
}
