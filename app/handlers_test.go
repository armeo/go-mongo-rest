package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/armeo/go-mongo-rest/app"
	"github.com/stretchr/testify/assert"
)

// === Setup ====
type MockDb struct {
	err   error
	note  *app.Note
	notes []app.Note
}

func (db *MockDb) GetAll() ([]app.Note, error) {
	return db.notes, db.err
}

func (db *MockDb) Create(note *app.Note) (*app.Note, error) {
	db.note = &app.Note{Title: "test", Description: "test"}

	return db.note, db.err
}

func (db *MockDb) GetByCode(code string) (*app.Note, error) {
	return db.note, db.err
}

func mockRoute(mockDb MockDb, method string, endpoint string, body io.Reader) *httptest.ResponseRecorder {
	mux := app.NewRoute(&mockDb)

	req, _ := http.NewRequest(method, endpoint, body)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	return res
}

// === Test Case ====
func TestHomeHandle(t *testing.T) {
	res := mockRoute(MockDb{}, "GET", "/", nil)

	assert.Equal(t, 200, res.Code)

	expected := map[string]string{"message": "Hello"}
	var actual map[string]string
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func TestNotesHandle(t *testing.T) {
	res := mockRoute(MockDb{}, "GET", "/api/v1/notes", nil)

	var expected app.NotesResource
	var actual app.NotesResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func TestNoteByCodeHandle(t *testing.T) {
	now := time.Now()
	mockDb := MockDb{
		note: &app.Note{Title: "test", Description: "test", CreatedOn: now},
	}

	res := mockRoute(mockDb, "GET", "/api/v1/notes/test1", nil)

	expected := app.NoteResource{Note: app.Note{Title: "test", Description: "test", CreatedOn: now}}
	var actual app.NoteResource
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}

func TestNoteByCodeHandleNotFound(t *testing.T) {
	mockDb := MockDb{
		err: mgo.ErrNotFound,
	}

	res := mockRoute(mockDb, "GET", "/api/v1/notes/test1", nil)

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

	res := mockRoute(mockDb, "GET", "/api/v1/notes/test1", nil)

	assert.Equal(t, res.Code, http.StatusInternalServerError)
	expected := map[string]interface{}{"error_code": float64(http.StatusInternalServerError), "error_msg": "Internal Server Error"}
	var actual map[string]interface{}
	json.NewDecoder(res.Body).Decode(&actual)
	assert.Equal(t, expected, actual)
}

func TestCreateNoteHandle(t *testing.T) {
	var jsonStr = []byte(`{"note":{"title":"test", "description":"test"}}`)
	res := mockRoute(MockDb{}, "POST", "/api/v1/notes", bytes.NewBuffer(jsonStr))

	n := app.Note{Title: "test", Description: "test"}
	expected := app.NoteResource{Note: n}

	var actual app.NoteResource
	json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(t, expected.Note.Title, actual.Note.Title)
	assert.Equal(t, expected.Note.Description, actual.Note.Description)
}
