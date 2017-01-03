package server

import (
	"encoding/json"
	"github.com/kellydunn/golang-geo"
)

type Configuration struct {
	Start                      *geo.Point `json:"start"`
	Finish                     *geo.Point `json:"finish"`
	WaypointsAllowed           int        `json:"waypointsAllowed"`
	DefaultLegDistanceInMeters float64    `json:"DefaultLegDistanceInMeters"`
}

func NewConfiguration() *Configuration {
	c := Configuration{}
	return &c
}

func (c *Configuration) Marshal() *[]byte {
	encoded, _ := json.Marshal(c)
	return &encoded
}

type Skipper struct {
	Class  string    `json:"class"`
	Wpt    int       `json:"wpt"`
	Latlng geo.Point `json:"latlng"`
	User   string    `json:"user"`
}

func (s *Skipper) Marshal() *[]byte {
	encoded, _ := json.Marshal(s)
	return &encoded

}
