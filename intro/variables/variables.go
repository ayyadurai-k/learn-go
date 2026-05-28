package main
import ("fmt")

// TYPE OF DECLARATION 
/* 
var : 
	- can be used outside and inside the function 
	- Declaration and initialization can be done in separately.
	- can be used with or without initial value. If not initialized, it will take the default value of the type.
	- can be used with the type or without the type. If the type is not specified, it will be inferred from the initial value.

Short variable declaration (:=) :
	- can be used only inside the function 
	- Declaration and initialization must be done in the same line.
	- cannot be used without initial value. It must be initialized at the time of declaration.
*/

func main(){
	// int
	var a int = 10 
	var b int = 20
	c := a+b // can be only used inside the function 
	fmt.Println(c)

	// string + without initial value 
	var name string 	
	name = "Ayyadurai"
	fmt.Println("Welcome",name)

	// multiple variable declaration
	var x, y, z int = 1, 2, 3
	fmt.Println(x,y,z)

	// multiple variable declaration with short variable declaration
	p, q, r := 4, 5, 6
	fmt.Println(p,q,r)

	// default value of the type
	var d int 
	var e string 
	fmt.Println("Default value of int:",d)
	fmt.Println("Default value of string:",e)


}