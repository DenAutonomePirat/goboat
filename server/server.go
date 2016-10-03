package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denautonomepirat/goboat/boat"
	"log"
)

func Listen() {
	fmt.Println("Morning")
	flag.Parse()
	log.SetFlags(0)

	web := NewWeb()

	//users := NewStore()
	//defer users.db.Close()
	go func() {
		b := boat.NewBoat()

		for {

			msg := <-web.mux.Recieve
			json.Unmarshal(msg, &b)
			fmt.Printf("%s\n", b.Marshal())
			web.mux.Broadcast <- b
		}
	}()
	web.ListenAndServe()
}
