package home

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/steffbo/intercise/interval"
)

// Home delivers the homepage
func Home(w http.ResponseWriter, r *http.Request) {

	tplExample := pongo2.Must(pongo2.FromFile("template/index.html"))

	interval := interval.Interval{}.AddExercise("pushups", 60).AddPause(20).AddExercise("burpees", 60)

	ctx := pongo2.Context{
		"interval": interval,
	}

	// Execute the template per HTTP request
	err := tplExample.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
