package main

import "fmt"

func main(){
	var s1 = []int{10,20,30} // slice literal
	s2 := []string{"Hi","hello"} // slice literal

	fmt.Println("Slice 1 : ",s1)
	fmt.Println("Slice 2 : ",s2)

	fmt.Println("Length of slice 1 : ",len(s1))
	fmt.Println("Capacity of slice 1 : ",cap(s1))

	fmt.Println("Element at index 2 : ",s1[2])

	s1[2] = 200
	fmt.Println("Updated Slice 1 : ",s1)
	
	var s3 = []int{}

	fmt.Println("Slice 3 : ",s3)

	// APPEND 

	s3 = append(s3,10)

	fmt.Println("Slice 3 after append : ",s3)

	s3 = append(s3,20,30)

	fmt.Println("Slice 3 after append : ",s3)

	s3 = append(s3,s1...)

	fmt.Println("Slice 3 after append : ",s3)


}