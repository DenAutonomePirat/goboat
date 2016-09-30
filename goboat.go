package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "46.101.213.117:8080", "http service address")
var port = flag.String("port", "/dev/ttyACM0", "ingest port /dev/someport")

func main() {
	fmt.Println("Morning")

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})

	ingestChannel := make(chan boat.Muxable)

	go boat.Ingest(*port, ingestChannel)
	go boat.IngestGPSD(ingestChannel)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			_ = message
			//log.Printf("recv: %s", message)
		}
	}()

	var muxable boat.Muxable
	var msg *[]byte

	for {
		select {

		case muxable = <-ingestChannel:
			msg = muxable.Marshal()
			printable := new(bytes.Buffer)
			json.Indent(printable, *msg, "", "    ")
			fmt.Println(printable)
			err := c.WriteMessage(websocket.TextMessage, *msg)
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
