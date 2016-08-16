package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// JSONHandle ...
func JSONHandle(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

// HomeHandle ...
func HomeHandle(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(map[string]string{"message": "Hello"})
}

// NotesHandle ...
func NotesHandle(db DbManager) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		json.NewEncoder(res).Encode(NotesResource{Notes: db.GetAll()})
	}
}

// NoteByCodeHandle ...
func NoteByCodeHandle(db DbManager) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		code := vars["code"]
		log.Printf("Note Code: %s", code)

		json.NewEncoder(res).Encode(NoteResource{Note: db.GetByCode(code)})
	}
}

// CreateNoteHandle ...
func CreateNoteHandle(db DbManager) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var noteResource NoteResource
		err := json.NewDecoder(req.Body).Decode(&noteResource)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(res).Encode(NoteResource{Note: db.Create(noteResource.Note)})
	}
}

func main() {
	db := NewMongoManager("localhost")
	r := mux.NewRouter()
	r.HandleFunc("/", JSONHandle(HomeHandle)).Methods("GET")
	r.Handle("/api/v1/notes", JSONHandle(NotesHandle(db))).Methods("GET")
	r.Handle("/api/v1/notes/{code}", JSONHandle(NoteByCodeHandle(db))).Methods("GET")
	r.Handle("/api/v1/notes", JSONHandle(CreateNoteHandle(db))).Methods("POST")
	http.Handle("/", r)

	log.Println("Listening on 8000")
	http.ListenAndServe(":8000", nil)
}
