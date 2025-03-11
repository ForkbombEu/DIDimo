package main

import (
	"log"
	"os"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP/workflow"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	godotenv.Load()
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = "localhost:7233"
	}

	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "openid-test-task-queue", worker.Options{})

	w.RegisterWorkflow(workflow.OpenIDTestWorkflow)
	w.RegisterActivity(workflow.GenerateYAMLActivity)
	w.RegisterActivity(workflow.RunStepCIJSProgramActivity)
	w.RegisterActivity(workflow.SendMailActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
