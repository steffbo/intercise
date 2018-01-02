package home

import (
	"net/http"

	"github.com/flosch/pongo2"
)

// Home delivers the homepage
func Home(w http.ResponseWriter, r *http.Request) {

	tplExample := pongo2.Must(pongo2.FromFile("template/index.html"))

	ctx := pongo2.Context{
		"interval": "stefan",
	}

	// Execute the template per HTTP request
	err := tplExample.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
