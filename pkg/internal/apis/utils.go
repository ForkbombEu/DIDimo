// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func decodeJSON(r io.Reader, dest interface{}) error {
	if err := json.NewDecoder(r).Decode(dest); err != nil {
		return apis.NewBadRequestError("invalid JSON input", err)
	}
	return nil
}


func normalizeProtocolAndAuthor(protocol, author string) (string, string) {
	switch protocol {
	case "openid4vp_wallet":
		protocol = "OpenID4VP_Wallet"
	case "openid4vci_wallet":
		protocol = "OpenID4VCI_Wallet"
	}
	if author == "openid_foundation" {
		author = "OpenID_foundation"
	}
	return protocol, author
}

func readTemplateFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to open template file: %w", err)
	}
	return string(data), nil
}

func getUserNamespace(app core.App, userId string) (string, error) {
	orgAuthCollection, err := app.FindCollectionByNameOrId("orgAuthorizations")
	if err != nil {
		return "", apis.NewInternalServerError("failed to find orgAuthorizations collection", err)
	}

	authOrgRecords, err := app.FindRecordsByFilter(orgAuthCollection.Id, "user={:user}", "", 0, 0, dbx.Params{"user": userId})
	if err != nil {
		return "", apis.NewInternalServerError("failed to find orgAuthorizations records", err)
	}
	if len(authOrgRecords) == 0 {
		return "", apis.NewInternalServerError("user is not authorized to access any organization", nil)
	}

	ownerRoleRecord, err := app.FindFirstRecordByFilter("orgRoles", "name='owner'")
	if err != nil {
		return "", apis.NewInternalServerError("failed to find owner role", err)
	}

	if len(authOrgRecords) > 1 {
		for _, record := range authOrgRecords {
			if record.GetString("role") == ownerRoleRecord.Id {
				return record.GetString("organization"), nil
			}
		}
	}
	if authOrgRecords[0].GetString("role") == ownerRoleRecord.Id {
		return authOrgRecords[0].GetString("organization"), nil
	}
	return "default", nil
}

func notifyLogsUpdate(app core.App, subscription string, data []map[string]any) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	message := subscriptions.Message{
		Name: subscription,
		Data: rawData,
	}
	clients := app.SubscriptionsBroker().Clients()
	for _, client := range clients {
		if client.HasSubscription(subscription) {
			client.Send(message)
		}
	}
	return nil
}

func WriteAPIError(w http.ResponseWriter, code int, domain, reason, message string) {
	errorResponse := APIErrorResponse{
		APIVersion: "2.0",
		Error: APIError{
			Code:    code,
			Message: message,
			Errors: []APIErrorDetail{
				{
					Domain:  domain,
					Reason:  reason,
					Message: message,
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse)
}
