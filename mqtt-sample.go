package main

import "io/ioutil"

import "fmt"
import "crypto/tls"
import "crypto/x509"
import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	fmt.Println("Importing RootCA file")
	pemCerts, err := ioutil.ReadFile("c:/samplecerts/rootCA.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	fmt.Println("Importing Client certs")
	cert, err := tls.LoadX509KeyPair("c:/samplecerts/keycert.pem", "c:/samplecerts/privatekey.pem")
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
	sensorl := gpio.NewGroveLightSensorDriver(board, "sensor", "0")
	sensort := gpio.NewGroveTemperatureSensorDriver(board, "sensor", "1")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			fmt.Println("current temp (c): ", sensort.Temperature())
			fmt.Println("current light : ", sensorl.Event("data"))
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensorl,sensort},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()

	jsonData := `
		{
		"msg": "Neel Mitra from Go"
		}
		`
	fmt.Println("Message published to the topic")
	c.Publish("/go-mqtt/sample", 0, false, jsonData)


	c.Disconnect(250)
}
