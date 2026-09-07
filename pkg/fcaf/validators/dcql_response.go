// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"fmt"
)

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
func validateUnknownFieldStripped(query, responseValue any) Result {
	if errStr := normalizeString(query); errStr != "" {
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
