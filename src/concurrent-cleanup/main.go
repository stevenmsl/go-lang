package main

import (
	"fmt"
)

func main() {
	runDoWorkLeak()
}

func doWorkLeak(strings <-chan string) <-chan interface{} {
	fmt.Println("In doWorkLeak ...")
	completed := make(chan interface{})
	go func() {
		fmt.Println("In go routine inside doWorkLeak ...") //This line will not get executed if the strings chan is nil.
		defer fmt.Println("doWork exited")                 //will never get executed
		defer close(completed)                             //will never get executed

		for s := range strings {
			fmt.Println(s)
		}
		fmt.Println("In doWorkLeak, done printing ...")
	}()
	return completed
}

func runDoWorkLeak() {
	fmt.Println("In runDoWorkLeak ...")
	/*
		strings := make(chan string)
		doWorkLeak(strings)
		strings <- "test"
	*/
	doWorkLeak(nil)
}
