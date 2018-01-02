package interval

import (
	"net/http"

	"github.com/flosch/pongo2"
)

// Interval struct
type Interval struct {
	Exercises []Exercise
	Pauses    []Pause
}

// Exercise struct
type Exercise struct {
	ID       int
	Name     string
	Duration int
}

// Pause struct
type Pause struct {
	ID       int
	Duration int
}

var intervals []Interval
var exercises []Exercise
var pauses []Pause

func init() {
	exercises = append(exercises, Exercise{ID: 1, Name: "pushups", Duration: 60})
	exercises = append(exercises, Exercise{ID: 3, Name: "burpees", Duration: 60})
	pauses = append(pauses, Pause{ID: 2, Duration: 15})
	intervals = append(intervals, Interval{Exercises: exercises, Pauses: pauses})
}

func addExercise(name string, duration int) Exercise {
	return Exercise{1, name, duration}
}

// Duration returns the summed interval time
func (i Interval) Duration() int {
	var duration int

	for _, e := range i.Exercises {
		duration += e.Duration
	}

	for _, e := range i.Pauses {
		duration += e.Duration
	}

	return duration
}

// HandleRequest returns the requested interval
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	tplExample := pongo2.Must(pongo2.FromFile("template/index.html"))

	ctx := pongo2.Context{
		"interval": intervals,
	}

	// Execute the template per HTTP request
	err := tplExample.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
