package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"log"

	"github.com/denautonomepirat/goboat/boat"
	rednet "github.com/denautonomepirat/goboat/udp"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
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

	currentBoat := boat.Boat{
		Class: "Boat",
		Navigation: boat.Nav{
			Position: conf.Start.Coordinate,
		},
	}

	goal := &s2.Point{}

	db := NewStore()
	/*
		u := NewUser()
		u.UserName = ""
		u.SetPassword("")
		db.AddUser(u)
	*/
	udp := rednet.NewUdpServer("10001")
	web := NewWeb(db)
	skippers := make(map[string]*Skipper)

	decisions := make(chan *boat.Command)

	go func() {

		for {
			select {
			case uname := <-web.skipperReg:
				delete(skippers, uname)

			case c := <-decisions:
				s := &Skipper{}
				if c.Waypoint.Name == "first" {

					s.first[0] = c.Waypoint.Coordinate.Lat
					s.first[1] = c.Waypoint.Coordinate.Lng

					nv := s2.LatLngFromDegrees(s.first[0], s.first[1])

					d := &boat.Ding{
						Class:    "ding",
						User:     c.User,
						Position: &nv,
					}
					web.mux.Broadcast <- d

					skippers[c.User] = s

					newp := s2.Point{}
					for _, s = range skippers {
						p := s2.PointFromLatLng(s2.LatLngFromDegrees(s.first[0], s.first[1]))
						newp.Vector = newp.Add(p.Vector)
					}
					newp.Normalize()
					newlatlng := s2.LatLngFromPoint(newp)
					d = &boat.Ding{
						Class:    "dong",
						User:     "boat",
						Position: &newlatlng,
					}
					goal = &newp

					web.mux.Broadcast <- d
				}
			}

		}
	}()
	go func() {

		for {
			select {
			case b := <-udp.Recieve:
				web.mux.Broadcast <- b

			case msg := <-web.mux.Recieve:
				var c map[string]interface{}
				fmt.Printf("%s", msg.user)
				json.Unmarshal(msg.payload, &c)

				if c["class"] == "User" {
					u := NewUser()
					json.Unmarshal(msg.payload, &u)
					fmt.Printf("The user %s send data\n", u.UserName)
					fmt.Printf("payload %v", u)
				}
				if c["class"] == "command" {
					c := &boat.Command{}
					json.Unmarshal(msg.payload, c)
					c.User = msg.user
					decisions <- c
				}

				if c["class"] == "Boat" {
					b := boat.NewBoat()
					fmt.Printf("Message recieved\n")
					json.Unmarshal(msg.payload, &b)
					db.AddTrack(b)
					web.mux.Broadcast <- b
				}
			}

		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(skippers) > 0 {
					var speed s1.Angle = 0.00001
					p := s2.PointFromLatLng(s2.LatLngFromDegrees(currentBoat.Navigation.Position[0], currentBoat.Navigation.Position[1]))
					polyline := s2.Polyline{}
					polyline = append(polyline, p)
					polyline = append(polyline, *goal)
					dist := (speed / polyline.Length()).Radians()
					newpos, _ := polyline.Interpolate(dist)
					currentBoat.Navigation.Position[0] = s2.LatLngFromPoint(newpos).Lat.Degrees()
					currentBoat.Navigation.Position[1] = s2.LatLngFromPoint(newpos).Lng.Degrees()
					web.mux.Broadcast <- &currentBoat

				}
			}
		}
	}()

	web.ListenAndServe(conf)
}

type Skipper struct {
	first  boat.Point
	second boat.Point
}
