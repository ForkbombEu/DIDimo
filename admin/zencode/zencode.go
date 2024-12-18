package zencode

import (
	_ "embed"
	"encoding/json"
	"errors"
	"pb/config"

	slangroom "github.com/dyne/slangroom-exec/bindings/go"
)

//go:embed zenflows-crypto/src/keypairoomServer-6-7.zen
var KEYPAIROOM_ZENCODE string

//go:embed zenflows-crypto/src/pubkeys-request-signed-02-24.zen
var PUBKEYS_REQUEST_SIGNED_ZENCODE string

func KeypairoomServer(conf *config.KeypairoomConfig, data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"userData":       data,
		"serverSideSalt": conf.Salt,
	})
	if err != nil {
		return "", err
	}

	result, success := slangroom.Exec(slangroom.SlangroomInput{Contract: KEYPAIROOM_ZENCODE, Keys: string(jsonData)})
	if success != nil {
		return "", errors.New(result.Logs)
	}
	var zenroomResult struct {
		HMAC string `json:"seedServerSideShard.HMAC"`
	}

	err = json.Unmarshal([]byte(result.Output), &zenroomResult)
	if err != nil {
		return "", err
	}

	return zenroomResult.HMAC, nil
}

func PubkeysRequestSigned(didRequest map[string]interface{}) (map[string]interface{}, error) {
	var err error
	var zenroomResult map[string]interface{}

	jsonData, err := json.Marshal(didRequest)
	if err != nil {
		return zenroomResult, err
	}

	result, success := slangroom.Exec(slangroom.SlangroomInput{Contract: PUBKEYS_REQUEST_SIGNED_ZENCODE, Keys: string(jsonData)})
	if success != nil {
		return zenroomResult, errors.New(result.Logs)
	}

	err = json.Unmarshal([]byte(result.Output), &zenroomResult)
	if err != nil {
		return zenroomResult, err
	}

	return zenroomResult, nil
}
