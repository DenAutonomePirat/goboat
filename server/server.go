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

	counter := 0

	go func() {

		for {

			msg := <-web.mux.Recieve
			var c map[string]interface{}
			json.Unmarshal(msg, &c)

			if c["class"] == "User" {
				u := NewUser()
				json.Unmarshal(msg, &u)
				fmt.Printf("User %d send data\n", u.Id)
			}
			if c["class"] == "Boat" {
				b := boat.NewBoat()
				counter++
				fmt.Printf("%d messages recieved\r", counter)
				json.Unmarshal(msg, &b)
				web.mux.Broadcast <- b
			}

		}
	}()
	web.ListenAndServe()
}
