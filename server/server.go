package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	rednet "github.com/denautonomepirat/goboat/udp"
	"github.com/kellydunn/golang-geo"
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

	s.configuration.Start = geo.NewPoint(56.72161, 8.21222)

	s.configuration.Finish = geo.NewPoint(56.96487, 10.36663)
	fmt.Printf("Total distance: %4f\n", s.configuration.Start.GreatCircleDistance(s.configuration.Finish))

	s.configuration.WaypointsAllowed = 4

	s.configuration.DefaultLegDistanceInMeters = 500

	fmt.Printf("%s\n", *s.configuration.Marshal())

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

				if c["class"] == "skipper" {
					var sk Skipper
					if err := json.Unmarshal(msg, &sk); err != nil {
						panic(err)
					}
					fmt.Println(sk.Latlng.GreatCircleDistance(s.configuration.Finish))
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
