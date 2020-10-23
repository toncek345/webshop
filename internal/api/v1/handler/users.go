package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/senko/clog"
)

func (app *App) adminRouter(r chi.Router) {
	r.Post("/user/login",
		func(w http.ResponseWriter, r *http.Request) {
			var obj struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			if err := app.JSONDecode(r, &obj); err != nil {
				clog.Warningf("%s", err)
				app.JSONRespond(w, r, http.StatusBadRequest, err)
				return
			}

			user, err := app.models.Users.GetByUsername(obj.Username)
			if err != nil {
				clog.Warningf("%s", err)
				app.JSONRespond(w, r, http.StatusNotFound, "finding user error")
				return
			}

			auth, err := app.models.Auth.AuthUser(user, obj.Password)
			if err != nil {
				clog.Warningf("%s", err)
				app.JSONRespond(w, r, http.StatusBadRequest, err)
				return
			}

			app.JSONRespond(w, r, http.StatusOK, auth)
		})

	r.Post("/user/logout",
		func(w http.ResponseWriter, r *http.Request) {
			token := app.getAuthHeader(r)
			if err := app.models.Auth.RemoveAuth(token); err != nil {
				clog.Error(err.Error())
				app.JSONRespond(w, r, http.StatusInternalServerError, err)
				return
			}

			app.JSONRespond(w, r, http.StatusOK, nil)
		})
}
