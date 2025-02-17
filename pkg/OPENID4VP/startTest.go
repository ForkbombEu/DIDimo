package main

import (
	"context"
	"fmt"
	"log"

	"github.com/forkbombeu/didimo/pkg/OPENID4VP/testdata"
	"github.com/forkbombeu/didimo/pkg/OPENID4VP/workflow"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	// Initialize the Temporal client using Dial (recommended method)
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// Prepare input for the workflow
	input := workflow.WorkflowInput{
		Variant: `{"credential_format":"sd_jwt_vc","client_id_scheme":"did","request_method":"request_uri_signed","response_mode":"direct_post"}`,
		JSONPayload: testdata.JSONPayload{
			Alias:       "TEST_from_rest",
			Description: "TEST FROM StepCI Workflow",
			Server: testdata.Server{
				AuthorizationEndpoint: "openid-vc://",
			},
			Client: testdata.Client{
				ClientID: "did:web:app.altme.io:issuer",
				PresentationDefinition: testdata.PresentationDefinition{
					ID: "two_sd_jwt",
					InputDescriptors: []testdata.InputDescriptor{
						{
							ID: "pid_credential",
							Constraints: testdata.Constraints{
								Fields: []testdata.Field{
									{
										Path: []string{"$.vct"},
										Filter: map[string]string{
											"type":  "string",
											"const": "urn:eu.europa.ec.eudi:pid:1",
										},
									},
								},
							},
							Format: testdata.Format{
								VCSDJWT: map[string][]string{
									"kb-jwt_alg_values": {"ES256", "ES256K", "EdDSA"},
									"sd-jwt_alg_values": {"ES256", "ES256K", "EdDSA"},
								},
							},
						},
					},
				},
				JWKS: testdata.JWKS{
					Keys: []testdata.JWKKey{
						{
							Kty: "EC",
							Alg: "ES256",
							Crv: "P-256",
							D:   "GSbo9TpmGaLgxxO6RNx6QnvcfykQJS7vUVgTe8vy9W0",
							X:   "m5uKsE35t3sP7gjmirUewufx2Gt2n6J7fSW68apB2Lo",
							Y:   "-V54TpMI8RbpB40hbAocIjnaHX5WP6NHjWkHfdCSAyU",
						},
					},
				},
			},
		},
	}

	// Define the options for the workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "OpenIDTestWorkflow" + uuid.NewString(), // Unique ID for this workflow
		TaskQueue: "openid-test-task-queue",                // The task queue where the worker listens
	}

	// Start the workflow execution
	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.OpenIDTestWorkflow, input)
	if err != nil {
		log.Fatalf("failed to start workflow: %v", err)
	}

	// Wait for the workflow result
	var result string
	if err := workflowRun.Get(context.Background(), &result); err != nil {
		log.Fatalf("failed to get workflow result: %v", err)
	}

	// Print the result from the workflow
	fmt.Printf("Workflow result: %s\n", result)
}
