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

package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"net/url"
)

type KeypairoomConfig struct {
	Salt string
}

type DidConfig struct {
	Keyring    map[string]interface{}
	DidURL     *url.URL
	Spec       string
	SignerSpec string
	Identity   string
}

func FetchKeypairoomConfig(app core.App) (*KeypairoomConfig, error) {
	c := &KeypairoomConfig{}
	conf, err := newConfig(app, "keypairoom")
	if err != nil {
		return c, err
	}

	if conf["SALT"] == "" {
		return c, fmt.Errorf("SALT must be provided")
	}
	c.Salt = conf["SALT"]

	return c, nil
}

func FetchDidConfig(app core.App) (*DidConfig, error) {
	c := &DidConfig{}
	conf, err := newConfig(app, "DID")
	if err != nil {
		return c, err
	}

	if conf["DID_SPEC"] == "" {
		return c, fmt.Errorf("DID_SPEC must be provided")
	}
	c.Spec = conf["DID_SPEC"]

	if conf["DID_SIGNER_SPEC"] == "" {
		return c, fmt.Errorf("DID_SIGNER_SPEC must be provided")
	}
	c.SignerSpec = conf["DID_SIGNER_SPEC"]

	if conf["DID_IDENTITY"] == "" {
		return c, fmt.Errorf("DID_IDENTITY must be provided")
	}
	c.Identity = conf["DID_IDENTITY"]

	u, err := fetchURL(conf["DID_URL"])
	if err != nil {
		return c, err
	}
	c.DidURL = u

	d, err := fetchDict(conf["DID_KEYRING"])
	if err != nil {
		return c, err
	}
	c.Keyring = d

	return c, nil
}

func fetchURL(env string) (*url.URL, error) {
	u, err := url.Parse(env)
	if err != nil {
		return nil, fmt.Errorf("%q is malformed: %w", env, err)
	}

	// not a uri, possibly a url
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("%q is malformed: not a url", env)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("%q is malformed: invalid scheme; must be http(s)", env)
	}

	// normalize it: take only what we need
	u = &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   "/",
	}

	return u, nil
}

// Parse a JSON dictionary encoded as base64
func fetchDict(env string) (map[string]interface{}, error) {
	rawDecodedDict, err := base64.StdEncoding.DecodeString(env)
	if err != nil {
		return nil, fmt.Errorf("%q is malformed: %w", err)
	}

	var body map[string]interface{}
	json.Unmarshal(rawDecodedDict, &body)

	return body, nil
}

func newConfig(app core.App, feature string) (map[string]string, error) {
	var envConfig map[string]string
	record, err := app.Dao().
		FindFirstRecordByData("features", "name", feature)
	envString, err := record.Get("envVariables").(types.JsonRaw).MarshalJSON()
	if err != nil {
		return envConfig, err
	}
	err = json.Unmarshal(envString, &envConfig)
	if err != nil {
		return envConfig, err
	}
	return envConfig, nil
}
