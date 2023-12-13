package did

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"pb/config"
	"pb/zencode"
	"strconv"
	"time"
)

const clientTimeout = 10 * time.Second

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}
var transport = &http.Transport{
	DisableKeepAlives:     true,
	Proxy:                 http.ProxyFromEnvironment,
	DialContext:           dialer.DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

type DidAgent struct {
	BitcoinPublicKey string
	EcdhPublicKey    string
	EddsaPublicKey   string
	EthereumAddress  string
	ReflowPublicKey  string
}

type DidResult struct {
	Created bool                   `json:"created"`
	Did     map[string]interface{} `json:"did"`
}

// from https://pkg.go.dev/net/http#pkg-overview
// Clients and Transports are safe for concurrent use by multiple goroutines
// and for efficiency should only be created once and re-used.
// TODO: Look at https://mauricio.github.io/golang-proxies
var client = &http.Client{
	Timeout: clientTimeout,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func restroomRequest(conf *config.DidConfig, contract string, body map[string]interface{}) (io.ReadCloser, error) {
	url := *conf.DidURL
	url.Path = fmt.Sprintf("/api/v1/sandbox/%s", contract)
	datakeys, _ := json.Marshal(map[string]map[string]interface{}{
		"keys": {},
		"data": body,
	})
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(datakeys))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func didId(conf *config.DidConfig, agent *DidAgent) string {
	return fmt.Sprintf("did:dyne:%s:%s", conf.Spec, agent.EddsaPublicKey)
}

func GetDid(conf *config.DidConfig, agent *DidAgent) (*DidResult, error) {
	url := *conf.DidURL
	url.Path = fmt.Sprintf("/dids/%s", didId(conf, agent))
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status code %d", url.String(), res.StatusCode)
	}

	bytesBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	err = json.Unmarshal(bytesBody, &body)

	result := DidResult{
		Created: false,
		Did:     body,
	}

	return &result, nil
}

func RequestNewDid(conf *config.DidConfig, agent *DidAgent) (*DidResult, error) {
	didRequest := map[string]interface{}{
		"proof": map[string]interface{}{
			"type":         "EcdsaSecp256k1Signature2019",
			"proofPurpose": "assertionMethod",
		},
		"@context": []interface{}{
			"https://www.w3.org/ns/did/v1",
			"https://w3id.org/security/suites/ed25519-2018/v1",
			"https://w3id.org/security/suites/secp256k1-2019/v1",
			"https://w3id.org/security/suites/secp256k1-2020/v1",
			"https://dyne.github.io/W3C-DID/specs/ReflowBLS12381.json",
			map[string]interface{}{
				"description": "https://schema.org/description",
				"identifier":  "https://schema.org/identifier",
			},
		},
		"did_spec":           conf.Spec,
		"signer_did_spec":    conf.SignerSpec,
		"identity":           conf.Identity,
		"bitcoin_public_key": agent.BitcoinPublicKey,
		"ecdh_public_key":    agent.EcdhPublicKey,
		"eddsa_public_key":   agent.EddsaPublicKey,
		"ethereum_address":   agent.EthereumAddress,
		"reflow_public_key":  agent.ReflowPublicKey,
		"timestamp":          strconv.FormatInt(time.Now().UnixMilli(), 10),
		"ifacer_id":          map[string]interface{}{"identifier": "43"},
	}
	for k, v := range conf.Keyring {
		didRequest[k] = v
	}

	did, err := zencode.PubkeysRequestSigned(didRequest)
	if err != nil {
		return nil, err
	}
	didSigned, err := restroomRequest(conf, "pubkeys-accept.chain", did)
	if err != nil {
		return nil, err
	}
	bytesBody, err := io.ReadAll(didSigned)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	err = json.Unmarshal(bytesBody, &body)
	if err != nil {
		return nil, err
	}

	result := DidResult{
		Created: true,
		Did:     body,
	}
	return &result, nil
}

func ClaimDid(conf *config.DidConfig, agent *DidAgent) (*DidResult, error) {
	did, err := GetDid(conf, agent)
	if err == nil {
		return did, nil
	}
	did, err = RequestNewDid(conf, agent)
	if err == nil {
		return did, nil
	}
	return nil, err
}
