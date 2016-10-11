package server

import (
	"encoding/json"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"time"
)

type User struct {
	connection     *Conn
	Class          string           `json:"class"`
	Id             uint32           `json:"id"`
	Name           string           `json:"name"`
	Online         bool             `json:"online"`
	ConnectedAt    time.Time        `json:"connected_at"`
	OnlineDuration time.Duration    `json:"onlineDuration"`
	Waypoints      [3]boat.Waypoint `json:"waypoints"`
}

func NewUser() *User {
	u := User{}
	u.Class = "User"
	return &u
}

func (u *User) Marshal() *[]byte {
	encoded, _ := json.Marshal(u)
	return &encoded
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
