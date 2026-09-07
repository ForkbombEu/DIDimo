// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeJSONObjectAcceptsJSONText(t *testing.T) {
	value, ok := normalizeJSONObject(`{"dcql_query":{"credentials":[]}}`)

	require.True(t, ok)
	require.Equal(t, map[string]any{
		"dcql_query": map[string]any{"credentials": []any{}},
	}, value)
}

func TestValidateDCQLCredentialQueriesRejectsDuplicateIDs(t *testing.T) {
	credential := map[string]any{
		"id":     "pid",
		"format": "dc+sd-jwt",
		"meta":   map[string]any{"vct_values": []any{"urn:eudi:pid:1"}},
	}

	err := validateDCQLCredentialQueries([]any{credential, credential})

	require.EqualError(t, err, `credentials[1].id "pid" is duplicated`)
}
