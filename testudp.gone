package main

import (
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"github.com/denautonomepirat/goboat/udp"
	"time"
)

func main() {

	s := udp.NewUdpServer("10005")
	go func() {

		for {
			r := boat.NewBoat()
			r = <-s.Recieve
			fmt.Printf("Er det mig %v\n", r.Navigation.HeadingMagnetic)
		}
	}()

	c := udp.NewUdpClient("10.0.0.11", "10001")
	t := boat.NewBoat()
	t.Navigation.HeadingMagnetic = 0.0
	for {

		c.Send <- t
		t.Navigation.HeadingMagnetic += 1.1
		time.Sleep(time.Millisecond * 100)
	}
}
