package routes

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

const testDataDir = "./../../pb_data"

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

func TestRoutes(t *testing.T) {
	recordToken, err := generateToken("users", "userA@example.org")
	if err != nil {
		t.Fatal(err)
	}

	superuserToken, err := generateToken(core.CollectionNameSuperusers, "admin@example.org")
	if err != nil {
		t.Fatal(err)
	}

	// setup the test ApiScenario app instance
	setupTestApp := func(t testing.TB) *tests.TestApp {
		testApp, err := tests.NewTestApp(testDataDir)
		if err != nil {
			t.Fatal(err)
		}
		bindAppHooks(testApp)

		return testApp
	}

	KeypairoomServerScenarios := []tests.ApiScenario{
		{
			Name:            "try as guest (aka. no Authorization header)",
			Method:          http.MethodPost,
			URL:             "/api/keypairoom-server",
			ExpectedStatus:  401,
			ExpectedContent: []string{"\"data\":{}"},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "try with wrong method",
			Method:          http.MethodGet,
			URL:             "/api/keypairoom-server",
			ExpectedStatus:  401,
			ExpectedContent: []string{"\"data\":{}"},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:   "try as authenticated app user",
			Method: http.MethodPost,
			URL:    "/api/keypairoom-server",
			Headers: map[string]string{
				"Authorization": recordToken,
			},
			Body: strings.NewReader(`{
				"userData": {
					"email": "test@te.com"
				}
			}`),
			ExpectedStatus:  200,
			ExpectedContent: []string{`{"hmac":"gEgInamN7yKK7BJ33Uqg6Jw/eUv43xw99KNsMen6tug="}`},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:   "try as authenticated admin",
			Method: http.MethodPost,
			URL:    "/api/keypairoom-server",
			Headers: map[string]string{
				"Authorization": superuserToken,
			},
			Body: strings.NewReader(`{
				"userData": {
					"email": "test@te.com"
				}
			}`),
			ExpectedStatus:  200,
			ExpectedContent: []string{`{"hmac":"gEgInamN7yKK7BJ33Uqg6Jw/eUv43xw99KNsMen6tug="}`},
			TestAppFactory:  setupTestApp,
		},
	}

	for _, scenario := range KeypairoomServerScenarios {
		scenario.Test(t)
	}
}
