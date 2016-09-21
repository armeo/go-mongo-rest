package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

// JSONHandle ...
func jsonHandle(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

// NewRoute ...
func NewRoute(db DbManager) http.Handler {
	m := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/", jsonHandle(HomeHandle)).Methods("GET")
	r.Handle("/api/v1/notes", jsonHandle(NotesHandle(db))).Methods("GET")
	r.Handle("/api/v1/notes/{code}", jsonHandle(NoteByCodeHandle(db))).Methods("GET")
	r.Handle("/api/v1/notes", jsonHandle(CreateNoteHandle(db))).Methods("POST")
	m.Handle("/", r)

	return m
}
