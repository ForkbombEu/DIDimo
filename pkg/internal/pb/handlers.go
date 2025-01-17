// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2022-2023 Dyne.org foundation <foundation@dyne.org>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package pb

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func ProxyHandler(req *core.RequestEvent) error {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:5100",
	})
	proxy.ServeHTTP(req.Response, req.Request)
	return nil
}

// KeypairoomServerHandler handles the `/api/keypairoom-server` route.
func KeypairoomServerHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var body map[string]map[string]interface{}

		conf, err := FetchKeypairoomConfig(app)
		if err != nil {
			return err
		}

		err = json.NewDecoder(e.Request.Body).Decode(&body)
		if err != nil {
			return err
		}
		hmac, err := KeypairoomServer(conf, body["userData"])
		if err != nil {
			return err
		}

		return e.JSON(http.StatusOK, map[string]string{"hmac": hmac})
	}
}

// DidHandler handles the `/api/did` route.
func DidHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {

		publicKeys, err := app.FindFirstRecordByFilter("users_public_keys", "owner = {:owner_id}", dbx.Params{"owner_id": e.Auth.Id})
		if err != nil {
			return apis.NewForbiddenError("Only users with public keys can access this endpoint", nil)
		}

		conf, err := FetchDidConfig(app)
		if err != nil {
			return err
		}

		did, err := ClaimDid(conf, &DidAgent{
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
	}
}
