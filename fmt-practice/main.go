package main

import "fmt"

func main() {
	for {
		fmt.Printf("하실 것을 적으세연\n1. 유저 보기\n2.유저 생성\n3. 유저 업데이트\n4. 유저 삭제\n그외. 프로그램 종료\n선택 : ")
		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		switch choice {
		case 1:
			fmt.Println("1")
		case 2:
			fmt.Println("2")
		case 3:
			fmt.Println("3")
		case 4:
			fmt.Println("4")
		default:
			fmt.Println("끝!")
			return
		}
	}
}
