package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/products", handler).Methods("POST")
	r.HandleFunc("/articles", handler).Methods("GET")
	r.HandleFunc("/articles/{id}", handler).Methods("GET", "PUT")
	r.HandleFunc("/authors", handler).Queries("surname", "{surname}")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		qt, err := route.GetQueriesTemplates()
		if err != nil {
			return nil
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		p, err := route.GetPathRegexp()
		if err != nil {
			return nil
		}
		// qr will contain a list of regular expressions with the same semantics as GetPathRegexp,
		// just applied to the Queries pairs instead, e.g., 'Queries("surname", "{surname}") will return
		// {"^surname=(?P<v0>.*)$}. Where each combined query pair will have an entry in the list.
		qr, err := route.GetQueriesRegexp()
		if err != nil {
			return nil
		}
		m, err := route.GetMethods()
		if err != nil {
			return nil
		}
		fmt.Println(strings.Join(m, ","), strings.Join(qt, ","), strings.Join(qr, ","), t, p)
		return nil
	})
	http.Handle("/", r)
}
