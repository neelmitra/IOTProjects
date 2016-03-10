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
		gobot.On(sensorl.Event("data"), func(data string) {
			fmt.Println("current temp (c): ", sensort.Temperature())
			fmt.Println("sensorl", data)
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
