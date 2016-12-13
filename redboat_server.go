package main

import (
	"github.com/denautonomepirat/goboat/server"
)

func main() {
	master := server.NewServer()
	go master.Listen()
	select {}
}
