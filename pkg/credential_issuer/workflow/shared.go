package workflow

import credentialissuer "github.com/forkbombeu/didimo/pkg/credential_issuer"

const FidesIssuersUrl = "https://credential-catalog.fides.community/api/public/credentialtype?includeAllDetails=false&size=200"

type FetchIssuersActivityResponse struct{ Issuers []string }

type FidesResponse struct {
	Content []struct {
		IssuanceURL               string `json:"issuanceUrl"`
		CredentialConfigurationID string `json:"credentialConfigurationId"`
		IssuePortalURL            string `json:"issuePortalUrl,omitempty"`
	} `json:"content"`
	Page struct {
		Size          int `json:"size"`
		Number        int `json:"number"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
	} `json:"page"`
}

type CredentialWorkflowInput struct {
	BaseURL  string // Base URL for the credential issuer
	IssuerID string // ID of the credentials issuer from PB
}

type CredentialWorkflowResponse struct {
	Message string
}

type CreateCredentialIssuersInput struct {
	Issuers []string
	DBPath  string
}

type Credential struct {
	CredentialDefinition                 *credentialissuer.CredentialDefinition                      `json:"credential_definition,omitempty"`
	CredentialSigningAlgValuesSupported  []credentialissuer.CredentialSigningAlgValuesSupportedElem  `json:"credential_signing_alg_values_supported,omitempty"`
	CryptographicBindingMethodsSupported []credentialissuer.CryptographicBindingMethodsSupportedElem `json:"cryptographic_binding_methods_supported,omitempty"`
	Display                              []credentialissuer.DisplayElem_1                            `json:"display,omitempty"`
	Format                               string                                                      `json:"format"`
	ProofTypesSupported                  credentialissuer.ProofTypesSupported                        `json:"proof_types_supported,omitempty"`
	Scope                                *string                                                     `json:"scope,omitempty"`
}

type StoreCredentialsActivityInput struct {
	IssuerData *credentialissuer.OpenidCredentialIssuerSchemaJson
	IssuerID   string
	DBPath     string
	CredKey   string
	IssuerName string
	Credential Credential
}

const FetchIssuersTaskQueue = "FetchIssuersTaskQueue"

const wellKnownJSON = `{
  "credential_issuer": "https://dev.issuer1.forkbomb.eu/credential_issuer",
  "credential_endpoint": "https://dev.issuer1.forkbomb.eu/credential_issuer/credential",
  "authorization_servers": [
    "https://dev.authz-server1.forkbomb.eu/authz_server"
  ],
  "display": [
    {
      "name": "Forkbomb Test Issuer",
      "locale": "en-US"
    }
  ],
  "jwks": {
    "keys": [
      {
        "kid": "did:dyne:sandbox.genericissuer:3suepGGjNHJmGDBebsCmapkdfBfXwFZzEQcEAMu7EdwA#es256_public_key",
        "crv": "P-256",
        "alg": "ES256",
        "kty": "EC"
      }
    ]
  },
  "credential_configurations_supported": {
    "discount_from_voucher": {
      "format": "vc+sd-jwt",
      "cryptographic_binding_methods_supported": [
        "jwk",
        "did:dyne:sandbox.signroom"
      ],
      "credential_signing_alg_values_supported": [
        "ES256"
      ],
      "proof_types_supported": {
        "jwt": {
          "proof_signing_alg_values_supported": [
            "ES256"
          ]
        }
      },
      "display": [
        {
          "name": "Get discount from Voucher dev",
          "locale": "en-US",
          "logo": {
            "url": "https://avatars.githubusercontent.com/u/96812851?s=200&v=4",
            "alt_text": "Get discount from Voucher dev logo",
            "uri": "https://avatars.githubusercontent.com/u/96812851?s=200&v=4"
          },
          "background_color": "#12107c",
          "text_color": "#FFFFFF",
          "description": "Get a special discount for all plans of DIDroom! Enter your voucher and get a discount credential."
        }
      ],
      "vct": "discount_from_voucher",
      "claims": {
        "has_discount_from_voucher": {
          "mandatory": true,
          "display": [
            {
              "locale": "en-US",
              "name": "Has a discount from Voucher"
            }
          ]
        }
      }
    }
  }
}`

