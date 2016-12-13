package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	rednet "github.com/denautonomepirat/goboat/udp"
	"log"
)

type Server struct {
	configuration *Configuration
	db            *Store
	udp           *rednet.UdpServer
	web           *Web
}

func NewServer() *Server {
	s := Server{
		configuration: NewConfiguration(),
		db:            NewStore(),
		udp:           rednet.NewUdpServer("10001"),
	}
	s.web = NewWeb(s.db)

	return &s
}

func (s *Server) Listen() {
	fmt.Println("Morning")
	flag.Parse()
	log.SetFlags(0)

	s.configuration.Start.Coordinate[0] = 56.72161
	s.configuration.Start.Coordinate[1] = 8.21222
	s.configuration.Start.Name = "start"

	s.configuration.Finish.Coordinate[0] = 56.96487
	s.configuration.Finish.Coordinate[1] = 10.36663
	s.configuration.Finish.Name = "finish"

	s.configuration.WaypointsAllowed = 4

	s.configuration.DefaultLegDistanceInMeters = 500

	user := NewUser()
	user.UserName = "Thomas"
	user.SetPassword("password")
	s.db.AddUser(user)

	go func() {

		for {
			select {
			case b := <-s.udp.Recieve:
				s.web.mux.Broadcast <- b

			case msg := <-s.web.mux.Recieve:
				var c map[string]interface{}
				json.Unmarshal(msg, &c)

				if c["class"] == "user" {
					fmt.Printf("User %s changed waypoint no. %s to:\nLatitude: \t%s\n", c["user"], c["wpt"], c["lat"])
				}

				if c["class"] == "Boat" {
					b := boat.NewBoat()
					fmt.Printf("Message recieved\n")
					json.Unmarshal(msg, &b)
					s.db.AddTrack(b)
					s.web.mux.Broadcast <- b
				}
			}

		}
	}()
	s.web.ListenAndServe(s.configuration)
}
