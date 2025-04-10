package worker_engine

import (
	"log"
	"sync"

	credentialWorkflow "github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"
	temporalclient "github.com/forkbombeu/didimo/pkg/internal/temporal_client"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/activities"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/workflows"
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
	c, err := temporalclient.GetTemporalClient()
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	var wg sync.WaitGroup

	var OpenIDNetWorkflow = workflows.OpenIDNetWorkflow{}
	var OpenIDNetLogsWorkflow = workflows.OpenIDNetLogsWorkflow{}
	workers := []WorkerConfig{
		{
			TaskQueue: workflows.OpenIDTestTaskQueue,
			Workflows: []interface{}{
				OpenIDNetWorkflow.Workflow,
				OpenIDNetLogsWorkflow.SubWorkflow,
			},
			Activities: []interface{}{
				&activities.StepCIWorkflowActivity{},
				&activities.SendMailActivity{},
				&activities.HTTPActivity{},
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
