package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

// MongoManager ...
type MongoManager struct {
	session *mgo.Session
}

// DbManager ...
type DbManager interface {
	GetAll() []Note
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
