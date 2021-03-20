package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/denautonomepirat/goboat/boat"
	"github.com/denautonomepirat/goboat/udp"
)

var addr = flag.String("addr", "46.101.213.117:8080", "http service address")
var port = flag.String("port", "/dev/ttyACM0", "ingest port /dev/someport")

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("Morning")

	flag.Parse()
	log.SetFlags(0)

	ingestChannel := make(chan *boat.Boat)
	//broadcatsChannel := make(chan boat.Muxable)

	go boat.Ingest(*port, ingestChannel)
	go boat.IngestGPSD(ingestChannel)
	//go boat.Connect(broadcatsChannel, interrupt, addr)
	c := udp.NewUdpClient("10.0.0.1", "10001")

	for {
		select {
		case b := <-ingestChannel:
			c.Send <- b
		case i := <-interrupt:
			log.Println("Stopping goroutines")
			interrupt <- i
			select {
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}
