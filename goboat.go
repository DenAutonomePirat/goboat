package main

import (
	"bufio"
	"encoding/json"
	"github.com/tarm/serial"
	"log"
)

type Boat struct {
	Id          int32   `json:"id"`
	Rudder      int32   `json:"rudder"`
	Depth       float32 `json:"depth"`
	MainSail    int32   `json:"mainsail"`
	Jib         int32   `json:"jib"`
	Volts       float32 `json:"volts`
	Amperes     float32 `json:"amperes"`
	JoulesTotal float32 `json:"joules_total"`
	joulesTrip  float32 `json:"joules_trip"`
	Heading     float32 `json:"heading"`
	Pitch       float32 `json:"pitch,omitempty"`
	Roll        float32 `json:"roll"`
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
