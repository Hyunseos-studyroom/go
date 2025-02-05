package main

import "fmt"

func main() {
	defer fmt.Println("defer statement")

	fmt.Println("Hello, world!")
	return
}
