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
