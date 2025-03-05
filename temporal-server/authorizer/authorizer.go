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

// ONLY FOR TEST
type myAuthorizer struct{}

func NewMyAuthorizer() authorization.Authorizer {
	return &myAuthorizer{}
}

var decisionAllow = authorization.Result{Decision: authorization.DecisionAllow}
var decisionDeny = authorization.Result{Decision: authorization.DecisionDeny}

func (a *myAuthorizer) Authorize(_ context.Context, claims *authorization.Claims,
	target *authorization.CallTarget) (authorization.Result, error) {

	// Allow all operations within "temporal-system" namespace
	// DON'T DO THIS IN A PRODUCTION ENVIRONMENT
	// IN PRODUCTION, only allow calls from properly authenticated and authorized callers
	// We are taking a shortcut in the sample because we don't have TLS or a auth token
	if target.Namespace == "temporal-system" {
		return decisionAllow, nil
	}

	// Allow all calls except UpdateNamespace through when claim mapper isn't invoked
	// Claim mapper is skipped unless TLS is configured or an auth token is passed
	if claims == nil && target.APIName != "UpdateNamespace" {
		return decisionAllow, nil
	}

	// Allow all operations for system-level admins and writers
	if claims.System&(authorization.RoleAdmin|authorization.RoleWriter) != 0 {
		return decisionAllow, nil
	}

	// For other namespaces, deny "UpdateNamespace" API unless the caller has a writer role in it
	if target.APIName == "UpdateNamespace" {
		if claims.Namespaces[target.Namespace]&authorization.RoleWriter != 0 {
			return decisionAllow, nil
		} else {
			return decisionDeny, nil
		}
	}

	// Allow all other requests
	return decisionAllow, nil
}
