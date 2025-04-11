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
	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
	"github.com/forkbombeu/didimo/pkg/workflow_engine/workflows"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"

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
				TaskQueue: credential_workflow.CredentialsTaskQueue,
			}
			c, err := temporalclient.GetTemporalClient()
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
			template, err := os.ReadFile(workflows.OpenIDNetStepCITemplatePath)
			if err != nil {
				return apis.NewBadRequestError("failed to open template file: %w", err)
			}
			// Start the workflow
			input := workflowengine.WorkflowInput{
				Payload: map[string]any{
					"variant":   string(req.Input.Variant),
					"form":      req.Input.Form,
					"user_mail": req.UserMail,
					"app_url":   appURL,
				},
				Config: map[string]any{
					"template": string(template),
				},
			}
			var w workflows.OpenIDNetWorkflow
			_, err = w.Start(input)
			if err != nil {
				return apis.NewBadRequestError("failed to start openidnet wallet workflow", err)
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
					stepCItemplate, err := os.ReadFile(workflows.OpenIDNetStepCITemplatePath)
					if err != nil {
						return apis.NewBadRequestError("failed to open template file: %w", err)
					}
					// Start the workflow
					input := workflowengine.WorkflowInput{
						Payload: map[string]any{
							"variant":   string(parsedData.Variant),
							"form":      parsedData.Form,
							"user_mail": "test@credimi.io",
							"app_url":   appURL,
						},
						Config: map[string]any{
							"template": string(stepCItemplate),
						},
					}
					var w workflows.OpenIDNetWorkflow
					_, err = w.Start(input)
					if err != nil {
						return apis.NewBadRequestError("failed to start openidnet wallet for test "+testName, err)
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
					stepCItemplate, err := os.ReadFile(workflows.OpenIDNetStepCITemplatePath)
					if err != nil {
						return apis.NewBadRequestError("failed to open template file: %w", err)
					}
					// Start the workflow
					input := workflowengine.WorkflowInput{
						Payload: map[string]any{
							"variant":   string(parsedVariant.Variant),
							"form":      parsedVariant.Form,
							"user_mail": "test@credimi.io",
							"app_url":   appURL,
						},
						Config: map[string]any{
							"template": string(stepCItemplate),
						},
					}
					var w workflows.OpenIDNetWorkflow
					_, err = w.Start(input)

					if err != nil {
						return apis.NewBadRequestError("failed to start openidnet wallet for test "+testName, err)
					}

				} else {
					return apis.NewBadRequestError("unsupported format for test "+testName, nil)
				}
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
			c, err := temporalclient.GetTemporalClient()
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

func RouteWorkflowList(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/workflows", func(e *core.RequestEvent) error {
			namespace := e.Request.URL.Query().Get("namespace")
			if namespace == "" {
				return apis.NewBadRequestError("namespace is required", nil)
			}

			authRecord := e.Auth

			orgRecord, err := e.App.FindFirstRecordByFilter("organizations", "name={:name}", dbx.Params{"name": namespace})
			if err != nil || orgRecord == nil {
				return apis.NewBadRequestError("Organization not found", err)
			}

			orgAuthRecord, err := e.App.FindRecordsByFilter("orgAuthorizations", "user={:user} && organization={:organization}", "", 0, 0, dbx.Params{"user": authRecord.Id, "organization": orgRecord.Id})
			if err != nil || orgAuthRecord == nil {
				return apis.NewUnauthorizedError("User is not authorized to access this organization", err)
			}
			if len(orgAuthRecord) > 1 {
				return apis.NewUnauthorizedError("User is not authorized to access this organization", nil)
			}

			c, err := temporalclient.GetTemporalClient()
			if err != nil {
				log.Fatalln("unable to create client", err)
			}
			list, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
				Namespace: namespace,
			})
			if err != nil {
				return apis.NewInternalServerError("failed to list workflows", err)
			}

			return e.JSON(http.StatusOK, list)
		}).Bind(apis.RequireAuth())
		return se.Next()
	})
}
