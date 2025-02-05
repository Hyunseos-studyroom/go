package main

import (
	"fmt"
	"time"
)

func main() {
	go Count("go")
	go Count("lang")
	time.Sleep(time.Second * 5)
}

func Count(person string) {
	for i := 0; i < 5; i++ {
		fmt.Println(person, i)
		time.Sleep(time.Second)
	}
}
