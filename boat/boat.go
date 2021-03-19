package boat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
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

type Command struct {
	Class    string `json:"class"`
	Waypoint struct {
		Coordinate struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"position"`
		Name string `json:"name"`
	} `json:"waypoint"`
}

type Boat struct {
	Class      string     `json:"class" bson:"-"`
	TimeStamp  int64      `json:"timestamp" bson:"timestamp"`
	Navigation Nav        `json:"navigation,omitempty" bson:"navigation,omitempty"`
	Power      Electrical `json:"power,omitempty" bson:"power,omitempty"`
	Route      []Route    `json:"route,omitempty" bson:"route,omitempty"`
}

func NewBoat() *Boat {
	b := Boat{}
	b.Class = "Boat"
	b.TimeStamp = int64(time.Now().UnixNano() / 1000 / 1000)
	return &b
}

func (b *Boat) Marshal() *[]byte {
	encoded, _ := json.Marshal(b)
	return &encoded
}

type Nav struct {
	Position         Point   `json:"position,omitempty" bson:"position,omitempty"`
	SpeedOverGround  float32 `json:"speedGPS,omitempty" bson:"speedGPS,omitempty"`
	CourseOverGround float32 `json:"courseGPS,omitempty" bson:"courseGPS,omitempty"`
	HeadingMagnetic  float32 `json:"heading,omitempty" bson:"heading,omitempty"`
	Log              float32 `json:"log,omitempty" bson:"log,omitempty"`
	Depth            float32 `json:"depth,omitempty" bson:"depth,omitempty"`
	MainSail         int32   `json:"mainsail,omitempty" bson:"mainsail,omitempty"`
	Jib              int32   `json:"jib,omitempty" bson:"jib,omitempty"`
	Rudder           int32   `json:"rudder,omitempty" bson:"rudder,omitempty"`
	Pitch            float32 `json:"pitch,omitempty" bson:"pitch,omitempty"`
	Roll             float32 `json:"roll,omitempty" bson:"roll,omitempty"`
	Rotation         float32 `json:"rot,omitempty" bson:"rot,omitempty"`
}

type Point [2]float64

type Electrical struct {
	Volts       float32 `json:"volts,omitempty" bson:"volts,omitempty"`
	Amperes     float32 `json:"amperes,omitempty" bson:"amperes,omitempty"`
	JoulesTotal float32 `json:"joules_total,omitempty" bson:"joules_total,omitempty"`
}
type Route struct {
	Waypoints []Waypoint `json:"waypoints,omitempty" bson:"waypoints,omitempty"`
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
	Name       string `json:"name,omitempty" bson:"name,omitempty"`
	Type       int    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinate Point  `json:"coordinate,omitempty" bson:"coordinate,omitempty"`
	Message    string `json:"message,omitempty" bson:"message,omitempty"`
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
