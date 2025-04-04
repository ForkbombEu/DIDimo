package pb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP"
	openid4vp_workflow "github.com/forkbombeu/didimo/pkg/OpenID4VP/workflow"
	credential_workflow "github.com/forkbombeu/didimo/pkg/credential_issuer/workflow"
	temporalclient "github.com/forkbombeu/didimo/pkg/internal/temporal_client"
	engine "github.com/forkbombeu/didimo/pkg/template_engine"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"google.golang.org/protobuf/types/known/durationpb"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
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
			c, err := temporalclient.GetTemporalClient("default")

			if err != nil {
				return err
			}
			we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, credential_workflow.CredentialWorkflow, workflowInput)
			if err != nil {
				return fmt.Errorf("error starting workflow for URL %s: %v", req.URL, err)
			}

			var result credential_workflow.CredentialWorkflowResponse
			err = we.Get(context.Background(), &result)
			if err != nil {
				return fmt.Errorf("error running workflow for URL %s: %v", req.URL, err)
			}

			log.Printf("Workflow completed successfully for URL %s: %s", req.URL, result.Message)

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
		se.Router.POST("/api/openid4vp-test", func(e *core.RequestEvent) error {
			var req OpenID4VPRequest
			appURL := app.Settings().Meta.AppURL
			if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
				return apis.NewBadRequestError("invalid JSON input", err)
			}

			// Start the workflow
			err := OpenID4VP.StartWorkflow(req.Input, req.UserMail, appURL, "default")
			if err != nil {
				return apis.NewBadRequestError("failed to start OpenID4VP workflow", err)
			}

			return e.JSON(http.StatusOK, map[string]bool{
				"started": true,
			})
		})

		se.Router.POST("/api/{protocol}/{author}/save-variables-and-start", func(e *core.RequestEvent) error {
			var req map[string]struct {
				Format string      `json:"format"`
				Data   interface{} `json:"data"`
			}
			appURL := app.Settings().Meta.AppURL

			User := e.Auth.Id
			email := e.Auth.GetString("email")

			authOrgCollection, err := e.App.FindCollectionByNameOrId("orgAuthorizations")
			if err != nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations collection", err)
			}
			if authOrgCollection == nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations collection", nil)
			}
			authOrgRecord, err := e.App.FindFirstRecordByFilter(authOrgCollection.Id, "user={:user}", dbx.Params{"user": User})
			if err != nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations record", err)
			}
			if authOrgRecord == nil {
				return apis.NewBadRequestError("user is not authorized to access this organization", nil)
			}
			organization := authOrgRecord.GetString("organization")
			if organization == "" {
				return apis.NewBadRequestError("organization is empty", nil)
			}

			protocol := e.Request.PathValue("protocol")
			author := e.Request.PathValue("author")

			if protocol == "" || author == "" {
				return apis.NewBadRequestError("protocol and author are required", nil)
			}
			if protocol == "openid4vp_wallet" {
				protocol = "OpenID4VP_Wallet"
			}
			if protocol == "openid4vci_wallet" {
				protocol = "OpenID4VCI_Wallet"
			}
			if author == "openid_foundation" {
				author = "OpenID_foundation"
			}
			filepath := os.Getenv("ROOT_DIR") + "/config_templates/" + protocol + "/" + author + "/"

			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				return apis.NewBadRequestError("directory does not exist for test "+protocol+"/"+author, err)
			}

			if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
				return apis.NewBadRequestError("invalid JSON input", err)
			}

			for testName, testData := range req {
				if testData.Format == "json" {
					jsonData, ok := testData.Data.(string)
					if !ok {
						return apis.NewBadRequestError("invalid JSON format for test "+testName, nil)
					}

					var parsedData OpenID4VP.OpenID4VPTestInputFile
					if err := json.Unmarshal([]byte(jsonData), &parsedData); err != nil {
						return apis.NewBadRequestError("failed to parse JSON for test "+testName, err)
					}

					err := OpenID4VP.StartWorkflow(parsedData, User, appURL, organization)
					if err != nil {
						return apis.NewBadRequestError("failed to start OpenID4VP workflow for test "+testName, err)
					}
				} else if testData.Format == "variables" {
					variables, ok := testData.Data.(map[string]interface{})
					values := make(map[string]interface{})
					config_values, err := app.FindCollectionByNameOrId("config_values")
					if err != nil {
						return err
					}
					if !ok {
						return apis.NewBadRequestError("invalid variables format for test "+testName, nil)
					}
					for credimiId, variable := range variables {
						v, ok := variable.(map[string]interface{})
						if !ok {
							return apis.NewBadRequestError("invalid variable format for test "+testName, nil)
						}

						fieldName, ok := v["fieldName"].(string)
						if !ok {
							return apis.NewBadRequestError("invalid fieldName format for test "+testName, nil)
						}

						record := core.NewRecord(config_values)
						record.Set("credimi_id", credimiId)
						record.Set("value", v["value"])
						record.Set("field_name", fieldName)
						record.Set("template_path", testName)
						if err := app.Save(record); err != nil {
							log.Println("Failed to create related record:", err)
						}
						values[fieldName] = v["value"]
					}

					template, err := os.Open(filepath + testName)
					if err != nil {
						return apis.NewBadRequestError("failed to open template for test "+testName, err)
					}
					defer template.Close()

					templateFile, err := os.Open(filepath + testName)
					if err != nil {
						return apis.NewBadRequestError("failed to open template for test "+testName, err)
					}
					defer templateFile.Close()

					renderedTemplate, err := engine.RenderTemplate(templateFile, values)
					if err != nil {
						return apis.NewInternalServerError("failed to render template for test "+testName, err)
					}

					var parsedVariant OpenID4VP.OpenID4VPTestInputFile
					errParsing := json.Unmarshal([]byte(renderedTemplate), &parsedVariant)
					if errParsing != nil {
						return apis.NewBadRequestError("failed to unmarshal JSON for test "+testName, errParsing)
					}

					err = OpenID4VP.StartWorkflow(OpenID4VP.OpenID4VPTestInputFile{
						Variant: json.RawMessage(parsedVariant.Variant),
						Form:    parsedVariant.Form,
					}, email, appURL, organization)
					if err != nil {
						return apis.NewBadRequestError("failed to start OpenID4VP workflow for test "+testName, err)
					}

				} else {
					return apis.NewBadRequestError("unsupported format for test "+testName, nil)
				}
			}

			return e.JSON(http.StatusOK, map[string]bool{
				"started": true,
			})
		}).Bind(apis.RequireAuth())

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
			c, err := temporalclient.GetTemporalClient("default")
			if err != nil {
				return err
			}
			err = c.SignalWorkflow(context.Background(), request.WorkflowID, "", "wallet-test-signal", data)
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
			c, err := temporalclient.GetTemporalClient("default")
			if err != nil {
				return err
			}
			err = c.SignalWorkflow(context.Background(), request.WorkflowID, "", "wallet-test-signal", data)
			if err != nil {
				return apis.NewBadRequestError("Failed to send failure signal", err)
			}

			return e.JSON(http.StatusOK, map[string]string{"message": "Test failed", "reason": request.Reason})
		})
		return se.Next()
	})
}

