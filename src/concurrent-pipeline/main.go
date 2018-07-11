/*
A pipeline consists of stages. Each stage will take the data in,
perform a transformation on it, and then send the data back out.

Stage
- A stage consumes and returns the same type.
- A stage must be reified by the language so that it may be passed around.
  (Need to do more research on this. (nomad, high order functions))

Batch processing
- Each stage will process chunks of data at once

Stream processing
- Each stage will receive and emit one element at a time



*/

package main

import (
	"fmt"
	"math/rand"
	"stage"
	"testing"
	"time"
)

func main() {
	runPipelineString()
	//runPipelineTakeFn()
	//runPipelineTake()
	//runPipelineC()
	//runPipelineS()
	//runPipelineBP()
}

func runBenchmark() {

}

//BenchmarkGeneric ...
func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range stage.ToString(done, take(done, repeat(done, "a"), b.N)) {

	}
}

func runPipelineString() {
	done := make(chan interface{})
	defer func() {
		close(done)
		fmt.Println("runPipelineTakeFn closed")
		time.Sleep(3 * time.Second)
	}()

	var message string
	for token := range stage.ToString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...\n", message)

}

func runPipelineTakeFn() {
	done := make(chan interface{})
	defer func() {
		close(done)
		fmt.Println("runPipelineTakeFn closed")
		time.Sleep(3 * time.Second) //Wait a bit so you can see the message printed when a stage is closed.
	}()
	rand := func() interface{} { return rand.Int() }
	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Printf("%v\n", num)
	}

}

func repeatFn(
	done <-chan interface{},
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer func() {
			fmt.Println("stage repeatFn closed")
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

func runPipelineTake() {
	done := make(chan interface{})
	defer func() {
		close(done)
	}()

	for num := range take(done, repeat(done, 1, 2), 10) {
		fmt.Printf("%v ", num)
	}
	fmt.Println()
}

//stage: generate a stream of data
func repeat(
	done <-chan interface{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer fmt.Println("repeat closed")
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

//stage: limit the pipeline
func take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer func() {
			fmt.Println("stage take closed")
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

/*
When you look at the output, it looks quite random:
Generated the number 1
1 + 2 = 2 done by stage Added by 2
3 x 3 = 9 done by stage Multiplied by 3
1 x 1 = 1 done by stage Multiplied by 1

In the above example, it appears to be in iteration 0 stage “Added by 2” finished before stage “Multiplied by 1”,
which it’s not possible as the later stage “Added by 2” is depending on the outcome of the previous stage “Multiplied by 1”.
I think the reason for this inconsistence is due to the time needed between you pushed the result to the channel
and you printed the result.
Don’t forget all stages are executing concurrently, and the next stage probably already finished
by the time you finished printing the result in the current stage.

*/
func runPipelineC() {
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	//each stage of the pipeline is executing concurrently
	pipeline := multiplyC(done, addC(done, multiplyC(done, intStream, 1), 2), 3)
	index := 0
	//If the range expression is a channel, at most one iteration variable is permitted.
	for v := range pipeline {
		fmt.Printf("Iteration %v Completed. Produced the number %v\n", index, v)
		index++
	}
}

//t huhe generator function converts a discrete set of values into a stream of data on a channel
func generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer fmt.Println("Generator closed")
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
				fmt.Printf("Generated the number %v\n", i)
			}
		}
	}()
	return intStream
}

func multiplyC(
	done <-chan interface{},
	intStream <-chan int,
	multiplier int,
) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer fmt.Printf("stage Multiplied by %v closed\n", multiplier)
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

func addC(
	done <-chan interface{},
	intStream <-chan int,
	additive int,
) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer fmt.Printf("stage Added by %v closed\n", additive)
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

func runPipelineS() {
	ints := []int{1, 2, 3, 4}
	for _, v := range ints {
		//We instantiate the pipeline for each iteration: 4 times in this case
		fmt.Println(multiplyS(addS(multiplyS(v, 2), 1), 2))
	}
}

//stream processing stage
func multiplyS(value, multiplier int) int {
	return value * multiplier
}

//stream processing
func addS(value, additive int) int {
	return value + additive
}

func runPipelineBP() {
	ints := []int{1, 2, 3, 4}
	for _, v := range addBP(multiplyBP(ints, 2), 1) {
		fmt.Println(v)
	}
}

//Batch processing stage
func multiplyBP(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}
	return multipliedValues
}

//Batch processing stage
func addBP(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {
		addedValues[i] = v + additive
	}
	return addedValues
}
