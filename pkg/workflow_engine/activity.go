package workflowengine

import (
	"context"
	"errors"
)

type ActivityInput struct {
	Payload map[string]any
	Config  map[string]string
}

type ActivityResult struct {
	Output any
	Errors []error
	Log    []string
}

type ExecutableActivity interface {
	Execute(ctx context.Context, input ActivityInput) (result ActivityResult, err error)
}

type ConfigurableActivity interface {
	Configure(input *ActivityInput) error
}

func Fail(result *ActivityResult, msg string) (ActivityResult, error) {
	err := errors.New(msg)
	result.Errors = append(result.Errors, err)
	return *result, err
}
