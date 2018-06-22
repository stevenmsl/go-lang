package main

/*
Lexical confinement involves using lexical scope to expose only the correct
data and concurrency primitives for multiple concurrent processes to use.
*/

import (
	"fmt"
)

func main() {
	results := chanOwner()
	consumer(results)
}

//The returned channel is one way – for receiving only.
//This confines usage of the channel within the consume function to only reads.
func chanOwner() <-chan int {
	//Instantiate the channel within the lexical scope of the chanOwner function
	results := make(chan int, 5)
	go func() {
		defer close(results)
		for i := 0; i <= 5; i++ {
			results <- i
		}
	}()
	return results
}
func consumer(results <-chan int) {
	for result := range results {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}
