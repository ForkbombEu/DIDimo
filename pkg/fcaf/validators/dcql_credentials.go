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
		if (mode == "credential_sets_options_empty" || mode == "credential_sets_options_non_array") &&
			errorValue != invalidRequestError {
			return Result{
				Status:  StatusFail,
				Message: "wallet did not return invalid_request for an invalid credential_sets.options query",
			}
		}
		return Result{
			Status:  StatusPass,
			Message: "wallet rejected invalid credential_sets.options",
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
