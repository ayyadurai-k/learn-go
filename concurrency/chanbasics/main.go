package main

import "fmt"

// func main(){
// 	ch := make(chan int)

// 	go func(){
// 		ch <- 42
// 	}()

// 	x := <- ch
// 	fmt.Println("Got : ",x)
// }



// func main(){
// 	ch := make(chan int)
// 	go func ()  {
// 		for i:=1;i<=5;i++ {
// 			ch <- i
// 		}
// 		close(ch)
// 	}()

// 	for n:= range ch{
// 		fmt.Println("got :",n)
// 	}

// }


func main(){
	ch := make(chan string)

	go func(){ ch <- "hello from goroutine A" }()
	go func(){ ch <- "hello from goroutine B" }()

	fmt.Println(<-ch)
	fmt.Println(<-ch)

}