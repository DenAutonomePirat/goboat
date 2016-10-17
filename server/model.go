package server

import (
	"encoding/json"
	"github.com/denautonomepirat/goboat/boat"
)

type GameSetup struct {
	Class                      string        `json:"class"`
	Start                      boat.Waypoint `json:"start"`
	Finish                     boat.Waypoint `json:"finish"`
	WaypointsAllowed           int           `json:"waypointsAllowed"`
	DefaultLegDistanceInMeters float64       `json:"DefaultLegDistanceInMeters"`
}

func NewGameSetup() *GameSetup {
	g := GameSetup{}
	g.Class = "GameSetup"
	return &g
}

func (g *GameSetup) Marshal() *[]byte {
	encoded, _ := json.Marshal(g)
	return &encoded
}
