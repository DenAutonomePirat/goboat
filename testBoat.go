package main

import (
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

func main() {
	influxConfig := client.UDPConfig{
		Addr: "10.0.0.1:8089",
	}
	// Make client
	c, err := client.NewUDPClient(influxConfig)
	if err != nil {
		panic(err.Error())
	}

	b := boat.NewBoat()
	for {
		b.Power.Amperes = 1.2001
		b.Power.Volts = 38.0
		b.Power.JoulesTotal = 175.34235
		b.Power.TimeStamp = time.Now()
		fmt.Printf("%s\n", *b.Marshal())

		// Create a new point batch
		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Precision: "ms",
		})

		bp.AddPoint(b.Influx())

		// Write the batch
		c.Write(bp)
		time.Sleep(time.Millisecond * 250)
	}

}
