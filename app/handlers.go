package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

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
