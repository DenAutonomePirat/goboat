package main

import (
	//"fmt"
	//"time"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
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

	return &Store{db: session}
}
