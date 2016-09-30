package main

import (
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
)

func main() {
	fmt.Println("tester boat.go")
	currentBoat := boat.NewBoat()
	currentBoat.Navigation.Rudder = 21
	fmt.Printf("%s", currentBoat.Marshal())

	newWaypoint := Waypoint{}
	newWaypoint.Type = 2
	fmt.Printf("%s", newWaypoint)

}
