// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package credentialissuer

import "encoding/json"
import "fmt"
import "reflect"

type CredentialDefinition struct {
	// CredentialSubject corresponds to the JSON schema field "credentialSubject".
	CredentialSubject CredentialDefinitionCredentialSubject `json:"credentialSubject,omitempty" yaml:"credentialSubject,omitempty" mapstructure:"credentialSubject,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type []string `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
}

type CredentialDefinitionCredentialSubject map[string]struct {
	// Display corresponds to the JSON schema field "display".
	Display []DisplayElem `json:"display,omitempty" yaml:"display,omitempty" mapstructure:"display,omitempty"`

	// Mandatory corresponds to the JSON schema field "mandatory".
	Mandatory *bool `json:"mandatory,omitempty" yaml:"mandatory,omitempty" mapstructure:"mandatory,omitempty"`
}

type CredentialSigningAlgValuesSupportedElem string

const CredentialSigningAlgValuesSupportedElemES256 CredentialSigningAlgValuesSupportedElem = "ES256"
const CredentialSigningAlgValuesSupportedElemES256K CredentialSigningAlgValuesSupportedElem = "ES256K"
const CredentialSigningAlgValuesSupportedElemEdDSA CredentialSigningAlgValuesSupportedElem = "EdDSA"
const CredentialSigningAlgValuesSupportedElemRS256 CredentialSigningAlgValuesSupportedElem = "RS256"
const CredentialSigningAlgValuesSupportedElemRSA CredentialSigningAlgValuesSupportedElem = "RSA"
const CredentialSigningAlgValuesSupportedElemRsaSignature2018 CredentialSigningAlgValuesSupportedElem = "RsaSignature2018"

var enumValues_CredentialSigningAlgValuesSupportedElem = []interface{}{
	"ES256",
	"EdDSA",
	"RS256",
	"ES256K",
	"RSA",
	"RsaSignature2018",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CredentialSigningAlgValuesSupportedElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_CredentialSigningAlgValuesSupportedElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_CredentialSigningAlgValuesSupportedElem, v)
	}
	*j = CredentialSigningAlgValuesSupportedElem(v)
	return nil
}

type CryptographicBindingMethodsSupportedElem string

const CryptographicBindingMethodsSupportedElemCoseKey CryptographicBindingMethodsSupportedElem = "cose_key"
const CryptographicBindingMethodsSupportedElemDid CryptographicBindingMethodsSupportedElem = "did"
const CryptographicBindingMethodsSupportedElemDidDyne CryptographicBindingMethodsSupportedElem = "did:dyne"
const CryptographicBindingMethodsSupportedElemDidDyneSandboxSignroom CryptographicBindingMethodsSupportedElem = "did:dyne:sandbox.signroom"
const CryptographicBindingMethodsSupportedElemDidEbsi CryptographicBindingMethodsSupportedElem = "did:ebsi"
const CryptographicBindingMethodsSupportedElemDidJwk CryptographicBindingMethodsSupportedElem = "did:jwk"
const CryptographicBindingMethodsSupportedElemDidKey CryptographicBindingMethodsSupportedElem = "did:key"
const CryptographicBindingMethodsSupportedElemDidWeb CryptographicBindingMethodsSupportedElem = "did:web"
const CryptographicBindingMethodsSupportedElemJWK CryptographicBindingMethodsSupportedElem = "JWK"
const CryptographicBindingMethodsSupportedElemJwk CryptographicBindingMethodsSupportedElem = "jwk"

var enumValues_CryptographicBindingMethodsSupportedElem = []interface{}{
	"JWK",
	"jwk",
	"did",
	"did:web",
	"did:ebsi",
	"did:jwk",
	"did:dyne",
	"did:dyne:sandbox.signroom",
	"did:key",
	"cose_key",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CryptographicBindingMethodsSupportedElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_CryptographicBindingMethodsSupportedElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_CryptographicBindingMethodsSupportedElem, v)
	}
	*j = CryptographicBindingMethodsSupportedElem(v)
	return nil
}

