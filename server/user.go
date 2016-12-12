package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"time"
)

type User struct {
	UserName       string    `json:"userName"bson:"userName"`
	HashedPassword []byte    `bson:"hashedPassword"`
	Salt           []byte    `bson:"salt"`
	Created        time.Time `json:"created"bson:"created"`
	Online         bool      `json:"online"`
	connection     *Conn
}

func NewUser() *User {
	u := User{}
	u.Created = time.Now()
	u.HashedPassword = make([]byte, 64)
	u.Salt = make([]byte, 32)
	return &u
}

func (u *User) Marshal() *[]byte {
	encoded, _ := json.Marshal(u)
	return &encoded
}

func (u *User) SetPassword(password string) {
	_, err := io.ReadFull(rand.Reader, u.Salt)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := scrypt.Key([]byte(password), u.Salt, 1<<14, 8, 1, 64)
	if err != nil {
		log.Fatal(err)
	}
	u.HashedPassword = hash
}

func (u *User) CheckPassword(passwordToTest string) bool {

	hashToCompare, err := scrypt.Key([]byte(passwordToTest), u.Salt, 1<<14, 8, 1, 64)

	if err != nil {
		log.Fatal(err)
	}

	if subtle.ConstantTimeCompare(hashToCompare, u.HashedPassword) == 1 {
		return true
	}
	return false

}

type SkipperWatch struct {
	UserName       string        `json:"UserName"bson:"userName`
	Id             []byte        `json:"id"bson:"id"`
	ConnectedAt    time.Time     `json:"connectedAt"bson:"connectedAt"`
	OnlineDuration time.Duration `json:"onlineDuration"`
	Expiry         time.Time
}

func newSkipperWatch() *SkipperWatch {
	s := SkipperWatch{}
	return &s

}
func name() {

}
