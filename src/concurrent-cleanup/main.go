package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	runNewRandStream()
	//runNewRandStreamLeak()
	//runDoWork()
	//runDoWorkLeak()
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

func doWork(
	//allow caller to signal cancellation.
	//By convention, this signal is usually a read-only channel named done and
	//is the first parameter.
	done <-chan interface{},
	strings <-chan string) <-chan interface{} {
	terminated := make(chan interface{})
	fmt.Println("In doWork ...")
	go func() {
		defer fmt.Println("doWork exited")
		defer close(terminated)
		for {
			select {
			case s := <-strings:
				fmt.Println(s)
			case <-done:
				fmt.Println("In doWork cancellation received")
				return
			}
		}
	}()
	return terminated
}

func runDoWork() {
	fmt.Println("In runDoWork ...")
	done := make(chan interface{})
	terminated := doWork(done, nil)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	<-terminated
	fmt.Println("runDoWork done.")
}

//in this function there is no way of telling the producer it can stop.
func newRandStreamLeak() <-chan int {
	randStream := make(chan int)
	go func() {
		defer fmt.Println("newRandStreamLeak closure exited.") //will never be executed
		defer close(randStream)                                //will never be executed
		for {
			//the go routine blocks once the channel is no longer being read from
			randStream <- rand.Int()
		}
	}()
	return randStream
}

func runNewRandStreamLeak() {
	fmt.Println("In runNewRandStreamLeak...")
	randStream := newRandStreamLeak()
	fmt.Println(" 3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

//provide the producer goroutine with a channel informing it to exit
func newRandStream(done <-chan interface{}) <-chan int {
	randStream := make(chan int)
	go func() {
		defer fmt.Println("newRandStream closure exited.") //will never be executed
		defer close(randStream)                            //will never be executed
		for {
			select {
			case randStream <- rand.Int():
			case <-done:
				return
			}
		}
	}()
	return randStream
}

func runNewRandStream() {
	fmt.Println("In runNewRandStream...")
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println(" 3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	//simulate there is more work to do
	time.Sleep(1 * time.Second)
	fmt.Println("runNewRandStream done")
}
