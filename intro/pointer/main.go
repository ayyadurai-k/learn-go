package main

import "fmt"

func update(num *int){
	*num = 200
}

func main() {
	var a int = 10

	var ptr *int = &a
	var zero_ptr *int

	fmt.Println("Ptr : ",ptr)

	fmt.Println("Value in pointer :",*ptr)

	// modify the value through
	*ptr = 20

	fmt.Println("Modified value :",a)

	fmt.Println("zero pointer : ",zero_ptr)

	update(&a)

	fmt.Println("updated value :",a)

	

}