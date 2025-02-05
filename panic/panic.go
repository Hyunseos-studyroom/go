package main

import "fmt"

func errorFunc() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("패닉 이유:", r)
		}
	}()

	panic("패닉")
}

func main() {
	fmt.Println("프로그램 시작")
	errorFunc()
	fmt.Println("프로그램 진행...")
}
