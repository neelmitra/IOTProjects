package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/mqtt"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://test.mosquitto.org:1883", "pinger")

	work := func() {
		mqttAdaptor.On("hello", func(data []byte) {
			fmt.Println("hello")
		})
		mqttAdaptor.On("hola", func(data []byte) {
			fmt.Println("hola")
		})
		data := []byte("o")
		gobot.Every(1*time.Second, func() {
			mqttAdaptor.Publish("hello", data)
		})
		gobot.Every(5*time.Second, func() {
			mqttAdaptor.Publish("hola", data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
