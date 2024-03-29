package http

import (
	"net/http"

	"github.com/cloudtrust/common-service/v2/security"
)

// MakeRightsHandler makes a HTTP handler that returns information about the rights of the user.
func MakeRightsHandler(authorizationManager security.AuthorizationManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = EncodeReply(r.Context(), w, authorizationManager.GetRightsOfCurrentUser(r.Context()))
	})
}
