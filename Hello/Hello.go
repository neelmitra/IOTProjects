package main

import "fmt"

func test() {
	fmt.Println("This is a freaking test")
}

func main() {
	var a int = 5
	fmt.Println("Hello world")
	fmt.Println("The value tested here is ", a)
	test()
}
