package udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/denautonomepirat/goboat/boat"
)

type UdpServer struct {
	Recieve chan *boat.Boat
	service string
}

func NewUdpServer(port string) *UdpServer {

	udpServer := UdpServer{
		Recieve: make(chan *boat.Boat),
		//		service: GetIpOfRednet() + ":" + port,
		service: "127.0.0.1:3030",
	}
	go udpServer.loop()
	return &udpServer
}

func (u *UdpServer) loop() {

	udpAddr, err := net.ResolveUDPAddr("udp4", u.service)

	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UDP server up and listening on %v\n", ln.LocalAddr())

	defer ln.Close()

	for {
		handleUDPConnection(ln, u.Recieve)
	}

}

func handleUDPConnection(conn *net.UDPConn, r chan *boat.Boat) {

	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	if err != nil {
		log.Fatal(err)
	}

	b := boat.NewBoat()
	err = json.Unmarshal(buffer[:n], &b)

	if err != nil {
		log.Println(err)
		log.Println(string(buffer[:n]))
		_, err = conn.WriteToUDP([]byte("nack"), addr)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	r <- b
	_, err = conn.WriteToUDP([]byte("ack"), addr)
	if err != nil {
		log.Fatal(err)
	}
}

// Nasty shit here
func GetIpOfRednet() string {
	addr, _ := net.InterfaceByName("rednet")
	ip := fmt.Sprint(addr.Addrs())
	ip = ip[1:]
	for p := range ip {
		if ip[p] == '/' {
			ip = ip[0:p]
			break
		}
	}
	return ip
}
