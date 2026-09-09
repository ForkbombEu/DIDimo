// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/forkbombeu/credimi/pkg/fcaf/evidence"
)

func validateTrustedAuthoritiesMatch(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet vp_token is not an object keyed by credential query ID",
		}
	}
	for credentialIndex, rawCredential := range credentials {
		credential, _ := normalizeJSONObject(rawCredential)
		id, _ := credential["id"].(string)
		presentations, ok := response[id].([]any)
		if !ok || len(presentations) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
			}
		}
		authorities, hasTA := credential["trusted_authorities"].([]any)
		if !hasTA || len(authorities) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d] does not contain trusted_authorities",
					credentialIndex,
				),
			}
		}
		for presentationIndex, rawPresentation := range presentations {
			token, ok := rawPresentation.(string)
			if !ok || token == "" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not an SD-JWT presentation",
						id,
						presentationIndex,
					),
				}
			}
			presentation, err := evidence.ParseSDJWTPresentation(token)
			if err != nil {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not a valid SD-JWT: %v",
						id,
						presentationIndex,
						err,
					),
				}
			}
			if !credentialMatchesTrustedAuthorities(presentation, authorities) {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] issuer does not match any trusted_authority",
						id,
						presentationIndex,
					),
				}
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "every returned credential issuer matches at least one trusted_authority",
	}
}
func validateTrustedAuthoritiesNoMatch(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}
	for index, rawCredential := range credentials {
		credential, _ := normalizeJSONObject(rawCredential)
		authorities, ok := credential["trusted_authorities"].([]any)
		if !ok || len(authorities) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] does not contain trusted_authorities", index),
			}
		}
		for authorityIndex, rawAuthority := range authorities {
			authority, ok := normalizeJSONObject(rawAuthority)
			if !ok || authority["type"] != "aki" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].trusted_authorities[%d] is not a valid aki authority",
						index,
						authorityIndex,
					),
				}
			}
			values, ok := authority["values"].([]any)
			if !ok || len(values) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].trusted_authorities[%d].values is empty",
						index,
						authorityIndex,
					),
				}
			}
			for valueIndex, rawValue := range values {
				value, ok := rawValue.(string)
				if !ok || value == "" {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].values[%d] is not a string",
							index,
							authorityIndex,
							valueIndex,
						),
					}
				}
				decoded, err := base64.RawURLEncoding.DecodeString(value)
				if err != nil || len(decoded) == 0 {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].values[%d] is not base64url",
							index,
							authorityIndex,
							valueIndex,
						),
					}
				}
			}
		}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for an unmatched trusted_authorities query",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned no credential for valid unmatched trusted_authorities",
	}
}
func credentialMatchesTrustedAuthorities(
	presentation *evidence.SDJWTPresentation,
	authorities []any,
) bool {
	for _, rawAuthority := range authorities {
		authority, ok := normalizeJSONObject(rawAuthority)
		if !ok {
			continue
		}
		authType, _ := authority["type"].(string)
		if authType == "" {
			continue
		}
		values, _ := authority["values"].([]any)
		if len(values) == 0 {
			continue
		}
		switch authType {
		case "aki":
			if sdjwtMatchesAKI(presentation, values) {
				return true
			}
		default:
			if sdjwtMatchesIssuerClaim(presentation, values) {
				return true
			}
		}
	}
	return false
}
func sdjwtMatchesAKI(presentation *evidence.SDJWTPresentation, values []any) bool {
	rawChain, ok := presentation.ProtectedHeaders["x5c"].([]any)
	if !ok || len(rawChain) == 0 {
		return false
	}
	encoded, ok := rawChain[0].(string)
	if !ok || encoded == "" {
		return false
	}
	der, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return false
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return false
	}
	if len(cert.AuthorityKeyId) == 0 {
		return false
	}
	encodedAKI := base64.RawURLEncoding.EncodeToString(cert.AuthorityKeyId)
	for _, rawValue := range values {
		value, ok := rawValue.(string)
		if ok && value == encodedAKI {
			return true
		}
	}
	return false
}
func sdjwtMatchesIssuerClaim(presentation *evidence.SDJWTPresentation, values []any) bool {
	iss, _ := presentation.IssuerPayload["iss"].(string)
	if iss == "" {
		return false
	}
	for _, rawValue := range values {
		value, ok := rawValue.(string)
		if ok && value == iss {
			return true
		}
	}
	return false
}
