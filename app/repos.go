package app

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
	GetAll() ([]Note, error)
	Create(note *Note) (*Note, error)
	GetByCode(code string) (*Note, error)
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
func (db *MongoManager) GetAll() ([]Note, error) {
	var notes []Note
	session := db.session.Clone()
	defer session.Close()
	collection := session.DB("notesdb").C("notes")
	err := collection.Find(nil).All(&notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

// GetByCode ...
func (db *MongoManager) GetByCode(code string) (*Note, error) {
	var note Note
	session := db.session.Clone()
	defer session.Close()
	collection := session.DB("notesdb").C("notes")
	err := collection.Find(bson.M{"note_code": code}).One(&note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

// Create ...
func (db *MongoManager) Create(note *Note) (*Note, error) {
	session := db.session.Clone()
	defer session.Close()
	collection := session.DB("notesdb").C("notes")

	objID := bson.NewObjectId()
	note.ID = objID
	note.CreatedOn = time.Now()

	err := collection.Insert(note)
	if err != nil {
		return nil, err
	}

	log.Printf("Inserted new Note %s with name %s", note.ID, note.Title)

	return note, err
}