func HookUpdateCredentialsIssuers(app *pocketbase.PocketBase) {
	app.OnRecordAfterUpdateSuccess().BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Collection().Name != "features" || e.Record.Get("name") != "updateIssuers" {
			return nil
		}
		if e.Record.Get("active") == false {
			return nil
		}
		envVariables := e.Record.Get("envVariables")
		if envVariables == nil {
			return nil
		}
		result := struct {
			Interval string `json:"interval"`
		}{}
		errJson := json.Unmarshal(e.Record.Get("envVariables").(types.JSONRaw), &result)
		if errJson != nil {
			log.Fatal(errJson)
		}
		if result.Interval == "" {
			return nil
		}
		var interval time.Duration
		switch result.Interval {
		case "every_minute":
			interval = time.Minute
		case "hourly":
			interval = time.Hour
		case "daily":
			interval = time.Hour * 24
		case "weekly":
			interval = time.Hour * 24 * 7
		case "monthly":
			interval = time.Hour * 24 * 30
		default:
			interval = time.Hour
		}
		workflowID := "schedule_workflow_id" + fmt.Sprintf("%d", time.Now().Unix())
		scheduleID := "schedule_id" + fmt.Sprintf("%d", time.Now().Unix())
		ctx := context.Background()

		temporalClient, err := client.Dial(client.Options{
			HostPort: client.DefaultHostPort,
		})
		if err != nil {
			log.Fatalln("Unable to create Temporal Client", err)
		}
		defer temporalClient.Close()

		scheduleHandle, err := temporalClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: scheduleID,
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: interval,
					},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        workflowID,
				Workflow:  credential_workflow.FetchIssuersWorkflow,
				TaskQueue: credential_workflow.FetchIssuersTaskQueue,
			},
		})
		if err != nil {
			log.Fatalln("Unable to create schedule", err)
		}
		_, _ = scheduleHandle.Describe(ctx)

		return nil
	})
}

