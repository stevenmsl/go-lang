/*
tee-channel - split values coming in from a channel to send them off into
two separate areas of the codebase.
*/

package main

import (
	"fmt"
	"stage"
)

func main() {
	done := make(chan interface{})
	defer close(done)
	out1, out2 := tee(done, stage.Take(done, stage.Repeat(done, 1, 2), 4))

	for val := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val, <-out2)
	}

}

func tee(
	done <-chan interface{},
	in <-chan interface{},
) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for value := range orDone(done, in) { //The iteration over “in” channel cannot continue until both out1 and out2 have been written to
			var out1, out2 = out1, out2 //Use local version by shadowing the original variables - Need to find out why this needed
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- value:
					out1 = nil //further writes will block
				case out2 <- value:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
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
