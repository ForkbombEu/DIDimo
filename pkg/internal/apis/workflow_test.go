// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package apis

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/tests"
	"go.temporal.io/sdk/testsuite"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
        suite.Suite
        testsuite.WorkflowTestSuite

        env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
        s.env = s.NewTestWorkflowEnvironment()
}

const testDataDir = "../../../test_pb_data"
// const rootDir = "/Users/alcibiade/dyne/temp/DIDimo"

func generateToken(collectionNameOrId string, email string) (string, error) {
	app, err := tests.NewTestApp(testDataDir)
	if err != nil {
		return "", err
	}
	defer app.Cleanup()

	record, err := app.FindAuthRecordByEmail(collectionNameOrId, email)
	if err != nil {
		return "", err
	}

	return record.NewAuthToken()
}

func TestAddOpenID4VPTestEndpoints_RoutesRegistered(t *testing.T) {
	godotenv.Load("../../../.env")

	app, err := tests.NewTestApp(testDataDir)
	defer app.Cleanup()
	require.NoError(t, err)

	setupTestApp := func(t testing.TB) *tests.TestApp {
		testApp, err := tests.NewTestApp(testDataDir)
		if err != nil {
			t.Fatal(err)
		}
		AddComplianceChecks(testApp)

		return testApp
	}

	authToken, err := generateToken("users", "userA@example.org")

	scenarios := []tests.ApiScenario{
		{
			Name:           "OpenID4VP Test - Valid Request",
			Method:         "POST",
			URL:            "/api/compliance/check/OpenID4VP_Wallet/OpenID_Foundation/save-variables-and-start",
			Body:           strings.NewReader(`{"sd_jwt_vc:did:request_uri_signed:direct_post.json":{"format":"variables","data":{"oid_description":{"type":"string","value":"jikj","fieldName":"description"},"oid_alias":{"type":"string","value":"knnkn","fieldName":"testalias"},"oid_client_id":{"type":"string","value":"did:web:app.altme.io:issuer","fieldName":"client_id"},"oid_client_jwks":{"type":"object","value":"{\n    \"keys\": [\n        {\n            \"kty\": \"EC\",\n            \"alg\": \"ES256\",\n            \"crv\": \"P-256\",\n            \"d\": \"GSbo9TpmGaLgxxO6RNx6QnvcfykQJS7vUVgTe8vy9W0\",\n            \"x\": \"m5uKsE35t3sP7gjmirUewufx2Gt2n6J7fSW68apB2Lo\",\n            \"y\": \"-V54TpMI8RbpB40hbAocIjnaHX5WP6NHjWkHfdCSAyU\"\n        }\n    ]\n}","fieldName":"jwks"},"oid_client_presentation_definition":{"type":"object","value":"{\n    \"id\": \"two_sd_jwt\",\n    \"input_descriptors\": [\n        {\n            \"constraints\": {\n                \"fields\": [\n                    {\n                        \"filter\": {\n                            \"const\": \"urn:eu.europa.ec.eudi:pid:1\",\n                            \"type\": \"string\"\n                        },\n                        \"path\": [\n                            \"$.vct\"\n                        ]\n                    }\n                ]\n            },\n            \"format\": {\n                \"vc+sd-jwt\": {\n                    \"kb-jwt_alg_values\": [\n                        \"ES256\",\n                        \"ES256K\",\n                        \"EdDSA\"\n                    ],\n                    \"sd-jwt_alg_values\": [\n                        \"ES256\",\n                        \"ES256K\",\n                        \"EdDSA\"\n                    ]\n                }\n            },\n            \"id\": \"pid_credential\"\n        }\n    ]\n}","fieldName":"presentation_definition"}}}}`),
			Headers:        map[string]string{"Content-Type": "application/json", "Authorization": authToken},
			Delay:          0,
			Timeout:        5 * time.Second,
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				"start",
			},
			NotExpectedContent: []string{"error"},
			// ExpectedEvents:     map[string]int{"OpenIDTestEvent": 1},
			TestAppFactory:     setupTestApp,
		},
		// {
		// 	Name:   "Save Variables and Start - Missing Auth",
		// 	Method: "POST",
		// 	// using dynamic URL path parameters
		// 	URL:     "/api/http/johndoe/save-variables-and-start",
		// 	Body:    strings.NewReader(`{"var1":"value1"}`),
		// 	Headers: map[string]string{"Content-Type": "application/json"},
		// 	Delay:   0,
		// 	Timeout: 5 * time.Second,
		// 	// Expecting unauthorized error due to missing auth header.
		// 	ExpectedStatus: http.StatusUnauthorized,
		// 	ExpectedContent: []string{
		// 		"unauthorized",
		// 	},
		// 	NotExpectedContent: []string{"success"},
		// 	ExpectedEvents:     map[string]int{"AuthRequired": 1},
		// 	TestAppFactory:     setupTestApp,
		// },
		// {
		// 	Name:   "Save Variables and Start - With Valid Auth",
		// 	Method: "POST",
		// 	URL:    "/api/http/janedoe/save-variables-and-start",
		// 	Body:   strings.NewReader(`{"var1":"value1", "var2":"value2"}`),
		// 	Headers: map[string]string{
		// 		"Content-Type":  "application/json",
		// 		"Authorization": "Bearer validtoken",
		// 	},
		// 	Delay:          0,
		// 	Timeout:        5 * time.Second,
		// 	ExpectedStatus: http.StatusOK,
		// 	ExpectedContent: []string{
		// 		"started",
		// 		"variables saved",
		// 	},
		// 	NotExpectedContent: []string{"error", "unauthorized"},
		// 	ExpectedEvents:     map[string]int{"SaveVariablesEvent": 1},
		// 	TestAppFactory:     setupTestApp,
		// },
		// {
		// 	Name:           "Wallet Test - Confirm Success",
		// 	Method:         "POST",
		// 	URL:            "/wallet-test/confirm-success",
		// 	Body:           strings.NewReader(`{"walletId":"12345","status":"confirmed"}`),
		// 	Headers:        map[string]string{"Content-Type": "application/json"},
		// 	Delay:          1 * time.Second,
		// 	Timeout:        5 * time.Second,
		// 	ExpectedStatus: http.StatusOK,
		// 	ExpectedContent: []string{
		// 		"confirmation received",
		// 	},
		// 	NotExpectedContent: []string{"failure"},
		// 	ExpectedEvents:     map[string]int{"WalletSuccessEvent": 1},
		// 	TestAppFactory:     setupTestApp,
		// },
		// {
		// 	Name:           "Wallet Test - Notify Failure",
		// 	Method:         "POST",
		// 	URL:            "/wallet-test/notify-failure",
		// 	Body:           strings.NewReader(`{"walletId":"12345","error":"timeout occurred"}`),
		// 	Headers:        map[string]string{"Content-Type": "application/json"},
		// 	Delay:          500 * time.Millisecond,
		// 	Timeout:        5 * time.Second,
		// 	ExpectedStatus: http.StatusOK,
		// 	ExpectedContent: []string{
		// 		"failure notified",
		// 	},
		// 	NotExpectedContent: []string{"success"},
		// 	ExpectedEvents:     map[string]int{"WalletFailureEvent": 1},
		// 	TestAppFactory:     setupTestApp,
		// },
		// {
		// 	Name:           "Wallet Test - Send Log Update",
		// 	Method:         "POST",
		// 	URL:            "/wallet-test/send-log-update",
		// 	Body:           strings.NewReader(`{"walletId":"12345","log":"log message update"}`),
		// 	Headers:        map[string]string{"Content-Type": "application/json"},
		// 	Delay:          500 * time.Millisecond,
		// 	Timeout:        5 * time.Second,
		// 	ExpectedStatus: http.StatusOK,
		// 	ExpectedContent: []string{
		// 		"log updated",
		// 	},
		// 	NotExpectedContent: []string{"error"},
		// 	ExpectedEvents:     map[string]int{"LogUpdateEvent": 1},
		// 	TestAppFactory:     setupTestApp,
		// },
		// {
		// 	Name:           "Wallet Test - Send Log Update Start",
		// 	Method:         "POST",
		// 	URL:            "/wallet-test/send-log-update-start",
		// 	Body:           strings.NewReader(`{"walletId":"12345","initLog":"beginning log update"}`),
		// 	Headers:        map[string]string{"Content-Type": "application/json"},
		// 	Delay:          500 * time.Millisecond,
		// 	Timeout:        5 * time.Second,
		// 	ExpectedStatus: http.StatusOK,
		// 	ExpectedContent: []string{
		// 		"log update started",
		// 	},
		// 	NotExpectedContent: []string{"error", "failure"},
		// 	ExpectedEvents:     map[string]int{"LogUpdateStartEvent": 1},
		// 	TestAppFactory:     setupTestApp,
		// },
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}

}
