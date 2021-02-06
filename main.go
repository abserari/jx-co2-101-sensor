package main

import (
	"bufio"
	"log"
	"time"

	"github.com/tarm/serial"
)

var Logmode = DebugLog

// pi3 should open uart and communicate with device: /dev/ttyAMA0 | /dev/serial0
func main() {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	sensor := NewCO2Sensor(s)

	for {

		data, _, err := sensor.ReadLine()
		if err != nil {
			LogError(err)
		}
		LogInfo(string(data))
	}

}

type LogMode int

const DebugLog LogMode = 0
const InfoLog LogMode = 1
const ErrLog LogMode = 2

func LogDebug(v ...interface{}) {
	if Logmode > 0 {
		return
	}
	log.Println(v...)
}

func LogInfo(v ...interface{}) {
	if Logmode > 1 {
		return
	}
	log.Println(v...)
}
func LogError(v ...interface{}) {
	if Logmode > 2 {
		return
	}
	log.Println(v...)
}

// CRC not implement from: https://github.com/sigurn/crc8
func CRC(b []byte) {
	return
}

var ActiveModeChange = []byte{0xff, 0x01, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0xfc}
var QueryModeChange = []byte{0xff, 0x01, 0x03, 0x02, 0x00, 0x00, 0x00, 0x00, 0xfb}
var QueryPPM = []byte{0xff, 0x01, 0x03, 0x03, 0x01, 0x00, 0x00, 0x00, 0xf9}
var Correct = []byte{0xff, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xfc}

type CO2Sensor struct {
	*serial.Port
	*bufio.Reader
}

func NewCO2Sensor(s *serial.Port) *CO2Sensor {
	sensor := &CO2Sensor{
		s,
		bufio.NewReader(s),
	}
	return sensor
}

func (s *CO2Sensor) SendCorrect() error {
	n, err := s.Write(Correct)
	if err != nil {
		LogError(err)
		return err
	}
	if n == len(Correct) {
		LogInfo("send Correct successful")
	}
	data, _, _ := s.ReadLine()
	LogInfo(string(data))

	return nil
}

func (s *CO2Sensor) SendActiveModeChange() error {
	n, err := s.Write(ActiveModeChange)
	if err != nil {
		LogError(err)
		return err
	}
	if n == len(ActiveModeChange) {
		LogInfo("send ActiveModeChange successful")
	}
	data, _, _ := s.ReadLine()
	LogInfo(string(data))

	return nil
}

func (s *CO2Sensor) SendQueryModeChange() error {
	n, err := s.Write(QueryModeChange)
	if err != nil {
		LogError(err)
		return err
	}
	if n == len(QueryModeChange) {
		LogInfo("send QueryModeChange successful")
	}

	return nil
}

func (s *CO2Sensor) SendQuery() (int, error) {
	n, err := s.Write(QueryPPM)
	if err != nil {
		LogError(err)
		return 0, err
	}
	if n == len(QueryPPM) {
		LogInfo("send QueryPPM successful")
	}

	return 0, nil
}
