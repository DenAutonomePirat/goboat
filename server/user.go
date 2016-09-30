package server

import (
	"time"
)

type User struct {
	ID          uint32
	Online      bool
	OnlineHours time.Duration
	Waypoints   []Waypoint
}
