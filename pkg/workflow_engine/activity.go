// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

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
	Execute(ctx context.Context, input ActivityInput) (ActivityResult, error)
	Name() string
}

type ConfigurableActivity interface {
	Configure(ctx context.Context, input *ActivityInput) error
}

func Fail(result *ActivityResult, msg string) (ActivityResult, error) {
	err := errors.New(msg)
	result.Errors = append(result.Errors, err)
	return *result, err
}
