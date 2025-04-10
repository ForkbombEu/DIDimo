// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"

	"github.com/forkbombeu/credimi/pkg/OpenID4VP/workflow"
	temporalclient "github.com/forkbombeu/credimi/pkg/internal/temporal_client"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/worker"
)

func main() {
	godotenv.Load()
	c, err := temporalclient.GetTemporalClient()
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	w := worker.New(c, workflow.OpenIDTestTaskQueue, worker.Options{})

	w.RegisterWorkflow(workflow.OpenIDTestWorkflow)
	w.RegisterWorkflow(workflow.LogSubWorkflow)
	w.RegisterActivity(workflow.GenerateYAMLActivity)
	w.RegisterActivity(workflow.RunStepCIJSProgramActivity)
	w.RegisterActivity(workflow.SendMailActivity)
	w.RegisterActivity(workflow.GetLogsActivity)
	w.RegisterActivity(workflow.TriggerLogsUpdateActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
