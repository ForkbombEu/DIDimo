// main entry point for the didimo backend server

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	_ "github.com/forkbombeu/didimo/migrations"
	"github.com/forkbombeu/didimo/pocketbase/did"
	"github.com/forkbombeu/didimo/pocketbase/feature"
	"github.com/forkbombeu/didimo/pocketbase/hooks"
	"github.com/forkbombeu/didimo/pocketbase/webauthn"
	"github.com/forkbombeu/didimo/pocketbase/zencode"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()
	app.RootCmd.Short = "\033[38;2;255;100;0m      dP oo       dP oo                     \033[0m\n" +
		"\033[38;2;255;71;43m      88          88                        \033[0m\n" +
		"\033[38;2;255;43;86m.d888b88 dP .d888b88 dP 88d8b.d8b. .d8888b. \033[0m\n" +
		"\033[38;2;255;14;129m88'  `88 88 88'  `88 88 88'`88'`88 88'  `88 \033[0m\n" +
		"\033[38;2;236;0;157m88.  .88 88 88.  .88 88 88  88  88 88.  .88 \033[0m\n" +
		"\033[38;2;197;0;171m`88888P8 dP `88888P8 dP dP  dP  dP `88888P' \033[0m\n" +
		"\033[38;2;159;0;186m                                             \033[0m\n" +
		"                   \033[48;2;0;0;139m\033[38;2;255;255;255m           :(){ :|:& };: \033[0m\n" + // Forkbomb with padding
		"                   \033[48;2;0;0;139m\033[38;2;255;255;255m by The Forkbomb Company \033[0m\n" // Company name aligned to right

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "localhost:5100"})
		e.Router.Any("/*", echo.WrapHandler(proxy))
		e.Router.Any("/", echo.WrapHandler(proxy))

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/keypairoom-server",
			Handler: func(c echo.Context) error {
				var body map[string]map[string]interface{}

				conf, err := feature.FetchKeypairoomConfig(app)
				if err != nil {
					return err
				}

				err = json.NewDecoder(c.Request().Body).Decode(&body)
				if err != nil {
					return err
				}
				hmac, err := zencode.KeypairoomServer(conf, body["userData"])
				if err != nil {
					return err
				}

				return c.JSON(http.StatusOK, map[string]string{"hmac": hmac})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/did",
			Handler: func(c echo.Context) error {
				authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
				if authRecord == nil {
					return apis.NewForbiddenError("Only auth records can access this endpoint", nil)
				}

				publicKeys, err := app.Dao().FindFirstRecordByFilter("users_public_keys", "owner = {:owner_id}", dbx.Params{"owner_id": authRecord.Id})
				if err != nil {
					return apis.NewForbiddenError("Only users with public keys can access this endpoint", nil)
				}

				conf, err := feature.FetchDidConfig(app)
				if err != nil {
					return err
				}

				did, err := did.ClaimDid(conf, &did.DidAgent{
					BitcoinPublicKey: publicKeys.Get("bitcoin_public_key").(string),
					EcdhPublicKey:    publicKeys.Get("ecdh_public_key").(string),
					EddsaPublicKey:   publicKeys.Get("eddsa_public_key").(string),
					EthereumAddress:  publicKeys.Get("ethereum_address").(string),
					ReflowPublicKey:  publicKeys.Get("reflow_public_key").(string),
					Es256PublicKey:   publicKeys.Get("es256_public_key").(string),
				})
				if err != nil {
					return err
				}

				return c.JSON(http.StatusOK, did)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	webauthn.Register(app)
	hooks.Register(app)
	jsvm.MustRegister(app, jsvm.Config{
		HooksWatch: true,
	})
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangJS,
		Automigrate:  true,
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