func RouteWorkflow(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/workflows", func(e *core.RequestEvent) error {
			authRecord := e.Auth

			log.Println("AuthRecord: ", authRecord)

			ownerRoleRecord, err := e.App.FindFirstRecordByFilter("orgRoles", "name='owner'")
			if err != nil {
				return apis.NewInternalServerError("failed to find owner role", err)
			}

			orgAuthRecord, err := e.App.FindRecordsByFilter("orgAuthorizations", "user={:user} && role={:role}", "", 0, 0, dbx.Params{"user": authRecord.Id, "role": ownerRoleRecord.Id})
			if err != nil || orgAuthRecord == nil {
				return apis.NewUnauthorizedError("User is not authorized to access this organization", err)
			}
			namespace := orgAuthRecord[0].GetString("organization")

			c, err := temporalclient.GetTemporalClient(namespace)
			if err != nil {
				return apis.NewInternalServerError("unable to create client", err)
			}
			list, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
				Namespace: namespace,
			})
			log.Printf("ListWorkflowExecutions: %v", list)
			if err != nil {
				log.Println("Error listing workflows:", err)
				return apis.NewInternalServerError("failed to list workflows", err)
			}

			// listJSON, err := json.Marshal(list)
			// if err != nil {
			// 	return apis.NewInternalServerError("failed to marshal workflow list", err)
			// }
			return e.JSON(http.StatusOK, list)
		}).Bind(apis.RequireAuth())

		se.Router.GET("/api/workflows/{workflowId}/{runId}", func(e *core.RequestEvent) error {
			workflowId := e.Request.PathValue("workflowId")
			if workflowId == "" {
				return apis.NewBadRequestError("workflowId is required", nil)
			}
			runId := e.Request.PathValue("runId")
			if runId == "" {
				return apis.NewBadRequestError("runId is required", nil)
			}
			authRecord := e.Auth
			orgAuthCollection, err := e.App.FindCollectionByNameOrId("orgAuthorizations")
			if err != nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations collection", err)
			}
			if orgAuthCollection == nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations collection", nil)
			}
			orgAuthRecord, err := e.App.FindFirstRecordByFilter(orgAuthCollection.Id, "user={:user}", dbx.Params{"user": authRecord.Id})
			if err != nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations record", err)
			}
			if orgAuthRecord == nil {
				return apis.NewBadRequestError("user is not authorized to access this organization", nil)
			}
			namespace := orgAuthRecord.GetString("organization")
			if namespace == "" {
				return apis.NewBadRequestError("organization is empty", nil)
			}
			c, err := temporalclient.GetTemporalClient(namespace)
			if err != nil {
				return apis.NewInternalServerError("unable to create client", err)
			}
			workflowExecution, err := c.DescribeWorkflowExecution(context.Background(), workflowId, runId)
			if err != nil {
				return apis.NewInternalServerError("failed to describe workflow execution", err)
			}
			if workflowExecution == nil {
				return apis.NewNotFoundError("workflow execution not found", nil)
			}

			workflowExecutionJSON, err := json.Marshal(workflowExecution)
			if err != nil {
				return apis.NewInternalServerError("failed to marshal workflow execution", err)
			}

			return e.JSON(http.StatusOK, workflowExecutionJSON)
		}).Bind(apis.RequireAuth())

		se.Router.GET("/api/workflows/{workflowId}/{runId}/history", func(e *core.RequestEvent) error {
			authRecord := e.Auth
			orgAuthCollection, err := e.App.FindCollectionByNameOrId("orgAuthorizations")
			if err != nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations collection", err)
			}
			if orgAuthCollection == nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations collection", nil)
			}
			orgAuthRecord, err := e.App.FindFirstRecordByFilter(orgAuthCollection.Id, "user={:user}", dbx.Params{"user": authRecord.Id})
			if err != nil {
				return apis.NewBadRequestError("failed to find orgAuthorizations record", err)
			}
			if orgAuthRecord == nil {
				return apis.NewBadRequestError("user is not authorized to access this organization", nil)
			}
			namespace := orgAuthRecord.GetString("organization")
			if namespace == "" {
				return apis.NewBadRequestError("organization is empty", nil)
			}
			workflowId := e.Request.PathValue("workflowId")
			if workflowId == "" {
				return apis.NewBadRequestError("workflowId is required", nil)
			}
			runId := e.Request.PathValue("runId")
			if runId == "" {
				return apis.NewBadRequestError("runId is required", nil)
			}
			c, err := temporalclient.GetTemporalClient(namespace)
			if err != nil {
				return apis.NewInternalServerError("unable to create client", err)
			}
			history := c.GetWorkflowHistory(context.Background(), workflowId, runId, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)

			historyJSON, err := json.Marshal(history)
			if err != nil {
				return apis.NewInternalServerError("failed to marshal workflow history", err)
			}

			return e.JSON(http.StatusOK, string(historyJSON))
		}).Bind(apis.RequireAuth())

		return se.Next()
	})
}

