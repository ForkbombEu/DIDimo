// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
)

const invalidRequestError = "invalid_request"

var dcqlIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type DCQLResponseConstraintsValidator struct{}

func (DCQLResponseConstraintsValidator) ID() string {
	return "dcql.response_satisfies_constraints"
}

//nolint:gocyclo // DCQL modes remain explicit so each conformance branch is auditable.
func (DCQLResponseConstraintsValidator) Validate(_ context.Context, input Input) Result {
	params, err := DecodeParams[struct {
		Mode           string  `json:"mode"`
		ForbiddenPaths [][]any `json:"forbidden_paths"`
		Property       string  `json:"property"`
		ExpectedType   string  `json:"expected_type"`
		Valid          bool    `json:"valid"`
		ExpectedValue  any     `json:"expected_value"`
	}](input.Params)
	if err != nil {
		return Result{Status: StatusError, Message: err.Error()}
	}
	switch params.Mode {
	case "credential_sets",
		"claims_present",
		"claims_subset",
		"claims_union",
		"claims_path_no_match",
		"claims_values_no_match",
		"claim_id_missing_with_claim_sets",
		"claims_without_id_without_claim_sets",
		"duplicate_claim_ids",
		"empty_claim_id",
		"invalid_claim_id_characters",
		"claim_path_missing",
		"claim_path_empty",
		"claim_path_non_array",
		"claim_path_allowed_components",
		"claims_without_values",
		"credential_sets_options_missing",
		"credential_sets_options_empty",
		"credential_sets_options_non_array",
		"credential_sets_options_valid_references",
		"credential_sets_options_invalid_references",
		"credential_sets_required_true_match",
		"credential_sets_required_true_no_match",
		"credential_sets_required_omitted",
		"credential_sets_required_false_with_match",
		"credentials_match",
		"without_credential_sets",
		"without_trusted_authorities",
		"without_claims",
		"empty_claims",
		"empty_array",
		"property_type",
		"property_equals",
		"trusted_authority_property_type",
		"trusted_authority_array_item_type",
		"trusted_authority_empty_string_item",
		"multiple_default_false",
		"multiple_true",
		"no_match",
		"request_rejected",
		"trusted_authorities_match",
		"trusted_authorities_no_match",
		"claim_sets":
	default:
		return Result{
			Status:  StatusError,
			Message: "mode must be credential_sets, credentials_match, without_credential_sets, without_trusted_authorities, without_claims, empty_claims, empty_array, property_type, property_equals, trusted_authority_property_type, trusted_authority_array_item_type, trusted_authority_empty_string_item, multiple_default_false, multiple_true, no_match, request_rejected, trusted_authorities_match, trusted_authorities_no_match, claim_sets, claim_path_member_type_error, wallet_error_expected, invalid_scope, unknown_field_stripped, vp_formats_not_supported, transaction_data_error, invalid_client, invalid_request_generic, access_denied, or jwe_enc_verified",
		}
	}

	root, ok := normalizeJSONObject(input.Value)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: fmt.Sprintf("DCQL evidence is %T, expected object", input.Value),
		}
	}
	queryValue, found := findObjectKey(root, "dcql_query")
	if !found {
		return Result{Status: StatusFail, Message: "captured evidence does not contain dcql_query"}
	}
	query, ok := normalizeJSONObject(queryValue)
	if !ok {
		return Result{Status: StatusFail, Message: "captured dcql_query is not an object"}
	}

	responseValue, _ := findObjectKey(root, "vp_token")
	errorValue, _ := findObjectKey(root, "error")
	switch params.Mode {
	case "credential_sets":
		sets, ok := query["credential_sets"].([]any)
		if !ok || len(sets) == 0 {
			return Result{
				Status:  StatusFail,
				Message: "dcql_query does not contain non-empty credential_sets",
			}
		}
	case "claims_present":
		return validateClaimsPresent(query, responseValue)
	case "claims_subset":
		return validateClaimsSubset(query, responseValue, params.ForbiddenPaths)
	case "claims_union":
		return validateClaimsUnion(query, responseValue, params.ForbiddenPaths)
	case "claims_path_no_match":
		return validateClaimsPathNoMatch(query, responseValue)
	case "claims_values_no_match":
		return validateClaimsValuesNoMatch(query, responseValue)
	case "claim_id_missing_with_claim_sets":
		return validateMissingClaimIDWithClaimSets(query, responseValue)
	case "claims_without_id_without_claim_sets":
		return validateClaimsWithoutIDWithoutClaimSets(query, responseValue)
	case "duplicate_claim_ids":
		return validateDuplicateClaimIDs(query, responseValue, errorValue)
	case "empty_claim_id":
		return validateEmptyClaimID(query, responseValue, errorValue)
	case "invalid_claim_id_characters":
		return validateInvalidClaimIDCharacters(query, responseValue, errorValue)
	case "claim_path_missing":
		return validateMissingClaimPath(query, responseValue, errorValue)
	case "claim_path_empty":
		return validateEmptyClaimPath(query, responseValue, errorValue)
	case "claim_path_non_array":
		return validateNonArrayClaimPath(query, responseValue, errorValue)
	case "claim_path_allowed_components":
		return validateAllowedClaimPathComponents(query, responseValue)
	case "claims_without_values":
		return validateClaimsWithoutValues(query, responseValue)
	case "credential_sets_options_missing":
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		sets, ok := query["credential_sets"].([]any)
		if !ok || len(sets) == 0 {
			return Result{
				Status:  StatusFail,
				Message: "dcql_query does not contain non-empty credential_sets",
			}
		}
		for index, rawSet := range sets {
			set, ok := normalizeJSONObject(rawSet)
			if !ok {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credential_sets[%d] is not an object", index),
				}
			}
			if _, exists := set["options"]; exists {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credential_sets[%d].options is present", index),
				}
			}
		}
		if !isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: "wallet returned a vp_token for a query missing credential_sets.options",
			}
		}
		if errorValue != invalidRequestError {
			return Result{
				Status:  StatusFail,
				Message: "wallet did not return invalid_request for a query missing credential_sets.options",
			}
		}
		return Result{
			Status:  StatusPass,
			Message: "wallet returned invalid_request for credential_sets without options",
		}
	case "credential_sets_options_empty",
		"credential_sets_options_non_array",
		"credential_sets_options_valid_references",
		"credential_sets_options_invalid_references":
		return validateCredentialSetsOptions(query, responseValue, errorValue, params.Mode)
	case "credential_sets_required_true_match",
		"credential_sets_required_true_no_match",
		"credential_sets_required_omitted",
		"credential_sets_required_false_with_match":
		return validateCredentialSetsRequired(query, responseValue, params.Mode)
	case "credentials_match",
		"without_credential_sets",
		"without_trusted_authorities",
		"without_claims",
		"multiple_default_false",
		"multiple_true":
		if params.Mode == "without_credential_sets" {
			if _, exists := query["credential_sets"]; exists {
				return Result{
					Status:  StatusFail,
					Message: "dcql_query contains credential_sets",
				}
			}
		}
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
		for index, rawCredential := range credentials {
			credential, ok := normalizeJSONObject(rawCredential)
			if !ok {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credentials[%d] is not an object", index),
				}
			}
			if params.Mode == "without_trusted_authorities" {
				if _, exists := credential["trusted_authorities"]; exists {
					return Result{
						Status:  StatusFail,
						Message: fmt.Sprintf("credentials[%d] contains trusted_authorities", index),
					}
				}
			}
			if params.Mode == "without_claims" {
				if _, exists := credential["claims"]; exists {
					return Result{
						Status:  StatusFail,
						Message: fmt.Sprintf("credentials[%d] contains claims", index),
					}
				}
			}
			id, _ := credential["id"].(string)
			if id == "" {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credentials[%d] has no id", index),
				}
			}
			presentation := response[id]
			if isEmptyDCQLValue(presentation) {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token has no presentation for credential query %q",
						id,
					),
				}
			}
			if params.Mode == "multiple_default_false" {
				if _, exists := credential["multiple"]; exists {
					return Result{
						Status:  StatusFail,
						Message: fmt.Sprintf("credentials[%d] contains multiple", index),
					}
				}
				presentations, ok := presentation.([]any)
				if !ok || len(presentations) != 1 {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token must contain exactly one presentation for credential query %q",
							id,
						),
					}
				}
			}
			if params.Mode == "multiple_true" {
				multiple, ok := credential["multiple"].(bool)
				if !ok || !multiple {
					return Result{
						Status:  StatusFail,
						Message: fmt.Sprintf("credentials[%d].multiple is not true", index),
					}
				}
				presentations, ok := presentation.([]any)
				if !ok || len(presentations) < 2 {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token must contain multiple presentations for credential query %q",
							id,
						),
					}
				}
			}
		}
	case "empty_claims", "empty_array":
		property := "claims"
		if params.Mode == "empty_array" {
			property = params.Property
			if property == "" {
				return Result{
					Status:  StatusError,
					Message: "property is required for empty_array mode",
				}
			}
		}
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		for index, rawCredential := range credentials {
			credential, ok := normalizeJSONObject(rawCredential)
			if !ok {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credentials[%d] is not an object", index),
				}
			}
			value, exists := credential[property]
			if !exists {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credentials[%d] does not contain %s", index, property),
				}
			}
			items, ok := value.([]any)
			if !ok || len(items) != 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].%s is not an empty array",
						index,
						property,
					),
				}
			}
		}
		if !isEmptyDCQLValue(responseValue) {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"wallet returned a credential for a query with empty %s",
					property,
				),
			}
		}
	case "no_match", "request_rejected":
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		if err := validateDCQLCredentialQueries(credentials); err != nil {
			return Result{Status: StatusFail, Message: err.Error()}
		}
		if errorText, _ := errorValue.(string); errorText == invalidRequestError {
			break
		}
		if !isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: "wallet returned a credential for a no-match DCQL query",
			}
		}
	case "trusted_authorities_match":
		return validateTrustedAuthoritiesMatch(query, responseValue)
	case "trusted_authorities_no_match":
		return validateTrustedAuthoritiesNoMatch(query, responseValue)
	case "claim_path_member_type_error":
		return validateClaimPathMemberTypeError(responseValue, errorValue)
	case "wallet_error_expected":
		return validateWalletErrorExpected(responseValue, errorValue, params.ExpectedValue)
	case "invalid_scope":
		return validateErrorCode(responseValue, errorValue, "invalid_scope")
	case "unknown_field_stripped":
		return validateUnknownFieldStripped(query, responseValue)
	case "vp_formats_not_supported":
		return validateErrorCode(responseValue, errorValue, "vp_formats_not_supported")
	case "transaction_data_error":
		return validateErrorCode(responseValue, errorValue, "invalid_transaction_data")
	case "invalid_client":
		return validateErrorCode(responseValue, errorValue, "invalid_client")
	case "invalid_request_generic":
		return validateErrorCode(responseValue, errorValue, invalidRequestError)
	case "access_denied":
		return validateErrorCode(responseValue, errorValue, "access_denied")
	case "jwe_enc_verified":
		return validateJWEEncVerified(responseValue)
	case "session_encryption":
		return validateSessionEncryption(responseValue)
	case "encoding":
		return validateEncoding(responseValue)
	case "interaction_engagement":
		return validateInteractionCompleted(responseValue)
	case "interaction_protocol_flow":
		return validateInteractionCompleted(responseValue)
	case "interaction_supportive":
		return validateInteractionCompleted(responseValue)
	case "rp_integrity_positive":
		return validateEvidencePresent(responseValue, errorValue)
	case "issuer_integrity_positive":
		return validateEvidencePresent(responseValue, errorValue)
	case "property_type":
		if params.Property == "" {
			return Result{
				Status:  StatusError,
				Message: "property is required for property_type mode",
			}
		}
		if _, exists := input.Params["valid"]; !exists {
			return Result{Status: StatusError, Message: "valid is required for property_type mode"}
		}
		if !supportedJSONType(params.ExpectedType) {
			return Result{
				Status:  StatusError,
				Message: "expected_type must be boolean, string, number, integer, array, object, or null",
			}
		}
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		for index, rawCredential := range credentials {
			credential, ok := normalizeJSONObject(rawCredential)
			if !ok {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credentials[%d] is not an object", index),
				}
			}
			value, exists := credential[params.Property]
			if !exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d] does not contain %s",
						index,
						params.Property,
					),
				}
			}
			matches := matchesJSONType(value, params.ExpectedType)
			if matches != params.Valid {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].%s type validity is %t, expected %t",
						index,
						params.Property,
						matches,
						params.Valid,
					),
				}
			}
		}
		if params.Valid && isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("wallet returned no credential for valid %s", params.Property),
			}
		}
		if !params.Valid && !isEmptyDCQLValue(responseValue) {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"wallet returned a credential for invalid %s type",
					params.Property,
				),
			}
		}
	case "property_equals":
		if params.Property == "" {
			return Result{
				Status:  StatusError,
				Message: "property is required for property_equals mode",
			}
		}
		if _, exists := input.Params["expected_value"]; !exists {
			return Result{
				Status:  StatusError,
				Message: "expected_value is required for property_equals mode",
			}
		}
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		for index, rawCredential := range credentials {
			credential, ok := normalizeJSONObject(rawCredential)
			if !ok {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("credentials[%d] is not an object", index),
				}
			}
			value, exists := credential[params.Property]
			if !exists || !reflect.DeepEqual(value, params.ExpectedValue) {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].%s does not equal the expected value",
						index,
						params.Property,
					),
				}
			}
		}
		if isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("wallet returned no credential for %s", params.Property),
			}
		}
	case "trusted_authority_array_item_type":
		if params.Property == "" {
			return Result{
				Status:  StatusError,
				Message: "property is required for trusted_authority_array_item_type mode",
			}
		}
		if _, exists := input.Params["valid"]; !exists {
			return Result{
				Status:  StatusError,
				Message: "valid is required for trusted_authority_array_item_type mode",
			}
		}
		if !supportedJSONType(params.ExpectedType) || params.ExpectedType != "array" {
			return Result{
				Status:  StatusError,
				Message: "expected_type must be array for trusted_authority_array_item_type mode",
			}
		}
		itemExpectedType, ok := input.Params["item_expected_type"].(string)
		if !ok || !supportedJSONType(itemExpectedType) {
			return Result{
				Status:  StatusError,
				Message: "item_expected_type must be boolean, string, number, integer, array, object, or null",
			}
		}
		itemValid, ok := input.Params["item_valid"].(bool)
		if !ok {
			return Result{Status: StatusError, Message: "item_valid must be a boolean"}
		}
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		if err := validateDCQLCredentialQueries(credentials); err != nil {
			return Result{Status: StatusFail, Message: err.Error()}
		}
		for credentialIndex, rawCredential := range credentials {
			credential, _ := normalizeJSONObject(rawCredential)
			authorities, ok := credential["trusted_authorities"].([]any)
			if !ok || len(authorities) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].trusted_authorities is not a non-empty array",
						credentialIndex,
					),
				}
			}
			for authorityIndex, rawAuthority := range authorities {
				authority, ok := normalizeJSONObject(rawAuthority)
				if !ok {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d] is not an object",
							credentialIndex,
							authorityIndex,
						),
					}
				}
				if authorityType, ok := authority["type"].(string); !ok || authorityType == "" {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].type is not a non-empty string",
							credentialIndex,
							authorityIndex,
						),
					}
				}
				values, ok := authority[params.Property].([]any)
				if !ok || len(values) == 0 {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].%s is not a non-empty array",
							credentialIndex,
							authorityIndex,
							params.Property,
						),
					}
				}
				foundInvalidItem := false
				for itemIndex, item := range values {
					matches := matchesJSONType(item, itemExpectedType)
					if !itemValid && !matches {
						foundInvalidItem = true
					}
					if itemValid && !matches {
						return Result{
							Status: StatusFail,
							Message: fmt.Sprintf(
								"credentials[%d].trusted_authorities[%d].%s[%d] type validity is %t, expected %t",
								credentialIndex,
								authorityIndex,
								params.Property,
								itemIndex,
								matches,
								itemValid,
							),
						}
					}
				}
				if !itemValid && !foundInvalidItem {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].%s contains no invalid item",
							credentialIndex,
							authorityIndex,
							params.Property,
						),
					}
				}
			}
		}
		if !params.Valid && !isEmptyDCQLValue(responseValue) {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"wallet returned a credential for invalid trusted authority %s item type",
					params.Property,
				),
			}
		}
	case "trusted_authority_property_type":
		if params.Property == "" {
			return Result{
				Status:  StatusError,
				Message: "property is required for trusted_authority_property_type mode",
			}
		}
		if _, exists := input.Params["valid"]; !exists {
			return Result{
				Status:  StatusError,
				Message: "valid is required for trusted_authority_property_type mode",
			}
		}
		if !supportedJSONType(params.ExpectedType) {
			return Result{
				Status:  StatusError,
				Message: "expected_type must be boolean, string, number, integer, array, object, or null",
			}
		}
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		if err := validateDCQLCredentialQueries(credentials); err != nil {
			return Result{Status: StatusFail, Message: err.Error()}
		}
		for credentialIndex, rawCredential := range credentials {
			credential, _ := normalizeJSONObject(rawCredential)
			authorities, ok := credential["trusted_authorities"].([]any)
			if !ok || len(authorities) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].trusted_authorities is not a non-empty array",
						credentialIndex,
					),
				}
			}
			for authorityIndex, rawAuthority := range authorities {
				authority, ok := normalizeJSONObject(rawAuthority)
				if !ok {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d] is not an object",
							credentialIndex,
							authorityIndex,
						),
					}
				}
				value, exists := authority[params.Property]
				if !exists {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d] does not contain %s",
							credentialIndex,
							authorityIndex,
							params.Property,
						),
					}
				}
				if params.Property == "type" && !nonEmptyStringArray(authority["values"]) {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].values is not a non-empty string array",
							credentialIndex,
							authorityIndex,
						),
					}
				}
				if params.Property == "values" {
					authorityType, ok := authority["type"].(string)
					if !ok || authorityType == "" {
						return Result{
							Status: StatusFail,
							Message: fmt.Sprintf(
								"credentials[%d].trusted_authorities[%d].type is not a non-empty string",
								credentialIndex,
								authorityIndex,
							),
						}
					}
				}
				matches := matchesJSONType(value, params.ExpectedType)
				if matches != params.Valid {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].%s type validity is %t, expected %t",
							credentialIndex,
							authorityIndex,
							params.Property,
							matches,
							params.Valid,
						),
					}
				}
			}
		}
		if params.Valid && isEmptyDCQLValue(responseValue) {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"wallet returned no credential for valid trusted authority %s",
					params.Property,
				),
			}
		}
		if !params.Valid && !isEmptyDCQLValue(responseValue) {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"wallet returned a credential for invalid trusted authority %s type",
					params.Property,
				),
			}
		}
	case "trusted_authority_empty_string_item":
		if params.Property == "" {
			return Result{
				Status:  StatusError,
				Message: "property is required for trusted_authority_empty_string_item mode",
			}
		}
		credentials, ok := query["credentials"].([]any)
		if !ok || len(credentials) == 0 {
			return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
		}
		if err := validateDCQLCredentialQueries(credentials); err != nil {
			return Result{Status: StatusFail, Message: err.Error()}
		}
		for credentialIndex, rawCredential := range credentials {
			credential, _ := normalizeJSONObject(rawCredential)
			authorities, ok := credential["trusted_authorities"].([]any)
			if !ok || len(authorities) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].trusted_authorities is not a non-empty array",
						credentialIndex,
					),
				}
			}
			for authorityIndex, rawAuthority := range authorities {
				authority, ok := normalizeJSONObject(rawAuthority)
				if !ok {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d] is not an object",
							credentialIndex,
							authorityIndex,
						),
					}
				}
				if authorityType, ok := authority["type"].(string); !ok || authorityType == "" {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].type is not a non-empty string",
							credentialIndex,
							authorityIndex,
						),
					}
				}
				values, ok := authority[params.Property].([]any)
				if !ok || len(values) == 0 {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].%s is not a non-empty array",
							credentialIndex,
							authorityIndex,
							params.Property,
						),
					}
				}
				foundEmpty := false
				for itemIndex, item := range values {
					itemString, ok := item.(string)
					if !ok {
						return Result{
							Status: StatusFail,
							Message: fmt.Sprintf(
								"credentials[%d].trusted_authorities[%d].%s[%d] is not a string",
								credentialIndex,
								authorityIndex,
								params.Property,
								itemIndex,
							),
						}
					}
					if itemString == "" {
						foundEmpty = true
					}
				}
				if !foundEmpty {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].trusted_authorities[%d].%s contains no empty string item",
							credentialIndex,
							authorityIndex,
							params.Property,
						),
					}
				}
			}
		}
		if !isEmptyDCQLValue(responseValue) {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"wallet returned a credential for trusted authority %s containing an empty string",
					params.Property,
				),
			}
		}
	case "claim_sets":
		credentials, ok := query["credentials"].([]any)
		if !ok || !containsClaimSets(credentials) {
			return Result{
				Status:  StatusFail,
				Message: "dcql_query credentials contain no claim_sets",
			}
		}
		if isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: "wallet response contains no vp_token fclaim_sets, claim_path_member_type_error, wallet_error_expected, invalid_scope, unknown_field_stripped, vp_formats_not_supported, transaction_data_error, invalid_client, invalid_request_generic, access_denied, or jwe_enc_verified",
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: fmt.Sprintf("wallet response satisfies DCQL %s constraints", params.Mode),
	}
}

//nolint:gocyclo // Each credential-set shape has distinct conformance semantics.

// validateClaimsSubset proves both sides of a user-controlled claim selection:
// requested claims are disclosed, while explicitly unchecked paths are absent.

// validateClaimPathMemberTypeError checks that wallet rejects DCQL queries with invalid claim-path
// member types (boolean, negative integer, unsupported object types).

// validateWalletErrorExpected checks that wallet returns an expected error code.
// If expected is nil, any error code is accepted.

// validateErrorCode checks that wallet returns a specific OAuth2/OID4VP error code.

// validateUnknownFieldStripped checks that wallet processes request with unknown fields stripped.

// validateJWEEncVerified checks JWE encryption parameters in the wallet response.

// validateSessionEncryption checks session encryption evidence.

// validateEncoding checks encoding validation evidence.

// validateInteractionCompleted checks that wallet interaction completed successfully.

// validateEvidencePresent checks that evidence exists and no error occurred.
