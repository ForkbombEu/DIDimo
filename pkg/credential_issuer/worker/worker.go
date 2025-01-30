package main

import (
	"log"

	"github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create a Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create a worker that listens to a task queue
	w := worker.New(c, "CredimiTaskQueue", worker.Options{})

	// Register the workflow and activities
	w.RegisterWorkflow(workflow.CredentialWorkflow)
	w.RegisterActivity(workflow.FetchCredentialIssuerActivity)
	w.RegisterActivity(workflow.StoreCredentialsActivity)

	// Start the worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}
