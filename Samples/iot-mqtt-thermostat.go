package main

import "io/ioutil"

import "fmt"
import "crypto/tls"
import "crypto/x509"
import (
	"encoding/json"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"math"
)

//NewTLSConfig SSL config for MQTT
func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	fmt.Println("Importing RootCA file")
	pemCerts, err := ioutil.ReadFile("/home/root/certs/rootCA.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	fmt.Println("Importing Client certs")
	cert, err := tls.LoadX509KeyPair("/home/root/certs/keycert.pem", "/home/root/certs/privatekey.pem")
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	//fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	//AWS MQTT Broker SSL Configuration
	tlsconfig := NewTLSConfig()
	fmt.Println("TLSConfig initiation Completed")
	opts := MQTT.NewClientOptions()
	opts.AddBroker("ssl://AEV5KR4BW3J9L.iot.us-east-1.amazonaws.com:8883")
	opts.SetClientID("iot-sample").SetTLSConfig(tlsconfig)
	fmt.Println("Invoking Publish Handler method ")
	opts.SetDefaultPublishHandler(f)

	// Start the connection
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// _ = "breakpoint"
	fmt.Println("MQTT Connection established")

	c.Subscribe("/go-mqtt/sample", 0, nil)

	// Gobot initiation
	gbot := gobot.NewGobot()
	board := edison.NewEdisonAdaptor("board")
	sensort := gpio.NewGroveTemperatureSensorDriver(board, "tempsensor", "0")
	screen := i2c.NewGroveLcdDriver(board, "screen")

	// Struct to hold sensor data
	type Sensord struct {
		Temp float64 `json:"temp"`
	}

	work := func() {
		screen.Write("Thermostat is On !!")
		screen.SetRGB(255, 0, 0)

		gobot.Every(5*time.Second, func() {
			fmt.Println("current temp (c): ", math.Floor(sensort.Temperature()))

			// LCD showing the temperature
			screen.Clear()
			screen.Home()
			screen.SetRGB(0,255,0)
			//screen.SetCustomChar(0, i2c.CustomLCDChars["smiley"])


			//Update the struct with sensor data, Jsonify & Convert to string
			res1Z := Sensord{
				Temp: math.Floor(sensort.Temperature()),
			}

			jData, err := json.Marshal(res1Z)
			if err != nil {
				fmt.Println(err)
				return
			}

			s := string(jData)

			// Writes to the LCD Screen & publish to AWS MQTT Broker.

			screen.Write(s)
			gobot.Every(2*time.Second, func() {
				screen.Scroll(false)
			})

			fmt.Println("The json data to be published in IOT topic is", s)
			c.Publish("/go-mqtt/sample", 0, false, s)
			//c.Disconnect(250)

			//screen.Home()
			//<-time.After(2 * time.Second)
			//screen.SetRGB(0, 0, 255)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensort,screen},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()

}
