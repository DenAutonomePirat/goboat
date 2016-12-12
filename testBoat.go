package main

import (
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
)

func main() {
	b := boat.NewBoat()
	fmt.Println(b.Marshal())
}
