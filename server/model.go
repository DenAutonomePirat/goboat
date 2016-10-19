package server

import (
	"encoding/json"
	"github.com/denautonomepirat/goboat/boat"
)

type Configuration struct {
	Start                      boat.Waypoint `json:"start"`
	Finish                     boat.Waypoint `json:"finish"`
	WaypointsAllowed           int           `json:"waypointsAllowed"`
	DefaultLegDistanceInMeters float64       `json:"DefaultLegDistanceInMeters"`
}

func NewConfiguration() *Configuration {
	c := Configuration{}
	return &c
}

func (c *Configuration) Marshal() *[]byte {
	encoded, _ := json.Marshal(c)
	return &encoded
}
