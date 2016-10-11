package boat

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"time"
)

func Connect(dataChannel chan Muxable, interrupt chan os.Signal, addr *string) {

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})

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
		}
	}()

	var muxable Muxable
	var msg *[]byte

	for {
		select {

		case muxable = <-dataChannel:
			msg = muxable.Marshal()
			err := c.WriteMessage(websocket.TextMessage, *msg)
			if err != nil {
				log.Println("write:", err)
				return
			}

		case i := <-interrupt:
			log.Println("Disconneting websocket")
			interrupt <- i
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
