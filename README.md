# jx-co2-101-sensor

树莓派 Raspberry model 3B+ 和 JX-CO2-101 传感器. [更多文档](https://www.yuque.com/abser/solutions/telmsy)

```go
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
```

![](https://cdn.nlark.com/yuque/0/2021/png/176280/1612499197873-7c58f77f-f8b9-40a5-8061-1aeb4f9ba4e8.png)
![](https://cdn.nlark.com/yuque/0/2021/png/176280/1612496135723-4376f1f6-188e-45a3-b9de-beea7f2761e3.png?x-oss-process=image%2Fresize%2Cw_738)
