// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

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
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/encoding/protojson"
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
				newRecord.Set("owner", e.Auth.Id)
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
				TaskQueue: credential_workflow.CredentialsTaskQueue,
			}

			c, err := temporalclient.GetTemporalClient()
			defer c.Close()
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
			//
			// providers, err := app.FindCollectionByNameOrId("services")
			// if err != nil {
			// 	return err
			// }
			//
			// newRecord := core.NewRecord(providers)
			// newRecord.Set("credential_issuers", issuerID)
			// newRecord.Set("name", "TestName")
			// // Save the new record in providers
			// if err := app.Save(newRecord); err != nil {
			// 	return err
			// }
			return e.JSON(http.StatusOK, map[string]string{
				"credentialIssuerUrl": req.URL,
			})
		}).Bind(apis.RequireAuth())
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
			err := OpenID4VP.StartWorkflow(req.Input, req.UserMail, appURL)
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

			namespace, err := getUserNamespace(app, User)
			if err != nil {
				return apis.NewBadRequestError("failed to get user namespace", err)
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
				memo := map[string]interface{}{
					"test":     testName,
					"standard": protocol,
					"author":   author,
				}
				if testData.Format == "json" {
					jsonData, ok := testData.Data.(string)
					if !ok {
						return apis.NewBadRequestError("invalid JSON format for test "+testName, nil)
					}

					var parsedData OpenID4VP.OpenID4VPTestInputFile
					if err := json.Unmarshal([]byte(jsonData), &parsedData); err != nil {
						return apis.NewBadRequestError("failed to parse JSON for test "+testName, err)
					}

					err := OpenID4VP.StartWorkflowWithNamespaceAndMemo(parsedData, User, appURL, namespace, memo)
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
					defer template.Close()
					if err != nil {
						return apis.NewBadRequestError("failed to open template for test "+testName, err)
					}

					templateFile, err := os.Open(filepath + testName)
					defer templateFile.Close()
					if err != nil {
						return apis.NewBadRequestError("failed to open template for test "+testName, err)
					}

					renderedTemplate, err := engine.RenderTemplate(templateFile, values)
					if err != nil {
						return apis.NewInternalServerError("failed to render template for test "+testName, err)
					}

					var parsedVariant OpenID4VP.OpenID4VPTestInputFile
					errParsing := json.Unmarshal([]byte(renderedTemplate), &parsedVariant)
					if errParsing != nil {
						return apis.NewBadRequestError("failed to unmarshal JSON for test "+testName, errParsing)
					}

					err = OpenID4VP.StartWorkflowWithNamespaceAndMemo(OpenID4VP.OpenID4VPTestInputFile{
						Variant: json.RawMessage(parsedVariant.Variant),
						Form:    parsedVariant.Form,
					}, email, appURL, namespace, memo)
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

			c, err := temporalclient.GetTemporalClient()
			defer c.Close()
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

			c, err := temporalclient.GetTemporalClient()
			defer c.Close()

			if err != nil {
				return err
			}

			err = c.SignalWorkflow(context.Background(), request.WorkflowID, "", "wallet-test-signal", data)
			if err != nil {
				return apis.NewBadRequestError("Failed to send failure signal", err)
			}

			return e.JSON(http.StatusOK, map[string]string{"message": "Test failed", "reason": request.Reason})
		})

		se.Router.POST("/wallet-test/send-log-update", func(e *core.RequestEvent) error {
			var logData openid4vp_workflow.LogUpdateRequest
			if err := json.NewDecoder(e.Request.Body).Decode(&logData); err != nil {
				return apis.NewBadRequestError("invalid JSON input", err)
			}
			if err := notifyLogsUpdate(app, logData.WorkflowID+"openid4vp-wallet-logs", logData.Logs); err != nil {
				return apis.NewBadRequestError("failed to send real-time log update", err)
			}

			return e.JSON(http.StatusOK, map[string]string{
				"message": "Log update sent successfully",
			})
		})
		se.Router.POST("/wallet-test/send-log-update-start", func(e *core.RequestEvent) error {
			var request struct {
				WorkflowID string `json:"workflow_id"`
			}
			if err := json.NewDecoder(e.Request.Body).Decode(&request); err != nil {
				return apis.NewBadRequestError("Invalid JSON input", err)
			}

			c, err := temporalclient.GetTemporalClient()
			defer c.Close()
			if err != nil {
				return err
			}

			err = c.SignalWorkflow(context.Background(), request.WorkflowID+"-log", "", "wallet-test-start-log-update", struct{}{})
			if err != nil {
				return apis.NewBadRequestError("Failed to send start logs update signal", err)
			}

			return e.JSON(http.StatusOK, map[string]string{"message": "Realtime Logs update started successfully"})
		})

		return se.Next()
	})
}
func notifyLogsUpdate(app core.App, subscription string, data []map[string]any) error {

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	message := subscriptions.Message{
		Name: subscription,
		Data: rawData,
	}
	clients := app.SubscriptionsBroker().Clients()
	for _, client := range clients {
		if client.HasSubscription(subscription) {
			client.Send(message)
		}
	}
	return nil
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
		defer temporalClient.Close()

		if err != nil {
			log.Fatalln("Unable to create Temporal Client", err)
		}
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
			namespace, err := getUserNamespace(e.App, authRecord.Id)
			if err != nil {
				return apis.NewUnauthorizedError("User is not authorized to access this organization", err)
			}
			c, err := temporalclient.GetTemporalClientWithNamespace(namespace)
			defer c.Close()
			if err != nil {
				return apis.NewInternalServerError("unable to create client", err)
			}
			list, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
				Namespace: namespace,
			})
			if err != nil {
				log.Println("Error listing workflows:", err)
				return apis.NewInternalServerError("failed to list workflows", err)
			}
			listJson, err := protojson.Marshal(list)
			if err != nil {
				return apis.NewInternalServerError("failed to marshal workflow list", err)
			}
			finalJson := make(map[string]interface{})
			err = json.Unmarshal(listJson, &finalJson)
			if err != nil {
				return apis.NewInternalServerError("failed to unmarshal workflow list", err)
			}
			if finalJson["executions"] == nil {
				finalJson["executions"] = []map[string]interface{}{}
			}
			return e.JSON(http.StatusOK, finalJson)
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

			namespace, err := getUserNamespace(e.App, authRecord.Id)
			if err != nil {
				return apis.NewUnauthorizedError("User is not authorized to access this organization", err)
			}
			if namespace == "" {
				return apis.NewBadRequestError("organization is empty", nil)
			}

			c, err := temporalclient.GetTemporalClientWithNamespace(namespace)
			defer c.Close()
			if err != nil {
				return apis.NewInternalServerError("unable to create client", err)
			}
			workflowExecution, err := c.DescribeWorkflowExecution(context.Background(), workflowId, runId)
			if err != nil {
				return apis.NewInternalServerError("failed to describe workflow execution", err)
			}
			weJson, err := protojson.Marshal(workflowExecution)
			if err != nil {
				return apis.NewInternalServerError("failed to marshal workflow execution", err)
			}
			finalJson := make(map[string]interface{})
			err = json.Unmarshal(weJson, &finalJson)
			if err != nil {
				return apis.NewInternalServerError("failed to unmarshal workflow execution", err)
			}
			return e.JSON(http.StatusOK, finalJson)
		}).Bind(apis.RequireAuth())

		se.Router.GET("/api/workflows/{workflowId}/{runId}/history", func(e *core.RequestEvent) error {
			authRecord := e.Auth

			namespace, err := getUserNamespace(e.App, authRecord.Id)
			if err != nil {
				return apis.NewBadRequestError("failed to get user namespace", err)
			}

			workflowId := e.Request.PathValue("workflowId")
			if workflowId == "" {
				return apis.NewBadRequestError("workflowId is required", nil)
			}
			runId := e.Request.PathValue("runId")
			if runId == "" {
				return apis.NewBadRequestError("runId is required", nil)
			}

			c, err := temporalclient.GetTemporalClientWithNamespace(namespace)
			defer c.Close()
			if err != nil {
				return apis.NewInternalServerError("unable to create client", err)
			}

			historyIterator := c.GetWorkflowHistory(context.Background(), workflowId, runId, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
			var history []map[string]interface{}
			for historyIterator.HasNext() {
				event, err := historyIterator.Next()
				if err != nil {
					return apis.NewInternalServerError("failed to iterate workflow history", err)
				}
				eventData, err := protojson.Marshal(event)
				if err != nil {
					return apis.NewInternalServerError("failed to marshal history event", err)
				}
				var eventMap map[string]interface{}
				if err := json.Unmarshal(eventData, &eventMap); err != nil {
					return apis.NewInternalServerError("failed to unmarshal history event", err)
				}
				history = append(history, eventMap)
			}

			return e.JSON(http.StatusOK, history)
		}).Bind(apis.RequireAuth())

		return se.Next()
	})
}