type DisplayElem struct {
	// Locale corresponds to the JSON schema field "locale".
	Locale *string `json:"locale,omitempty" yaml:"locale,omitempty" mapstructure:"locale,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

type DisplayElemLogo struct {
	// AltText corresponds to the JSON schema field "alt_text".
	AltText *string `json:"alt_text,omitempty" yaml:"alt_text,omitempty" mapstructure:"alt_text,omitempty"`

	// Uri corresponds to the JSON schema field "uri".
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DisplayElemLogo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in DisplayElemLogo: required")
	}
	type Plain DisplayElemLogo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DisplayElemLogo(plain)
	return nil
}

type DisplayElem_1 struct {
	// Locale corresponds to the JSON schema field "locale".
	Locale *string `json:"locale,omitempty" yaml:"locale,omitempty" mapstructure:"locale,omitempty"`

	// Logo corresponds to the JSON schema field "logo".
	Logo *DisplayElemLogo `json:"logo,omitempty" yaml:"logo,omitempty" mapstructure:"logo,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DisplayElem_1) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in DisplayElem_1: required")
	}
	type Plain DisplayElem_1
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DisplayElem_1(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DisplayElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in DisplayElem: required")
	}
	type Plain DisplayElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DisplayElem(plain)
	return nil
}

type KeyAttestationsRequired struct {
	// KeyStorage corresponds to the JSON schema field "key_storage".
	KeyStorage []string `json:"key_storage,omitempty" yaml:"key_storage,omitempty" mapstructure:"key_storage,omitempty"`

	// UserAuthentication corresponds to the JSON schema field "user_authentication".
	UserAuthentication []string `json:"user_authentication,omitempty" yaml:"user_authentication,omitempty" mapstructure:"user_authentication,omitempty"`
}

type OpenidCredentialIssuerSchemaJson struct {
	// Array of OAuth 2.0 Authorization Server identifiers.
	AuthorizationServers []string `json:"authorization_servers,omitempty" yaml:"authorization_servers,omitempty" mapstructure:"authorization_servers,omitempty"`

	// BatchCredentialIssuance corresponds to the JSON schema field
	// "batch_credential_issuance".
	BatchCredentialIssuance *OpenidCredentialIssuerSchemaJsonBatchCredentialIssuance `json:"batch_credential_issuance,omitempty" yaml:"batch_credential_issuance,omitempty" mapstructure:"batch_credential_issuance,omitempty"`

	// CredentialConfigurationsSupported corresponds to the JSON schema field
	// "credential_configurations_supported".
	CredentialConfigurationsSupported OpenidCredentialIssuerSchemaJsonCredentialConfigurationsSupported `json:"credential_configurations_supported" yaml:"credential_configurations_supported" mapstructure:"credential_configurations_supported"`

	// URL of the Credential Issuer's Credential Endpoint.
	CredentialEndpoint string `json:"credential_endpoint" yaml:"credential_endpoint" mapstructure:"credential_endpoint"`

	// The Credential Issuer's identifier
	CredentialIssuer string `json:"credential_issuer" yaml:"credential_issuer" mapstructure:"credential_issuer"`

	// CredentialResponseEncryption corresponds to the JSON schema field
	// "credential_response_encryption".
	CredentialResponseEncryption *OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption `json:"credential_response_encryption,omitempty" yaml:"credential_response_encryption,omitempty" mapstructure:"credential_response_encryption,omitempty"`

	// URL of the Credential Issuer's Deferred Credential Endpoint.
	DeferredCredentialEndpoint *string `json:"deferred_credential_endpoint,omitempty" yaml:"deferred_credential_endpoint,omitempty" mapstructure:"deferred_credential_endpoint,omitempty"`

	// Display corresponds to the JSON schema field "display".
	Display []OpenidCredentialIssuerSchemaJsonDisplayElem `json:"display,omitempty" yaml:"display,omitempty" mapstructure:"display,omitempty"`

	// URL of the Credential Issuer's Nonce Endpoint.
	NonceEndpoint *string `json:"nonce_endpoint,omitempty" yaml:"nonce_endpoint,omitempty" mapstructure:"nonce_endpoint,omitempty"`

	// URL of the Credential Issuer's Notification Endpoint.
	NotificationEndpoint *string `json:"notification_endpoint,omitempty" yaml:"notification_endpoint,omitempty" mapstructure:"notification_endpoint,omitempty"`

	// JWT containing Credential Issuer metadata parameters as claims.
	SignedMetadata *string `json:"signed_metadata,omitempty" yaml:"signed_metadata,omitempty" mapstructure:"signed_metadata,omitempty"`
}

