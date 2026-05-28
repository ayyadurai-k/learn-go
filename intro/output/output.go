package main
import ("fmt")

func main() {
	var word1,word2 string = "Hello","World"

	// PRINT

	fmt.Print("Print() : \n\n")

	fmt.Print(word1,word2) // print with default format
	fmt.Print("\n") // print a new line
	fmt.Print(word1,"\n") // print with a new line
	fmt.Print(word2,"\n") // print with a new line
	fmt.Print(word1," ",word2) // print with a space

	// PRINTLN

	fmt.Print("\n\n Println() : \n\n")


	fmt.Println(word1)
	fmt.Println(word2)

	// PRINTF

	fmt.Printf("Value : %v Type: %T",word1,word1)


}