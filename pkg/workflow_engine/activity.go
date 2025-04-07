package workflowengine

import "context"

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