func HookAtUserCreation(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		err := addUserToDefaultOrganization(e)
		if err != nil {
			return err
		}
		return e.Next()
	})
}

func getUserNamespace(app core.App, userId string) (string, error) {
	orgAuthCollection, err := app.FindCollectionByNameOrId("orgAuthorizations")
	if err != nil {
		return "", apis.NewInternalServerError("failed to find orgAuthorizations collection", err)
	}

	authOrgRecords, err := app.FindRecordsByFilter(orgAuthCollection.Id, "user={:user}", "", 0, 0, dbx.Params{"user": userId})
	if err != nil {
		return "", apis.NewInternalServerError("failed to find orgAuthorizations records", err)
	}
	if len(authOrgRecords) == 0 {
		return "", apis.NewInternalServerError("user is not authorized to access any organization", nil)
	}

	ownerRoleRecord, err := app.FindFirstRecordByFilter("orgRoles", "name='owner'")
	if err != nil {
		return "", apis.NewInternalServerError("failed to find owner role", err)
	}

	if len(authOrgRecords) > 1 {
		for _, record := range authOrgRecords {
			if record.GetString("role") == ownerRoleRecord.Id {
				return record.GetString("organization"), nil
			}
		}
	}
	if authOrgRecords[0].GetString("role") == ownerRoleRecord.Id {
		return authOrgRecords[0].GetString("organization"), nil
	}
	return "default", nil
}

