package server

import (
	//"flag"
	"log"
	//"net/url"
	//"os"
	//"os/signal"
	//"time"
	"github.com/denautonomepirat/goboat/boat"
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
			log.Printf("Client registered: %p, %d total.", conn, len(m.connections))

		case conn = <-m.unregister:
			//remove connection
			delete(m.connections, conn)
			close(conn.Output)
			log.Printf("Client unregistered: %p, %d total.", conn, len(m.connections))

		case muxable = <-m.Broadcast:
			msg = muxable.Marshal()
			for conn := range m.connections {
				conn.Output <- *msg
			}
		}
	}
}
