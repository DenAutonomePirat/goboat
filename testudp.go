package main

import (
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"github.com/denautonomepirat/goboat/udp"
	"time"
)

func main() {

	s := udp.NewUdpServer("10001")
	go func() {

		for {
			r := boat.NewBoat()
			r = <-s.Recieve
			fmt.Printf("Er det mig %v\n", r.Power.Amperes)
		}
	}()

	c := udp.NewUdpClient("10.0.0.11", "10001")
	t := boat.NewBoat()
	t.Power.Amperes = 20
	for {

		c.Send <- t
		time.Sleep(time.Second)
	}
}
