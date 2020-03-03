package urls

import (
	"net/http"

	"github.com/toncek345/webshop/models"
)

// Middleware for checking if the client is authorized. Returns 403 if not.
func authenticationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := getAuthHeader(r)

		ok := models.IsAuth(authToken)
		if !ok {
			respond(w, r, http.StatusForbidden, "Not authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
