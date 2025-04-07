// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"
	"os"

	"github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create a Temporal client
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = "localhost:7233"
	}
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create a worker that listens to a task queue
	w := worker.New(c, workflow.FetchIssuersTaskQueue, worker.Options{})

	// Register the workflow and activities
	w.RegisterWorkflow(workflow.FetchIssuersWorkflow)
	w.RegisterActivity(workflow.FetchIssuersActivity)
	w.RegisterActivity(workflow.CreateCredentialIssuersActivity)

	// Start the worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}
