package main

import (
	"encoding/json"
	"fmt"
)

//Box create astruct to hold values
type Box struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Color  string `json:"color"`
	Open   bool   `json:"open"`
}

func main() {
	// Initiatlise some variables
	var wid = 35
	var hi = 40
	// Create an instance of the Box struct.
	box := Box{
		Width:  wid,
		Height: hi,
		Color:  "blue",
		Open:   false,
	}
	// Create JSON from the instance data.
	// ... Ignore errors.
	b, _ := json.Marshal(box)
	// Convert bytes to string.
	s := string(b)
	fmt.Println(s)
}