func HookAtUserCreation(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess().BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Collection().Name != "users" {
			return nil
		}
		log.Println("User created:", e.Record.Id)
		user := e.Record
		return createNamespaceForUser(e, user)
	})
}

func createNamespaceForUser(e *core.RecordEvent, user *core.Record) error {
	log.Println("Creating namespace for user:", user.Id)
	err := e.App.RunInTransaction(func(txApp core.App) error {
		orgCollection, err := txApp.FindCollectionByNameOrId("organizations")
		log.Println("OrgCollection: ", orgCollection.Id)
		if err != nil {
			return apis.NewInternalServerError("failed to find organizations collection", err)
		}

		newOrg := core.NewRecord(orgCollection)
		newOrg.Set("name", user.Id)
		txApp.Save(newOrg)

		ownerRoleRecord, err := txApp.FindFirstRecordByFilter("orgRoles", "name='owner'")
		if err != nil {
			return apis.NewInternalServerError("failed to find owner role", err)
		}

		orgAuthCollection, err := txApp.FindCollectionByNameOrId("orgAuthorizations")
		if err != nil {
			return apis.NewInternalServerError("failed to find orgAuthorizations collection", err)
		}
		newOrgAuth := core.NewRecord(orgAuthCollection)
		newOrgAuth.Set("user", e.Record.Id)
		newOrgAuth.Set("organization", newOrg.Id)
		newOrgAuth.Set("role", ownerRoleRecord.Id)
		txApp.Save(newOrgAuth)
		retention := durationpb.New(24 * time.Hour)

		namespace := newOrg.Id
		client, err := client.NewNamespaceClient(client.Options{})

		if err != nil {
			log.Println("Error creating Temporal client:", err)
			return apis.NewInternalServerError("failed to create Temporal client", err)
		}
		defer client.Close()
		// Check if the namespace already exists
		_, err = client.Describe(context.Background(), namespace)
		if err == nil {
			log.Println("Namespace already exists for user:", user.Id)
			return nil
		}
		err = client.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
			Namespace:                        namespace,
			WorkflowExecutionRetentionPeriod: retention,
		})
		if err != nil {
			log.Println("Error registering namespace:", err)
			return apis.NewInternalServerError("failed to register namespace", err)
		}
		log.Println("Namespace created successfully for user:", user.Id)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
