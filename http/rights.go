package http

import (
	"context"
	"net/http"

	"github.com/cloudtrust/common-service/security"
)

// MakeRigtsHandler makes a HTTP handler that returns information about the rights of the user.
func MakeRightsHandler(authorizationManager security.AuthorizationManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		EncodeReply(context.TODO(), w, authorizationManager.GetRightsOfCurrentUser(r.Context()))
	})
}
