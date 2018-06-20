package main

/*
Channel Allow different concurrent parts of the program to communicate
*/

import (
	"fmt"
	"strconv"
	"time"
)

func main() {

	bufferedChanResolveDeadlock()

	//The The following will throw a fatal error: all goroutines are asleep - deadlock!
	//Failed to continue.
	//bufferedChanDeadlock()

	bufferedChannel()
	runSum()
	buyItems()
}

func bufferedChannel() {
	fmt.Println("In bufferedChannel ...")
	bc := make(chan int, 2) //If you provide the buffer length in the second parameter this channel become a buffered channel
	bc <- 1
	go delayInput(bc, 2)
	fmt.Println(<-bc)
	fmt.Println(<-bc) //This line will wait for 3 seconds until the integer 2 is sent to the channel
}

func bufferedChanDeadlock() {
	fmt.Println("In bufferedChanDeadlock ...")
	bc := make(chan int, 2)
	bc <- 1
	bc <- 2
	bc <- 3 //Overfill the buffer without letting the code a chance to read/remove a value from the channel
	fmt.Println(<-bc)
}

func bufferedChanResolveDeadlock() {
	fmt.Println("In bufferedChanResolveDeadlock...")
	bc := make(chan int, 2)
	bc <- 1
	bc <- 2

	add3 := func() {
		bc <- 3
	}

	go add3() //The goroutine is being called before the channel is being emptied, but that is fine, the goroutine will wait until the channel is available.
	//It doesn’t block the main thread.

	fmt.Println(<-bc)
	fmt.Println(<-bc)
	fmt.Println(<-bc)

}

func delayInput(bc chan<- int, input int) {
	fmt.Println("In delayInput...Sleep for 3 secs")
	time.Sleep(3 * time.Second)
	bc <- input
	fmt.Printf("%v sent to channel\n", input)
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
