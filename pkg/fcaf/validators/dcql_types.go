// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"encoding/json"
	"fmt"
)

func supportedJSONType(expected string) bool {
	switch expected {
	case "boolean", "string", "number", "integer", "array", "object", "null":
		return true
	default:
		return false
	}
}
func matchesJSONType(value any, expected string) bool {
	switch expected {
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			return true
		default:
			return false
		}
	case "integer":
		switch typed := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		case float32:
			return typed == float32(int64(typed))
		case float64:
			return typed == float64(int64(typed))
		default:
			return false
		}
	case "array":
		_, ok := value.([]any)
		return ok
	case "object":
		_, ok := normalizeJSONObject(value)
		return ok
	case "null":
		return value == nil
	default:
		return false
	}
}
func validateDCQLCredentialQueries(credentials []any) error {
	ids := make(map[string]struct{}, len(credentials))
	for index, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			return fmt.Errorf("credentials[%d] is not an object", index)
		}
		id, _ := credential["id"].(string)
		if !dcqlIDPattern.MatchString(id) {
			return fmt.Errorf("credentials[%d].id is not a valid DCQL identifier", index)
		}
		if _, duplicate := ids[id]; duplicate {
			return fmt.Errorf("credentials[%d].id %q is duplicated", index, id)
		}
		ids[id] = struct{}{}

		format, _ := credential["format"].(string)
		if format == "" {
			return fmt.Errorf("credentials[%d].format is missing", index)
		}
		meta, ok := normalizeJSONObject(credential["meta"])
		if !ok {
			return fmt.Errorf("credentials[%d].meta is not an object", index)
		}
		switch format {
		case "dc+sd-jwt":
			if !nonEmptyStringArray(meta["vct_values"]) {
				return fmt.Errorf(
					"credentials[%d].meta.vct_values is not a non-empty string array",
					index,
				)
			}
		case "mso_mdoc":
			docType, _ := meta["doctype_value"].(string)
			if docType == "" {
				return fmt.Errorf("credentials[%d].meta.doctype_value is missing", index)
			}
		default:
			return fmt.Errorf("credentials[%d].format %q is not supported", index, format)
		}
		if claims, exists := credential["claims"]; exists {
			items, ok := claims.([]any)
			if !ok || len(items) == 0 {
				return fmt.Errorf("credentials[%d].claims is not a non-empty array", index)
			}
			for claimIndex, rawClaim := range items {
				claim, ok := normalizeJSONObject(rawClaim)
				if !ok || !nonEmptyStringArray(claim["path"]) {
					return fmt.Errorf(
						"credentials[%d].claims[%d].path is invalid",
						index,
						claimIndex,
					)
				}
			}
		}
	}
	return nil
}
func nonEmptyStringArray(value any) bool {
	items, ok := value.([]any)
	if !ok || len(items) == 0 {
		return false
	}
	for _, item := range items {
		text, ok := item.(string)
		if !ok || text == "" {
			return false
		}
	}
	return true
}
func normalizeJSONObject(value any) (map[string]any, bool) {
	if object, ok := value.(map[string]any); ok {
		return object, true
	}
	text, ok := value.(string)
	if !ok {
		return nil, false
	}
	var object map[string]any
	if err := json.Unmarshal([]byte(text), &object); err != nil {
		return nil, false
	}
	return object, true
}
func findObjectKey(value any, key string) (any, bool) {
	object, ok := normalizeJSONObject(value)
	if !ok {
		return nil, false
	}
	if found, exists := object[key]; exists {
		return found, true
	}
	for _, child := range object {
		if found, exists := findObjectKey(child, key); exists {
			return found, true
		}
		if array, ok := child.([]any); ok {
			for _, item := range array {
				if found, exists := findObjectKey(item, key); exists {
					return found, true
				}
			}
		}
	}
	return nil, false
}
func containsClaimSets(credentials []any) bool {
	for _, rawCredential := range credentials {
		credential, ok := normalizeJSONObject(rawCredential)
		if !ok {
			continue
		}
		claimSets, ok := credential["claim_sets"].([]any)
		if ok && len(claimSets) > 0 {
			return true
		}
	}
	return false
}
func isEmptyDCQLValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return typed == ""
	case []any:
		return len(typed) == 0
	case map[string]any:
		return len(typed) == 0
	default:
		return false
	}
}
func containsCredentialFormat(formats map[string]string, expected string) bool {
	for _, actual := range formats {
		if actual == expected {
			return true
		}
	}
	return false
}
func validatedVPTokenPresentations(
	query map[string]any,
	responseValue any,
) (map[string]any, *Result) {
	credentials, ok := query["credentials"].([]any)
	if !ok || len(credentials) == 0 {
		return nil, &Result{Status: StatusFail, Message: "dcql_query does not contain credentials"}
	}
	if err := validateDCQLCredentialQueries(credentials); err != nil {
		return nil, &Result{Status: StatusFail, Message: err.Error()}
	}
	response, ok := normalizeJSONObject(responseValue)
	if !ok {
		return nil, &Result{
			Status:  StatusFail,
			Message: "wallet response vp_token is not a JSON object",
		}
	}
	for _, queryID := range queryCredentialIDs(query) {
		presentations, ok := response[queryID].([]any)
		if !ok || len(presentations) == 0 {
			return nil, &Result{
				Status:  StatusFail,
				Message: fmt.Sprintf("vp_token[%q] is not a non-empty presentation array", queryID),
			}
		}
	}
	return response, nil
}
func queryCredentialIDs(query map[string]any) []string {
	credentials, _ := query["credentials"].([]any)
	ids := make([]string, 0, len(credentials))
	for _, rawCredential := range credentials {
		credential, _ := normalizeJSONObject(rawCredential)
		id, _ := credential["id"].(string)
		ids = append(ids, id)
	}
	return ids
}
func normalizeString(v any) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
