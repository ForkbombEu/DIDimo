package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/forkbombeu/didimo/pkg/internal/pb"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func Setup(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.Any("/*", proxyHandler)
		se.Router.Any("/", proxyHandler)
		se.Router.POST("/api/keypairoom-server", func(e *core.RequestEvent) error {
			var body map[string]map[string]interface{}

			conf, err := pb.FetchKeypairoomConfig(app)
			if err != nil {
				return err
			}

			err = json.NewDecoder(e.Request.Body).Decode(&body)
			if err != nil {
				return err
			}
			hmac, err := pb.KeypairoomServer(conf, body["userData"])
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

			conf, err := pb.FetchDidConfig(app)
			if err != nil {
				return err
			}

			did, err := pb.ClaimDid(conf, &pb.DidAgent{
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

	pb.DatabaseHooks(app)
	pb.Register(app)

	// ** BAD LINE **
	// hooks.Register(app)

	jsvm.MustRegister(app, jsvm.Config{
		HooksWatch: true,
	})
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangJS,
		Automigrate:  true,
	})
}

func proxyHandler(req *core.RequestEvent) error {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:5100",
	})
	proxy.ServeHTTP(req.Response, req.Request)
	return nil
}
