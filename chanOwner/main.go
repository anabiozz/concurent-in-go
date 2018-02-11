package main

import (
	"fmt"
)

func main() {
	chanOwer := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i < 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwer()
	for result := range resultStream {
		fmt.Printf("received: %d\n", result)
	}

	fmt.Println("done")
}
