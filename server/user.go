package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"time"
)

type User struct {
	Name           string           `json:"name"`
	HashedPassword []byte           `bson:"hashedPassword"`
	Salt           []byte           `bson:"salt"`
	connection     *Conn            `json:"-"bson:"-"`
	Online         bool             `json:"online"`
	ConnectedAt    time.Time        `json:"connected_at"`
	OnlineDuration time.Duration    `json:"onlineDuration"`
	Waypoints      [3]boat.Waypoint `json:"waypoints"`
}

func NewUser() *User {
	u := User{}
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
	fmt.Printf("%x\n", hash)

}

func (u *User) CheckPassword(passwordToTest string) error {

	hashToCompare, err := scrypt.Key([]byte(passwordToTest), u.Salt, 1<<14, 8, 1, 64)

	fmt.Printf("%x\n", u.HashedPassword)
	fmt.Printf("%x\n", hashToCompare)

	if err != nil {
		log.Fatal(err)
	}

	if subtle.ConstantTimeCompare(hashToCompare, u.HashedPassword) == 1 {
		return nil
	}
	return errors.New("Password didn't match")

}

func (u *User) WaypointReached() {
	u.Waypoints[0] = u.Waypoints[1]
	u.Waypoints[1] = u.Waypoints[2]
	u.Waypoints[2] = boat.Waypoint{}
}

func (u *User) SetWaypoint(n int, w *boat.Waypoint) {
	u.Waypoints[n] = *w
	fmt.Printf("User %s changed waypoint %d\n", u.Name, n)
}
