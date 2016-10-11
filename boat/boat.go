package boat

import (
	"encoding/json"
	"fmt"
	"log"
)

func CheckGracefull(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func Check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Boat struct {
	Class      string     `json:"class"`
	Navigation Nav        `json:"navigation,omitempty"`
	Power      Electrical `json:"power,omitempty"`
	Routes     []Route
}

func NewBoat() *Boat {
	b := Boat{}
	b.Class = "Boat"
	return &b
}

func (b *Boat) Marshal() *[]byte {
	encoded, _ := json.Marshal(b)
	return &encoded
}

type Nav struct {
	Position         Point   `json:"position,omitempty"`
	SpeedOverGround  float32 `json:"speedGPS,omitempty"`
	CourseOverGround float32 `json:"courseGPS,omitempty"`
	HeadingMagnetic  float32 `json:"heading,omitempty"`
	Log              float32 `json:"log,omitempty"`
	Depth            float32 `json:"depth,omitempty"`
	MainSail         int32   `json:"mainsail,omitempty"`
	Jib              int32   `json:"jib,omitempty"`
	Rudder           int32   `json:"rudder,omitempty"`
	Pitch            float32 `json:"pitch,omitempty"`
	Roll             float32 `json:"roll,omitempty"`
	Rotation         float32 `json:"rot,omitempty"`
}

type Point [2]float64

type Electrical struct {
	Volts       float32 `json:"volts,omitempty"`
	Amperes     float32 `json:"amperes,omitempty"`
	JoulesTotal float32 `json:"joules_total,omitempty"`
}
type Route struct {
	Waypoints []Waypoint `json:"waypoints"`
}

func NewNav() *Nav {
	n := Nav{}
	return &n
}

func (n *Nav) Marshal() *[]byte {
	encoded, _ := json.Marshal(n)
	return &encoded
}

type Waypoint struct {
	Name       string `json:"name,omitempty"`
	Type       int    `json:"type,omitempty"`
	Coordinate Point  `json:"coordinate,omitempty"`
	Message    string `json:"message,omitempty"`
}

func NewWaypoint() *Waypoint {
	w := Waypoint{}
	return &w
}

func (w *Waypoint) Marshal() *[]byte {
	encoded, _ := json.Marshal(w)
	return &encoded
}

type Muxable interface {
	Marshal() *[]byte
}
