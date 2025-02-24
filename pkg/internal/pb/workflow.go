package pb

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP"
	"github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"
	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"go.temporal.io/sdk/client"
)

type OpenID4VPRequest struct {
	Input      OpenID4VP.OpenID4VPTestInputFile `json:"input"`
	UserMail   string                           `json:"user_mail"`
	WorkflowID string                           `json:"workflow_id"`
	TestName   string                           `json:"test_name"`
}

func HookCredentialWorkflow(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess("credential_issuers").BindFunc(func(e *core.RecordEvent) error {

		hostPort := os.Getenv("TEMPORAL_ADDRESS")
		if hostPort == "" {
			hostPort = "localhost:7233"
		}
		c, err := client.Dial(client.Options{
			HostPort: hostPort,
		})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer c.Close()

		workflowInput := workflow.WorkflowInput{
			BaseURL:  e.Record.Get("url").(string),
			IssuerID: e.Record.Id,
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
		return e.Next()
	})
}

func AddOpenID4VPTestEndpoint(app *pocketbase.PocketBase) {

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.POST("/api/openid4vp-test", func(e *core.RequestEvent) error {
			var req OpenID4VPRequest
			if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
				return apis.NewBadRequestError("invalid JSONE input", err)
			}

			// Start the workflow
			err := OpenID4VP.StartWorkflow(req.Input, req.UserMail)
			if err != nil {
				return apis.NewBadRequestError("failed to start OpenID4VP workflow", err)
			}

			return e.JSON(http.StatusOK, map[string]string{
				"message": "Workflow started successfully",
			})
		})
		return se.Next()
	})
}
