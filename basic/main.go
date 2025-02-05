package main

import "fmt"

func main() {
	golang := "golang"
	golangPointer := &golang

	fmt.Println("원본 고랭: ", golang)
	fmt.Println("포인터가 가리키는 주소: ", golangPointer)
	fmt.Println("포인터를 따라간 고랭 데이터: ", *golangPointer)

	*golangPointer = "Hello Golang"
	fmt.Println("바뀐 고랭 원본 데이터:", golang)
}
