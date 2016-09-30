package boat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/stratoberry/go-gpsd"
	"github.com/tarm/serial"
	"log"
)

type Boat struct {
	Navigation Nav
	Power      Electrical
}

func NewBoat() *Boat {
	b := Boat{}
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

func NewNav() *Nav {
	n := Nav{}
	return &n
}

func (n *Nav) Marshal() *[]byte {
	encoded, _ := json.Marshal(n)
	return &encoded
}

type Waypoint struct {
	Name       string
	Type       WaypointType
	Coordinate Point
}

type WaypointType int

const (
	NorthCardinalBuoy WaypointType = iota
	SouthCardinalBuoy
	EastCardinalBuoy
	WestCardinalBuoy
	FairwayBuoy
)

type Muxable interface {
	Marshal() *[]byte
}

func checkGracefull(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

//-------------------------------
// Ingest json from serial port
//-------------------------------

func Ingest(s string, message chan Muxable) {

	config := &serial.Config{
		Name: s,
		Baud: 115200,
	}
	arduino, err := serial.OpenPort(config)
	check(err)
	defer arduino.Close()

	reader := bufio.NewReader(arduino)
	var token []byte

	for {
		token, _, err = reader.ReadLine()
		check(err)
		currentBoat := NewNav()
		err = json.Unmarshal(token, currentBoat)
		checkGracefull(err)
		if err == nil {
			message <- currentBoat
		}

	}
}

//--------------------------------------------
// Ingest data from gpsd running on local host
//--------------------------------------------

func IngestGPSD(message chan Muxable) {

	var gps *gpsd.Session
	var err error

	if gps, err = gpsd.Dial(gpsd.DefaultAddress); err != nil {
		panic(fmt.Sprintf("Failed to connect to GPSD: ", err))
	}

	gps.AddFilter("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		currentBoat := NewBoat()
		currentBoat.Navigation.Position[0] = tpv.Lat
		currentBoat.Navigation.Position[1] = tpv.Lon
		currentBoat.Navigation.SpeedOverGround = float32(tpv.Speed)
		currentBoat.Navigation.CourseOverGround = float32(tpv.Track)

		message <- currentBoat
	})

	done := gps.Watch()
	<-done

}
