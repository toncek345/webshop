package urls

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/senko/clog"
	"github.com/toncek345/webshop/models"
)

func adminUrls(r chi.Router) {
	r.Post("/user/login",
		func(w http.ResponseWriter, r *http.Request) {
			var obj models.User

			err := decode(r, &obj)
			if err != nil {
				clog.Warningf("%s", err)
				respond(w, r, http.StatusBadRequest, err)
				return
			}

			auth, err := models.AuthUser(obj.Username, obj.Password)
			if err != nil {
				clog.Warningf("%s", err)
				respond(w, r, http.StatusBadRequest, err)
				return
			}

			respond(w, r, http.StatusOK, auth)
		})

	r.Post("/user/logout",
		func(w http.ResponseWriter, r *http.Request) {
			token := getAuthHeader(r)
			err := models.RemoveAuth(token)
			if err != nil {
				respond(w, r, http.StatusInternalServerError, err)
				return
			}

			respond(w, r, http.StatusOK, nil)
		})
}
