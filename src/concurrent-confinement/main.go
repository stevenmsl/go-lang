package main

/*
Lexical confinement involves using lexical scope to expose only the correct
data and concurrency primitives for multiple concurrent processes to use.
*/

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	printString("golang")
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

func printData(wg *sync.WaitGroup, data []byte) {
	defer wg.Done()
	var buff bytes.Buffer
	for _, b := range data {
		fmt.Fprintf(&buff, "%c", b)
	}
	fmt.Println(buff.String())
}

func printString(value string) {
	fmt.Println("In printString...")
	var wg sync.WaitGroup //WaitGroup waits for a collection of goroutines to finish
	wg.Add(2)
	data := []byte(value)
	//Constrain the go routines to only work on certain the part of the slice.
	//Using lexical scope to avoid the need of implementing synchronization
	//the program will be simpler, and the performance would be better.
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait() //"Wait" blocks until the WaitGroup counter is zero.
}
