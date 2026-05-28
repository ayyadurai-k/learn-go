package main

import "fmt"

func main()  {
	var numbers = [3]int{10,20,30}
	names := [2]string{"Hi","hello"}

	fmt.Println("Numbers : ",numbers)
	fmt.Println("Names : ",names)

	fmt.Println("Numbers of 2 :",numbers[2])

	numbers[2] = 200
	fmt.Println("Updated Numbers : ",numbers)

	var arr1 = [3]int{} // uninitialized
	var arr2 = [4]int{2,2} // partially uninitialized
	var arr3 = [2]int{10,20} // fully initialized

	fmt.Println("Array 1 : ",arr1)
	fmt.Println("Array 2 : ",arr2)
	fmt.Println("Array 3 : ",arr3)

	fmt.Println("Length of array :",len(arr1))

}