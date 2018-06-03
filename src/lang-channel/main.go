package main

/*
Channel Allow different concurrent parts of the program to communicate
*/

import (
	"fmt"
	"strconv"
)

func main() {
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
