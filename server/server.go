package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"log"
)

func Listen() {
	fmt.Println("Morning")
	flag.Parse()
	log.SetFlags(0)

	game := NewGameSetup()

	game.Start.Coordinate[0] = 56.72161
	game.Start.Coordinate[1] = 8.21222
	game.Start.Name = "start"

	game.Finish.Coordinate[0] = 56.96487
	game.Finish.Coordinate[1] = 10.36663
	game.Finish.Name = "finish"

	game.WaypointsAllowed = 3

	game.DefaultLegDistanceInMeters = 500

	web := NewWeb()
	db := NewStore()

	u := NewUser()
	u.SetPassword("jeg elsker at lege med kokmputere")
	fmt.Println(u.HashedPassword)

	fmt.Println(u.CheckPassword("jeg elsker at lege med komputere"))
	go func() {

		for {

			msg := <-web.mux.Recieve
			var c map[string]interface{}
			json.Unmarshal(msg, &c)

			if c["class"] == "User" {
				u := NewUser()
				json.Unmarshal(msg, &u)
				fmt.Printf("The user %s send data\n", u.Name)
			}

			if c["class"] == "Boat" {
				b := boat.NewBoat()
				fmt.Printf("Message recieved\n")
				json.Unmarshal(msg, &b)
				db.AddTrack(b)
				web.mux.Broadcast <- b
			}

		}
	}()
	web.ListenAndServe(game)
}
