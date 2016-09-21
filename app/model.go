package app

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Note ...
type Note struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	NoteCode    string        `bson:"note_code" json:"note_code"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	CreatedOn   time.Time     `json:"createdon"`
}

// NoteResource ...
type NoteResource struct {
	Note Note `json:"note"`
}

// NotesResource ...
type NotesResource struct {
	Notes []Note `json:"notes"`
}
