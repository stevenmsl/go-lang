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
	runPipelineS()
	//runPipelineBP()
}

func runPipelineS() {
	ints := []int{1, 2, 3, 4}
	for _, v := range ints {
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
