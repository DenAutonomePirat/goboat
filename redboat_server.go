package main

import (
	"github.com/denautonomepirat/goboat/server"
)

func main() {
	go server.Listen()
	select {}
}
