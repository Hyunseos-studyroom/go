package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := 0; i < 5; i++ {
		fmt.Printf("Producing: %d\n", i)
		ch <- i
		time.Sleep(time.Second)
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for value := range ch {
		fmt.Printf("Consuming: %d\n", value)
	}
}

func main() {
	ch := make(chan int)

	go producer(ch)
	consumer(ch)
}
