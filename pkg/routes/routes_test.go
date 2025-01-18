package routes

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

const testDataDir = "./../../fixtures/test_pb_data"

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

func TestKeypairoomRoute(t *testing.T) {
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
			ExpectedContent: []string{`{"data":{},"message":"The request requires valid record authorization token.","status":401}`},
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
                    "email": "test@example.com"
                }
            }`),
			ExpectedStatus:  200,
			ExpectedContent: []string{`{"hmac":"cFV3s5YaSriFM9eF/yDqkwF0snxhqK45sffFN0SwQSo="}`},
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
                    "email": "test@example.com"
                }
            }`),
			ExpectedStatus:  200,
			ExpectedContent: []string{`{"hmac":"cFV3s5YaSriFM9eF/yDqkwF0snxhqK45sffFN0SwQSo="}`},
			TestAppFactory:  setupTestApp,
		},
	}

	for _, scenario := range KeypairoomServerScenarios {
		scenario.Test(t)
	}
}

func TestDidRoute(t *testing.T) {
	recordTokenA, err := generateToken("users", "userA@example.org")
	if err != nil {
		t.Fatal(err)
	}

	recordTokenB, err := generateToken("users", "userB@example.org")
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

	DIDServerScenarios := []tests.ApiScenario{
		{
			Name:            "try as guest (aka. no Authorization header)",
			Method:          http.MethodGet,
			URL:             "/api/did",
			ExpectedStatus:  401,
			ExpectedContent: []string{"\"data\":{}"},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:   "try as authenticated app user with key",
			Method: http.MethodGet,
			URL:    "/api/did",
			Headers: map[string]string{
				"Authorization": recordTokenA,
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{`{"created":false,"did":{"@context":"https://w3id.org/did-resolution/v1","didDocument":{"@context":["https://www.w3.org/ns/did/v1","https://w3id.org/security/suites/ed25519-2018/v1","https://w3id.org/security/suites/secp256k1-2019/v1","https://w3id.org/security/suites/secp256k1-2020/v1","https://dyne.github.io/W3C-DID/specs/ReflowBLS12381.json","https://dyne.github.io/W3C-DID/specs/EcdsaSecp256r1.json",{"description":"https://schema.org/description","identifier":"https://schema.org/identifier"}],"description":"DIDimo user","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","identifier":"43","proof":{"created":"1737044873293","jws":"eyJhbGciOiJFUzI1NksiLCJiNjQiOnRydWUsImNyaXQiOiJiNjQifQ..wmsQ-rjrWvZv_R1o5VZUZ5uKiS7McnECNDHU3uDY4pC0g-Qe0sPLDrGhN8T82WPQtUDdr4-vJ_b4GChI-5FsPw","proofPurpose":"assertionMethod","type":"EcdsaSecp256k1Signature2019","verificationMethod":"did:dyne:sandbox.didimo_A:GTkLUcMs4zdntg9U4hb2FJMCvuQMa1swqqLxKxcVXDKa#ecdh_public_key"},"verificationMethod":[{"controller":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV#ecdh_public_key","publicKeyBase58":"RC5aweQQeXAgbtr8HC2PHQsJ1Hi3p1VNCqUCwr2RJCwHntT3VRBZ4jhK5hAaZJxx6Ed8BJoCbnWSeVFLy5nrWMqy","type":"EcdsaSecp256k1VerificationKey2019"},{"controller":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV#reflow_public_key","publicKeyBase58":"ABKyacZJXACEUCsXNjA2UbeVLqpfeuymA4Hn63B6skJb6sEsPrtfCLRht6o9dNiYF2XKFyc4Q3uSUM42xvxnnDB2q1kYLWknKgF4gFbxw5YLFghArLcVU9WVBoYRujSmk4e5SPvqQMDLbAZzk3fz8jt9WA5Kp7cjtkXLpvAAVbMNyHKafsFqfk24KGV6GXwnJsK3tejmvGiKqCHbgNaoDWup125VKtVVSGjBu4y9e3N9xnC9rN7dV7CX2ukTyZ1DWtMXzL","type":"ReflowBLS12381VerificationKey"},{"controller":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV#bitcoin_public_key","publicKeyBase58":"fGFjnyf9NE8wcKSBVSi4v7g6aXMfJA8R7HV4WF54uo9y","type":"EcdsaSecp256k1VerificationKey2019"},{"controller":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV#eddsa_public_key","publicKeyBase58":"6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","type":"Ed25519VerificationKey2018"},{"controller":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV#es256_public_key","publicKeyBase58":"454GSTNiKekt68pDFe83xE5zBv9p2tJiz5pWA6tbL1ifj8ZqMdwXL9gqvDJ6hAgGhE1WQsdeZSwiEKdgUTxibHJC","type":"EcdsaSecp256r1VerificationKey"},{"blockchainAccountId":"eip155:1:0xb172b9a3faa978a00dfe1648631097d194e9b7ee","controller":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV","id":"did:dyne:sandbox.didimo:6XmGn1Fv9BEWUB9dr5v3GpYhUj4bJd9zZVbpM2f3q7VV#ethereum_address","type":"EcdsaSecp256k1RecoveryMethod2020"}]},"didDocumentMetadata":{"created":"1737044873293","deactivated":"false"}}}`},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:   "try as user without key",
			Method: http.MethodGet,
			URL:    "/api/did",
			Headers: map[string]string{
				"Authorization": recordTokenB,
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`{"data":{},"message":"Only users with public keys can access this endpoint.","status":403}`},
			TestAppFactory:  setupTestApp,
		},
	}
	for _, scenario := range DIDServerScenarios {
		scenario.Test(t)
	}
}
