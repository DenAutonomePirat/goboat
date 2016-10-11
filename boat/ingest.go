package boat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/stratoberry/go-gpsd"
	"github.com/tarm/serial"
)

//-------------------------------
// Ingest json from serial port
//-------------------------------

func Ingest(s string, message chan *Boat) {

	config := &serial.Config{
		Name: s,
		Baud: 115200,
	}
	arduino, err := serial.OpenPort(config)
	Check(err)
	defer arduino.Close()

	reader := bufio.NewReader(arduino)
	var token []byte

	for {
		token, _, err = reader.ReadLine()
		Check(err)
		currentBoat := NewBoat()
		err = json.Unmarshal(token, currentBoat)
		CheckGracefull(err)
		if err == nil {
			message <- currentBoat
		}

	}
}

//--------------------------------------------
// Ingest data from gpsd running on local host
//--------------------------------------------

func IngestGPSD(message chan *Boat) {

	var gps *gpsd.Session
	var err error

	if gps, err = gpsd.Dial(gpsd.DefaultAddress); err != nil {
		panic(fmt.Sprintf("Failed to connect to GPSD: ", err))
	}

	gps.AddFilter("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		currentBoat := NewBoat()
		currentBoat.Navigation.Position[0] = tpv.Lat
		currentBoat.Navigation.Position[1] = tpv.Lon
		currentBoat.Navigation.SpeedOverGround = float32(tpv.Speed)
		currentBoat.Navigation.CourseOverGround = float32(tpv.Track)

		message <- currentBoat
	})

	done := gps.Watch()
	<-done

}
