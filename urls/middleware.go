package urls

import (
	"net/http"
	"time"
	"webshop/models"

	"github.com/senko/clog"
)

// Route logging middleware. Prints in console http version, method,
// request path and duration.
func logRoute(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()
		handle.ServeHTTP(w, r)
		endTime := time.Now()

		duration := endTime.Sub(initTime)
		clog.Debugf("%s %s %s %dms",
			r.Proto,
			r.Method,
			r.RequestURI,
			duration.Nanoseconds()/int64(time.Millisecond))
	})
}

// Middleware for checking if the client is authorized. Returns 403 if not.
func authenticationRequired(handle http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := getAuthHeader(r)

		ok := models.IsAuth(authToken)
		if !ok {
			respond(w, r, http.StatusForbidden, "Not authorized")
			return
		}

		handle.ServeHTTP(w, r)
	})
}
