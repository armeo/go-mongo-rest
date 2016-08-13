package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDb struct{}

func (db MockDb) GetAll() []Note {
	return nil
}

func TestHomeHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	HomeHandler(res, req)

	assert.Equal(t, 200, res.Code)
	expected := map[string]string{"message": "Hello"}
	var actual map[string]string
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}

func TestNotesHandler(t *testing.T) {
	mockDb := MockDb{}
	notesHandle := NotesHandler(mockDb)
	req, _ := http.NewRequest("GET", "/api/v1/notes", nil)
	res := httptest.NewRecorder()
	notesHandle.ServeHTTP(res, req)
	var expected NotesResource
	var actual NotesResource
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}
