package testdata

type Server struct {
	AuthorizationEndpoint string `yaml:"authorization_endpoint" json:"authorization_endpoint"`
}

type Field struct {
	Path   []string          `yaml:"path" json:"path"`
	Filter map[string]string `yaml:"filter" json:"filter"`
}

type Constraints struct {
	Fields []Field `yaml:"fields" json:"fields"`
}

type Format struct {
	VCSDJWT map[string][]string `yaml:"vc+sd-jwt" json:"vc+sd-jwt"`
}

type InputDescriptor struct {
	ID          string      `yaml:"id" json:"id"`
	Constraints Constraints `yaml:"constraints" json:"constraints"`
	Format      Format      `yaml:"format" json:"format"`
}

type PresentationDefinition struct {
	ID               string            `yaml:"id" json:"id"`
	InputDescriptors []InputDescriptor `yaml:"input_descriptors" json:"input_descriptors"`
}

type Client struct {
	ClientID               string                 `yaml:"client_id" json:"client_id"`
	PresentationDefinition PresentationDefinition `yaml:"presentation_definition" json:"presentation_definition"`
	JWKS                   JWKS                   `yaml:"jwks" json:"jwks"`
}

type JWKS struct {
	Keys []JWKKey `yaml:"keys" json:"keys"`
}

type JWKKey struct {
	Kty string `yaml:"kty" json:"kty"`
	Alg string `yaml:"alg" json:"alg"`
	Crv string `yaml:"crv" json:"crv"`
	D   string `yaml:"d" json:"d"`
	X   string `yaml:"x" json:"x"`
	Y   string `yaml:"y" json:"y"`
}

type Form struct {
	Alias       string `yaml:"alias" json:"alias"`
	Description string `yaml:"description" json:"description"`
	Server      Server `yaml:"server" json:"server"`
	Client      Client `yaml:"client" json:"client"`
}
