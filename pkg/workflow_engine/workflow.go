package workflowengine

import (
	"go.temporal.io/sdk/workflow"
)

type WorkflowInput struct {
	Payload        map[string]any
	Config         map[string]any
	ActvityOptions *workflow.ActivityOptions
}

type WorkflowResult struct {
	Message string
	Errors  []error
	Log     any
}

type Workflow interface {
	Workflow(ctx workflow.Context, input WorkflowInput) (result WorkflowResult, err error)
	Start(input WorkflowInput) (result WorkflowResult, err error)
}
