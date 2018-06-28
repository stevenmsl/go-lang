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
		fmt.Printf("chan with duration %v sleeping...\n", after)
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
		<-1
		<-2
		<-3
		<-orDone2

		go routine 2
		<-4
		<-5
		<-orDone1 (pass down from the parent)
		<-orDone2

		Once 1 is closed this will end the go routine 1, which orDone2 will be closed immediately,
		which will in turn end the go routine 2 and close the orDone1.

	*/
	<-or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)

	//Will hit this line in roughly one second. And the chan with 1 second duration will be closed.
	fmt.Printf("done after %v\n", time.Since(start))
	//Wait a bit longer to see all other 4 channels got closed.
	time.Sleep(8 * time.Second)

}
