package server

import (
	"github.com/denautonomepirat/goboat/boat"
	"time"
)

type User struct {
	ID          uint32
	Online      bool
	OnlineHours time.Duration
	Waypoints   []boat.Waypoint
}
