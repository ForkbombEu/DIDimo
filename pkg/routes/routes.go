package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/forkbombeu/didimo/pkg/internal/pb"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func Setup(app *pocketbase.PocketBase) {
	routes := map[string]string{
		"/workflows/{path...}":  getEnv("ADDRESS_TEMPORAL", "http://localhost:8080"),
		"/monitoring/{path...}": getEnv("ADDRESS_GRAFANA", "http://localhost:8085"),
		"/{path...}":            getEnv("ADDRESS_UI", "http://localhost:5100"),
	}

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		for path, target := range routes {
			se.Router.Any(path, createReverseProxy(target))
		}

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
		}).Bind(apis.RequireAuth())

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
		}).Bind(apis.RequireAuth())

		return se.Next()
	})

	pb.HookNamespaceOrgs(app)
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

func createReverseProxy(target string) func(r *core.RequestEvent) error {
	return func(r *core.RequestEvent) error {
		targetURL, err := url.Parse(target)
		if err != nil {
			return err
		}
		log.Printf("Proxying request: %s -> %s%s", r.Request.URL.Path, targetURL.String(), r.Request.URL.Path)

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.Header.Set("Host", targetURL.Host)
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
			if origin := req.Header.Get("Origin"); origin != "" {
				req.Header.Set("Origin", origin)
			}
			if referer := req.Header.Get("Referer"); referer != "" {
				req.Header.Set("Referer", referer)
			}
		}
		proxy.ServeHTTP(r.Response, r.Request)
		return nil
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
