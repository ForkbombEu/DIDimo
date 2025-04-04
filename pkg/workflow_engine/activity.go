package workflowengine

type ActivityInput struct {
	Payload map[string]any
	Config  map[string]string
}

type ActivityResult struct {
	Output any
	Errors []error
	Log []string
}

type ExecutableActivity interface {
	Execute(input ActivityInput) (result ActivityResult,error)
}

type ConfigurableActivity interface {
	Configure(input *ActivityInput)
}
