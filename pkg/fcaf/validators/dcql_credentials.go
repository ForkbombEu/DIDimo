// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import "fmt"

func validateCredentialSetsOptions(
	query map[string]any,
	responseValue, errorValue any,
	mode string,
) Result {
	credentials, ok := query["credentials"].([]any)
	sets, setsOK := query["credential_sets"].([]any)
	if !ok || len(credentials) == 0 || !setsOK || len(sets) == 0 {
		return Result{
			Status:  StatusFail,
			Message: "dcql_query must contain credentials and credential_sets",
		}
	}
	ids := make(map[string]struct{}, len(credentials))
	for _, raw := range credentials {
		credential, ok := normalizeJSONObject(raw)
		if !ok {
			return Result{Status: StatusFail, Message: "dcql credential is not an object"}
		}
		id, ok := credential["id"].(string)
		if !ok || id == "" {
			return Result{Status: StatusFail, Message: "dcql credential id is invalid"}
		}
		ids[id] = struct{}{}
	}
	invalid := false
	for _, raw := range sets {
		set, ok := normalizeJSONObject(raw)
		if !ok {
			invalid = true
			continue
		}
		options, exists := set["options"]
		if !exists {
			invalid = true
			continue
		}
		groups, ok := options.([]any)
		if mode == "credential_sets_options_non_array" {
			if ok {
				return Result{Status: StatusFail, Message: "credential_sets.options is an array"}
			}
			invalid = true
			continue
		}
		if !ok || len(groups) == 0 {
			invalid = true
			continue
		}
		for _, rawGroup := range groups {
			group, ok := rawGroup.([]any)
			if !ok || len(group) == 0 {
				invalid = true
				continue
			}
			for _, rawID := range group {
				id, ok := rawID.(string)
				if !ok {
					invalid = true
					continue
				}
				if _, found := ids[id]; !found {
					invalid = true
				}
			}
		}
	}
	if mode == "credential_sets_options_valid_references" && invalid {
		return Result{
			Status:  StatusFail,
			Message: "credential_sets.options contains invalid references",
		}
	}
	if mode == "credential_sets_options_invalid_references" && !invalid {
		return Result{
			Status:  StatusFail,
			Message: "credential_sets.options contains no invalid references",
		}
	}
	if mode == "credential_sets_options_empty" && !invalid {
		return Result{Status: StatusFail, Message: "credential_sets.options is non-empty"}
	}
	if mode == "credential_sets_options_non_array" || mode == "credential_sets_options_empty" ||
		mode == "credential_sets_options_invalid_references" {
		if !isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: "wallet returned a vp_token for an invalid credential_sets.options query",
			}
		}
		if mode == "credential_sets_options_invalid_references" && normalizeString(errorValue) == "" {
			return Result{
				Status:  StatusFail,
				Message: "wallet did not return a privacy-preserving error for invalid credential_sets.options references",
			}
		}
		if (mode == "credential_sets_options_empty" || mode == "credential_sets_options_non_array") && errorValue != invalidRequestError {
			return Result{
				Status:  StatusFail,
				Message: "wallet did not return invalid_request for an invalid credential_sets.options query",
			}
		}
		return Result{
			Status:  StatusPass,
			Message: "wallet returned a privacy-preserving error for invalid credential_sets.options",
		}
	}
	if isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned no vp_token for valid credential_sets.options references",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet processed valid credential_sets.options references",
	}
}
func validateCredentialSetsRequired(query map[string]any, responseValue any, mode string) Result {
	sets, ok := query["credential_sets"].([]any)
	if !ok || len(sets) == 0 {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credential_sets"}
	}
	response, responseOK := normalizeJSONObject(responseValue)
	for index, rawSet := range sets {
		set, ok := normalizeJSONObject(rawSet)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credential_sets[%d] is not an object", index),
			}
		}
		required, exists := set["required"]
		if mode == "credential_sets_required_true_match" && (!exists || required != true) {
			return Result{Status: StatusFail, Message: "required is not true"}
		}
		if mode == "credential_sets_required_true_no_match" && (!exists || required != true) {
			return Result{Status: StatusFail, Message: "required is not true"}
		}
		if mode == "credential_sets_required_omitted" && exists {
			return Result{Status: StatusFail, Message: "required is present"}
		}
		if mode == "credential_sets_required_false_with_match" && required != false {
			return Result{Status: StatusFail, Message: "required is not false"}
		}
	}
	if mode == "credential_sets_required_true_match" ||
		mode == "credential_sets_required_omitted" ||
		mode == "credential_sets_required_false_with_match" {
		if !responseOK || isEmptyDCQLValue(response) {
			return Result{
				Status:  StatusFail,
				Message: "wallet returned no vp_token for a satisfiable credential set",
			}
		}
		return Result{Status: StatusPass, Message: "wallet presented the credential set"}
	}
	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a presentation for a missing required credential set",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet stopped without presenting a missing required credential set",
	}
}
func validateRequiredCredentialsNoPartialPresentation(
	query map[string]any,
	responseValue any,
	errorValue any,
) Result {
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

	credentialIDs := make([]string, len(credentials))
	for index, rawCredential := range credentials {
		credential, _ := normalizeJSONObject(rawCredential)
		credentialIDs[index], _ = credential["id"].(string)
	}

	sets, ok := query["credential_sets"].([]any)
	if !ok || len(sets) != 1 || !credentialSetHasExactOption(sets[0], credentialIDs...) {
		return Result{
			Status:  StatusFail,
			Message: "credential_sets must contain one option requiring both credential queries",
		}
	}
	set, _ := normalizeJSONObject(sets[0])
	if set["required"] != true {
		return Result{Status: StatusFail, Message: "credential_sets option is not required"}
	}

	if !isEmptyDCQLValue(responseValue) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned a partial presentation for an unsatisfied required credential set",
		}
	}
	if errorMessage, ok := errorValue.(string); !ok || errorMessage == "" {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return a requirement-not-met error for an unsatisfied required credential set",
		}
	}

	return Result{
		Status:  StatusPass,
		Message: "wallet returned a requirement-not-met error without presenting a partial credential set",
	}
}
func validateCredentialSetInteraction(query map[string]any, responseValue any, mode string) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 2 {
		return Result{
			Status:  StatusFail,
			Message: "dcql_query must contain exactly two credential queries",
		}
	}
	ids := make([]string, 2)
	for index, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] is not an object", index),
			}
		}
		id, _ := credential["id"].(string)
		if id == "" {
			return Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("credentials[%d] has no id", index),
			}
		}
		ids[index] = id
	}
	sets, ok := query["credential_sets"].([]any)
	if !ok {
		return Result{Status: StatusFail, Message: "dcql_query does not contain credential_sets"}
	}
	if mode == "credential_sets_required_optional" {
		if len(sets) != 2 || !credentialSetHasExactOption(sets[0], ids[0]) ||
			!credentialSetHasExactOption(sets[1], ids[1]) {
			return Result{
				Status:  StatusFail,
				Message: "credential_sets must contain separate options for the available required and unavailable optional queries",
			}
		}
		first, _ := normalizeJSONObject(sets[0])
		second, _ := normalizeJSONObject(sets[1])
		if first["required"] != true || second["required"] != false {
			return Result{
				Status:  StatusFail,
				Message: "credential_sets must mark the available option required and unavailable option optional",
			}
		}
	} else {
		if len(sets) != 1 {
			return Result{
				Status:  StatusFail,
				Message: "credential_sets must contain exactly one set",
			}
		}
		expected := []string{ids[0]}
		if mode == "credential_sets_combined_option_no_match" {
			expected = append(expected, ids[1])
		}
		if !credentialSetHasExactOption(sets[0], expected...) {
			return Result{
				Status:  StatusFail,
				Message: "credential_sets option does not match the requested credential queries",
			}
		}
	}
	response, responseOK := normalizeJSONObject(responseValue)
	if mode == "credential_sets_combined_option_no_match" {
		if !isEmptyDCQLValue(responseValue) {
			return Result{
				Status:  StatusFail,
				Message: "wallet returned a presentation for an unsatisfied combined credential-set option",
			}
		}
		return Result{
			Status:  StatusPass,
			Message: "wallet did not return an unsatisfied combined credential-set option",
		}
	}
	if !responseOK || isEmptyDCQLValue(response[ids[0]]) {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return the available credential-set option",
		}
	}
	if !isEmptyDCQLValue(response[ids[1]]) {
		return Result{
			Status:  StatusFail,
			Message: "wallet returned the unavailable credential query",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned only the available credential-set option",
	}
}
func validateOptionalCredentialSetNoMatch(query map[string]any, responseValue any) Result {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) != 1 {
		return Result{
			Status:  StatusFail,
			Message: "dcql_query must contain exactly one optional credential query",
		}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return Result{Status: StatusFail, Message: err.Error()}
	}
	credential, _ := normalizeJSONObject(credentials[0])
	credentialID, _ := credential["id"].(string)
	if credentialID == "" {
		return Result{Status: StatusFail, Message: "optional credential query has no id"}
	}

	sets, ok := query["credential_sets"].([]any)
	if !ok || len(sets) != 1 || !credentialSetHasExactOption(sets[0], credentialID) {
		return Result{
			Status:  StatusFail,
			Message: "dcql_query must contain one credential set for the optional credential query",
		}
	}
	set, _ := normalizeJSONObject(sets[0])
	if set["required"] != false {
		return Result{
			Status:  StatusFail,
			Message: "credential set must mark the unavailable credential query optional",
		}
	}

	vpToken, ok := normalizeJSONObject(responseValue)
	if !ok || len(vpToken) != 0 {
		return Result{
			Status:  StatusFail,
			Message: "wallet did not return an empty vp_token for the unavailable optional credential query",
		}
	}
	return Result{
		Status:  StatusPass,
		Message: "wallet returned an empty vp_token for the unavailable optional credential query",
	}
}
func credentialSetHasExactOption(rawSet any, expected ...string) bool {
	set, ok := normalizeJSONObject(rawSet)
	if !ok {
		return false
	}
	options, ok := set["options"].([]any)
	if !ok || len(options) != 1 {
		return false
	}
	option, ok := options[0].([]any)
	if !ok || len(option) != len(expected) {
		return false
	}
	for index, id := range expected {
		if option[index] != id {
			return false
		}
	}
	return true
}