type OpenidCredentialIssuerSchemaJsonBatchCredentialIssuance struct {
	// BatchSize corresponds to the JSON schema field "batch_size".
	BatchSize int `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OpenidCredentialIssuerSchemaJsonBatchCredentialIssuance) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["batch_size"]; raw != nil && !ok {
		return fmt.Errorf("field batch_size in OpenidCredentialIssuerSchemaJsonBatchCredentialIssuance: required")
	}
	type Plain OpenidCredentialIssuerSchemaJsonBatchCredentialIssuance
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if 1 > plain.BatchSize {
		return fmt.Errorf("field %s: must be >= %v", "batch_size", 1)
	}
	*j = OpenidCredentialIssuerSchemaJsonBatchCredentialIssuance(plain)
	return nil
}

type OpenidCredentialIssuerSchemaJsonCredentialConfigurationsSupported map[string]struct {
	// CredentialDefinition corresponds to the JSON schema field
	// "credential_definition".
	CredentialDefinition *CredentialDefinition `json:"credential_definition,omitempty" yaml:"credential_definition,omitempty" mapstructure:"credential_definition,omitempty"`

	// CredentialSigningAlgValuesSupported corresponds to the JSON schema field
	// "credential_signing_alg_values_supported".
	CredentialSigningAlgValuesSupported []CredentialSigningAlgValuesSupportedElem `json:"credential_signing_alg_values_supported,omitempty" yaml:"credential_signing_alg_values_supported,omitempty" mapstructure:"credential_signing_alg_values_supported,omitempty"`

	// CryptographicBindingMethodsSupported corresponds to the JSON schema field
	// "cryptographic_binding_methods_supported".
	CryptographicBindingMethodsSupported []CryptographicBindingMethodsSupportedElem `json:"cryptographic_binding_methods_supported,omitempty" yaml:"cryptographic_binding_methods_supported,omitempty" mapstructure:"cryptographic_binding_methods_supported,omitempty"`

	// Display corresponds to the JSON schema field "display".
	Display []DisplayElem_1 `json:"display,omitempty" yaml:"display,omitempty" mapstructure:"display,omitempty"`

	// Format corresponds to the JSON schema field "format".
	Format string `json:"format" yaml:"format" mapstructure:"format"`

	// ProofTypesSupported corresponds to the JSON schema field
	// "proof_types_supported".
	ProofTypesSupported ProofTypesSupported `json:"proof_types_supported,omitempty" yaml:"proof_types_supported,omitempty" mapstructure:"proof_types_supported,omitempty"`

	// Scope corresponds to the JSON schema field "scope".
	Scope *string `json:"scope,omitempty" yaml:"scope,omitempty" mapstructure:"scope,omitempty"`
}

type OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption struct {
	// AlgValuesSupported corresponds to the JSON schema field "alg_values_supported".
	AlgValuesSupported []OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem `json:"alg_values_supported" yaml:"alg_values_supported" mapstructure:"alg_values_supported"`

	// EncValuesSupported corresponds to the JSON schema field "enc_values_supported".
	EncValuesSupported []OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem `json:"enc_values_supported" yaml:"enc_values_supported" mapstructure:"enc_values_supported"`

	// EncryptionRequired corresponds to the JSON schema field "encryption_required".
	EncryptionRequired bool `json:"encryption_required" yaml:"encryption_required" mapstructure:"encryption_required"`
}

type OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem string

const OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElemES256 OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem = "ES256"
const OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElemEdDSA OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem = "EdDSA"
const OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElemRS256 OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem = "RS256"

var enumValues_OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem = []interface{}{
	"ES256",
	"EdDSA",
	"RS256",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem, v)
	}
	*j = OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionAlgValuesSupportedElem(v)
	return nil
}

type OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem string

const OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElemA128CBCHS256 OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem = "A128CBC-HS256"
const OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElemA128GCM OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem = "A128GCM"

var enumValues_OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem = []interface{}{
	"A128CBC-HS256",
	"A128GCM",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem, v)
	}
	*j = OpenidCredentialIssuerSchemaJsonCredentialResponseEncryptionEncValuesSupportedElem(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["alg_values_supported"]; raw != nil && !ok {
		return fmt.Errorf("field alg_values_supported in OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption: required")
	}
	if _, ok := raw["enc_values_supported"]; raw != nil && !ok {
		return fmt.Errorf("field enc_values_supported in OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption: required")
	}
	if _, ok := raw["encryption_required"]; raw != nil && !ok {
		return fmt.Errorf("field encryption_required in OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption: required")
	}
	type Plain OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OpenidCredentialIssuerSchemaJsonCredentialResponseEncryption(plain)
	return nil
}

type OpenidCredentialIssuerSchemaJsonDisplayElem struct {
	// Locale corresponds to the JSON schema field "locale".
	Locale *string `json:"locale,omitempty" yaml:"locale,omitempty" mapstructure:"locale,omitempty"`

	// Logo corresponds to the JSON schema field "logo".
	Logo *OpenidCredentialIssuerSchemaJsonDisplayElemLogo `json:"logo,omitempty" yaml:"logo,omitempty" mapstructure:"logo,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
}

