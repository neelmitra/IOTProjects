package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	board := edison.NewEdisonAdaptor("board")
	sensorl := gpio.NewGroveLightSensorDriver(board, "sensor", "0")
	sensort := gpio.NewGroveTemperatureSensorDriver(board, "sensor", "1")

	work := func() {
		gobot.On(sensorl.Event("data"), func(data interface{}) {
			fmt.Println("current temp (c): ", sensort.Temperature())
			fmt.Println("current light : ", sensorl.Event("data"))
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensorl, sensort},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}