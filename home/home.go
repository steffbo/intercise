package home

import (
	"net/http"

	"github.com/flosch/pongo2"
)

// Home Comment
func Home(w http.ResponseWriter, r *http.Request) {

	template := pongo2.FromFile("template/index.html")
	tplExample := pongo2.Must(template)
		
	ctx := pongo2.Context{
		"destination": "stefan",
	}

	// Execute the template per HTTP request
	err := tplExample.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
