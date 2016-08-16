package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
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
		notes, err := db.GetAll()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(res, `{"error_code":%d, "error_msg":"%s"}`, http.StatusInternalServerError, err)
			return
		}

		json.NewEncoder(res).Encode(NotesResource{Notes: notes})
	}
}

// NoteByCodeHandle ...
func NoteByCodeHandle(db DbManager) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		code := vars["code"]
		log.Printf("Note Code: %s", code)

		note, err := db.GetByCode(code)
		if err == mgo.ErrNotFound {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, `{"error_code":%d, "error_msg":"%s"}`, http.StatusNotFound, err)
			return
		} else if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(res, `{"error_code":%d, "error_msg":"%s"}`, http.StatusInternalServerError, err)
			return
		}

		json.NewEncoder(res).Encode(NoteResource{Note: *note})
	}
}

// CreateNoteHandle ...
func CreateNoteHandle(db DbManager) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var noteResource NoteResource
		err := json.NewDecoder(req.Body).Decode(&noteResource)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, `{"error_code":%d, "error_msg":"%s"}`, http.StatusBadRequest, err)
			return
		}

		note, err := db.Create(&noteResource.Note)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(res, `{"error_code":%d, "error_msg":"%s"}`, http.StatusInternalServerError, err)
			return
		}

		json.NewEncoder(res).Encode(NoteResource{Note: *note})
	}
}

func Route(db DbManager) http.Handler {
	m := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/", JSONHandle(HomeHandle)).Methods("GET")
	r.Handle("/api/v1/notes", JSONHandle(NotesHandle(db))).Methods("GET")
	r.Handle("/api/v1/notes/{code}", JSONHandle(NoteByCodeHandle(db))).Methods("GET")
	r.Handle("/api/v1/notes", JSONHandle(CreateNoteHandle(db))).Methods("POST")
	m.Handle("/", r)
	return m
}

func main() {
	log.Println("Listening on 8000")
	http.ListenAndServe(":8000", Route(NewMongoManager("localhost")))
}
