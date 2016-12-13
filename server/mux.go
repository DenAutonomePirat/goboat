package server

import (
	//"flag"
	//"fmt"
	"log"
	"net/http"
	//"net/url"
	//"os"
	//"os/signal"
	//"time"
	"github.com/denautonomepirat/goboat/boat"
	"github.com/gorilla/websocket"
)

type Mux struct {
	connections map[*Conn]bool

	Broadcast chan boat.Muxable

	Recieve chan []byte

	register chan *Conn

	unregister chan *Conn
}

func NewMux() *Mux {

	mux := Mux{connections: make(map[*Conn]bool),
		Broadcast:  make(chan boat.Muxable),
		Recieve:    make(chan []byte),
		register:   make(chan *Conn),
		unregister: make(chan *Conn),
	}
	go mux.loop()

	return &mux
}

func (m *Mux) loop() {
	var conn *Conn
	var muxable boat.Muxable
	var msg *[]byte
	for {
		select {
		case conn = <-m.register:
			//register new connection
			m.connections[conn] = true
			log.Printf("Client registered: %s, %d total.", conn.User, len(m.connections))

		case conn = <-m.unregister:
			//remove connection
			delete(m.connections, conn)
			close(conn.Output)
			log.Printf("Client unregistered: %s, %d total.", conn.User, len(m.connections))

		case muxable = <-m.Broadcast:
			msg = muxable.Marshal()
			for conn := range m.connections {
				conn.Output <- *msg
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func (m *Mux) Handle(w http.ResponseWriter, r *http.Request, u string) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Could not upgrade http request: %s", err.Error())
		return
	}
	conn := NewConn(m, ws, u) //set user pointer
	m.register <- conn
}
