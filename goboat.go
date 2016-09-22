package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/stratoberry/go-gpsd"
	"github.com/tarm/serial"
	"log"
)

type Boat struct {
	Id          int32   `json:"id",omitempty`
	Rudder      int32   `json:"rudder",omitempty`
	Depth       float32 `json:"depth",omitempty`
	MainSail    int32   `json:"mainsail",omitempty`
	Jib         int32   `json:"jib",omitempty`
	Volts       float32 `json:"volts",omitempty`
	Amperes     float32 `json:"amperes,omitempty"`
	JoulesTotal float32 `json:"joules_total,omitempty"`
	joulesTrip  float32 `json:"joules_trip,omitempty"`
	Heading     float32 `json:"heading,omitempty"`
	Pitch       float32 `json:"pitch,omitempty"`
	Roll        float32 `json:"roll,omitempty"`
	Lat         float32 `json:"lat,omitempty"`
	Lon         float32 `json:"lon,omitempty"`
}

func (b *Boat) Marshal() *[]byte {
	encoded, _ := json.Marshal(b)
	return &encoded
}

func NewBoat() *Boat {
	b := Boat{}
	return &b
}

func Ingest(s string, message chan Muxable) {

	config := &serial.Config{
		Name: s,
		Baud: 115200,
	}
	arduino, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer arduino.Close()

	reader := bufio.NewReader(arduino)
	var token []byte

	for {
		token, _, err = reader.ReadLine()
		if err != nil {
			panic(err)
		}
		currentBoat := NewBoat()
		err = json.Unmarshal(token, currentBoat)
		if err == nil {
			message <- currentBoat
		}
	}
}

func IngestGPSD(message chan Muxable) {

	var gps *gpsd.Session
	var err error

	if gps, err = gpsd.Dial(gpsd.DefaultAddress); err != nil {
		panic(fmt.Sprintf("Failed to connect to GPSD: ", err))
	}

	gps.AddFilter("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		currentBoat := NewBoat()
		currentBoat.Lat = float32(tpv.Lat)
		currentBoat.Lon = float32(tpv.Lon)
		message <- currentBoat
	})

	done := gps.Watch()
	<-done

}
