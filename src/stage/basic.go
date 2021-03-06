package stage

import (
	"fmt"
	"sync"
)

func isPrime(value int) bool {
	if value <= 0 {
		return false
	}
	if value == 2 || value == 3 {
		return true
	}

	for i := 2; i < value-1; i++ {
		if value%i == 0 {
			return false
		}
	}

	return true
}

//FanIn stage - multiplexing or joining together multiple streams of data into a single stream
func FanIn(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))

	//spinning up one goroutine for each incoming channel
	for _, c := range channels {
		go multiplex(c)
	}
	//spinning up another goroutine to wait for everyone
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

//PrimeFinder ...
func PrimeFinder(
	done <-chan interface{},
	valueStream <-chan int,
) <-chan interface{} {
	primeStream := make(chan interface{})
	go func() {
		defer close(primeStream)
		for v := range valueStream {

			if !isPrime(v) {
				continue
			} else {
				select {
				case <-done:
					return
				case primeStream <- v:
				}
			}
		}

	}()
	return primeStream
}

//ToInt stage
func ToInt(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}

	}()
	return intStream
}

//ToString stage:type assertion, Change the data type of the channel from interface {} to string
func ToString(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}

	}()
	return stringStream
}

//RepeatFn ...
func RepeatFn(
	done <-chan interface{},
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer func() {
			//fmt.Println("stage repeatFn closed")
			close(valueStream)
		}()
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

//Repeat stage: generate a stream of data
func Repeat(
	done <-chan interface{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		//defer fmt.Println("repeat closed")
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

//RepeatS stage (typed) : generate a stream of data
func RepeatS(
	done <-chan interface{},
	values ...string,
) <-chan string {
	valueStream := make(chan string)
	go func() {
		//defer fmt.Println("repeat closed")
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

//Take stage: limit the pipeline
func Take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer func() {
			//fmt.Println("stage take closed")
			close(takeStream)
		}()

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

//TakeS stage: limit the pipeline
func TakeS(
	done <-chan interface{},
	valueStream <-chan string,
	num int,
) <-chan string {
	takeStream := make(chan string)
	go func() {
		defer func() {
			//fmt.Println("stage take closed")
			close(takeStream)
		}()

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

//Generator stage : converts a discrete set of values into a stream of data on a channel
func Generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		//defer fmt.Println("Generator closed")
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
				//fmt.Printf("Generated the number %v\n", i)
			}
		}
	}()
	return intStream
}

//MultiplyC stage
func MultiplyC(
	done <-chan interface{},
	intStream <-chan int,
	multiplier int,
) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		//defer fmt.Printf("stage Multiplied by %v closed\n", multiplier)
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
				fmt.Printf("%v x %v = %v done by stage Multiplied by %v \n",
					i, multiplier, i*multiplier, multiplier)
			}
		}
	}()
	return multipliedStream
}

//AddC stage
func AddC(
	done <-chan interface{},
	intStream <-chan int,
	additive int,
) <-chan int {
	addedStream := make(chan int)
	go func() {
		//defer fmt.Printf("stage Added by %v closed\n", additive)
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i + additive:
				fmt.Printf("%v + %v = %v done by stage Added by %v \n",
					i, additive, i*additive, additive)
			}
		}
	}()
	return addedStream
}

//MultiplyS stage : stream processing stage
func MultiplyS(value, multiplier int) int {
	return value * multiplier
}

//AddS stage: stream processing
func AddS(value, additive int) int {
	return value + additive
}

//MultiplyBP stage: Batch processing stage
func MultiplyBP(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}
	return multipliedValues
}

//AddBP stage: batch processing stage
func AddBP(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {
		addedValues[i] = v + additive
	}
	return addedValues
}