type OpenidCredentialIssuerSchemaJsonDisplayElemLogo struct {
	// AltText corresponds to the JSON schema field "alt_text".
	AltText *string `json:"alt_text,omitempty" yaml:"alt_text,omitempty" mapstructure:"alt_text,omitempty"`

	// Uri corresponds to the JSON schema field "uri".
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OpenidCredentialIssuerSchemaJsonDisplayElemLogo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in OpenidCredentialIssuerSchemaJsonDisplayElemLogo: required")
	}
	type Plain OpenidCredentialIssuerSchemaJsonDisplayElemLogo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OpenidCredentialIssuerSchemaJsonDisplayElemLogo(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OpenidCredentialIssuerSchemaJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["credential_configurations_supported"]; raw != nil && !ok {
		return fmt.Errorf("field credential_configurations_supported in OpenidCredentialIssuerSchemaJson: required")
	}
	if _, ok := raw["credential_endpoint"]; raw != nil && !ok {
		return fmt.Errorf("field credential_endpoint in OpenidCredentialIssuerSchemaJson: required")
	}
	if _, ok := raw["credential_issuer"]; raw != nil && !ok {
		return fmt.Errorf("field credential_issuer in OpenidCredentialIssuerSchemaJson: required")
	}
	type Plain OpenidCredentialIssuerSchemaJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OpenidCredentialIssuerSchemaJson(plain)
	return nil
}

type ProofSigningAlgValuesSupportedElem string

const ProofSigningAlgValuesSupportedElemES256 ProofSigningAlgValuesSupportedElem = "ES256"
const ProofSigningAlgValuesSupportedElemEdDSA ProofSigningAlgValuesSupportedElem = "EdDSA"

var enumValues_ProofSigningAlgValuesSupportedElem = []interface{}{
	"ES256",
	"EdDSA",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProofSigningAlgValuesSupportedElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_ProofSigningAlgValuesSupportedElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_ProofSigningAlgValuesSupportedElem, v)
	}
	*j = ProofSigningAlgValuesSupportedElem(v)
	return nil
}

type ProofTypesSupported map[string]struct {
	// KeyAttestationsRequired corresponds to the JSON schema field
	// "key_attestations_required".
	KeyAttestationsRequired *KeyAttestationsRequired `json:"key_attestations_required,omitempty" yaml:"key_attestations_required,omitempty" mapstructure:"key_attestations_required,omitempty"`

	// ProofSigningAlgValuesSupported corresponds to the JSON schema field
	// "proof_signing_alg_values_supported".
	ProofSigningAlgValuesSupported []ProofSigningAlgValuesSupportedElem `json:"proof_signing_alg_values_supported" yaml:"proof_signing_alg_values_supported" mapstructure:"proof_signing_alg_values_supported"`
}
