package server

import (
	//"fmt"
	//"time"

	"github.com/denautonomepirat/goboat/boat"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store struct {
	db *mgo.Session
}

func NewStore() *Store {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err.Error())
	}

	index := mgo.Index{
		Key:        []string{"+name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	session.DB("redboat").C("users").EnsureIndex(index)

	index = mgo.Index{
		Key:        []string{"+timeStamp"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	session.DB("redboat").C("track").EnsureIndex(index)

	return &Store{db: session}
}

func (s *Store) AddUser(u *User) error {
	_, err := s.db.DB("redboat").C("users").Upsert(bson.M{"name": u.Name}, u)
	return err
}

func (s *Store) AddTrack(b *boat.Boat) error {
	err := s.db.DB("redboat").C("track").Insert(b)
	return err
}