func addUserToDefaultOrganization(e *core.RecordEvent) error {
	user := e.Record
	errTx := e.App.RunInTransaction(func(txApp core.App) error {
		orgCollection, err := txApp.FindCollectionByNameOrId("organizations")
		if err != nil {
			return apis.NewInternalServerError("failed to find organizations collection", err)
		}
		defaultOrgRecord, err := txApp.FindFirstRecordByFilter(orgCollection.Id, "name='default'")
		if err != nil {
			return apis.NewInternalServerError("failed to find default organization", err)
		}
		if defaultOrgRecord == nil {
			return apis.NewInternalServerError("default organization not found", nil)
		}
		orgAuthCollection, err := txApp.FindCollectionByNameOrId("orgAuthorizations")
		if err != nil {
			return apis.NewInternalServerError("failed to find orgAuthorizations collection", err)
		}
		newOrgAuth := core.NewRecord(orgAuthCollection)
		newOrgAuth.Set("user", user.Id)
		newOrgAuth.Set("organization", defaultOrgRecord.Id)
		memberRoleRecord, err := txApp.FindFirstRecordByFilter("orgRoles", "name='member'")
		if err != nil {
			return apis.NewInternalServerError("failed to find owner role", err)
		}
		newOrgAuth.Set("role", memberRoleRecord.Id)
		err = txApp.Save(newOrgAuth)
		if err != nil {
			return apis.NewInternalServerError("failed to save orgAuthorization record", err)
		}
		return nil
	})

	if errTx != nil {
		return apis.NewInternalServerError("failed to add user to default organization", errTx)
	}
	return nil
}

// This function will be used when user will claim the organization
func createNamespaceForUser(e *core.RecordEvent, user *core.Record) error {

	err := e.App.RunInTransaction(func(txApp core.App) error {
		orgCollection, err := txApp.FindCollectionByNameOrId("organizations")
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
		newOrgAuth.Set("user", user.Id)
		newOrgAuth.Set("organization", newOrg.Id)
		newOrgAuth.Set("role", ownerRoleRecord.Id)
		txApp.Save(newOrgAuth)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
