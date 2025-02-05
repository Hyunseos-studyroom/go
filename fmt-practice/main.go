package main

import "fmt"

type User struct {
	name string
	age  int
	job  string
}

var user []User

func main() {
	for {
		fmt.Printf("하실 것을 적으세연\n1. 유저 보기\n2. 유저 생성\n3. 유저 업데이트\n4. 유저 삭제\n그외 프로그램 종료\n선택 : ")
		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		switch choice {
		case 1:
			fmt.Println("-----------------------------------------------------------------------------")
			if len(user) > 0 {
				for i := 0; i < len(user); i++ {
					fmt.Printf("%d. %s (%d, %s)\n", i+1, user[i].name, user[i].age, user[i].job)
				}
			} else {
				fmt.Println("유저가 존재하지 않아요.")
			}
			fmt.Println("-----------------------------------------------------------------------------")
		case 2:
			fmt.Println("-----------------------------------------------------------------------------")
			var newUser User
			fmt.Print("유저의 이름: ")
			fmt.Scan(&newUser.name)
			fmt.Print("유저의 나이: ")
			fmt.Scan(&newUser.age)
			fmt.Print("유저의 직업: ")
			fmt.Scan(&newUser.job)
			user = append(user, newUser)
			fmt.Println("유저가 생성되었어요.")
			fmt.Println(user)
			fmt.Println("-----------------------------------------------------------------------------")
		case 3:
			fmt.Println("-----------------------------------------------------------------------------")
			if len(user) > 0 {
				for i := 0; i < len(user); i++ {
					fmt.Printf("%d. %s (%d, %s)\n", i+1, user[i].name, user[i].age, user[i].job)
				}
				fmt.Println("어떤 유저를 지우실래요? 번호를 적어주세요")
				var index int
				fmt.Scan(&index)
				if index <= len(user) && index > 0 {

				} else {
					fmt.Println("유저가 존재하지 않아요.")
				}
			} else {
				fmt.Println("유저가 존재하지 않아요.")
			}
			fmt.Println("-----------------------------------------------------------------------------")
		case 4:
			fmt.Println("4")
		default:
			fmt.Println("끝!")
			return
		}
	}
}
