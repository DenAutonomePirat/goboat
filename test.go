package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Morning")

	ingestChannel := make(chan *Boat)
	//currentBoat := NewBoat()

	go Ingest("/dev/ttyACM1", ingestChannel)

	go func() {
		for {

			update := <-ingestChannel
			fmt.Println(update.Heading)
		}
	}()
	for {
		time.Sleep(time.Second)
	}

}
