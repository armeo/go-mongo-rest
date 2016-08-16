package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoManager ...
type MongoManager struct {
	session *mgo.Session
}

// DbManager ...
type DbManager interface {
	GetAll() []Note
	Create(note Note) Note
	GetByCode(code string) Note
}

// NewMongoManager ...
func NewMongoManager(dbPath string) *MongoManager {
	log.Println("Starting mongodb session")
	session, err := mgo.Dial(dbPath)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return &MongoManager{session: session}
}

// GetAll ...
func (db *MongoManager) GetAll() []Note {
	var notes []Note
	session := db.session.Clone()
	defer session.Close()
	collection := session.DB("notesdb").C("notes")
	iter := collection.Find(nil).Iter()

	result := Note{}
	for iter.Next(&result) {
		notes = append(notes, result)
	}

	return notes
}

// GetByCode ...
func (db *MongoManager) GetByCode(code string) Note {
	var note Note
	session := db.session.Clone()
	defer session.Close()
	collection := session.DB("notesdb").C("notes")
	err := collection.Find(bson.M{"note_code": code}).One(&note)

	if err != nil {
		log.Printf("Note %s not found !!", code)
	}

	return note
}

// Create ...
func (db *MongoManager) Create(note Note) Note {
	session := db.session.Clone()
	defer session.Close()
	collection := session.DB("notesdb").C("notes")

	objID := bson.NewObjectId()
	note.ID = objID
	note.CreatedOn = time.Now()

	err := collection.Insert(&note)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Inserted new Note %s with name %s", note.ID, note.Title)
	}

	return note
}
