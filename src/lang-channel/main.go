package main

/*
Channel Allow different concurrent parts of the program to communicate
*/

import (
	"fmt"
	"strconv"
)

func main() {
	runSum()
	buyItems()
}

func runSum() {
	fmt.Println("In runSum...")
	input := []int{1, 2, 3, 4, 5, 6}
	c := make(chan int)
	go sum(input[:len(input)/2], c)
	go sum(input[len(input)/2:], c)
	x, y := <-c, <-c
	fmt.Printf("x: %v, y: %v, x+y: %v\n", x, y, x+y)
}
func sum(input []int, c chan int) {
	fmt.Printf("Sum up %v\n", input)
	sum := 0
	for _, v := range input {
		sum += v
	}
	c <- sum
	fmt.Printf("result %v pushed to channel\n", sum)
}

func buyItems() {
	fmt.Println("In buyItems...")
	out := make(chan string)
	in := make(chan string)
	go buy(in, out, 1)
	go buy(in, out, 2)
	go buy(in, out, 3)

	/* there is no guarantee as to which goroutine will
	accept which input, or which goroutine will return an output first */

	in <- "item 1"
	in <- "item 2"
	in <- "item 3"

	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

//<-chan string means you can only put stuff into the channel
func buy(in <-chan string, out chan<- string, id int) {
	item := <-in
	fmt.Println("Buyer " + strconv.Itoa(id) + " purchased " + item)
	out <- item + " purchased by buyer " + strconv.Itoa(id)
}
