// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"fmt"
	"reflect"

	"github.com/forkbombeu/credimi/pkg/fcaf/evidence"
)

func validateCredentialFormatPresentation(
	query map[string]any,
	responseValue any,
	expectedFormat string,
) Result {
	if expectedFormat != "mso_mdoc" && expectedFormat != "dc+sd-jwt" {
		return Result{
			Status:  StatusError,
			Message: "expected_format must be mso_mdoc or dc+sd-jwt",
		}
	}

	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 1 {
		return Result{
			Status:  StatusFail,
			Message: "dcql_query must contain exactly one credential query",
		}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}
	credential, _ := normalizeJSONObject(credentials[0])
	queryID, _ := credential["id"].(string)
	actualFormat, _ := credential["format"].(string)
	if actualFormat != expectedFormat {
		return Result{
			Status: StatusFail,
			Message: fmt.Sprintf(
				"credential query format is %q, expected %q",
				actualFormat,
				expectedFormat,
			),
		}
	}

	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet response vp_token is not an object",
		}
	}
	presentations, ok := response[queryID].([]any)
	if !ok || len(presentations) == 0 {
		return Result{
			Status:  StatusFail,
			Message: fmt.Sprintf("vp_token has no presentation for query %q", queryID),
		}
	}
	return Result{
		Status: StatusPass,
		Message: fmt.Sprintf(
			"wallet returned %s presentation for query %q",
			expectedFormat,
			queryID,
		),
	}
}
func validateMDocClaimPathPresentation(
	query map[string]any,
	responseValue any,
	expectedPath []any,
) Result {
	if len(expectedPath) != 2 {
		return Result{Status: StatusError, Message: "expected_claim_path must contain namespace and element"}
	}
	namespace, namespaceOK := expectedPath[0].(string)
	element, elementOK := expectedPath[1].(string)
	if !namespaceOK || namespace == "" || !elementOK || element == "" {
		return Result{Status: StatusError, Message: "expected_claim_path namespace and element must be non-empty strings"}
	}

	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 1 {
		return Result{Status: StatusFail, Message: "dcql_query must contain exactly one credential query"}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}
	credential, _ := normalizeJSONObject(credentials[0])
	if credential["format"] != "mso_mdoc" {
		return Result{Status: StatusFail, Message: "credential query format must be mso_mdoc"}
	}
	queryID, _ := credential["id"].(string)
	claims, ok := credential["claims"].([]any)
	if !ok {
		return Result{Status: StatusFail, Message: "mdoc credential query has no claims"}
	}
	pathFound := false
	for _, rawClaim := range claims {
		claim, ok := normalizeJSONObject(rawClaim)
		if !ok {
			continue
		}
		path, ok := claim["path"].([]any)
		if ok && reflect.DeepEqual(path, expectedPath) {
			pathFound = true
			break
		}
	}
	if !pathFound {
		return Result{Status: StatusFail, Message: "mdoc credential query does not contain the expected namespace and element path"}
	}

	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{Status: StatusFail, Message: "wallet response vp_token is not an object"}
	}
	presentations, ok := response[queryID].([]any)
	if !ok || len(presentations) == 0 {
		return Result{Status: StatusFail, Message: fmt.Sprintf("vp_token has no presentation for query %q", queryID)}
	}
	for index, rawPresentation := range presentations {
		token, ok := rawPresentation.(string)
		if !ok || token == "" {
			return Result{Status: StatusFail, Message: fmt.Sprintf("vp_token[%q][%d] is not an mdoc presentation", queryID, index)}
		}
		presentation, err := evidence.ParseMDocPresentation(token)
		if err != nil {
			return Result{Status: StatusFail, Message: fmt.Sprintf("vp_token[%q][%d] is not a valid mdoc presentation: %v", queryID, index, err)}
		}
		if _, found := presentation.Element(namespace, element); !found {
			return Result{Status: StatusFail, Message: fmt.Sprintf("vp_token[%q][%d] does not contain mdoc element %q in namespace %q", queryID, index, element, namespace)}
		}
	}
	return Result{Status: StatusPass, Message: "wallet returned the requested mdoc namespace and element in a CBOR presentation"}
}
func validateMDocClaimPathNoMatch(query map[string]any, expectedPath []any) Result {
	if len(expectedPath) != 2 {
		return Result{Status: StatusError, Message: "expected_claim_path must contain namespace and element"}
	}
	namespace, namespaceOK := expectedPath[0].(string)
	element, elementOK := expectedPath[1].(string)
	if !namespaceOK || namespace == "" || !elementOK || element == "" {
		return Result{Status: StatusError, Message: "expected_claim_path namespace and element must be non-empty strings"}
	}

	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 1 {
		return Result{Status: StatusFail, Message: "dcql_query must contain exactly one credential query"}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}
	credential, _ := normalizeJSONObject(credentials[0])
	if credential["format"] != "mso_mdoc" {
		return Result{Status: StatusFail, Message: "credential query format must be mso_mdoc"}
	}
	claims, ok := credential["claims"].([]any)
	if !ok {
		return Result{Status: StatusFail, Message: "mdoc credential query has no claims"}
	}
	for _, rawClaim := range claims {
		claim, ok := normalizeJSONObject(rawClaim)
		if !ok {
			continue
		}
		path, ok := claim["path"].([]any)
		if ok && reflect.DeepEqual(path, expectedPath) {
			return Result{Status: StatusPass, Message: "mdoc credential query contains the expected absent namespace path"}
		}
	}
	return Result{Status: StatusFail, Message: "mdoc credential query does not contain the expected absent namespace path"}
}
func validateVPTokenSignedPresentation(
	root map[string]any,
	query map[string]any,
	responseValue any,
) Result {
	responseTypes := collectObjectFieldValues(root, "response_type")
	if len(responseTypes) != 1 || responseTypes[0] != "vp_token" {
		return Result{
			Status:  StatusFail,
			Message: "captured authorization request response_type must equal vp_token",
		}
	}
	response, result := validatedVPTokenPresentations(query, responseValue)
	if result != nil {
		return *result
	}
	for queryID, rawPresentations := range response {
		presentations, ok := rawPresentations.([]any)
		if !ok {
			continue
		}
		for index, rawPresentation := range presentations {
			token, ok := rawPresentation.(string)
			if !ok || token == "" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not an SD-JWT presentation",
						queryID,
						index,
					),
				}
			}
			parsed, err := evidence.ParseSDJWTPresentation(token)
			if err != nil {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not a valid SD-JWT presentation: %v",
						queryID,
						index,
						err,
					),
				}
			}
			algorithm, _ := parsed.ProtectedHeaders["alg"].(string)
			if algorithm == "" || algorithm == "none" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not a signed JWT presentation",
						queryID,
						index,
					),
				}
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "authorization response contains signed SD-JWT vp_token presentations",
	}
}
func validateVPTokenJSONObject(responseValue any) Result {
	if _, ok := normalizeJSONObject(responseValue); !ok {
		return Result{Status: StatusFail, Message: "wallet response vp_token is not a JSON object"}
	}
	return Result{Status: StatusPass, Message: "wallet response vp_token is a JSON object"}
}
func validateVPTokenQueryIDs(query map[string]any, responseValue any) Result {
	response, result := validatedVPTokenPresentations(query, responseValue)
	if result != nil {
		return *result
	}
	if len(response) != len(queryCredentialIDs(query)) {
		return Result{
			Status:  StatusFail,
			Message: "vp_token keys do not exactly match credential query IDs",
		}
	}
	for _, queryID := range queryCredentialIDs(query) {
		if _, found := response[queryID]; !found {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no key for credential query %q", queryID),
			}
		}
	}
	return Result{Status: StatusPass, Message: "vp_token keys exactly match credential query IDs"}
}
func validateVPTokenPresentationArrays(query map[string]any, responseValue any) Result {
	_, result := validatedVPTokenPresentations(query, responseValue)
	if result != nil {
		return *result
	}
	return Result{
		Status:  StatusPass,
		Message: "vp_token entries are non-empty presentation arrays for each credential query",
	}
}
func validateTwoDistinctFormatPresentations(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 2 {
		return Result{
			Status:  StatusFail,
			Message: "dcql_query must contain exactly two credential queries",
		}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}

	formats := make(map[string]string, len(credentials))
	for _, rawCredential := range credentials {
		credential, _ := normalizeJSONObject(rawCredential)
		id, _ := credential["id"].(string)
		formats[id], _ = credential["format"].(string)
	}
	if len(formats) != 2 || !containsCredentialFormat(formats, "dc+sd-jwt") ||
		!containsCredentialFormat(formats, "mso_mdoc") {
		return Result{
			Status:  StatusFail,
			Message: "credential queries must contain one dc+sd-jwt and one mso_mdoc credential",
		}
	}

	response, result := validatedVPTokenPresentations(query, responseValue)
	if result != nil {
		return *result
	}
	if len(response) != len(formats) {
		return Result{
			Status:  StatusFail,
			Message: "vp_token keys do not exactly match the two credential query IDs",
		}
	}
	for queryID, format := range formats {
		presentations := response[queryID].([]any)
		for index, rawPresentation := range presentations {
			token, ok := rawPresentation.(string)
			if !ok || token == "" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not a presentation string",
						queryID,
						index,
					),
				}
			}
			switch format {
			case "dc+sd-jwt":
				parsed, err := evidence.ParseSDJWTPresentation(token)
				if err != nil {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token[%q][%d] is not a valid SD-JWT presentation: %v",
							queryID,
							index,
							err,
						),
					}
				}
				algorithm, _ := parsed.ProtectedHeaders["alg"].(string)
				if algorithm == "" || algorithm == "none" {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token[%q][%d] is not a signed SD-JWT presentation",
							queryID,
							index,
						),
					}
				}
			case "mso_mdoc":
				if _, err := evidence.ParseMDocPresentation(token); err != nil {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token[%q][%d] is not a valid mDoc presentation: %v",
							queryID,
							index,
							err,
						),
					}
				}
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "vp_token exactly maps distinct SD-JWT and mDoc credential queries to valid presentations",
	}
}
func validateClaimPathMemberTypeError(responseValue, errorValue any) Result {
	if errStr := normalizeString(errorValue); errStr != "" {
		if errStr == invalidRequestError {
			return Result{
				Status: StatusPass,
				Message: fmt.Sprintf(
					"wallet returned %s for invalid claim-path member type",
					errStr,
				),
			}
		}
		return Result{
			Status: StatusPass,
			Message: fmt.Sprintf(
				"wallet returned error %s for invalid claim-path member type",
				errStr,
			),
		}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned vp_token for query with invalid claim-path member type",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet did not return vp_token for invalid claim-path member type",
	}
}
func validateWalletErrorExpected(responseValue, errorValue, expected any) Result {
	if errStr := normalizeString(errorValue); errStr != "" {
		if expected != nil {
			expectedStr, ok := expected.(string)
			if ok && errStr == expectedStr {
				return Result{
					Status:  StatusPass,
					Message: fmt.Sprintf("wallet returned expected error %s", errStr),
				}
			}
			if ok {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("wallet returned %s, expected %s", errStr, expectedStr),
				}
			}
		}
		return Result{Status: StatusPass, Message: fmt.Sprintf("wallet returned error %s", errStr)}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{Status: StatusFail, Message: "wallet returned vp_token, expected error"}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet did not return vp_token (expected error case)",
	}
}
func validateWalletErrorRequired(responseValue, errorValue any) Result {
	if !isEmptyDCQLValue(responseValue) {
		return Result{Status: StatusFail, Message: "wallet returned vp_token, expected error"}
	}
	if errStr := normalizeString(errorValue); errStr != "" {
		return Result{Status: StatusPass, Message: fmt.Sprintf("wallet returned error %s", errStr)}
	}
	return Result{Status: StatusFail, Message: "wallet returned no vp_token but no error, expected error"}
}
func validateErrorCode(responseValue, errorValue any, expectedCode string) Result {
	if errStr := normalizeString(errorValue); errStr != "" {
		if errStr == expectedCode {
			return Result{
				Status:  StatusPass,
				Message: fmt.Sprintf("wallet returned expected error %s", expectedCode),
			}
		}
		return Result{
			Status:  StatusFail,
			Message: fmt.Sprintf("wallet returned error %s, expected %s", errStr, expectedCode),
		}
	}
	if isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusPass,
			Message: fmt.Sprintf("wallet did not return vp_token for %s case", expectedCode),
		}
	}
	return Result{
		Status:  StatusFail,
		Message: fmt.Sprintf("wallet returned vp_token, expected error %s", expectedCode),
	}
}
func validateRequiredErrorCode(responseValue, errorValue any, expectedCode string) Result {
	errStr := normalizeString(errorValue)
	if errStr == expectedCode && isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusPass,
			Message: fmt.Sprintf("wallet returned expected error %s", expectedCode),
		}
	}
	if errStr == "" {
		return Result{
			Status:  StatusFail,
			Message: fmt.Sprintf("wallet did not return required error %s", expectedCode),
		}
	}
	return Result{
		Status:  StatusFail,
		Message: fmt.Sprintf("wallet returned error %s, expected %s", errStr, expectedCode),
	}
}
func validateUnknownFieldStripped(query, responseValue any, property string) Result {
	if property == "" {
		return Result{
			Status:  StatusError,
			Message: "property is required for unknown_field_stripped mode",
		}
	}
	queryObject, ok := normalizeJSONObject(query)
	if !ok {
		return Result{Status: StatusFail, Message: "dcql_query is not an object"}
	}
	if _, exists := queryObject[property]; !exists {
		credentials, _ := queryObject["credentials"].([]any)
		found := false
		for _, rawCredential := range credentials {
			credential, ok := normalizeJSONObject(rawCredential)
			if ok {
				_, found = credential[property]
			}
			if found {
				break
			}
		}
		if !found {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("dcql_query does not contain unknown property %q", property),
			}
		}
	}
	if errStr := normalizeString(responseValue); errStr != "" {
		return Result{
			Status: StatusFail,
			Message: fmt.Sprintf(
				"wallet returned error %s for unknown field (should have been stripped)",
				errStr,
			),
		}
	}
	if isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned no vp_token for request with unknown fields",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet accepted request with unknown fields stripped",
	}
}
func validateJWEEncVerified(responseValue any) Result {
	resp, _ := normalizeJSONObject(responseValue)
	if len(resp) == 0 {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned empty vp_token, cannot verify JWE enc",
		}
	}
	if _, exists := resp["response"]; exists {
		return Result{Status: StatusPass, Message: "wallet response contains JWE response evidence"}
	}
	for _, v := range resp {
		if str, ok := v.(string); ok && len(str) > 0 {
			parts := 0
			for _, c := range str {
				if c == '.' {
					parts++
				}
			}
			if parts == 4 {
				return Result{Status: StatusPass, Message: "wallet response contains compact JWE"}
			}
		}
	}
	return Result{Status: StatusPass, Message: "wallet returned response evidence"}
}
func validateSessionEncryption(responseValue any) Result {
	if isEmptyDCQLValue(responseValue) {
		return Result{Status: StatusFail, Message: "wallet returned no session encryption evidence"}
	}
	return Result{Status: StatusPass, Message: "wallet returned session encryption evidence"}
}
func validateEncoding(responseValue any) Result {
	if isEmptyDCQLValue(responseValue) {
		return Result{Status: StatusFail, Message: "wallet returned no encoding evidence"}
	}
	return Result{Status: StatusPass, Message: "wallet returned encoding evidence"}
}
func validateInteractionCompleted(responseValue any) Result {
	if isEmptyDCQLValue(responseValue) {
		return Result{Status: StatusFail, Message: "wallet interaction did not complete"}
	}
	return Result{Status: StatusPass, Message: "wallet interaction completed successfully"}
}
func validateEvidencePresent(responseValue, errorValue any) Result {
	if errStr := normalizeString(errorValue); errStr != "" {
		return Result{Status: StatusFail, Message: fmt.Sprintf("wallet returned error: %s", errStr)}
	}
	if isEmptyDCQLValue(responseValue) {
		return Result{Status: StatusFail, Message: "wallet returned no evidence"}
	}
	return Result{Status: StatusPass, Message: "wallet evidence present"}
}
