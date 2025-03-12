package pb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP"
	openid4vp_workflow "github.com/forkbombeu/didimo/pkg/OpenID4VP/workflow"
	credential_workflow "github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"go.temporal.io/sdk/client"
)

type OpenID4VPRequest struct {
	Input    OpenID4VP.OpenID4VPTestInputFile `json:"input"`
	UserMail string                           `json:"user_mail"`
	TestName string                           `json:"test_name"`
}

type IssuerURL struct {
	URL string `json:"credentialIssuerUrl"`
}

func HookCredentialWorkflow(app *pocketbase.PocketBase) {
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = "localhost:7233"
	}
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalln("unable to create client", err)
	}

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.POST("/credentials_issuers/start-check", func(e *core.RequestEvent) error {
			var req IssuerURL

			if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
				return apis.NewBadRequestError("invalid JSON input", err)
			}
			// Check if a record with the given URL already exists
			collection, err := app.FindCollectionByNameOrId("credential_issuers")
			if err != nil {
				return err
			}

			existingRecords, err := app.FindRecordsByFilter(
				collection.Id,
				"url = {:url}",
				"",
				1,
				0,
				dbx.Params{"url": req.URL},
			)
			if err != nil {
				return err
			}
			var issuerID string

			if len(existingRecords) > 0 {
				issuerID = existingRecords[0].Id
			} else {
				// Create a new record
				newRecord := core.NewRecord(collection)
				newRecord.Set("url", req.URL)

				if err := app.Save(newRecord); err != nil {
					return err
				}

				issuerID = newRecord.Id
			}

			// Start the workflow
			workflowInput := credential_workflow.CredentialWorkflowInput{
				BaseURL:  req.URL,
				IssuerID: issuerID,
			}

			workflowOptions := client.StartWorkflowOptions{
				ID:        "credentials-workflow-" + uuid.New().String(),
				TaskQueue: "CredentialsTaskQueue",
			}

			we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, credential_workflow.CredentialWorkflow, workflowInput)
			if err != nil {
				return fmt.Errorf("error starting workflow for URL %s: %v", req.URL, err)

			}

			var result string
			err = we.Get(context.Background(), &result)
			if err != nil {
				return fmt.Errorf("error running workflow for URL %s: %v", req.URL, err)

			}

			log.Printf("Workflow completed successfully for URL %s: %s", req.URL, result)

			providers, err := app.FindCollectionByNameOrId("services")
			if err != nil {
				return err
			}

			newRecord := core.NewRecord(providers)
			newRecord.Set("credential_issuers", issuerID)
			newRecord.Set("name", "TestName")
			// Save the new record in providers
			if err := app.Save(newRecord); err != nil {
				log.Println("Failed to create related record:", err)
				return err
			}
			return e.JSON(http.StatusOK, map[string]string{
				"credentialIssuerUrl": req.URL,
			})
		})
		return se.Next()
	})
}

func AddOpenID4VPTestEndpoints(app *pocketbase.PocketBase) {

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
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

		se.Router.POST("/api/openid4vp-test", func(e *core.RequestEvent) error {
			var req OpenID4VPRequest
			appURL := app.Settings().Meta.AppURL
			if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
				return apis.NewBadRequestError("invalid JSON input", err)
			}

			// Start the workflow
			err := OpenID4VP.StartWorkflow(req.Input, req.UserMail, appURL)
			if err != nil {
				return apis.NewBadRequestError("failed to start OpenID4VP workflow", err)
			}

			return e.JSON(http.StatusOK, map[string]bool{
				"started": true,
			})
		})

		se.Router.POST("/wallet-test/confirm-success", func(e *core.RequestEvent) error {
			var request struct {
				WorkflowID string `json:"workflow_id"`
			}
			if err := json.NewDecoder(e.Request.Body).Decode(&request); err != nil {
				return apis.NewBadRequestError("Invalid JSON input", err)
			}
			data := openid4vp_workflow.SignalData{
				Success: true,
			}
			err := c.SignalWorkflow(context.Background(), request.WorkflowID, "", "wallet-test-signal", data)
			if err != nil {
				return apis.NewBadRequestError("Failed to send success signal", err)
			}

			return e.JSON(http.StatusOK, map[string]string{"message": "Workflow completed successfully"})
		})

		se.Router.POST("/wallet-test/notify-failure", func(e *core.RequestEvent) error {
			var request struct {
				WorkflowID string `json:"workflow_id"`
				Reason     string `json:"reason"`
			}
			if err := json.NewDecoder(e.Request.Body).Decode(&request); err != nil {
				return apis.NewBadRequestError("Invalid JSON input", err)
			}
			data := openid4vp_workflow.SignalData{
				Success: false,
				Reason:  request.Reason,
			}
			err := c.SignalWorkflow(context.Background(), request.WorkflowID, "", "wallet-test-signal", data)
			if err != nil {
				return apis.NewBadRequestError("Failed to send failure signal", err)
			}

			return e.JSON(http.StatusOK, map[string]string{"message": "Test failed", "reason": request.Reason})
		})
		return se.Next()
	})
}
