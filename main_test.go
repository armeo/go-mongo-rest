package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDb struct{}

func (db MockDb) GetAll() ([]Note, error) {
	return nil, nil
}

func (db MockDb) Create(note *Note) (*Note, error) {
	n := Note{Title: "test", Description: "test"}

	return &n, nil
}

func (db MockDb) GetByCode(code string) (*Note, error) {
	n := Note{Title: "test", Description: "test"}

	return &n, nil
}

func TestHomeHandle(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	HomeHandle(res, req)

	assert.Equal(t, 200, res.Code)

	expected := map[string]string{"message": "Hello"}
	var actual map[string]string
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func TestNotesHandle(t *testing.T) {
	mockDb := MockDb{}
	notesHandle := NotesHandle(mockDb)
	req, _ := http.NewRequest("GET", "/api/v1/notes", nil)
	res := httptest.NewRecorder()
	notesHandle.ServeHTTP(res, req)

	var expected NotesResource
	var actual NotesResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func TestCreateNoteHandle(t *testing.T) {
	mockDb := MockDb{}
	createNoteHandle := CreateNoteHandle(mockDb)
	var jsonStr = []byte(`{"note":{"title":"test", "description":"test"}}`)
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()
	createNoteHandle.ServeHTTP(res, req)

	n := Note{Title: "test", Description: "test"}
	expected := NoteResource{Note: n}

	var actual NoteResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected.Note.Title, actual.Note.Title)
	assert.Equal(t, expected.Note.Description, actual.Note.Description)
}

func TestNoteByCodeHandle(t *testing.T) {
	mockDb := MockDb{}
	noteByCodeHandle := NoteByCodeHandle(mockDb)
	req, _ := http.NewRequest("GET", "/api/v1/notes/test1", nil)
	res := httptest.NewRecorder()
	noteByCodeHandle.ServeHTTP(res, req)

	var expected NotesResource
	var actual NotesResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}
