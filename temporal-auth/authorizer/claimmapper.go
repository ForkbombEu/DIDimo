package authorizer

import (
	"fmt"
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
)

type CustomClaimMapper struct {
	app *pocketbase.PocketBase
}

func DidimoCustomClaimMapper(_ *config.Config, app *pocketbase.PocketBase) *CustomClaimMapper {
	return &CustomClaimMapper{app: app}
}

// GetClaims extracts user claims and retrieves organization names
func (m *CustomClaimMapper) GetClaims(
	authInfo *authorization.AuthInfo,
) (*authorization.Claims, error) {
	// Decode JWT Token to get user ID
	authRecord, err := m.app.FindAuthRecordByToken(authInfo.AuthToken, core.TokenTypeAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to extract find user: %w", err)
	}
	orgNames, err := m.getUserOrganizations(authRecord.Id)
	if err != nil {
		log.Printf("Error fetching organizations for user %s: %v", authRecord.Id, err)
		return nil, fmt.Errorf("failed to resolve user organizations")
	}

	// Create a map for Namespaces where the key is the organization name, and the value is a Role
	namespaces := make(map[string]authorization.Role)
	for _, orgName := range orgNames {
		namespaces[orgName] = authorization.RoleReader | authorization.RoleWriter
	}

	return &authorization.Claims{
		Subject:    authRecord.Id,
		Namespaces: namespaces,
	}, nil
}

// getUserOrganizations fetches the organizations and roles of a user
func (m *CustomClaimMapper) getUserOrganizations(userID string) ([]string, error) {
	var orgNames []string

	// Query user_organization table for user's organizations
	records, err := m.app.FindAllRecords("user_organization",
		dbx.NewExp("user = {:user}", dbx.Params{"user": userID}),
	)
	if err != nil {
		return nil, err
	}

	// Process each organization relationship
	for _, rec := range records {
		orgID := rec.GetString("organization")

		// Get organization details
		orgRec, err := m.app.FindRecordById("organizations", orgID)
		if err != nil {
			log.Printf("Warning: Organization ID %s not found, skipping.", orgID)
			continue // Skip invalid organizations
		}

		orgNames = append(orgNames, orgRec.GetString("name"))
	}

	return orgNames, nil

}
