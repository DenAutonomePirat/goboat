// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Morning")
	flag.Parse()
	log.SetFlags(0)

	web := NewWeb()
	//measurements := make(chan Muxable)
	//users := NewStore()
	//defer users.db.Close()
	go func() {
		b := NewBoat()

		for {

			msg := <-web.mux.Recieve
			json.Unmarshal(msg, &b)
			web.mux.Broadcast <- b
		}
	}()
	web.ListenAndServe()
}

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
