package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	goji "goji.io"
	"goji.io/pat"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DbName = "intercise"
const DbCol = "intervals"

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// Interval struct
type Interval struct {
	ID    int    `json:"id"`
	Items []Item `json:"items"`
}

// Item is an interface for both Exercise and Pause
type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/intervals"), allIntervals(session))
	mux.HandleFunc(pat.Post("/intervals"), addInterval(session))
	mux.HandleFunc(pat.Get("/intervals/:id"), intervalByID(session))
	mux.HandleFunc(pat.Put("/intervals/:id"), updateInterval(session))
	mux.HandleFunc(pat.Delete("/intervals/:id"), deleteInterval(session))
	http.ListenAndServe("localhost:8080", mux)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(DbName).C(DbCol)

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func allIntervals(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DbName).C(DbCol)

		var intervals []Interval
		err := c.Find(bson.M{}).All(&intervals)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all intervals: ", err)
			return
		}

		respBody, err := json.MarshalIndent(intervals, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func addInterval(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var interval Interval
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&interval)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB(DbName).C(DbCol)

		err = c.Insert(interval)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Interval with this ID already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert interval: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+strconv.Itoa(interval.ID))
		w.WriteHeader(http.StatusCreated)
	}
}

func intervalByID(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		idParam := pat.Param(r, "id")
		id, _ := strconv.Atoi(idParam)

		c := session.DB(DbName).C(DbCol)

		var interval Interval
		err := c.Find(bson.M{"id": id}).One(&interval)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find interval: ", err)
			return
		}

		respBody, err := json.MarshalIndent(interval, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func updateInterval(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")

		var interval Interval
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&interval)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB(DbName).C(DbCol)

		err = c.Update(bson.M{"id": id}, &interval)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed update interval: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Interval not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteInterval(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")

		c := session.DB(DbName).C(DbCol)

		err := c.Remove(bson.M{"id": id})
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed delete interval: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Interval not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
