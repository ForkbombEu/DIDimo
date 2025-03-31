package testdata

type Server struct {
	AuthorizationEndpoint string `yaml:"authorization_endpoint" json:"authorization_endpoint"`
}

type Client struct {
	ClientID                          string      `yaml:"client_id,omitempty" json:"client_id,omitempty"`
	AuthorizationEncryptedResponseAlg string      `yaml:"authorization_encrypted_response_alg,omitempty" json:"authorization_encrypted_response_alg,omitempty"`
	AuthorizationEncryptedResponseEnc string      `yaml:"authorization_encrypted_response_enc,omitempty" json:"authorization_encrypted_response_enc,omitempty"`
	PresentationDefinition            interface{} `yaml:"presentation_definition" json:"presentation_definition"`
	JWKS                              interface{} `yaml:"jwks" json:"jwks"`
}

// type Form struct {
// 	Alias       string `yaml:"alias" json:"alias"`
// 	Description string `yaml:"description,omitempty" json:"description,omitempty"`
// 	Server      Server `yaml:"server" json:"server"`
// 	Client      Client `yaml:"client" json:"client"`
// }
type Form any
