package server

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/denautonomepirat/goboat/boat"
	rednet "github.com/denautonomepirat/goboat/udp"
	"log"
)

func Listen() {
	fmt.Println("Morning")
	flag.Parse()
	log.SetFlags(0)

	conf := NewConfiguration()

	conf.Start.Coordinate[0] = 56.72161
	conf.Start.Coordinate[1] = 8.21222
	conf.Start.Name = "start"

	conf.Finish.Coordinate[0] = 56.96487
	conf.Finish.Coordinate[1] = 10.36663
	conf.Finish.Name = "finish"

	conf.WaypointsAllowed = 3

	conf.DefaultLegDistanceInMeters = 500

	db := NewStore()
	u := NewUser()
	u.UserName = "thomas"
	u.SetPassword("password")
	db.AddUser(u)
	udp := rednet.NewUdpServer("10001")
	web := NewWeb(db)

	go func() {

		for {
			select {
			case b := <-udp.Recieve:
				web.mux.Broadcast <- b

			case msg := <-web.mux.Recieve:
				var c map[string]interface{}
				json.Unmarshal(msg, &c)

				if c["class"] == "User" {
					u := NewUser()
					json.Unmarshal(msg, &u)
					fmt.Printf("The user %s send data\n", u.UserName)
				}

				if c["class"] == "Boat" {
					b := boat.NewBoat()
					fmt.Printf("Message recieved\n")
					json.Unmarshal(msg, &b)
					db.AddTrack(b)
					web.mux.Broadcast <- b
				}
			}

		}
	}()
	web.ListenAndServe(conf)
}
