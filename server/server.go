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
	udp := rednet.NewUdpServer("10001")
	web := NewWeb(db)

	user := NewUser()
	user.UserName = "Thomas"
	user.SetPassword("password")
	db.AddUser(user)

	go func() {

		for {
			select {
			case b := <-udp.Recieve:
				web.mux.Broadcast <- b

			case msg := <-web.mux.Recieve:
				var c map[string]interface{}
				json.Unmarshal(msg, &c)

				if c["class"] == "user" {
					fmt.Printf("User %s changed the %s waypoint to:\nLatitude: \t%s\n", c["user"], c["wpt"], c["lat"])
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
