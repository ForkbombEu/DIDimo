package webauthn

import (
	"encoding/json"
	"errors"

	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

type User struct {
	app   *pocketbase.PocketBase
	model *models.Record
}

func NewUser(app *pocketbase.PocketBase, m *models.Record) *User {

	user := &User{
		model: m,
		app:   app,
	}

	return user
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	return []byte(u.model.Id)
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.model.GetString("email")
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.model.GetString("email")
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential, description string) error {
	credentialsStore, err := u.app.Dao().FindCollectionByNameOrId("webauthnCredentials")
	if err != nil {
		return err
	}
	record := models.NewRecord(credentialsStore)
	record.Set("user", u.model.Id)
	record.Set("credential", cred)
	record.Set("description", description)

	if err := u.app.Dao().SaveRecord(record); err != nil {
		return err
	}
	return nil
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	var credentials []webauthn.Credential
	records, err := u.app.Dao().FindRecordsByExpr("webauthnCredentials",
		dbx.NewExp("user = {:user}", dbx.Params{"user": u.model.Id}))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	for i := 0; i < len(records); i++ {
		var credential webauthn.Credential
		err := json.Unmarshal([]byte(records[i].GetString("credential")), &credential)
		if err != nil {
			log.Printf("Should not be here, wrong data in JSON field: %s", err.Error())
			return nil
		} else {
			credentials = append(credentials, credential)
		}
	}
	return credentials
}

func NewWebAuthnFromEnv(app *pocketbase.PocketBase) (*webauthn.WebAuthn, error) {
	record, err := app.Dao().FindFirstRecordByData("features", "name", "webauthn")
	if err != nil {
		return nil, err
	}

	if !record.GetBool("active") {
		return nil, apis.NewNotFoundError("Webauthn not enabled", nil)
	}

	var envConfig struct {
		DisplayName string `json:"DISPLAY_NAME"`
		RPId        string `json:"RPID"`
		RPOrigin    string `json:"RPORIGINS"`
	}

	err = json.Unmarshal([]byte(record.GetString("envVariables")), &envConfig)
	if err != nil {
		return nil, err
	}

	if envConfig.DisplayName == "" {
		return nil, errors.New("Display name is empty")
	}

	if envConfig.RPId == "" {
		return nil, errors.New("Relying party not set")
	}

	if envConfig.RPOrigin == "" {
		return nil, errors.New("Relying party origin not set")
	}

	wconfig := &webauthn.Config{
		RPDisplayName: envConfig.DisplayName, // Display Name for your site
		RPID:          envConfig.RPId,        // Generally the FQDN for your site
		RPOrigins:     []string{envConfig.RPOrigin},
	}

	w, err := webauthn.New(wconfig)

	return w, nil
}

func storeSessionData(app *pocketbase.PocketBase, userRecord *models.Record, sessionData *webauthn.SessionData) error {
	// Remove old session data
	record, err := app.Dao().FindFirstRecordByData("sessionDataWebauthn", "user", userRecord.Id)
	if record != nil {
		if err := app.Dao().DeleteRecord(record); err != nil {
			return err
		}
	}

	// store session data as marshaled JSON
	sessionStore, err := app.Dao().FindCollectionByNameOrId("sessionDataWebauthn")
	if err != nil {
		return err
	}
	record = models.NewRecord(sessionStore)
	record.Set("user", userRecord.Id)
	record.Set("session", sessionData)

	if err := app.Dao().SaveRecord(record); err != nil {
		return err
	}
	return nil
}

func Register(app *pocketbase.PocketBase) error {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/webauthn/register/begin/:email",
			Handler: func(c echo.Context) error {
				w, err := NewWebAuthnFromEnv(app)
				if err != nil {
					return err
				}

				email := c.PathParam("email")

				authenticated := true
				userRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
				if userRecord == nil {
					authenticated = false
					// User not authenticated I have to create a new user,
					// but if a user exists I may go on if that user doesn't
					// have neither a password nor a credential
					userRecord, err = app.Dao().FindAuthRecordByEmail("users", email)
					if err != nil {
						// Could not fetch the user, try to create a new one
						collection, err := app.Dao().FindCollectionByNameOrId("users")
						if err != nil {
							return err
						}

						userRecord = models.NewRecord(collection)
						userRecord.Set("email", email)
						userRecord.Set("username", email)
						userRecord.RefreshTokenKey()
						if err := app.Dao().SaveRecord(userRecord); err != nil {
							return err
						}
					}
				} else if userRecord.Get("email") != email { // User is logged in
					return apis.NewForbiddenError("Wrong email", nil)
				}

				user := NewUser(app, userRecord)

				if !authenticated && (len(user.WebAuthnCredentials()) > 0 || userRecord.PasswordHash() != "") {
					return apis.NewForbiddenError("A user already exists with this email", nil)
				}

				options, sessionData, err := w.BeginRegistration(
					user,
				)

				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}

				err = storeSessionData(app, userRecord, sessionData)

				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}

				return c.JSON(http.StatusOK, options)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/webauthn/register/finish/:email",
			Handler: func(c echo.Context) error {
				w, err := NewWebAuthnFromEnv(app)
				if err != nil {
					return err
				}

				data := new(struct {
					Description string `json:"description" form:"description" query:"description"`
				})
				var b []byte

				// I have to read c.Request() twice.. :(
				if c.Request().Body != nil {
					// TODO: check that the body is not tooo big (it should not)
					b, _ = io.ReadAll(c.Request().Body)
					c.Request().Body = io.NopCloser(bytes.NewBuffer(b))

					if err := c.Bind(data); err != nil {
						return c.String(http.StatusBadRequest, err.Error())
					}
					c.Request().Body = io.NopCloser(bytes.NewBuffer(b))
				}
				fmt.Printf("data: %v\n", data)
				email := c.PathParam("email")

				userRecord, err := app.Dao().FindAuthRecordByEmail("users", email)
				if err != nil {
					log.Println(err)
					return err
				}
				user := NewUser(app, userRecord)
				record, err := app.Dao().FindFirstRecordByData("sessionDataWebauthn", "user", userRecord.Id)
				if err != nil {
					return err
				}
				var sessionData webauthn.SessionData
				json.Unmarshal([]byte(record.GetString("session")), &sessionData)

				credential, err := w.FinishRegistration(user, sessionData, c.Request())
				if err != nil {
					fmt.Println(c.Request())
					return err
				}
				user.AddCredential(*credential, data.Description)

				if err := app.Dao().SaveRecord(userRecord); err != nil {
					return err
				}
				if err := app.Dao().DeleteRecord(record); err != nil {
					return err
				}

				return c.JSON(http.StatusOK, make(map[string]interface{}))
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/webauthn/login/begin/:email",
			Handler: func(c echo.Context) error {
				w, err := NewWebAuthnFromEnv(app)
				if err != nil {
					return err
				}
				if w == nil {
					return apis.NewNotFoundError("Webauthn not enabled", nil)
				}

				email := c.PathParam("email")
				userRecord, err := app.Dao().FindAuthRecordByEmail("users", email)
				if err != nil {
					return err
				}
				user := NewUser(app, userRecord)

				// generate PublicKeyCredentialRequestOptions, session data
				options, sessionData, err := w.BeginLogin(user)
				if err != nil {
					return err
				}

				err = storeSessionData(app, userRecord, sessionData)

				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}

				return c.JSON(http.StatusOK, options)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/webauthn/login/finish/:email",
			Handler: func(c echo.Context) error {
				w, err := NewWebAuthnFromEnv(app)
				if err != nil {
					return err
				}
				if w == nil {
					return apis.NewNotFoundError("Webauthn not enabled", nil)
				}

				email := c.PathParam("email")
				userRecord, err := app.Dao().FindAuthRecordByEmail("users", email)
				if err != nil {
					return err
				}
				user := NewUser(app, userRecord)
				record, err := app.Dao().FindFirstRecordByData("sessionDataWebauthn", "user", userRecord.Id)
				if err != nil {
					return err
				}
				var sessionData webauthn.SessionData
				json.Unmarshal([]byte(record.GetString("session")), &sessionData)

				_, err = w.FinishLogin(user, sessionData, c.Request())
				if err != nil {
					return err
				}

				// generate an auth token and return an auth response
				// note: in the future the below will be simplified to just: return api.AuthResponse(c, user)
				token, tokenErr := tokens.NewRecordAuthToken(app, userRecord)
				if tokenErr != nil {
					return errors.New("Failed to create user token")
				}

				return c.JSON(http.StatusOK, map[string]any{
					"token": token,
					"user":  userRecord,
				})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		return nil
	})
	return nil
}
