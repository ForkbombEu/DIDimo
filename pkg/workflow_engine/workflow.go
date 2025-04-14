// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package workflowengine

import (
	"go.temporal.io/sdk/workflow"
)

type WorkflowInput struct {
	Payload map[string]any
	Config  map[string]any
}

type WorkflowResult struct {
	Message string
	Errors  []error
	Log     any
}

type Workflow interface {
	Workflow(ctx workflow.Context, input WorkflowInput) (WorkflowResult, error)
	Name() string
	GetOptions() workflow.ActivityOptions
}

type StartableWorkflow interface {
	Start(input WorkflowInput) (WorkflowResult, error)
}

type ConfigurableChild interface {
	Configure(ctx workflow.Context) workflow.Context
}
