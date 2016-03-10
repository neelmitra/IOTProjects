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

	lsensor := gpio.A0
	tsensor := gpio.A1
	//mysensor[0] := gpio.NewGroveLightSensorDriver(board, "sensor", "0")
	//mysensor[1] := gpio.NewGroveTemperatureSensorDriver(board, "sensor", "1")

	work := func() {
		gobot.On {
			fmt.Println("current temp (c): ", tsensor)
			fmt.Println("current light intensity", lsensor)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{lsensor,tsensor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
