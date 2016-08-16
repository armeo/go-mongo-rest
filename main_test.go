package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/stretchr/testify/assert"
)

type MockDb struct {
	err   error
	note  *Note
	notes []Note
}

func (db *MockDb) GetAll() ([]Note, error) {
	return db.notes, db.err
}

func (db *MockDb) Create(note *Note) (*Note, error) {
	db.note = &Note{Title: "test", Description: "test"}

	return db.note, db.err
}

func (db *MockDb) GetByCode(code string) (*Note, error) {
	return db.note, db.err
}

func TestHomeHandle(t *testing.T) {
	mockDb := MockDb{}
	mux := Route(&mockDb)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)

	expected := map[string]string{"message": "Hello"}
	var actual map[string]string
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func TestNotesHandle(t *testing.T) {
	mockDb := MockDb{}
	mux := Route(&mockDb)
	req, _ := http.NewRequest("GET", "/api/v1/notes", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	var expected NotesResource
	var actual NotesResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func TestCreateNoteHandle(t *testing.T) {
	mockDb := MockDb{}
	mux := Route(&mockDb)
	var jsonStr = []byte(`{"note":{"title":"test", "description":"test"}}`)
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	n := Note{Title: "test", Description: "test"}
	expected := NoteResource{Note: n}

	var actual NoteResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected.Note.Title, actual.Note.Title)
	assert.Equal(t, expected.Note.Description, actual.Note.Description)
}

func TestNoteByCodeHandle(t *testing.T) {
	now := time.Now()
	mockDb := MockDb{
		note: &Note{Title: "test", Description: "test", CreatedOn: now},
	}

	mux := Route(&mockDb)

	req, _ := http.NewRequest("GET", "/api/v1/notes/test1", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	expected := NoteResource{Note: Note{Title: "test", Description: "test", CreatedOn: now}}
	var actual NoteResource
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}

func TestNoteByCodeHandleNotFound(t *testing.T) {
	mockDb := MockDb{
		err: mgo.ErrNotFound,
	}

	mux := Route(&mockDb)

	req, _ := http.NewRequest("GET", "/api/v1/notes/test1", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusNotFound)
	expected := map[string]interface{}{"error_code": float64(http.StatusNotFound), "error_msg": "not found"}
	var actual map[string]interface{}
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}

func TestNoteByCodeHandleInternalServerError(t *testing.T) {
	mockDb := MockDb{
		err: fmt.Errorf("Internal Server Error"),
	}

	mux := Route(&mockDb)

	req, _ := http.NewRequest("GET", "/api/v1/notes/test1", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusInternalServerError)
	expected := map[string]interface{}{"error_code": float64(http.StatusInternalServerError), "error_msg": "Internal Server Error"}
	var actual map[string]interface{}
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}
