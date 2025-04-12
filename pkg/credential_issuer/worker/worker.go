// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"

	"github.com/forkbombeu/credimi/pkg/credential_issuer/workflow"
	temporalclient "github.com/forkbombeu/credimi/pkg/internal/temporal_client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create a Temporal client
	c, err := temporalclient.GetTemporalClient()
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create a worker that listens to a task queue
	w := worker.New(c, "CredentialsTaskQueue", worker.Options{})

	// Register the workflow and activities
	w.RegisterWorkflow(workflow.CredentialWorkflow)
	w.RegisterActivity(workflow.FetchCredentialIssuerActivity)
	w.RegisterActivity(workflow.StoreOrUpdateCredentialsActivity)
	w.RegisterActivity(workflow.CleanupCredentialsActivity)

	// Start the worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}
