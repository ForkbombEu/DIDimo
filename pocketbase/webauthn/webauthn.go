package webauthn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type User struct {
	app   *pocketbase.PocketBase
	model *core.Record
}

func NewUser(app *pocketbase.PocketBase, m *core.Record) *User {
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
	credentialsStore, err := u.app.FindCollectionByNameOrId("webauthnCredentials")
	if err != nil {
		return err
	}
	record := core.NewRecord(credentialsStore)
	record.Set("user", u.model.Id)
	record.Set("credential", cred)
	record.Set("description", description)

	if err := u.app.Save(record); err != nil {
		return err
	}
	return nil
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	var credentials []webauthn.Credential
	records, err := u.app.FindAllRecords("webauthnCredentials",
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
	record, err := app.FindFirstRecordByData("features", "name", "webauthn")
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

func storeSessionData(app *pocketbase.PocketBase, userRecord *core.Record, sessionData *webauthn.SessionData) error {
	// Remove old session data
	record, err := app.FindFirstRecordByData("sessionDataWebauthn", "user", userRecord.Id)
	if record != nil {
		if err := app.Delete(record); err != nil {
			return err
		}
	}

	// store session data as marshaled JSON
	sessionStore, err := app.FindCollectionByNameOrId("sessionDataWebauthn")
	if err != nil {
		return err
	}
	record = core.NewRecord(sessionStore)
	record.Set("user", userRecord.Id)
	record.Set("session", sessionData)

	if err := app.Save(record); err != nil {
		return err
	}
	return nil
}

func Register(app *pocketbase.PocketBase) error {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/webauthn/register/begin/:email", func(e *core.RequestEvent) error {
			w, err := NewWebAuthnFromEnv(app)
			if err != nil {
				return err
			}

			email := e.Request.PathValue("email")

			authenticated := true
			if e.Auth == nil {
				authenticated = false
				// User not authenticated I have to create a new user,
				// but if a user exists I may go on if that user doesn't
				// have neither a password nor a credential
				e.Auth, err = app.FindAuthRecordByEmail("users", email)
				if err != nil {
					// Could not fetch the user, try to create a new one
					collection, err := app.FindCollectionByNameOrId("users")
					if err != nil {
						return err
					}

					e.Auth = core.NewRecord(collection)
					e.Auth.Set("email", email)
					e.Auth.Set("username", email)
					e.Auth.RefreshTokenKey()
					if err := app.Save(e.Auth); err != nil {
						return err
					}
				}
			} else if e.Auth.Get("email") != email { // User is logged in
				return apis.NewForbiddenError("Wrong email", nil)
			}

			user := NewUser(app, e.Auth)

			if !authenticated && (len(user.WebAuthnCredentials()) > 0 || e.Auth.GetString("passwordHash") != "") {
				return apis.NewForbiddenError("A user already exists with this email", nil)
			}

			options, sessionData, err := w.BeginRegistration(
				user,
			)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, err.Error())
			}

			err = storeSessionData(app, e.Auth, sessionData)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, err.Error())
			}

			return e.JSON(http.StatusOK, options)
		})
		se.Router.POST("/api/webauthn/register/finish/:email", func(e *core.RequestEvent) error {
			w, err := NewWebAuthnFromEnv(app)
			if err != nil {
				return err
			}

			data := new(struct {
				Description string `json:"description" form:"description" query:"description"`
			})
			var b []byte

			// I have to read c.Request() twice.. :(
			if e.Request.Body != nil {
				// TODO: check that the body is not tooo big (it should not)
				b, _ = io.ReadAll(e.Request.Body)
				e.Request.Body = io.NopCloser(bytes.NewBuffer(b))

				if err := json.Unmarshal(b, data); err != nil {
					return e.String(http.StatusBadRequest, err.Error())
				}
				e.Request.Body = io.NopCloser(bytes.NewBuffer(b))
			}
			fmt.Printf("data: %v\n", data)
			email := e.Request.PathValue("email")

			userRecord, err := app.FindAuthRecordByEmail("users", email)
			if err != nil {
				log.Println(err)
				return err
			}
			user := NewUser(app, userRecord)
			record, err := app.FindFirstRecordByData("sessionDataWebauthn", "user", userRecord.Id)
			if err != nil {
				return err
			}
			var sessionData webauthn.SessionData
			json.Unmarshal([]byte(record.GetString("session")), &sessionData)

			credential, err := w.FinishRegistration(user, sessionData, e.Request)
			if err != nil {
				fmt.Println(e.Request)
				return err
			}
			user.AddCredential(*credential, data.Description)

			if err := app.Save(userRecord); err != nil {
				return err
			}
			if err := app.Delete(record); err != nil {
				return err
			}
			return e.JSON(http.StatusOK, make(map[string]interface{}))
		})
		se.Router.GET("/api/webauthn/login/begin/:email", func(e *core.RequestEvent) error {
			w, err := NewWebAuthnFromEnv(app)
			if err != nil {
				return err
			}
			if w == nil {
				return apis.NewNotFoundError("Webauthn not enabled", nil)
			}

			email := e.Request.PathValue("email")
			userRecord, err := app.FindAuthRecordByEmail("users", email)
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
				return e.JSON(http.StatusInternalServerError, err.Error())
			}

			return e.JSON(http.StatusOK, options)
		})
		se.Router.POST("/api/webauthn/login/finish/:email", func(e *core.RequestEvent) error {
			w, err := NewWebAuthnFromEnv(app)
			if err != nil {
				return err
			}
			if w == nil {
				return apis.NewNotFoundError("Webauthn not enabled", nil)
			}

			email := e.Request.PathValue("email")
			userRecord, err := app.FindAuthRecordByEmail("users", email)
			if err != nil {
				return err
			}
			user := NewUser(app, userRecord)
			record, err := app.FindFirstRecordByData("sessionDataWebauthn", "user", userRecord.Id)
			if err != nil {
				return err
			}
			var sessionData webauthn.SessionData
			json.Unmarshal([]byte(record.GetString("session")), &sessionData)

			_, err = w.FinishLogin(user, sessionData, e.Request)
			if err != nil {
				return err
			}

			// generate an auth token and return an auth response
			// note: in the future the below will be simplified to just: return api.AuthResponse(c, user)
			token, tokenErr := userRecord.NewAuthToken()
			if tokenErr != nil {
				return errors.New("Failed to create user token")
			}

			return e.JSON(http.StatusOK, map[string]any{
				"token": token,
				"user":  userRecord,
			})
		})
		return se.Next()
	})
	return nil
}
