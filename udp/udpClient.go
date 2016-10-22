package udp

import (
	"github.com/denautonomepirat/goboat/boat"
	"log"
	"net"
)

type UdpClient struct {
	Send        chan boat.Muxable
	service     string
	PacketCount int32
	PacketLoss  int32
}

func NewUdpClient(name, port string) *UdpClient {
	udpClient := UdpClient{
		Send:    make(chan boat.Muxable),
		service: name + ":" + port,
	}

	go udpClient.loop()

	return &udpClient
}

func (u *UdpClient) loop() {

	udpAddr, err := net.ResolveUDPAddr("udp4", u.service)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Established connection to %s \n", u.service)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	var msg *[]byte

	for {
		muxable := <-u.Send
		msg = muxable.Marshal()
		_, err = conn.Write(*msg)

		if err != nil {
			log.Println(err)
		}

	}

}
