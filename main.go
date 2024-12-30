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

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   "localhost:5100",
		})
		se.Router.Any("/*", func(req *core.RequestEvent) error {
			proxy.ServeHTTP(req.Response, req.Request)
			return nil
		})
		se.Router.Any("/", func(req *core.RequestEvent) error {
			proxy.ServeHTTP(req.Response, req.Request)
			return nil
		})

		se.Router.POST("/api/keypairoom-server", func(e *core.RequestEvent) error {
			var body map[string]map[string]interface{}

			conf, err := feature.FetchKeypairoomConfig(app)
			if err != nil {
				return err
			}

			err = json.NewDecoder(e.Request.Body).Decode(&body)
			if err != nil {
				return err
			}
			hmac, err := zencode.KeypairoomServer(conf, body["userData"])
			if err != nil {
				return err
			}

			return e.JSON(http.StatusOK, map[string]string{"hmac": hmac})
		})

		se.Router.GET("/api/did", func(e *core.RequestEvent) error {
			authRecord := e.Auth
			if authRecord == nil {
				return apis.NewForbiddenError("Only auth records can access this endpoint", nil)
			}

			publicKeys, err := app.FindFirstRecordByFilter("users_public_keys", "owner = {:owner_id}", dbx.Params{"owner_id": authRecord.Id})
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

			return e.JSON(http.StatusOK, did)
		})

		return se.Next()
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
