// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"fmt"
	"reflect"

	"github.com/forkbombeu/credimi/pkg/fcaf/evidence"
)

func validateClaimsPresent(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
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
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d].id is not a non-empty string", index),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d].claims is not a non-empty array", index),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						index,
						claimIndex,
					),
				}
			}
			path, ok := claim["path"].([]any)
			if !ok || len(path) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not a non-empty array",
						index,
						claimIndex,
					),
				}
			}
			for pathIndex, segment := range path {
				if _, ok := segment.(string); !ok {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].claims[%d].path[%d] is not a string",
							index,
							claimIndex,
							pathIndex,
						),
					}
				}
			}
		}
		if isEmptyDCQLValue(response[id]) {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
			}
		}
	}
	return Result{Status: StatusPass, Message: "wallet processed credential queries with claims"}
}
func validateClaimsSubset(query map[string]any, responseValue any, forbiddenPaths [][]any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	if len(forbiddenPaths) == 0 {
		return Result{Status: StatusFail, Message: "claims_subset requires forbidden_paths"}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet vp_token is not an object keyed by credential query ID",
		}
	}
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].id is not a non-empty string",
					credentialIndex,
				),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		presentations, ok := response[id].([]any)
		if !ok || len(presentations) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
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
						"vp_token[%q][%d] is not a valid SD-JWT presentation: %v",
						id,
						presentationIndex,
						err,
					),
				}
			}
			for claimIndex, rawClaim := range claims {
				claim, ok := normalizeJSONObject(rawClaim)
				if !ok {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].claims[%d] is not an object",
							credentialIndex,
							claimIndex,
						),
					}
				}
				path, ok := claim["path"].([]any)
				if !ok || len(path) == 0 || !claimPathResolves(presentation.Claims, path) {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token[%q][%d] does not disclose requested claims[%d].path",
							id,
							presentationIndex,
							claimIndex,
						),
					}
				}
			}
			for pathIndex, path := range forbiddenPaths {
				if len(path) == 0 {
					return Result{
						Status:  StatusFail,
						Message: fmt.Sprintf("forbidden_paths[%d] is empty", pathIndex),
					}
				}
				if claimPathResolves(presentation.Claims, path) {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token[%q][%d] discloses unchecked forbidden_paths[%d]",
							id,
							presentationIndex,
							pathIndex,
						),
					}
				}
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet disclosed requested claims and omitted unchecked claims",
	}
}
func validateClaimSetsPreferredOption(
	query map[string]any,
	responseValue, expectedValue any,
) Result {
	expectedFloat, ok := expectedValue.(float64)
	if !ok || expectedFloat < 0 || expectedFloat != float64(int(expectedFloat)) {
		return Result{
			Status:  StatusError,
			Message: "expected_value must be a non-negative claim_sets index",
		}
	}
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 1 {
		return Result{
			Status:  StatusFail,
			Message: "claim_sets preference requires exactly one credential query",
		}
	}
	credential, ok := normalizeJSONObject(credentials[0])
	if !ok {
		return Result{Status: StatusFail, Message: "credential query is not an object"}
	}
	id, _ := credential["id"].(string)
	claims, claimsOK := credential["claims"].([]any)
	sets, setsOK := credential["claim_sets"].([]any)
	if id == "" || !claimsOK || !setsOK || len(sets) != 3 || int(expectedFloat) >= len(sets) {
		return Result{
			Status:  StatusFail,
			Message: "request does not contain three claim_sets over identified claims",
		}
	}
	paths := map[string][]any{}
	for _, rawClaim := range claims {
		claim, ok := normalizeJSONObject(rawClaim)
		if !ok {
			return Result{Status: StatusFail, Message: "claim is not an object"}
		}
		claimID, _ := claim["id"].(string)
		path, pathOK := claim["path"].([]any)
		if claimID == "" || !pathOK || len(path) == 0 {
			return Result{Status: StatusFail, Message: "claim has no id or path"}
		}
		paths[claimID] = path
	}
	expected, ok := sets[int(expectedFloat)].([]any)
	if !ok || len(expected) == 0 {
		return Result{Status: StatusFail, Message: "preferred claim_set is empty or invalid"}
	}
	expectedIDs := map[string]struct{}{}
	for _, rawID := range expected {
		claimID, ok := rawID.(string)
		if !ok || paths[claimID] == nil {
			return Result{
				Status:  StatusFail,
				Message: "preferred claim_set references an unknown claim",
			}
		}
		expectedIDs[claimID] = struct{}{}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{Status: StatusFail, Message: "wallet response contains no vp_token"}
	}
	presentations, ok := response[id].([]any)
	if !ok || len(presentations) != 1 {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return exactly one selected presentation",
		}
	}
	token, ok := presentations[0].(string)
	if !ok || token == "" {
		return Result{Status: StatusFail, Message: "selected presentation is not an SD-JWT"}
	}
	presentation, err := evidence.ParseSDJWTPresentation(token)
	if err != nil {
		return Result{
			Status:  StatusFail,
			Message: fmt.Sprintf("selected presentation is not valid SD-JWT: %v", err),
		}
	}
	for claimID, path := range paths {
		_, expected := expectedIDs[claimID]
		disclosed := claimPathResolves(presentation.Claims, path)
		if expected && !disclosed {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("selected claim_set claim %q was not disclosed", claimID),
			}
		}
		if !expected && disclosed {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"claim %q outside the selected claim_set was disclosed",
					claimID,
				),
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet disclosed only the preferred satisfiable claim_set",
	}
}
func validateClaimSetsNoMatch(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 || !containsClaimSets(credentials) {
		return Result{
			Status:  StatusFail,
			Message: "request does not contain credential claims and claim_sets",
		}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned claims although no claim_set is satisfiable",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned no claims for unsatisfiable claim_sets",
	}
}
func validateClaimSetsWithoutClaims(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "request contains no credential query"}
	}
	for _, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{Status: StatusFail, Message: "credential query is not an object"}
		}
		if _, exists := credential["claims"]; exists {
			return Result{
				Status:  StatusFail,
				Message: "invalid request unexpectedly contains claims",
			}
		}
		sets, ok := credential["claim_sets"].([]any)
		if !ok || len(sets) == 0 {
			return Result{Status: StatusFail, Message: "invalid request contains no claim_sets"}
		}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for claim_sets without claims",
		}
	}
	return Result{Status: StatusPass, Message: "wallet rejected claim_sets without claims"}
}
func validateClaimsUnion(query map[string]any, responseValue any, forbiddenPaths [][]any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) < 2 {
		return Result{
			Status:  StatusFail,
			Message: "claims_union requires at least two credential queries",
		}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet vp_token is not an object keyed by credential query ID",
		}
	}
	requested := make([][]any, 0)
	presentations := make([]*evidence.SDJWTPresentation, 0)
	for index, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", index),
			}
		}
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d].id is not a non-empty string", index),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d].claims is not a non-empty array", index),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						index,
						claimIndex,
					),
				}
			}
			path, ok := claim["path"].([]any)
			if !ok || len(path) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not a non-empty array",
						index,
						claimIndex,
					),
				}
			}
			requested = append(requested, path)
		}
		values, ok := response[id].([]any)
		if !ok || len(values) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
			}
		}
		for presentationIndex, raw := range values {
			token, ok := raw.(string)
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
			parsed, err := evidence.ParseSDJWTPresentation(token)
			if err != nil {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"vp_token[%q][%d] is not a valid SD-JWT presentation: %v",
						id,
						presentationIndex,
						err,
					),
				}
			}
			presentations = append(presentations, parsed)
		}
	}
	for pathIndex, path := range requested {
		found := false
		for _, presentation := range presentations {
			if claimPathResolves(presentation.Claims, path) {
				found = true
				break
			}
		}
		if !found {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"union response does not disclose requested claims[%d].path",
					pathIndex,
				),
			}
		}
	}
	for pathIndex, path := range forbiddenPaths {
		for _, presentation := range presentations {
			if claimPathResolves(presentation.Claims, path) {
				return Result{
					Status:  StatusFail,
					Message: fmt.Sprintf("union response discloses forbidden_paths[%d]", pathIndex),
				}
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned the union of claims requested by multiple queries",
	}
}
func validateClaimsPathNoMatch(query map[string]any, responseValue any, expectedClaimPath []any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundExpectedClaimPath := len(expectedClaimPath) == 0
	for index, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", index),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d].claims is not a non-empty array", index),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						index,
						claimIndex,
					),
				}
			}
			path, ok := claim["path"].([]any)
			if !ok || len(path) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not a non-empty array",
						index,
						claimIndex,
					),
				}
			}
			if reflect.DeepEqual(path, expectedClaimPath) {
				foundExpectedClaimPath = true
			}
		}
	}
	if !foundExpectedClaimPath {
		return Result{Status: StatusFail, Message: "dcql_query does not contain the expected unmatched claim path"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for an unmatched claim path",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned no credential for the unmatched claim path",
	}
}
func validateClaimsValuesNoMatch(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			path, pathOK := claim["path"].([]any)
			values, valuesOK := claim["values"].([]any)
			if !pathOK || len(path) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not a non-empty array",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if !valuesOK || len(values) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].values is not a non-empty array",
						credentialIndex,
						claimIndex,
					),
				}
			}
		}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for mismatched claim values",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned no credential for mismatched claim values",
	}
}
func validateMissingClaimIDWithClaimSets(
	query map[string]any,
	responseValue any,
	errorValue any,
) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundMissingID := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claimSets, ok := credential["claim_sets"].([]any)
		if !ok || len(claimSets) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claim_sets is not a non-empty array",
					credentialIndex,
				),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for _, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				continue
			}
			if _, exists := claim["id"]; !exists {
				foundMissingID = true
			}
		}
	}
	if !foundMissingID {
		return Result{Status: StatusFail, Message: "claims contain no missing id"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for claims missing id with claim_sets",
		}
	}
	if errorValue != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for claims missing id with claim_sets",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned invalid_request for claims missing id with claim_sets",
	}
}
func validateClaimsWithoutIDWithoutClaimSets(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet vp_token is not an object keyed by credential query ID",
		}
	}
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		if _, exists := credential["claim_sets"]; exists {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] contains claim_sets", credentialIndex),
			}
		}
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].id is not a non-empty string",
					credentialIndex,
				),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if _, exists := claim["id"]; exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] contains id",
						credentialIndex,
						claimIndex,
					),
				}
			}
			path, ok := claim["path"].([]any)
			if !ok || len(path) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not a non-empty array",
						credentialIndex,
						claimIndex,
					),
				}
			}
			for pathIndex, segment := range path {
				if value, ok := segment.(string); !ok || value == "" {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"credentials[%d].claims[%d].path[%d] is not a non-empty string",
							credentialIndex,
							claimIndex,
							pathIndex,
						),
					}
				}
			}
		}
		if isEmptyDCQLValue(response[id]) {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
			}
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet matched claims without ids when claim_sets was absent",
	}
}
func validateDuplicateClaimIDs(query map[string]any, responseValue any, errorValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundDuplicate := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		seen := make(map[string]struct{}, len(claims))
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			id, ok := claim["id"].(string)
			if !ok || id == "" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].id is not a non-empty string",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if _, exists := seen[id]; exists {
				foundDuplicate = true
			}
			seen[id] = struct{}{}
		}
	}
	if !foundDuplicate {
		return Result{
			Status:  StatusFail,
			Message: "no credential claims array contains a duplicate id",
		}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for duplicate claim ids",
		}
	}
	if errorText, _ := errorValue.(string); errorText != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for duplicate claim ids",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet rejected duplicate claim ids with invalid_request",
	}
}
func validateEmptyClaimID(query map[string]any, responseValue any, errorValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundEmpty := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			idValue, exists := claim["id"]
			if !exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].id is missing",
						credentialIndex,
						claimIndex,
					),
				}
			}
			id, ok := idValue.(string)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].id is not a string",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if id == "" {
				foundEmpty = true
			}
		}
	}
	if !foundEmpty {
		return Result{Status: StatusFail, Message: "no claim id is empty"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for an empty claim id",
		}
	}
	if errorText, _ := errorValue.(string); errorText != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for an empty claim id",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet rejected an empty claim id with invalid_request",
	}
}
func validateInvalidClaimIDCharacters(
	query map[string]any,
	responseValue any,
	errorValue any,
) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundInvalid := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			idValue, exists := claim["id"]
			if !exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].id is missing",
						credentialIndex,
						claimIndex,
					),
				}
			}
			id, ok := idValue.(string)
			if !ok || id == "" {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].id is not a non-empty string",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if !dcqlIDPattern.MatchString(id) {
				foundInvalid = true
			}
		}
	}
	if !foundInvalid {
		return Result{Status: StatusFail, Message: "no claim id contains a forbidden character"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for a malformed claim id",
		}
	}
	if errorText, _ := errorValue.(string); errorText != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for a malformed claim id",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet rejected a malformed claim id with invalid_request",
	}
}
func validateMissingClaimPath(query map[string]any, responseValue any, errorValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundMissing := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if _, exists := claim["path"]; !exists {
				foundMissing = true
			}
		}
	}
	if !foundMissing {
		return Result{Status: StatusFail, Message: "no claim is missing path"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for a claim missing path",
		}
	}
	if errorText, _ := errorValue.(string); errorText != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for a claim missing path",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet rejected a claim missing path with invalid_request",
	}
}
func validateEmptyClaimPath(query map[string]any, responseValue any, errorValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundEmpty := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			pathValue, exists := claim["path"]
			if !exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is missing",
						credentialIndex,
						claimIndex,
					),
				}
			}
			path, ok := pathValue.([]any)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not an array",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if len(path) == 0 {
				foundEmpty = true
			}
		}
	}
	if !foundEmpty {
		return Result{Status: StatusFail, Message: "no claim path is empty"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for an empty claim path",
		}
	}
	if errorText, _ := errorValue.(string); errorText != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for an empty claim path",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet rejected an empty claim path with invalid_request",
	}
}
func validateNonArrayClaimPath(query map[string]any, responseValue any, errorValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	foundNonArray := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			pathValue, exists := claim["path"]
			if !exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is missing",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if _, ok := pathValue.([]any); !ok {
				foundNonArray = true
			}
		}
	}
	if !foundNonArray {
		return Result{Status: StatusFail, Message: "no claim path has a non-array value"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a credential for a non-array claim path",
		}
	}
	if errorText, _ := errorValue.(string); errorText != invalidRequestError {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return invalid_request for a non-array claim path",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet rejected a non-array claim path with invalid_request",
	}
}
func validateAllowedClaimPathComponents(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet vp_token is not an object keyed by credential query ID",
		}
	}
	seenString := false
	seenNull := false
	seenNonNegativeInteger := false
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].id is not a non-empty string",
					credentialIndex,
				),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		paths := make([][]any, 0, len(claims))
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			path, ok := claim["path"].([]any)
			if !ok || len(path) == 0 {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is not a non-empty array",
						credentialIndex,
						claimIndex,
					),
				}
			}
			paths = append(paths, path)
			for componentIndex, component := range path {
				switch typed := component.(type) {
				case string:
					if typed == "" {
						return Result{
							Status: StatusFail,
							Message: fmt.Sprintf(
								"credentials[%d].claims[%d].path[%d] is an empty string",
								credentialIndex,
								claimIndex,
								componentIndex,
							),
						}
					}
					seenString = true
				case nil:
					seenNull = true
				default:
					if !isNonNegativeInteger(component) {
						return Result{
							Status: StatusFail,
							Message: fmt.Sprintf(
								"credentials[%d].claims[%d].path[%d] is not a string, null, or non-negative integer",
								credentialIndex,
								claimIndex,
								componentIndex,
							),
						}
					}
					seenNonNegativeInteger = true
				}
			}
		}
		presentations, ok := response[id].([]any)
		if !ok || len(presentations) == 0 {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
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
						"vp_token[%q][%d] is not a valid SD-JWT presentation: %v",
						id,
						presentationIndex,
						err,
					),
				}
			}
			for pathIndex, path := range paths {
				if !claimPathResolves(presentation.Claims, path) {
					return Result{
						Status: StatusFail,
						Message: fmt.Sprintf(
							"vp_token[%q][%d] does not disclose a value resolved by claims[%d].path",
							id,
							presentationIndex,
							pathIndex,
						),
					}
				}
			}
		}
	}
	if !seenString || !seenNull || !seenNonNegativeInteger {
		return Result{
			Status:  StatusFail,
			Message: "claim paths do not cover string, null, and non-negative integer components",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet resolved claim paths with all allowed component types",
	}
}
func claimPathResolves(root any, path []any) bool {
	values := []any{root}
	for _, component := range path {
		next := make([]any, 0)
		for _, value := range values {
			switch typed := component.(type) {
			case string:
				object, ok := value.(map[string]any)
				if !ok {
					continue
				}
				if resolved, exists := object[typed]; exists {
					next = append(next, resolved)
				}
			case nil:
				array, ok := value.([]any)
				if ok {
					next = append(next, array...)
				}
			default:
				array, ok := value.([]any)
				if !ok {
					continue
				}
				index, ok := claimPathArrayIndex(typed, len(array))
				if ok {
					next = append(next, array[index])
				}
			}
		}
		if len(next) == 0 {
			return false
		}
		values = next
	}
	return len(values) > 0
}
func claimPathArrayIndex(value any, length int) (int, bool) {
	if !isNonNegativeInteger(value) {
		return 0, false
	}
	var index uint64
	switch typed := value.(type) {
	case int:
		index = uint64(typed)
	case int8:
		index = uint64(typed)
	case int16:
		index = uint64(typed)
	case int32:
		index = uint64(typed)
	case int64:
		index = uint64(typed)
	case uint:
		index = uint64(typed)
	case uint8:
		index = uint64(typed)
	case uint16:
		index = uint64(typed)
	case uint32:
		index = uint64(typed)
	case uint64:
		index = typed
	case float32:
		index = uint64(typed)
	case float64:
		index = uint64(typed)
	default:
		return 0, false
	}
	if index >= uint64(length) {
		return 0, false
	}
	return int(index), true
}
func isNonNegativeInteger(value any) bool {
	switch typed := value.(type) {
	case int:
		return typed >= 0
	case int8:
		return typed >= 0
	case int16:
		return typed >= 0
	case int32:
		return typed >= 0
	case int64:
		return typed >= 0
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32:
		return typed >= 0 && typed == float32(int64(typed))
	case float64:
		return typed >= 0 && typed == float64(int64(typed))
	default:
		return false
	}
}
func validateClaimsWithoutValues(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return Result{
			Status:  StatusFail,
			Message: "wallet vp_token is not an object keyed by credential query ID",
		}
	}
	for credentialIndex, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", credentialIndex),
			}
		}
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].id is not a non-empty string",
					credentialIndex,
				),
			}
		}
		claims, ok := credential["claims"].([]any)
		if !ok || len(claims) == 0 {
			return Result{
				Status: StatusFail,
				Message: fmt.Sprintf(
					"credentials[%d].claims is not a non-empty array",
					credentialIndex,
				),
			}
		}
		for claimIndex, rawClaim := range claims {
			claim, ok := normalizeJSONObject(rawClaim)
			if !ok {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] is not an object",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if _, exists := claim["values"]; exists {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d] contains values",
						credentialIndex,
						claimIndex,
					),
				}
			}
			if !nonEmptyStringArray(claim["path"]) {
				return Result{
					Status: StatusFail,
					Message: fmt.Sprintf(
						"credentials[%d].claims[%d].path is invalid",
						credentialIndex,
						claimIndex,
					),
				}
			}
		}
		if isEmptyDCQLValue(response[id]) {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token has no presentation for credential query %q", id),
			}
		}
	}
	return Result{Status: StatusPass, Message: "wallet matched claims without values"}
}
