package main

import "fmt"

type User struct {
	personName string
}

func updateName(user *User,name string){
	user.personName = name
}

func main(){
	user := User{
		personName: "Ayyadurai",
	}

	updateName(&user,"Ayyadurai K")

	fmt.Println("User name :",user.personName)


}