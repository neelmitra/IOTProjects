// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command label uses the Vision API's label detection capabilities to find a label
// based on an image's content.
//
//     go run visionlabel.go <path-to-image>
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
	"strings"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	m "github.com/neelmitra/IOTProjects/MQ"
	"encoding/json"
)



var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// run submits a label request on a single image by given file.
func run(file string) (error,string) {
	ctx := context.Background()

	// Authenticate to generate a vision service.
	client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
	if err != nil {
		return err,"Exception"
	}
	service, err := vision.New(client)
	if err != nil {
		return err,"Exception"
	}

	// Read the image.
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err,"Exception"
	}

	// Construct a label request, encoding the image in base64.
	req := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},
		Features: []*vision.Feature{{Type: "LABEL_DETECTION"}},
	}
	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}
	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		return err,"Exception"
	}

	// Parse annotations from responses
	if annotations := res.Responses[0].LabelAnnotations; len(annotations) > 0 {
		label := annotations[0].Description
		fmt.Printf("Found label: %s for %s\n", label, file)
		return nil, label
	}
	fmt.Printf("Not found label: %s\n", file)
	return nil,"Not found"
}


// create a struct to notify the threats
type Detect struct {
	Threat  string `json:"threat"`
	Type   string   `json:"label"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	threats := []string{
		"gun",
		"baseball",
		"person",
		"bat",
		"zombie",
	}

	err ,label := run(args[0])

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	//certs path initialisation
	rootCA := "/home/root/certs/rootCA.pem"
	keyCert := "/home/root/certs/keycert.pem"
	privateKey := "/home/root/certs/privatekey.pem"

	//AWS MQTT Broker SSL Configuration
	tlsconfig := m.NewTLSConfig(rootCA,keyCert,privateKey)
	fmt.Println("TLSConfig initiation Completed")
	opts := MQTT.NewClientOptions()
	opts.AddBroker("ssl://AEV5KR4BW3J9L.iot.us-east-1.amazonaws.com:8883")
	opts.SetClientID("ConnectedGoDev").SetTLSConfig(tlsconfig)
	fmt.Println("Invoking Publish Handler method ")
	opts.SetDefaultPublishHandler(f)

	// Start the connection
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// _ = "breakpoint"
	fmt.Println("MQTT Connection established")
	c.Subscribe("/go-mqtt/opencv", 0, nil)

	//Analysing the match
	for j := 0; j < len(threats); j++ {
		if(strings.Contains(label,threats[j])) {
			fmt.Println("Match found !! Its a", threats[j])
			// Create an instance of the Threat struct.
			box := Detect{
				Threat:  label,
				Type:   threats[j],
			}
			// Create JSON from the instance data.
			// ... Ignore errors.
			b, _ := json.Marshal(box)
			// Convert bytes to string.
			s := string(b)
			if(c.IsConnected()) {
				fmt.Println("Publishing to MQ the message : ", s)
				c.Publish("/go-mqtt/opencv", 0, false, s)
				c.Disconnect(100)
			}
		}
	}
}
