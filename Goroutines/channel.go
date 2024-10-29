package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan string)
	people := [2]string{"John", "Jessi"}
	for _, person := range people {
		go isSexy(person, c)
	}
	fmt.Println(<-c)
	fmt.Println(<-c)
}

func isSexy(person string, c chan string) {
	time.Sleep(time.Second * 10)
	fmt.Println(person)
	c <- person + "is sexy"
}
