package main

import (
	"fmt"
	"time"
)

func main() {
	defer fmt.Println("Goodbye cruel world")

	fmt.Println("Hello World")
	fmt.Println("Golang is a synchronous programming language")

	go fmt.Println("But it has threads too")

	fmt.Println("Sleeping 1 second")
	time.Sleep(1 * time.Second)
}
