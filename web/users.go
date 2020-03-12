package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/senko/clog"
	"github.com/toncek345/webshop/models"
)

func (app *App) adminRouter(r chi.Router) {
	r.Post("/user/login",
		func(w http.ResponseWriter, r *http.Request) {
			var obj models.User

			err := JSONDecode(r, &obj)
			if err != nil {
				clog.Warningf("%s", err)
				app.JSONRespond(w, r, http.StatusBadRequest, err)
				return
			}

			auth, err := app.models.Auth.AuthUser(obj.Username, obj.Password)
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
			err := app.models.Auth.RemoveAuth(token)
			if err != nil {
				app.JSONRespond(w, r, http.StatusInternalServerError, err)
				return
			}

			app.JSONRespond(w, r, http.StatusOK, nil)
		})
}
