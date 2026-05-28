package main
import ("fmt")

// - can be declared inside or outside of a function
// - cannot be declared using the := syntax
// - cannot be reassigned a new value


const PI = 3.14

func main() {
	// typed 
	const a int = 42
	fmt.Println(a)
	
	// untyped
	const b = 3.14
	fmt.Println(b)
	
  	fmt.Println(PI)
}