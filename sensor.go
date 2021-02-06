package sensor

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

var Logmode = DebugLog

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
func CRC(b []byte, check []byte) bool {
	return true
}

func ReadData(raw []byte) (int, error) {
	// bytes data            4444       ppm
	//           space space 4444 space ppm
	strs := strings.Split(string(raw), " ")
	LogDebug("receive raw", raw, "and split to ", strs)

	// convert 4444 string to int
	return strconv.Atoi(strs[2])
}

var ActiveModeChange = []byte{0xff, 0x05, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0xf8}
var QueryModeChange = []byte{0xff, 0x05, 0x03, 0x02, 0x00, 0x00, 0x00, 0x00, 0xf7}
var QueryPPM = []byte{0xff, 0x05, 0x03, 0x03, 0x01, 0x00, 0x00, 0x00, 0xf5}
var Correct = []byte{0xff, 0x05, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8}
var MODBUS_RTU = []byte{0x05, 0x03, 0x00, 0x05, 0x00, 0x01, 0x94, 0x07}

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

func (s *CO2Sensor) SendMODBUS_RTU() error {
	n, err := s.Write(MODBUS_RTU)
	if err != nil {
		LogError(err)
		return err
	}
	if n == len(MODBUS_RTU) {
		LogInfo("send MODBUS_RTU successful")
	}
	data, _, _ := s.ReadLine()
	LogInfo(string(data))

	return nil
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
		LogDebug("😀 send ActiveModeChange successful")
	}

	// read response
	var b = make([]byte, 8)
	n, err = s.Port.Read(b)
	if err != nil || n != 8 {
		LogError("😥 set ActiveModeChange failed", err, n, b)
		return err
	}

	var check = false
	// read check bit and crc8 check
	// check if crc[0] -1 != b[0] - b[1]- b[2]- b[3] - b[4] - b[5] - b[6] - b[7]
	var crc = make([]byte, 1)
	n, err = s.Port.Read(crc)
	if err != nil || n != 1 {
		LogError("😥 set ActiveModeChange failed", err, n, crc)
		return err
	}

	// need to check
	check = true
	if !check {
		LogError("😥 set ActiveModeChange failed", err, n, crc)
		return err
	}

	LogInfo("😀 set QueryModeChange successful")

	return nil
}

func (s *CO2Sensor) SendQueryModeChange() error {
	n, err := s.Write(QueryModeChange)
	if err != nil {
		LogError(err)
		return err
	}
	if n == len(QueryModeChange) {
		LogDebug("😀 send QueryModeChange successful")
	}

	// read response
	var b = make([]byte, 8)
	n, err = s.Port.Read(b)
	if err != nil || n != 8 {
		LogError("😥 set QueryModeChange failed", err, n, b)
		return err
	}

	var check bool = false
	// read check bit and crc8 check
	// check if crc[0] -1 != b[0] - b[1]- b[2]- b[3] - b[4] - b[5] - b[6] - b[7]
	var crc = make([]byte, 1)
	n, err = s.Port.Read(crc)
	if err != nil || n != 1 {
		LogError("😥 set QueryModeChange failed", err, n, crc)
		return err
	}

	// need to check
	check = true
	if !check {
		LogError("😥 set QueryModeChange failed", err, n, crc)
		return err
	}

	LogInfo("😀 set QueryModeChange successful")

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
