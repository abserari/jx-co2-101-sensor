package main

import (
	"log"
	"time"

	. "github.com/abserari/jx-co2-101-sensor"
	"github.com/tarm/serial"
)

// pi3 should open uart and communicate with device: /dev/ttyAMA0 | /dev/serial0
func main() {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	sensor := NewCO2Sensor(s)
	sensor.SendActiveModeChange()

	// go func() {
	// 	for {
	// 		sensor.SendQuery()
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()
	for {
		data, _, err := sensor.ReadLine()
		if err != nil {
			LogError(err)
		}
		value, err := ReadData(data)
		LogInfo(value)
	}

}
