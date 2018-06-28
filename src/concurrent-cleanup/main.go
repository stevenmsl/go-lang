/*
The "done" channel convention - If a goroutine is responsible for creating a goroutine,
it is also responsible for ensuring it can stop the goroutine.

*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	runSigs()
	//runNewRandStream()
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

//the or-channel pattern
func or(channels ...<-chan interface{}) <-chan interface{} {
	//termination criteria
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer fmt.Println("closing orDone chan")
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}

	}()
	fmt.Println("returning orDone chan")
	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer fmt.Printf("closing chan with duration %v\n", after)
		defer close(c)
		fmt.Printf("chan with duration %v kicked off...\n", after)
		time.Sleep(after)
	}()
	return c
}

func runSigs() {
	start := time.Now()

	/*
		So this is pretty much the same as create a certain number of go routines,
		each with a select statement that will receive from certain number of channels.
		go routine 1
		<-1 (1 sec)
		<-2
		<-3
		<-orDone2

		go routine 2
		<-4
		<-5
		<-orDone (passing down from the first call)
		<-orDone2

		Once 1 is closed the go routine 1 will proceed to close orDone,
		which in turns will trigger the select case <-orDon in the go routine 2. As a result,
		the go routine 2 will proceed to close the orDone2.

		2 to 5 will however close at a much later time after the orDone is closed. Will not this cause a problem?
	*/
	<-or2(
		sig(2*time.Second),
		sig(1*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)

	//Will hit this line in roughly one second. And the chan with 1 second duration will be closed.
	fmt.Printf("done after %v\n", time.Since(start))
	//Wait a bit longer to see all other 4 channels got closed.
	time.Sleep(8 * time.Second)

}

//This function is implemented to help understand how the ‘or’ function works.
//It only takes 5 channels so the process can be illustrated and explained more easily.
func or2(chan1, chan2, chan3, chan4, chan5 <-chan interface{}) <-chan interface{} {
	orDone := make(chan interface{})
	orDone2 := make(chan interface{})
	go func() {
		defer fmt.Println("closing orDone in go routine 1")
		defer close(orDone)
		//only one of these cases will be triggered depending on which one received the message first
		//The routine will then proceed to close orDone chan.
		select {
		case <-chan1:
			fmt.Println("Received from chan1 in go routine 1")
		case <-chan2:
			fmt.Println("Received from chan2 in go routine 1")
		case <-chan3:
			fmt.Println("Received from chan3 in go routine 1")
		case <-orDone2:
			fmt.Println("Received from orDone2 in go routine 1")
		}
	}()

	go func() {
		defer fmt.Println("closing orDone2 in go routine 2")
		defer close(orDone2)
		select {
		case <-chan4:
			fmt.Println("Received from chan4 in go routine 2")
		case <-chan5:
			fmt.Println("Received from chan5 in go routine 2")
		case <-orDone:
			fmt.Println("Received from orDone in go routine 2")
		case <-orDone2:
			fmt.Println("Received from orDone2 in go routine 2")
		}
	}()

	return orDone
}
