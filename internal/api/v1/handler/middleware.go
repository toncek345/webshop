package handler

import (
	"net/http"
)

// Middleware for checking if the client is authorized. Returns 403 if not.
func (app *App) authenticationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := app.getAuthHeader(r)

		ok, err := app.models.Auth.IsAuth(authToken)
		if err != nil {
			app.JSONRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		if !ok {
			app.JSONRespond(w, r, http.StatusForbidden, "Not authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
