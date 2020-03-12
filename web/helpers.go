package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/senko/clog"
)

// getAuthHeader gets auth uuid from header which is sent in x-auth.
func (app *App) getAuthHeader(r *http.Request) string {
	return r.Header.Get("x-auth")
}

// JSONDecode decodes given request to object that is given.
// Call this function with pointer to object that needs to be decoded.
func JSONDecode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

// JSONRespond generic respond func.
func (app *App) JSONRespond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		clog.Errorf("error marshalling obj: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		clog.Errorf("error returning data: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// parseUrlVarsInt parses variables from url to integer.
//
// TODO: consider removing this
func parseUrlVarsInt(r *http.Request, name string) (int, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, name), 10, 32)
	if err != nil {
		clog.Warningf("id parse error: %s", err)
	}
	return int(id), nil
}
