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
