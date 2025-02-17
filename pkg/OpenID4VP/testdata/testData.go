package testdata

type Server struct {
	AuthorizationEndpoint string `yaml:"authorization_endpoint"`
}

type Field struct {
	Path   []string          `yaml:"path"`
	Filter map[string]string `yaml:"filter"`
}

type Constraints struct {
	Fields []Field `yaml:"fields"`
}

type Format struct {
	VCSDJWT map[string][]string `yaml:"vc+sd-jwt"`
}

type InputDescriptor struct {
	ID          string      `yaml:"id"`
	Constraints Constraints `yaml:"constraints"`
	Format      Format      `yaml:"format"`
}

type PresentationDefinition struct {
	ID               string            `yaml:"id"`
	InputDescriptors []InputDescriptor `yaml:"input_descriptors"`
}

type Client struct {
	ClientID               string                 `yaml:"client_id"`
	PresentationDefinition PresentationDefinition `yaml:"presentation_definition"`
	JWKS                   JWKS                   `yaml:"jwks"`
}

type JWKS struct {
	Keys []JWKKey `yaml:"keys"`
}

type JWKKey struct {
	Kty string `yaml:"kty"`
	Alg string `yaml:"alg"`
	Crv string `yaml:"crv"`
	D   string `yaml:"d"`
	X   string `yaml:"x"`
	Y   string `yaml:"y"`
}

type JSONPayload struct {
	Alias       string `yaml:"alias"`
	Description string `yaml:"description"`
	Server      Server `yaml:"server"`
	Client      Client `yaml:"client"`
}
