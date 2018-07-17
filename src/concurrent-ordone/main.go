/* or-done channel

Prevent goroutine leaks but keep the code simple – hide the verbosity of implementing done channel

*/

package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	myChan := make(chan interface{})
	go func() {
		myChan <- "A"
		myChan <- "B"
	}()

	go func() {
		for val := range orDone(done, myChan) {
			s := fmt.Sprintf("%v", val)
			fmt.Println(s)
		}
	}()

	time.Sleep(time.Second * 5)
}

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valueStream <- v:
				case <-done:
				}

			}
		}
	}()

	return valueStream
}
