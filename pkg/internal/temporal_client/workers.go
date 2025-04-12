// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package temporalclient

import (
	"log"
	"sync"

	openidWorkflow "github.com/forkbombeu/credimi/pkg/OpenID4VP/workflow"
	credentialWorkflow "github.com/forkbombeu/credimi/pkg/credential_issuer/workflow"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// WorkerConfig defines a worker setup
type WorkerConfig struct {
	TaskQueue  string
	Workflows  []interface{}
	Activities []interface{}
}

// StartWorker initializes and runs a single Temporal worker
func StartWorker(client client.Client, config WorkerConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	w := worker.New(client, config.TaskQueue, worker.Options{})

	for _, wf := range config.Workflows {
		w.RegisterWorkflow(wf)
	}

	for _, act := range config.Activities {
		w.RegisterActivity(act)
	}

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Failed to start worker for %s: %v", config.TaskQueue, err)
	}
}

// StartAllWorkers initializes and starts multiple Temporal workers
func StartAllWorkers() {
	c, err := GetTemporalClient()
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	var wg sync.WaitGroup

	workers := []WorkerConfig{
		{
			TaskQueue: openidWorkflow.OpenIDTestTaskQueue,
			Workflows: []interface{}{
				openidWorkflow.OpenIDTestWorkflow,
				openidWorkflow.LogSubWorkflow,
			},
			// Get the activities via a function/method in a dynamic way
			Activities: []interface{}{
				openidWorkflow.GenerateYAMLActivity,
				openidWorkflow.RunStepCIJSProgramActivity,
				openidWorkflow.SendMailActivity,
				openidWorkflow.GetLogsActivity,
				openidWorkflow.TriggerLogsUpdateActivity,
			},
		},
		{
			TaskQueue: credentialWorkflow.CredentialsTaskQueue,
			Workflows: []interface{}{
				credentialWorkflow.CredentialWorkflow,
			},
			Activities: []interface{}{
				credentialWorkflow.FetchCredentialIssuerActivity,
				credentialWorkflow.StoreOrUpdateCredentialsActivity,
				credentialWorkflow.CleanupCredentialsActivity,
			},
		},
		{
			TaskQueue: credentialWorkflow.FetchIssuersTaskQueue,
			Workflows: []interface{}{
				credentialWorkflow.FetchIssuersWorkflow,
			},
			Activities: []interface{}{
				credentialWorkflow.FetchIssuersActivity,
				credentialWorkflow.CreateCredentialIssuersActivity,
			},
		},
	}

	for _, config := range workers {
		wg.Add(1)
		go StartWorker(c, config, &wg)
	}

	wg.Wait()
}

func WorkersHook(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		log.Println("Starting workers...")
		go StartAllWorkers()
		return se.Next()
	})
}

func StartUserWorker(namespace string) {
	c, err := GetTemporalClientWithNamespace(namespace)
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	worker := WorkerConfig{
		TaskQueue: openidWorkflow.OpenIDTestTaskQueue,
		Workflows: []interface{}{
			openidWorkflow.OpenIDTestWorkflow,
			openidWorkflow.LogSubWorkflow,
		},
		Activities: []interface{}{
			openidWorkflow.GenerateYAMLActivity,
			openidWorkflow.RunStepCIJSProgramActivity,
			openidWorkflow.SendMailActivity,
			openidWorkflow.GetLogsActivity,
			openidWorkflow.TriggerLogsUpdateActivity,
		},
	}

	go StartWorker(c, worker, &sync.WaitGroup{})
}
