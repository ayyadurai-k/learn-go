package main

import (
	"fmt"
	"time"
)

func main()  {
	// sayHello()
    fmt.Println("Start:", time.Now())

    // Launch 3 goroutines
    go printNumbers()
    go printLetters()

    // Wait for all to complete
    time.Sleep(2 * time.Second)

    fmt.Println("End:", time.Now())
}


func printNumbers(){
	for i:=1;i<=10;i++{
		fmt.Println(i)
		time.Sleep(500 * time.Millisecond)

	}
}

func printLetters() {
    for i := 0; i < 3; i++ {
        fmt.Println("Letter:", string(rune('A'+i)))
        time.Sleep(500 * time.Millisecond)
    }
}

func sayHello(){
	fmt.Println("Hello")
}