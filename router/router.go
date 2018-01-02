package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/steffbo/intercise/home"
)

func Router() {
	r := mux.NewRouter()
	r.HandleFunc("/", home.Home)
	http.Handle("/", r)
}
