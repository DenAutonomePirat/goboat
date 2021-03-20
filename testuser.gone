package main

import (
	"encoding/json"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"github.com/denautonomepirat/goboat/server"
)

func main() {
	User := server.NewUser()
	User.Name = "Thomas"
	User.Waypoints[0].Name = "waypoint 0"
	User.Waypoints[1].Name = "waypoint 1"
	User.Waypoints[2].Name = "waypoint 2"
	User.WaypointReached()
	fmt.Printf("%s\n", *User.Marshal())
	w := boat.NewWaypoint()
	w.Coordinate[0] = 12.54
	w.Coordinate[1] = 65.45
	w.Name = "lkj 34"
	User.SetWaypoint(2, w)
	fmt.Printf("%s\n", *User.Marshal())
}
