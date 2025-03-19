package main

import (
	"log"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP/workflow"
	temporalclient "github.com/forkbombeu/didimo/pkg/internal/temporal_client"
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

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
