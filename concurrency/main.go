package main

import (
	// "time"
	"math/rand"
	"sync"
)

func main(){
	wg := sync.WaitGroup{}
	orders := generateOrders(10)


	wg.Add(2)

	go func() {
		defer wg.Done()
		processOrders(orders)
	}()
	
	go func() {
		defer wg.Done()
		updateOrderStatuses(orders)
	}()
	wg.Wait()


	reportOrderStatuses(orders)

}

type Order struct {
	ID int
	Status string
}


func generateOrders(count int) ([]*Order) {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{ID: i + 1, Status: "pending"}
	}
	return orders
}

func processOrders(orders []*Order) {
	for _, order := range orders {
		// time.Sleep(1 * time.Second) // Simulate processing time
		order.Status = "processed"
		println("Processed order ID:", order.ID)
	}	
}

func updateOrderStatuses(orders []*Order) {
	for _, order := range orders {
		// time.Sleep(500 * time.Millisecond)
		status := []string{"pending", "processing", "shipped", "delivered"}[rand.Intn(4)]
		order.Status = status
		println("Updated order ID:", order.ID, "to status:", order.Status)
	}
}

func reportOrderStatuses(orders []*Order) {
	statusCount := make(map[string]int)
	for _, order := range orders {
		statusCount[order.Status]++
		// time.Sleep(200 * time.Millisecond) // Simulate time taken to count statuses
	}
	println("Order status report:")
	for status, count := range statusCount {
		println(status, ":", count)
		// time.Sleep(100 * time.Millisecond) // Simulate time taken to print each status
	}
}