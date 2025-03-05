package authorizer

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/server/common/authorization"
)

type CustomAuthorizer struct{}

// NewCustomAuthorizer initializes the authorizer
func DIDimoCustomAuthorizer() *CustomAuthorizer {
	return &CustomAuthorizer{}
}

// Authorize checks if a user can access a Temporal namespace
func (a *CustomAuthorizer) Authorize(
	ctx context.Context,
	claims *authorization.Claims,
	target *authorization.CallTarget,
) (authorization.Result, error) {

	// Ensure claims.Namespaces is a map[string]Role
	if claims.Namespaces == nil {
		log.Println("Invalid or missing Namespaces in claims")
		return authorization.Result{
			Decision: authorization.DecisionDeny,
			Reason:   "Invalid claims format",
		}, nil
	}

	// Extract user's organization names from namespaces
	_, exists := claims.Namespaces[target.Namespace]
	if !exists {
		log.Println("User does not have access to this namespace")
		return authorization.Result{
			Decision: authorization.DecisionDeny,
			Reason:   "User has no access to the requested namespace",
		}, nil
	}
	return authorization.Result{
		Decision: authorization.DecisionAllow,
		Reason:   fmt.Sprintf("User has access to namespace %s", target.Namespace),
	}, nil
}
