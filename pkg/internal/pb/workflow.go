package pb

import (
	"context"
	"log"

	"github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"
	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"go.temporal.io/sdk/client"
)

func HookCredentialWorkflow(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess("credential_issuers").BindFunc(func(e *core.RecordEvent) error {
		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer c.Close()

		workflowInput := workflow.WorkflowInput{
			BaseURL: e.Record.Get("url").(string),
		}

		// Execute the workflow
		workflowOptions := client.StartWorkflowOptions{
			ID:        "credentials-workflow" + uuid.New().String(),
			TaskQueue: "CredentialsTaskQueue",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.CredentialWorkflow, workflowInput)
		if err != nil {
			log.Fatalf("Error starting worflow for URL %s: %v", e.Record.Get("url").(string), err)
		}
		var result string
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatalf("Error running worflow for URL %s: %v", e.Record.Get("url").(string), err)
		}

		log.Default().Printf("Workflow completed successfully: %s\n", result)

		// Part of the code for Giovanni
		// TODO - Not the real implementation should be removed in future

		relatedID := e.Record.Id
		providers, err := app.FindCollectionByNameOrId("services")
		if err != nil {
			return err
		}

		newRecord := core.NewRecord(providers)
		newRecord.Set("credential_issuers", relatedID)
		newRecord.Set("name", "TestName")
		// Save the new record in providers
		if err := app.Save(newRecord); err != nil {
			log.Println("Failed to create related record:", err)
			return err
		}

		log.Println("Successfully created related record in services")

		return e.Next()
	})
}
