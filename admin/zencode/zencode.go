package zencode

import (
	_ "embed"
	"encoding/json"
	"errors"
	zenroom "github.com/dyne/Zenroom/bindings/golang/zenroom"
	"pb/config"
)

//go:embed zenflows-crypto/src/keypairoomServer-6-7.zen
var KEYPAIROOM_ZENCODE string

//go:embed zenflows-crypto/src/pubkeys-request-signed.zen
var PUBKEYS_REQUEST_SIGNED_ZENCODE string

func KeypairoomServer(conf *config.KeypairoomConfig, data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"userData":       data,
		"serverSideSalt": conf.Salt,
	})
	if err != nil {
		return "", err
	}

	result, success := zenroom.ZencodeExec(KEYPAIROOM_ZENCODE, "", string(jsonData), "")
	if !success {
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

	result, success := zenroom.ZencodeExec(PUBKEYS_REQUEST_SIGNED_ZENCODE, "", string(jsonData), "")
	if !success {
		return zenroomResult, errors.New(result.Logs)
	}

	err = json.Unmarshal([]byte(result.Output), &zenroomResult)
	if err != nil {
		return zenroomResult, err
	}

	return zenroomResult, nil
}
