package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var retryPolicy = &temporal.RetryPolicy{
	InitialInterval:    time.Second * 5,
	BackoffCoefficient: 2.0,
	MaximumInterval:    time.Minute,
	MaximumAttempts:    5,
}

var ActivityOptions = workflow.ActivityOptions{
	ScheduleToCloseTimeout: time.Minute * 10,
	StartToCloseTimeout:    time.Minute * 5,
	RetryPolicy:            retryPolicy,
}
