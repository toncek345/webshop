package urls

import (
	"net/http"

	"webshop/models"

	"github.com/gorilla/mux"
	"github.com/senko/clog"
)

func adminUrls(r *mux.Router) {
	r.HandleFunc("/user/login",
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
		}).Methods("POST")

	r.HandleFunc("/user/logout",
		func(w http.ResponseWriter, r *http.Request) {
			token := getAuthHeader(r)
			err := models.RemoveAuth(token)
			if err != nil {
				respond(w, r, http.StatusInternalServerError, err)
				return
			}

			respond(w, r, http.StatusOK, nil)
		}).Methods("POST")
}
