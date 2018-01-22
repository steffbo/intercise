package router

import (
	"net/http"

	"github.com/steffbo/intercise/home"

	"github.com/gorilla/mux"
	"github.com/steffbo/intercise/interval"
)

// Router initializes all possible routes
func Router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", home.Home).Methods("GET")
	r.HandleFunc("/interval", interval.GetIntervals).Methods("GET")
	r.HandleFunc("/interval/{id}", interval.GetInterval).Methods("GET")
	r.HandleFunc("/interval/{id}", interval.CreateInterval).Methods("POST")
	r.HandleFunc("/interval/{id}", interval.DeleteInterval).Methods("DELETE")
	return r
}
