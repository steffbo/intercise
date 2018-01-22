package interval

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Interval struct
type Interval struct {
	ID    int
	Items []Item `json:"items"`
}

// Item is an interface for both Exercise and Pause
type Item struct {
	ID       int
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

// MakeExercise creates a new Item for an exercise
func MakeExercise(name string, duration int) Item {
	return Item{-1, name, duration}
}

// MakePause creates a new Item for a pause
func MakePause(duration int) Item {
	return Item{-1, "pause", duration}
}

// AddItems adds multiple Items at once to an interval
func (i Interval) AddItems(items []Item) Interval {
	for _, e := range items {
		i.AddItem(e)
	}
	return i
}

// AddItem adds a single item to an interval
// The item is assigned an ID
func (i Interval) AddItem(item Item) Interval {
	item.ID = len(i.Items)
	if i.Items == nil {
		i.Items = []Item{item}
	} else {
		i.Items = append(i.Items, item)
	}
	return i
}

func (i Item) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, Duration: %d", i.ID, i.Name, i.Duration)
}

func (i Interval) String() string {
	str := ""
	for _, e := range i.Items {
		str = str + e.String()
	}
	return str
}

// AddExercise adds a new Item to the slice of items
func (i Interval) AddExercise(name string, duration int) Interval {

	nextID := len(i.Items)
	newItem := Item{nextID, name, duration}

	if i.Items == nil {
		i.Items = []Item{newItem}
	} else {
		i.Items = append(i.Items, newItem)
	}
	return i
}

// AddPause adds a new Item to the slice of items
func (i Interval) AddPause(duration int) Interval {

	nextID := len(i.Items)
	newItem := Item{nextID, "pause", duration}

	if i.Items == nil {
		i.Items = []Item{newItem}
	} else {
		i.Items = append(i.Items, newItem)
	}
	return i
}

// Duration returns the summed interval time
func (i Interval) Duration() int {
	var duration int

	for _, e := range i.Items {
		duration += e.Duration
	}

	return duration
}

var intervals []Interval

func init() {
	interval0 := Interval{0, nil}.AddExercise("pushups", 60).AddPause(20).AddExercise("burpees", 60)
	interval1 := Interval{1, nil}.AddExercise("jumping jacks", 60).AddPause(20).AddExercise("crunches", 60)
	intervals = append(intervals, interval0, interval1)
}

// GetIntervals returns all intervals
func GetIntervals(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(intervals)
}

// GetInterval returns a single interval
func GetInterval(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range intervals {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// CreateInterval adds a new interval
func CreateInterval(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var interval Interval
	_ = json.NewDecoder(r.Body).Decode(&interval)
	id, _ := strconv.Atoi(params["id"])
	interval.ID = id
	intervals = append(intervals, interval)
	json.NewEncoder(w).Encode(intervals)
}

// DeleteInterval removes an interval
func DeleteInterval(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range intervals {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			intervals = append(intervals[:index], intervals[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(intervals)
	}
}
