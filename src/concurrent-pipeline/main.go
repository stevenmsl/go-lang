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
)

func main() {
	runPipelineC()
	//runPipelineS()
	//runPipelineBP()
}

func runPipelineC() {
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	//each stage of the pipeline is executing concurrently
	pipeline := multiplyC(done, addC(done, multiplyC(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}

//t huhe generator function converts a discrete set of values into a stream of data on a channel
func generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
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
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
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
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i + additive:
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
