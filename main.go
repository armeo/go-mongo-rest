package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func JSONHandle(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

// HomeHandler ...
func HomeHandler(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(map[string]string{"message": "Hello"})
}

// NotesHandler ...
func NotesHandler(db DbManager) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		json.NewEncoder(res).Encode(NotesResource{Notes: db.GetAll()})
	}
}

func main() {
	db := NewMongoManager("localhost")
	r := mux.NewRouter()
	r.HandleFunc("/", JSONHandle(HomeHandler)).Methods("GET")
	r.Handle("/api/v1/notes", JSONHandle(NotesHandler(db))).Methods("GET")
	http.Handle("/", r)

	log.Println("Listening on 8000")
	http.ListenAndServe(":8000", nil)
}
